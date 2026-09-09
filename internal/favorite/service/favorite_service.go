package service

import (
	"errors"
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/favorite/model"
	"github.com/Mosteben/hotel-booking-system/internal/favorite/repository"
	"gorm.io/gorm"
)

var ErrFavoriteNotFound = errors.New("hotel is not in favorites")

type FavoriteService interface {
	AddFavorite(favorite *model.Favorite) error
	GetMyFavorites(userID string) ([]model.Favorite, error)
	RemoveFavorite(userID string, hotelID uint) error
}

type favoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(
	repo repository.FavoriteRepository,
) FavoriteService {
	return &favoriteService{
		repo: repo,
	}
}

// =========================
// Add Favorite
// =========================

func (s *favoriteService) AddFavorite(
	favorite *model.Favorite,
) error {

	if favorite == nil {
		return errors.New("favorite data is required")
	}

	if strings.TrimSpace(favorite.UserID) == "" {
		return errors.New("user id is required")
	}

	if favorite.HotelID == 0 {
		return errors.New("hotel id is required")
	}

	// Prevent duplicate favorites.
	_, err := s.repo.GetByUserAndHotel(
		favorite.UserID,
		favorite.HotelID,
	)

	if err == nil {
		return errors.New("hotel is already in favorites")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.repo.Create(favorite)
}

// =========================
// Get My Favorites
// =========================

func (s *favoriteService) GetMyFavorites(
	userID string,
) ([]model.Favorite, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetByUserID(userID)
}

// =========================
// Remove Favorite
// =========================

func (s *favoriteService) RemoveFavorite(
	userID string,
	hotelID uint,
) error {

	if strings.TrimSpace(userID) == "" {
		return errors.New("invalid user id")
	}

	if hotelID == 0 {
		return errors.New("invalid hotel id")
	}

	favorite, err := s.repo.GetByUserAndHotel(
		userID,
		hotelID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFavoriteNotFound
		}

		return err
	}

	return s.repo.Delete(favorite)
}
