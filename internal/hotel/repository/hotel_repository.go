package repository

import (
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"gorm.io/gorm"
)

type HotelRepository interface {
	Create(hotel *model.Hotel) error
	GetAll() ([]model.Hotel, error)
	GetByID(id uint) (*model.Hotel, error)
	Update(hotel *model.Hotel) error
	Delete(id uint) error
	Search(filters *model.SearchRequest) ([]model.Hotel, error)
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

func (r *hotelRepository) Search(
	filters *model.SearchRequest,
) ([]model.Hotel, error) {

	var hotels []model.Hotel

	query := r.db.
		Table("hotels").
		Select("DISTINCT hotels.*").
		Joins("JOIN rooms ON rooms.hotel_id = hotels.id")

	// Search by hotel name
	if strings.TrimSpace(filters.Name) != "" {
		query = query.Where(
			"LOWER(hotels.name) LIKE LOWER(?)",
			"%"+strings.TrimSpace(filters.Name)+"%",
		)
	}

	// Search by city
	if strings.TrimSpace(filters.City) != "" {
		query = query.Where(
			"LOWER(hotels.city) = LOWER(?)",
			strings.TrimSpace(filters.City),
		)
	}

	// Search by country
	if strings.TrimSpace(filters.Country) != "" {
		query = query.Where(
			"LOWER(hotels.country) = LOWER(?)",
			strings.TrimSpace(filters.Country),
		)
	}

	// Search by stars
	if filters.Stars > 0 {
		query = query.Where(
			"hotels.stars = ?",
			filters.Stars,
		)
	}

	// Minimum room price
	if filters.MinPrice > 0 {
		query = query.Where(
			"rooms.price_per_night >= ?",
			filters.MinPrice,
		)
	}

	// Maximum room price
	if filters.MaxPrice > 0 {
		query = query.Where(
			"rooms.price_per_night <= ?",
			filters.MaxPrice,
		)
	}

	// Search by room type
	if strings.TrimSpace(filters.RoomType) != "" {
		query = query.Where(
			"LOWER(rooms.type) = LOWER(?)",
			strings.TrimSpace(filters.RoomType),
		)
	}

	// Minimum room capacity
	if filters.MinCapacity > 0 {
		query = query.Where(
			"rooms.capacity >= ?",
			filters.MinCapacity,
		)
	}

	err := query.
		Order("hotels.id DESC").
		Find(&hotels).Error

	return hotels, err
}
