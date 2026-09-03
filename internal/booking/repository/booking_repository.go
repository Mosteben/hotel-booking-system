package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	GetAll() ([]model.Booking, error)
	GetByID(id uint) (*model.Booking, error)
	GetByUserID(userID string) ([]model.Booking, error)
	Update(booking *model.Booking) error
	Delete(id uint) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) GetAll() ([]model.Booking, error) {
	var bookings []model.Booking

	err := r.db.
		Order("id DESC").
		Find(&bookings).Error

	return bookings, err
}

func (r *bookingRepository) GetByID(id uint) (*model.Booking, error) {
	var booking model.Booking

	err := r.db.First(&booking, id).Error
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) GetByUserID(userID string) ([]model.Booking, error) {
	var bookings []model.Booking

	err := r.db.
		Where("user_id = ?", userID).
		Order("id DESC").
		Find(&bookings).Error

	return bookings, err
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&model.Booking{}, id).Error
}
