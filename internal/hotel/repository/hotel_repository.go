package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"gorm.io/gorm"
)

type HotelRepository interface {
	Create(hotel *model.Hotel) error
	GetAll() ([]model.Hotel, error)
	GetByID(id uint) (*model.Hotel, error)
	Update(hotel *model.Hotel) error
	Delete(id uint) error
}

type hotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{
		db: db,
	}
}

func (r *hotelRepository) Create(hotel *model.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *hotelRepository) GetAll() ([]model.Hotel, error) {
	var hotels []model.Hotel

	err := r.db.
		Order("id DESC").
		Find(&hotels).Error

	return hotels, err
}

func (r *hotelRepository) GetByID(id uint) (*model.Hotel, error) {
	var hotel model.Hotel

	err := r.db.First(&hotel, id).Error
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (r *hotelRepository) Update(hotel *model.Hotel) error {
	return r.db.Save(hotel).Error
}

func (r *hotelRepository) Delete(id uint) error {
	return r.db.Delete(&model.Hotel{}, id).Error
}