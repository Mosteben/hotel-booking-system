package service

import (
	"errors"
	"strings"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(booking *model.Booking) error
	GetAllBookings() ([]model.Booking, error)
	GetBookingByID(id uint) (*model.Booking, error)
	GetBookingsByUserID(userID string) ([]model.Booking, error)
	UpdateBooking(id uint, userID string, booking *model.Booking) error
	DeleteBooking(id uint) error
}

type bookingService struct {
	repo repository.BookingRepository
}

func NewBookingService(repo repository.BookingRepository) BookingService {
	return &bookingService{
		repo: repo,
	}
}

func (s *bookingService) CreateBooking(booking *model.Booking) error {
	if err := validateBooking(booking); err != nil {
		return err
	}

	booking.Status = "pending"

	return s.repo.Create(booking)
}

func (s *bookingService) GetAllBookings() ([]model.Booking, error) {
	return s.repo.GetAll()
}

func (s *bookingService) GetBookingByID(id uint) (*model.Booking, error) {
	if id == 0 {
		return nil, errors.New("invalid booking id")
	}

	return s.repo.GetByID(id)
}

func (s *bookingService) GetBookingsByUserID(userID string) ([]model.Booking, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetByUserID(userID)
}

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
		return errors.New("you are not allowed to update this booking")
	}

	existingBooking.RoomID = booking.RoomID
	existingBooking.CheckIn = booking.CheckIn
	existingBooking.CheckOut = booking.CheckOut
	existingBooking.Guests = booking.Guests
	existingBooking.TotalPrice = booking.TotalPrice

	// Status is controlled by the system.
	// We don't allow the customer to change it here.

	return s.repo.Update(existingBooking)
}

func (s *bookingService) DeleteBooking(id uint) error {
	if id == 0 {
		return errors.New("invalid booking id")
	}

	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func validateBooking(booking *model.Booking) error {
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
		return errors.New("check-out must be after check-in")
	}

	if booking.CheckIn.Before(time.Now()) {
		return errors.New("check-in date cannot be in the past")
	}

	if booking.Guests < 1 {
		return errors.New("guests must be at least 1")
	}

	if booking.TotalPrice <= 0 {
		return errors.New("total price must be greater than 0")
	}

	if booking.Status != "" &&
		booking.Status != "pending" &&
		booking.Status != "confirmed" &&
		booking.Status != "cancelled" {
		return errors.New("invalid booking status")
	}

	return nil
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}