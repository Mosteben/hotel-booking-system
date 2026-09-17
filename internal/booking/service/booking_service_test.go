package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
	"gorm.io/gorm"
)

// =========================
// Mocks
// =========================

type mockRoomLister struct {
	room    *roomModel.Room
	getErr  error
	byIDErr error
}

func (m *mockRoomLister) GetByID(id uint) (*roomModel.Room, error) {
	return m.room, m.getErr
}

func (m *mockRoomLister) GetByIDs(ids []uint) ([]roomModel.Room, error) {
	return nil, m.byIDErr
}

type mockHotelLister struct {
	hotel  *hotelModel.Hotel
	getErr error
}

func (m *mockHotelLister) GetByID(id uint) (*hotelModel.Hotel, error) {
	return m.hotel, m.getErr
}

func (m *mockHotelLister) GetByIDs(ids []uint) ([]hotelModel.Hotel, error) {
	return nil, nil
}

// =========================
// applyBookingSnapshot
// =========================

func TestApplyBookingSnapshot_Success(t *testing.T) {
	room := &roomModel.Room{
		ID:            52,
		HotelID:       12,
		RoomNumber:    "M-10-05",
		Type:          "suite",
		PricePerNight: 3000,
	}

	hotel := &hotelModel.Hotel{
		ID:   12,
		Name: "Mock Hotel 10",
	}

	svc := &bookingService{
		hotelRepo: &mockHotelLister{hotel: hotel},
	}

	booking := &model.Booking{}

	if err := svc.applyBookingSnapshot(booking, room); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if booking.HotelID != 12 {
		t.Errorf("expected HotelID 12, got %d", booking.HotelID)
	}
	if booking.HotelName != "Mock Hotel 10" {
		t.Errorf("expected HotelName 'Mock Hotel 10', got %q", booking.HotelName)
	}
	if booking.RoomNumber != "M-10-05" {
		t.Errorf("expected RoomNumber 'M-10-05', got %q", booking.RoomNumber)
	}
	if booking.RoomType != "suite" {
		t.Errorf("expected RoomType 'suite', got %q", booking.RoomType)
	}
	if booking.PricePerNight != 3000 {
		t.Errorf("expected PricePerNight 3000, got %v", booking.PricePerNight)
	}
}

func TestApplyBookingSnapshot_RejectsHotelMismatch(t *testing.T) {
	room := &roomModel.Room{
		ID:      52,
		HotelID: 12,
	}

	svc := &bookingService{
		hotelRepo: &mockHotelLister{hotel: &hotelModel.Hotel{ID: 12, Name: "Mock Hotel 10"}},
	}

	// Client claims the room belongs to a different hotel than it actually
	// does - must be rejected, never silently trusted or corrected.
	booking := &model.Booking{HotelID: 999}

	err := svc.applyBookingSnapshot(booking, room)
	if err == nil {
		t.Fatal("expected an error for mismatched hotel id, got nil")
	}
}

func TestApplyBookingSnapshot_AllowsMatchingHotelID(t *testing.T) {
	room := &roomModel.Room{ID: 52, HotelID: 12, RoomNumber: "M-10-05", Type: "suite", PricePerNight: 3000}
	hotel := &hotelModel.Hotel{ID: 12, Name: "Mock Hotel 10"}

	svc := &bookingService{
		hotelRepo: &mockHotelLister{hotel: hotel},
	}

	// Client-supplied hotel_id that DOES match the room's real hotel is
	// fine - it's still fully overwritten from the server-loaded hotel
	// afterward, never taken as-is.
	booking := &model.Booking{HotelID: 12}

	if err := svc.applyBookingSnapshot(booking, room); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if booking.HotelName != "Mock Hotel 10" {
		t.Errorf("expected HotelName from the server-loaded hotel, got %q", booking.HotelName)
	}
}

func TestApplyBookingSnapshot_HotelNotFound(t *testing.T) {
	room := &roomModel.Room{ID: 52, HotelID: 12}

	svc := &bookingService{
		hotelRepo: &mockHotelLister{getErr: gorm.ErrRecordNotFound},
	}

	booking := &model.Booking{}

	err := svc.applyBookingSnapshot(booking, room)
	if err == nil {
		t.Fatal("expected an error when the hotel can't be found, got nil")
	}
}

// =========================
// CreateBooking
// =========================

func TestCreateBooking_RoomNotFound(t *testing.T) {
	svc := &bookingService{
		roomRepo: &mockRoomLister{getErr: gorm.ErrRecordNotFound},
	}

	booking := &model.Booking{
		UserID:   "user-1",
		RoomID:   999,
		CheckIn:  time.Now().Add(24 * time.Hour),
		CheckOut: time.Now().Add(48 * time.Hour),
		Guests:   1,
	}

	err := svc.CreateBooking(booking)
	if err == nil || err.Error() != "room not found" {
		t.Fatalf("expected 'room not found', got %v", err)
	}
}

func TestCreateBooking_RoomNotAvailable(t *testing.T) {
	svc := &bookingService{
		roomRepo: &mockRoomLister{room: &roomModel.Room{ID: 1, Status: "maintenance", Capacity: 2}},
	}

	booking := &model.Booking{
		UserID:   "user-1",
		RoomID:   1,
		CheckIn:  time.Now().Add(24 * time.Hour),
		CheckOut: time.Now().Add(48 * time.Hour),
		Guests:   1,
	}

	err := svc.CreateBooking(booking)
	if err == nil || err.Error() != "room is not available" {
		t.Fatalf("expected 'room is not available', got %v", err)
	}
}

func TestCreateBooking_GuestsExceedCapacity(t *testing.T) {
	svc := &bookingService{
		roomRepo: &mockRoomLister{room: &roomModel.Room{ID: 1, Status: "available", Capacity: 2}},
	}

	booking := &model.Booking{
		UserID:   "user-1",
		RoomID:   1,
		CheckIn:  time.Now().Add(24 * time.Hour),
		CheckOut: time.Now().Add(48 * time.Hour),
		Guests:   5,
	}

	err := svc.CreateBooking(booking)
	if err == nil || err.Error() != "number of guests exceeds room capacity" {
		t.Fatalf("expected capacity error, got %v", err)
	}
}

// TestCreateBooking_PricingUnchangedAndSnapshotAppliedBeforeCommit verifies
// two things at once: (1) the server-side price calculation
// (room.PricePerNight * nights) still runs exactly as before this change,
// and (2) the snapshot is applied before the booking is ever handed to the
// DB transaction. It forces the flow to fail at the hotel lookup (a step
// that only exists after pricing runs) so it can assert on booking.TotalPrice
// without needing a real *gorm.DB/transaction.
func TestCreateBooking_PricingUnchangedAndSnapshotAppliedBeforeCommit(t *testing.T) {
	svc := &bookingService{
		roomRepo: &mockRoomLister{room: &roomModel.Room{
			ID:            1,
			HotelID:       7,
			Status:        "available",
			Capacity:      4,
			RoomNumber:    "101",
			Type:          "double",
			PricePerNight: 500,
		}},
		hotelRepo: &mockHotelLister{getErr: errors.New("boom")},
	}

	checkIn := time.Now().Add(24 * time.Hour)
	checkOut := checkIn.AddDate(0, 0, 3) // 3 nights

	booking := &model.Booking{
		UserID:   "user-1",
		RoomID:   1,
		CheckIn:  checkIn,
		CheckOut: checkOut,
		Guests:   2,
	}

	err := svc.CreateBooking(booking)
	if err == nil {
		t.Fatal("expected an error from the forced hotel-lookup failure, got nil")
	}

	if booking.TotalPrice != 1500 {
		t.Errorf("expected TotalPrice 500*3=1500 (unchanged pricing logic), got %v", booking.TotalPrice)
	}

	if booking.RoomNumber != "" || booking.HotelName != "" {
		t.Errorf("snapshot should not be applied when the hotel lookup fails, got RoomNumber=%q HotelName=%q", booking.RoomNumber, booking.HotelName)
	}
}

func TestCreateBooking_InvalidDuration(t *testing.T) {
	svc := &bookingService{
		roomRepo: &mockRoomLister{room: &roomModel.Room{ID: 1, Status: "available", Capacity: 2}},
	}

	checkIn := time.Now().Add(24 * time.Hour)

	booking := &model.Booking{
		UserID: "user-1",
		RoomID: 1,
		// Less than a full night apart - CheckOut is still after CheckIn
		// (passes validateBooking), but rounds down to zero whole nights.
		CheckIn:  checkIn,
		CheckOut: checkIn.Add(2 * time.Hour),
		Guests:   1,
	}

	err := svc.CreateBooking(booking)
	if err == nil || err.Error() != "invalid booking duration" {
		t.Fatalf("expected 'invalid booking duration', got %v", err)
	}
}

// =========================
// Historical resilience: GetBookingByID never depends on the room/hotel
// still existing - it only ever reads the booking row itself, snapshot
// fields included, regardless of whether RoomID/UserID currently resolve.
// =========================

// mockBookingRepoForGet implements the full repository.BookingRepository
// interface; only GetByID actually does anything, since that's all
// GetBookingByID needs.
type mockBookingRepoForGet struct {
	booking *model.Booking
	err     error
}

func (m *mockBookingRepoForGet) Create(b *model.Booking) error    { return nil }
func (m *mockBookingRepoForGet) GetAll() ([]model.Booking, error) { return nil, nil }
func (m *mockBookingRepoForGet) GetByID(id uint) (*model.Booking, error) {
	return m.booking, m.err
}
func (m *mockBookingRepoForGet) GetByIDs(ids []uint) ([]model.Booking, error) { return nil, nil }
func (m *mockBookingRepoForGet) GetByUserID(userID string) ([]model.Booking, error) {
	return nil, nil
}
func (m *mockBookingRepoForGet) Update(b *model.Booking) error { return nil }
func (m *mockBookingRepoForGet) Delete(id uint) error          { return nil }
func (m *mockBookingRepoForGet) IsRoomAvailable(roomID uint, checkIn, checkOut time.Time, excludeBookingID uint) (bool, error) {
	return true, nil
}
func (m *mockBookingRepoForGet) WithTx(tx *gorm.DB) repository.BookingRepository { return m }

func TestGetBookingByID_ReturnsSnapshotEvenWhenRoomIsGone(t *testing.T) {
	// This booking's RoomID (999) doesn't exist anywhere the mock knows
	// about - simulating a room deleted after the booking was made. The
	// point is that GetBookingByID never looks the room up at all, so it
	// can't be affected either way; the snapshot fields already on the row
	// are simply returned as-is.
	stored := &model.Booking{
		ID:            56,
		UserID:        "real-user",
		RoomID:        999,
		Status:        "cancelled",
		HotelID:       12,
		HotelName:     "Mock Hotel 10",
		RoomNumber:    "M-10-05",
		RoomType:      "suite",
		PricePerNight: 3000,
	}

	svc := &bookingService{repo: &mockBookingRepoForGet{booking: stored}}

	got, err := svc.GetBookingByID(56, "real-user", "customer")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.HotelName != "Mock Hotel 10" || got.RoomNumber != "M-10-05" {
		t.Errorf("expected snapshot fields intact, got HotelName=%q RoomNumber=%q", got.HotelName, got.RoomNumber)
	}
}
