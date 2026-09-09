package service

import (
	"errors"
	"strings"

	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/repository"
	"gorm.io/gorm"
)

type ReviewService interface {
	CreateReview(review *model.Review) error
	GetHotelReviews(hotelID uint) ([]model.Review, error)
	GetHotelAverageRating(hotelID uint) (float64, error)
}

type reviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(
	repo repository.ReviewRepository,
) ReviewService {
	return &reviewService{
		repo: repo,
	}
}

func (s *reviewService) CreateReview(
	review *model.Review,
) error {
	if review == nil {
		return errors.New("review data is required")
	}

	if review.HotelID == 0 {
		return errors.New("hotel id is required")
	}

	if strings.TrimSpace(review.UserID) == "" {
		return errors.New("user id is required")
	}

	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	review.Comment = strings.TrimSpace(review.Comment)

	_, err := s.repo.GetByUserAndHotel(
		review.UserID,
		review.HotelID,
	)

	if err == nil {
		return errors.New(
			"user has already reviewed this hotel",
		)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.repo.Create(review)
}

func (s *reviewService) GetHotelReviews(
	hotelID uint,
) ([]model.Review, error) {
	if hotelID == 0 {
		return nil, errors.New("invalid hotel id")
	}

	return s.repo.GetByHotelID(hotelID)
}

func (s *reviewService) GetHotelAverageRating(
	hotelID uint,
) (float64, error) {
	if hotelID == 0 {
		return 0, errors.New("invalid hotel id")
	}

	return s.repo.GetAverageRating(hotelID)
}