package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/favorite/model"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	Create(favorite *model.Favorite) error
	GetByUserID(userID string) ([]model.Favorite, error)
	GetByUserAndHotel(userID string, hotelID uint) (*model.Favorite, error)
	Delete(favorite *model.Favorite) error
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{
		db: db,
	}
}

func (r *favoriteRepository) Create(
	favorite *model.Favorite,
) error {
	return r.db.Create(favorite).Error
}

func (r *favoriteRepository) GetByUserID(
	userID string,
) ([]model.Favorite, error) {

	var favorites []model.Favorite

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&favorites).Error

	return favorites, err
}

func (r *favoriteRepository) GetByUserAndHotel(
	userID string,
	hotelID uint,
) (*model.Favorite, error) {

	var favorite model.Favorite

	err := r.db.
		Where("user_id = ? AND hotel_id = ?", userID, hotelID).
		First(&favorite).Error

	if err != nil {
		return nil, err
	}

	return &favorite, nil
}

func (r *favoriteRepository) Delete(
	favorite *model.Favorite,
) error {
	return r.db.Delete(favorite).Error
}
