package service

import (
	"errors"
	"strings"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	roomRepository "github.com/Mosteben/hotel-booking-system/internal/room/repository"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(booking *model.Booking) error
	GetAllBookings() ([]model.Booking, error)
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
	repo     repository.BookingRepository
	roomRepo roomRepository.RoomRepository
}

func NewBookingService(
	repo repository.BookingRepository,
	roomRepo roomRepository.RoomRepository,
) BookingService {
	return &bookingService{
		repo:     repo,
		roomRepo: roomRepo,
	}
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

	// Check room availability for the selected dates.
	available, err := s.repo.IsRoomAvailable(
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

	// New bookings always start as pending.
	booking.Status = "pending"

	return s.repo.Create(booking)
}

// =========================
// Get All Bookings
// =========================

func (s *bookingService) GetAllBookings() (
	[]model.Booking,
	error,
) {
	return s.repo.GetAll()
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

	// Check room availability.
	// Exclude the current booking so it doesn't
	// conflict with itself.
	available, err := s.repo.IsRoomAvailable(
		booking.RoomID,
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

	// Calculate number of nights.
	nights := int(
		booking.CheckOut.Sub(booking.CheckIn).Hours() / 24,
	)

	if nights <= 0 {
		return errors.New("invalid booking duration")
	}

	// Calculate total price on the server.
	totalPrice := room.PricePerNight * float64(nights)

	existingBooking.RoomID = booking.RoomID
	existingBooking.CheckIn = booking.CheckIn
	existingBooking.CheckOut = booking.CheckOut
	existingBooking.Guests = booking.Guests
	existingBooking.TotalPrice = totalPrice

	// Status is controlled by the system.
	// Customer cannot change booking status here.

	return s.repo.Update(existingBooking)
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
