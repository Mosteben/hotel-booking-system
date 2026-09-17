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
	GetDetails(id uint) (*model.Hotel, error)
	Update(hotel *model.Hotel) error
	Delete(id uint) error
	Search(filters *model.SearchRequest) ([]model.Hotel, error)

	// GetByIDs batch-fetches hotels (no image preload - name/id only) for
	// Admin DTO resolution, so those don't do one query per row.
	GetByIDs(ids []uint) ([]model.Hotel, error)

	CreateImage(image *model.HotelImage) error
	GetImageByID(id uint) (*model.HotelImage, error)
	GetImagesByHotelID(hotelID uint) ([]model.HotelImage, error)
	CountImagesByHotelID(hotelID uint) (int64, error)
	DeleteImage(id uint) error
	UnsetMainImage(hotelID uint) error
	SetMainImage(id uint) error

	// Create a repository that uses the provided transaction.
	WithTx(tx *gorm.DB) HotelRepository
}

type hotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{
		db: db,
	}
}

func (r *hotelRepository) WithTx(
	tx *gorm.DB,
) HotelRepository {
	return &hotelRepository{
		db: tx,
	}
}

func (r *hotelRepository) Create(hotel *model.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *hotelRepository) GetAll() ([]model.Hotel, error) {
	var hotels []model.Hotel

	err := r.db.
		Preload("Images").
		Order("id DESC").
		Find(&hotels).Error

	return hotels, err
}

func (r *hotelRepository) GetByIDs(ids []uint) ([]model.Hotel, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	var hotels []model.Hotel

	err := r.db.
		Where("id IN ?", ids).
		Find(&hotels).Error

	return hotels, err
}

func (r *hotelRepository) GetByID(id uint) (*model.Hotel, error) {
	var hotel model.Hotel

	err := r.db.
		Preload("Images").
		First(&hotel, id).Error

	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (r *hotelRepository) GetDetails(id uint) (*model.Hotel, error) {
	var hotel model.Hotel

	err := r.db.
		Preload("Rooms.Images").
		Preload("Images").
		First(&hotel, id).Error

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
		Preload("Images").
		Order("hotels.id DESC").
		Find(&hotels).Error

	return hotels, err
}

// =========================
// Hotel Images
// =========================

func (r *hotelRepository) CreateImage(image *model.HotelImage) error {
	return r.db.Create(image).Error
}

func (r *hotelRepository) GetImageByID(id uint) (*model.HotelImage, error) {
	var image model.HotelImage

	err := r.db.First(&image, id).Error
	if err != nil {
		return nil, err
	}

	return &image, nil
}

func (r *hotelRepository) GetImagesByHotelID(
	hotelID uint,
) ([]model.HotelImage, error) {

	var images []model.HotelImage

	err := r.db.
		Where("hotel_id = ?", hotelID).
		Order("id ASC").
		Find(&images).Error

	return images, err
}

func (r *hotelRepository) CountImagesByHotelID(
	hotelID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&model.HotelImage{}).
		Where("hotel_id = ?", hotelID).
		Count(&count).Error

	return count, err
}

func (r *hotelRepository) DeleteImage(id uint) error {
	return r.db.Delete(&model.HotelImage{}, id).Error
}

func (r *hotelRepository) UnsetMainImage(hotelID uint) error {
	return r.db.
		Model(&model.HotelImage{}).
		Where("hotel_id = ? AND is_main = ?", hotelID, true).
		Update("is_main", false).Error
}

func (r *hotelRepository) SetMainImage(id uint) error {
	return r.db.
		Model(&model.HotelImage{}).
		Where("id = ?", id).
		Update("is_main", true).Error
}
