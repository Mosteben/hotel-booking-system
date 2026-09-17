package service

import (
	"errors"
	"strings"

	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/repository"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"gorm.io/gorm"
)

// userLister/hotelLister are the minimal slices of UserRepository/
// HotelRepository this service needs for Admin DTO resolution - the real
// repositories already satisfy these, so tests can use small local fakes
// instead of mocking the full interfaces.
type userLister interface {
	GetByIDs(ids []string) ([]userModel.User, error)
}

type hotelLister interface {
	GetByIDs(ids []uint) ([]hotelModel.Hotel, error)
}

type ReviewService interface {
	CreateReview(review *model.Review) error
	// GetAllReviews is admin/manager only - it returns each review with
	// its author and hotel resolved server-side (see AdminReviewSummary).
	GetAllReviews() ([]AdminReviewSummary, error)
	GetHotelReviews(hotelID uint) ([]model.Review, error)
	GetHotelAverageRating(hotelID uint) (float64, error)
	DeleteReview(id uint) error
}

type reviewService struct {
	repo      repository.ReviewRepository
	userRepo  userLister
	hotelRepo hotelLister
}

func NewReviewService(
	repo repository.ReviewRepository,
	userRepo userLister,
	hotelRepo hotelLister,
) ReviewService {
	return &reviewService{
		repo:      repo,
		userRepo:  userRepo,
		hotelRepo: hotelRepo,
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

func (s *reviewService) GetAllReviews() ([]AdminReviewSummary, error) {
	reviews, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return s.resolveAdminReviewSummaries(reviews)
}

// DeleteReview is moderation-only - the route this backs is admin/manager
// gated, so no ownership check is needed here (unlike a user deleting
// their own review, which this codebase doesn't support yet).
func (s *reviewService) DeleteReview(id uint) error {
	if id == 0 {
		return errors.New("invalid review id")
	}

	if _, err := s.repo.GetByID(id); err != nil {
		return err
	}

	return s.repo.Delete(id)
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
