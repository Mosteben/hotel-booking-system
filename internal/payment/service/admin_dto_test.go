package service

import (
	"testing"
	"time"

	bookingModel "github.com/Mosteben/hotel-booking-system/internal/booking/model"
	bookingRepository "github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	paymentModel "github.com/Mosteben/hotel-booking-system/internal/payment/model"
	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeUserLister struct{ users []userModel.User }

func (f *fakeUserLister) GetByIDs(ids []string) ([]userModel.User, error) {
	return f.users, nil
}

type fakeRoomLister struct{ rooms []roomModel.Room }

func (f *fakeRoomLister) GetByIDs(ids []uint) ([]roomModel.Room, error) {
	return f.rooms, nil
}

type fakeHotelLister struct{ hotels []hotelModel.Hotel }

func (f *fakeHotelLister) GetByIDs(ids []uint) ([]hotelModel.Hotel, error) {
	return f.hotels, nil
}

// paymentService.bookingRepo stays the full BookingRepository interface
// (it's also used for the real payment/booking status transitions), so
// resolveAdminPaymentSummaries is tested against a stub implementing that
// full interface rather than the narrow fakeBookingGetter above.
type stubBookingRepository struct {
	bookings []bookingModel.Booking
}

func (s *stubBookingRepository) Create(*bookingModel.Booking) error { return nil }
func (s *stubBookingRepository) GetAll() ([]bookingModel.Booking, error) {
	return nil, nil
}
func (s *stubBookingRepository) GetByID(uint) (*bookingModel.Booking, error) {
	return nil, gorm.ErrRecordNotFound
}
func (s *stubBookingRepository) GetByIDs(ids []uint) ([]bookingModel.Booking, error) {
	return s.bookings, nil
}
func (s *stubBookingRepository) GetByUserID(string) ([]bookingModel.Booking, error) {
	return nil, nil
}
func (s *stubBookingRepository) Update(*bookingModel.Booking) error { return nil }
func (s *stubBookingRepository) Delete(uint) error                  { return nil }
func (s *stubBookingRepository) IsRoomAvailable(uint, time.Time, time.Time, uint) (bool, error) {
	return true, nil
}
func (s *stubBookingRepository) WithTx(*gorm.DB) bookingRepository.BookingRepository {
	return s
}

func TestResolveAdminPaymentSummaries_ResolvesRealRelationships(t *testing.T) {
	userID := uuid.New()
	hotel := hotelModel.Hotel{ID: 1, Name: "NileStay Cairo"}
	room := roomModel.Room{ID: 10, HotelID: hotel.ID, RoomNumber: "101", Type: "Deluxe"}
	booking := bookingModel.Booking{
		ID:       50,
		RoomID:   room.ID,
		Status:   "confirmed",
		CheckIn:  time.Now(),
		CheckOut: time.Now().Add(48 * time.Hour),
	}
	user := userModel.User{ID: userID, FirstName: "Ahmed", LastName: "Tester", Email: "ahmed@test.com", Phone: "+201000000000"}

	svc := &paymentService{
		userRepo:  &fakeUserLister{users: []userModel.User{user}},
		roomRepo:  &fakeRoomLister{rooms: []roomModel.Room{room}},
		hotelRepo: &fakeHotelLister{hotels: []hotelModel.Hotel{hotel}},
	}

	payments := []paymentModel.Payment{
		{ID: 1, BookingID: booking.ID, UserID: userID.String(), Amount: 500, PaymentMethod: "card", Status: "paid"},
	}

	// resolveAdminPaymentSummaries needs bookingsByID populated - since
	// bookingRepo is the full BookingRepository interface, test the
	// resolution against a manually-seeded map by calling the resolver
	// with a service whose bookingRepo returns our fixture booking.
	svc.bookingRepo = &stubBookingRepository{bookings: []bookingModel.Booking{booking}}

	summaries, err := svc.resolveAdminPaymentSummaries(payments)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	s := summaries[0]

	if s.User == nil || s.User.Name != "Ahmed Tester" || s.User.Email != "ahmed@test.com" {
		t.Errorf("expected resolved user Ahmed Tester, got %+v", s.User)
	}
	if s.Booking == nil || s.Booking.ID != booking.ID {
		t.Errorf("expected resolved booking #%d, got %+v", booking.ID, s.Booking)
	}
	if s.Room == nil || s.Room.RoomNumber != "101" {
		t.Errorf("expected resolved room 101, got %+v", s.Room)
	}
	if s.Hotel == nil || s.Hotel.Name != "NileStay Cairo" {
		t.Errorf("expected resolved hotel NileStay Cairo, got %+v", s.Hotel)
	}
}

func TestResolveAdminPaymentSummaries_OrphanedUserStaysNilNeverFabricated(t *testing.T) {
	svc := &paymentService{
		userRepo:    &fakeUserLister{users: nil},
		roomRepo:    &fakeRoomLister{rooms: nil},
		hotelRepo:   &fakeHotelLister{hotels: nil},
		bookingRepo: &stubBookingRepository{bookings: nil},
	}

	payments := []paymentModel.Payment{
		{ID: 1, BookingID: 999, UserID: "00000000-0000-0000-0000-000000001003", Amount: 100, PaymentMethod: "cash", Status: "pending"},
	}

	summaries, err := svc.resolveAdminPaymentSummaries(payments)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	s := summaries[0]
	if s.User != nil {
		t.Errorf("expected nil User for an unresolvable user_id, got %+v", s.User)
	}
	if s.Booking != nil {
		t.Errorf("expected nil Booking for an unresolvable booking_id, got %+v", s.Booking)
	}
	// The payment's own real fields must still come through untouched.
	if s.ID != 1 || s.Amount != 100 || s.Status != "pending" {
		t.Errorf("expected the payment's own fields to be preserved, got %+v", s)
	}
}
