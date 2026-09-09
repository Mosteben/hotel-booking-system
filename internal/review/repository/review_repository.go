package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *model.Review) error
	GetByHotelID(hotelID uint) ([]model.Review, error)
	GetByUserAndHotel(userID string, hotelID uint) (*model.Review, error)
	GetAverageRating(hotelID uint) (float64, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

func (r *reviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByHotelID(
	hotelID uint,
) ([]model.Review, error) {
	var reviews []model.Review

	err := r.db.
		Where("hotel_id = ?", hotelID).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, err
}

func (r *reviewRepository) GetByUserAndHotel(
	userID string,
	hotelID uint,
) (*model.Review, error) {
	var review model.Review

	err := r.db.
		Where("user_id = ? AND hotel_id = ?", userID, hotelID).
		First(&review).Error

	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *reviewRepository) GetAverageRating(
	hotelID uint,
) (float64, error) {
	var average float64

	err := r.db.
		Model(&model.Review{}).
		Where("hotel_id = ?", hotelID).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&average).Error

	return average, err
}
