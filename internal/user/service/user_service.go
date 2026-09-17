package service

import (
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/user/model"
)

// userLister is the minimal slice of UserRepository this service needs.
// *repository.UserRepository already satisfies this (structural typing),
// so the repository itself doesn't need to change to make this testable.
type userLister interface {
	GetAll() ([]model.User, error)
}

// AdminUserSummary is the safe, admin-facing view of a user - explicitly
// whitelisted fields only, so a password hash can never leak here even by
// accident (unlike returning model.User directly and relying on a json tag).
type AdminUserSummary struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	IsActive   bool   `json:"is_active"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
}

type UserService interface {
	GetAllUsers() ([]AdminUserSummary, error)
}

type userService struct {
	repo userLister
}

func NewUserService(repo userLister) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetAllUsers() ([]AdminUserSummary, error) {

	users, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	summaries := make([]AdminUserSummary, 0, len(users))

	for _, u := range users {
		summaries = append(summaries, AdminUserSummary{
			ID:         u.ID.String(),
			FirstName:  u.FirstName,
			LastName:   u.LastName,
			Email:      u.Email,
			Phone:      u.Phone,
			Role:       u.Role,
			IsActive:   u.IsActive,
			IsVerified: u.IsVerified,
			CreatedAt:  u.CreatedAt.Format(time.RFC3339),
		})
	}

	return summaries, nil
}
