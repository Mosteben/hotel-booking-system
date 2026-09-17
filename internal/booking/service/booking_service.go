package service

import (
	"errors"
	"strings"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"gorm.io/gorm"
)

// userLister/hotelLister/roomLister are the minimal slices of
// UserRepository/HotelRepository/RoomRepository this service needs. The
// real repositories already satisfy these (structural typing), so this
// keeps the service testable without mocking their full interfaces.
type userLister interface {
	GetByIDs(ids []string) ([]userModel.User, error)
}

type hotelLister interface {
	GetByIDs(ids []uint) ([]hotelModel.Hotel, error)
	GetByID(id uint) (*hotelModel.Hotel, error)
}

type roomLister interface {
	GetByID(id uint) (*roomModel.Room, error)
	GetByIDs(ids []uint) ([]roomModel.Room, error)
}

type BookingService interface {
	CreateBooking(booking *model.Booking) error
	// GetAllBookings is admin/manager only - it returns the booking with
	// its user/hotel/room resolved server-side (see AdminBookingSummary),
	// unlike every other method here which stays on the raw model.
	GetAllBookings() ([]AdminBookingSummary, error)
	GetBookingByID(id uint, userID string, role string) (*model.Booking, error)
	GetBookingsByUserID(userID string) ([]model.Booking, error)
	UpdateBooking(id uint, userID string, booking *model.Booking) error
	DeleteBooking(id uint, userID string) error

	UpdateBookingStatus(id uint, status string) error

	IsRoomAvailable(
		roomID uint,
		checkIn time.Time,
		checkOut time.Time,
	) (bool, error)
}

type bookingService struct {
	repo      repository.BookingRepository
	roomRepo  roomLister
	userRepo  userLister
	hotelRepo hotelLister
	db        *gorm.DB
}

func NewBookingService(
	repo repository.BookingRepository,
	roomRepo roomLister,
	userRepo userLister,
	hotelRepo hotelLister,
	db *gorm.DB,
) BookingService {
	return &bookingService{
		repo:      repo,
		roomRepo:  roomRepo,
		userRepo:  userRepo,
		hotelRepo: hotelRepo,
		db:        db,
	}
}

// =========================
// Booking Snapshot
// =========================

// applyBookingSnapshot loads the room's hotel and stamps booking with a
// point-in-time copy of the room/hotel details (hotel_id, hotel_name,
// room_number, room_type, price_per_night). This is what lets a booking
// stay a meaningful historical record even after its room or hotel is
// later deleted - the raw RoomID FK can go stale, but these columns don't.
//
// If the request supplied a hotel_id (booking.HotelID != 0), it's
// cross-checked against the room's real hotel and rejected on mismatch -
// but never trusted as the value actually stored: every snapshot field is
// always overwritten from the server-loaded room/hotel below, regardless
// of what the client sent.
func (s *bookingService) applyBookingSnapshot(
	booking *model.Booking,
	room *roomModel.Room,
) error {

	if booking.HotelID != 0 && booking.HotelID != room.HotelID {
		return errors.New("room does not belong to the specified hotel")
	}

	hotel, err := s.hotelRepo.GetByID(room.HotelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("hotel not found")
		}

		return err
	}

	booking.HotelID = hotel.ID
	booking.HotelName = hotel.Name
	booking.RoomNumber = room.RoomNumber
	booking.RoomType = room.Type
	booking.PricePerNight = room.PricePerNight

	return nil
}

// =========================
// Create Booking
// =========================

func (s *bookingService) CreateBooking(
	booking *model.Booking,
) error {

	if err := validateBooking(booking); err != nil {
		return err
	}

	// Get room and validate existence, status, and capacity.
	room, err := s.roomRepo.GetByID(booking.RoomID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("room not found")
		}

		return err
	}

	// Room must be available.
	if room.Status != "available" {
		return errors.New("room is not available")
	}

	// Guests cannot exceed room capacity.
	if booking.Guests > room.Capacity {
		return errors.New(
			"number of guests exceeds room capacity",
		)
	}

	// Calculate number of nights.
	nights := int(
		booking.CheckOut.Sub(booking.CheckIn).Hours() / 24,
	)

	if nights <= 0 {
		return errors.New("invalid booking duration")
	}

	// Calculate total price on the server.
	booking.TotalPrice = room.PricePerNight * float64(nights)

	if err := s.applyBookingSnapshot(booking, room); err != nil {
		return err
	}

	// New bookings always start as pending.
	booking.Status = "pending"

	// The availability check and the insert must happen inside the same
	// transaction and be serialized per room, otherwise two concurrent
	// requests for the same room/dates can both pass the check before
	// either commits (TOCTOU race). A plain transaction is not enough on
	// its own: under Postgres' default READ COMMITTED isolation, the
	// availability check is a "no matching rows" query, so two concurrent
	// transactions can each see zero conflicts and both insert. A
	// session-scoped advisory lock keyed on the room ID closes that gap -
	// the second concurrent transaction blocks on the lock until the
	// first commits or rolls back, so its availability check always sees
	// whatever the first transaction just created. Postgres releases the
	// lock automatically at the end of the transaction.
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"SELECT pg_advisory_xact_lock(?)",
			int64(booking.RoomID),
		).Error; err != nil {
			return err
		}

		txRepo := s.repo.WithTx(tx)

		available, err := txRepo.IsRoomAvailable(
			booking.RoomID,
			booking.CheckIn,
			booking.CheckOut,
			0,
		)

		if err != nil {
			return err
		}

		if !available {
			return errors.New(
				"room is not available for the selected dates",
			)
		}

		return txRepo.Create(booking)
	})
}

// =========================
// Get All Bookings
// =========================

func (s *bookingService) GetAllBookings() (
	[]AdminBookingSummary,
	error,
) {
	bookings, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.resolveAdminBookingSummaries(bookings)
}

// =========================
// Get Booking By ID
// =========================

func (s *bookingService) GetBookingByID(
	id uint,
	userID string,
	role string,
) (*model.Booking, error) {

	if id == 0 {
		return nil, errors.New("invalid booking id")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	booking, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Admin and manager can view any booking.
	if role == "admin" || role == "manager" {
		return booking, nil
	}

	// Customers can only view their own bookings.
	if booking.UserID != userID {
		return nil, errors.New(
			"you are not allowed to view this booking",
		)
	}

	return booking, nil
}

// =========================
// Get My Bookings
// =========================

func (s *bookingService) GetBookingsByUserID(
	userID string,
) ([]model.Booking, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetByUserID(userID)
}

// =========================
// Update Booking
// =========================

func (s *bookingService) UpdateBooking(
	id uint,
	userID string,
	booking *model.Booking,
) error {

	if id == 0 {
		return errors.New("invalid booking id")
	}

	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	if err := validateBooking(booking); err != nil {
		return err
	}

	existingBooking, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Make sure the booking belongs to the logged-in user.
	if existingBooking.UserID != userID {
		return errors.New(
			"you are not allowed to update this booking",
		)
	}

	// Get room and validate existence, status, and capacity.
	room, err := s.roomRepo.GetByID(booking.RoomID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("room not found")
		}

		return err
	}

	// Room must be available.
	if room.Status != "available" {
		return errors.New("room is not available")
	}

	// Guests cannot exceed room capacity.
	if booking.Guests > room.Capacity {
		return errors.New(
			"number of guests exceeds room capacity",
		)
	}

	// Calculate number of nights.
	nights := int(
		booking.CheckOut.Sub(booking.CheckIn).Hours() / 24,
	)

	if nights <= 0 {
		return errors.New("invalid booking duration")
	}

	// Calculate total price on the server.
	totalPrice := room.PricePerNight * float64(nights)

	// Refresh the snapshot too - if the room changed, the old snapshot
	// (hotel_name/room_number/room_type/price_per_night) would otherwise
	// keep describing the previous room while room_id already points at
	// the new one.
	if err := s.applyBookingSnapshot(booking, room); err != nil {
		return err
	}

	oldRoomID := existingBooking.RoomID
	newRoomID := booking.RoomID

	// The availability check and the update must happen inside the same
	// transaction and be serialized per room, for the same reason as
	// CreateBooking: under Postgres' default READ COMMITTED isolation, a
	// plain transaction alone doesn't stop two concurrent updates from
	// each observing "no conflicts" before either commits. If the room is
	// being changed, both the old and new room are locked - always in
	// ascending room-ID order, regardless of which one is "old" or "new"
	// - so two concurrent updates touching an overlapping pair of rooms
	// always acquire their locks in the same order and can never deadlock
	// on each other.
	return s.db.Transaction(func(tx *gorm.DB) error {
		lockRoomIDs := []uint{oldRoomID}

		if newRoomID != oldRoomID {
			if newRoomID < oldRoomID {
				lockRoomIDs = []uint{newRoomID, oldRoomID}
			} else {
				lockRoomIDs = append(lockRoomIDs, newRoomID)
			}
		}

		for _, roomID := range lockRoomIDs {
			if err := tx.Exec(
				"SELECT pg_advisory_xact_lock(?)",
				int64(roomID),
			).Error; err != nil {
				return err
			}
		}

		txRepo := s.repo.WithTx(tx)

		// Check room availability.
		// Exclude the current booking so it doesn't
		// conflict with itself.
		available, err := txRepo.IsRoomAvailable(
			newRoomID,
			booking.CheckIn,
			booking.CheckOut,
			id,
		)

		if err != nil {
			return err
		}

		if !available {
			return errors.New(
				"room is not available for the selected dates",
			)
		}

		existingBooking.RoomID = booking.RoomID
		existingBooking.CheckIn = booking.CheckIn
		existingBooking.CheckOut = booking.CheckOut
		existingBooking.Guests = booking.Guests
		existingBooking.TotalPrice = totalPrice
		existingBooking.HotelID = booking.HotelID
		existingBooking.HotelName = booking.HotelName
		existingBooking.RoomNumber = booking.RoomNumber
		existingBooking.RoomType = booking.RoomType
		existingBooking.PricePerNight = booking.PricePerNight

		// Status is controlled by the system.
		// Customer cannot change booking status here.

		return txRepo.Update(existingBooking)
	})
}

// =========================
// Check Room Availability
// =========================

func (s *bookingService) IsRoomAvailable(
	roomID uint,
	checkIn time.Time,
	checkOut time.Time,
) (bool, error) {

	if roomID == 0 {
		return false, errors.New("invalid room id")
	}

	if checkIn.IsZero() {
		return false, errors.New("check-in date is required")
	}

	if checkOut.IsZero() {
		return false, errors.New("check-out date is required")
	}

	if !checkOut.After(checkIn) {
		return false, errors.New(
			"check-out must be after check-in",
		)
	}

	// Make sure the room exists.
	_, err := s.roomRepo.GetByID(roomID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("room not found")
		}

		return false, err
	}

	return s.repo.IsRoomAvailable(
		roomID,
		checkIn,
		checkOut,
		0,
	)
}

// =========================
// Cancel Booking
// =========================

func (s *bookingService) DeleteBooking(
	id uint,
	userID string,
) error {

	if id == 0 {
		return errors.New("invalid booking id")
	}

	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	booking, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Make sure the booking belongs to the logged-in user.
	if booking.UserID != userID {
		return errors.New(
			"you are not allowed to cancel this booking",
		)
	}

	// Booking is already cancelled.
	if booking.Status == "cancelled" {
		return errors.New("booking is already cancelled")
	}

	// Cancel instead of deleting the booking.
	booking.Status = "cancelled"

	return s.repo.Update(booking)
}

// =========================
// Update Booking Status
// =========================

func (s *bookingService) UpdateBookingStatus(
	id uint,
	status string,
) error {

	if id == 0 {
		return errors.New("invalid booking id")
	}

	status = strings.ToLower(strings.TrimSpace(status))

	// Admin/Manager can only set these statuses
	// through the status management endpoint.
	if status != "confirmed" && status != "cancelled" {
		return errors.New("invalid booking status")
	}

	booking, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Cancelled bookings cannot be changed again.
	if booking.Status == "cancelled" {
		return errors.New(
			"cancelled booking cannot be updated",
		)
	}

	// Booking already has this status.
	if booking.Status == status {
		return errors.New(
			"booking already has this status",
		)
	}

	// Only pending bookings can be confirmed.
	if status == "confirmed" &&
		booking.Status != "pending" {
		return errors.New(
			"only pending bookings can be confirmed",
		)
	}

	// Pending or confirmed bookings can be cancelled.
	if status == "cancelled" &&
		booking.Status != "pending" &&
		booking.Status != "confirmed" {
		return errors.New(
			"booking cannot be cancelled",
		)
	}

	booking.Status = status

	return s.repo.Update(booking)
}

// =========================
// Booking Validation
// =========================

func validateBooking(
	booking *model.Booking,
) error {

	if booking == nil {
		return errors.New("booking data is required")
	}

	if strings.TrimSpace(booking.UserID) == "" {
		return errors.New("user id is required")
	}

	if booking.RoomID == 0 {
		return errors.New("room id is required")
	}

	if booking.CheckIn.IsZero() {
		return errors.New("check-in date is required")
	}

	if booking.CheckOut.IsZero() {
		return errors.New("check-out date is required")
	}

	if !booking.CheckOut.After(booking.CheckIn) {
		return errors.New(
			"check-out must be after check-in",
		)
	}

	if booking.CheckIn.Before(time.Now()) {
		return errors.New(
			"check-in date cannot be in the past",
		)
	}

	if booking.Guests < 1 {
		return errors.New(
			"guests must be at least 1",
		)
	}

	if booking.Status != "" &&
		booking.Status != "pending" &&
		booking.Status != "confirmed" &&
		booking.Status != "cancelled" {
		return errors.New(
			"invalid booking status",
		)
	}

	return nil
}

// =========================
// Not Found Helper
// =========================

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
