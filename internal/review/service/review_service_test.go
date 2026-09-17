package service

import (
	"testing"

	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"gorm.io/gorm"
)

type fakeUserLister struct {
	users []userModel.User
}

func (f *fakeUserLister) GetByIDs(ids []string) ([]userModel.User, error) {
	return f.users, nil
}

type fakeHotelLister struct {
	hotels []hotelModel.Hotel
}

func (f *fakeHotelLister) GetByIDs(ids []uint) ([]hotelModel.Hotel, error) {
	return f.hotels, nil
}

type mockReviewRepository struct {
	reviews      []model.Review
	getByIDErr   error
	deleteErr    error
	deleteCalled bool
	deletedID    uint
}

func (m *mockReviewRepository) Create(review *model.Review) error {
	return nil
}

func (m *mockReviewRepository) GetAll() ([]model.Review, error) {
	return m.reviews, nil
}

func (m *mockReviewRepository) GetByID(id uint) (*model.Review, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	return &model.Review{ID: id}, nil
}

func (m *mockReviewRepository) GetByHotelID(hotelID uint) ([]model.Review, error) {
	return nil, nil
}

func (m *mockReviewRepository) GetByUserAndHotel(userID string, hotelID uint) (*model.Review, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockReviewRepository) GetAverageRating(hotelID uint) (float64, error) {
	return 0, nil
}

func (m *mockReviewRepository) Delete(id uint) error {
	m.deleteCalled = true
	m.deletedID = id
	return m.deleteErr
}

func TestGetAllReviews_ReturnsRepositoryResults(t *testing.T) {
	repo := &mockReviewRepository{
		reviews: []model.Review{{ID: 1, HotelID: 5}, {ID: 2, HotelID: 6}},
	}
	svc := NewReviewService(repo, &fakeUserLister{}, &fakeHotelLister{})

	reviews, err := svc.GetAllReviews()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(reviews) != 2 {
		t.Fatalf("expected 2 reviews, got %d", len(reviews))
	}
}

func TestDeleteReview_InvalidID(t *testing.T) {
	repo := &mockReviewRepository{}
	svc := NewReviewService(repo, &fakeUserLister{}, &fakeHotelLister{})

	if err := svc.DeleteReview(0); err == nil {
		t.Fatal("expected an error for id 0, got nil")
	}
	if repo.deleteCalled {
		t.Error("Delete should not be called for an invalid id")
	}
}

func TestDeleteReview_NotFound(t *testing.T) {
	repo := &mockReviewRepository{getByIDErr: gorm.ErrRecordNotFound}
	svc := NewReviewService(repo, &fakeUserLister{}, &fakeHotelLister{})

	err := svc.DeleteReview(42)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if repo.deleteCalled {
		t.Error("Delete should not be called when the review doesn't exist")
	}
}

func TestDeleteReview_Success(t *testing.T) {
	repo := &mockReviewRepository{}
	svc := NewReviewService(repo, &fakeUserLister{}, &fakeHotelLister{})

	if err := svc.DeleteReview(7); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !repo.deleteCalled {
		t.Error("expected Delete to be called")
	}
	if repo.deletedID != 7 {
		t.Errorf("expected Delete called with id 7, got %d", repo.deletedID)
	}
}
