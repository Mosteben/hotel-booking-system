package service

import (
	"strings"
	"time"

	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	"github.com/Mosteben/hotel-booking-system/internal/review/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
)

// AdminReviewSummary is what GET /reviews (admin/manager only) returns -
// the raw review plus its author and hotel resolved server-side.
type AdminReviewSummary struct {
	ID        uint      `json:"id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Nil when the review's user_id/hotel_id doesn't resolve to an
	// existing record - never fabricated.
	User  *AdminReviewUser  `json:"user"`
	Hotel *AdminReviewHotel `json:"hotel"`
}

type AdminReviewUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AdminReviewHotel struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (s *reviewService) resolveAdminReviewSummaries(
	reviews []model.Review,
) ([]AdminReviewSummary, error) {

	userIDSet := make(map[string]struct{})
	hotelIDSet := make(map[uint]struct{})

	for _, r := range reviews {
		userIDSet[r.UserID] = struct{}{}
		hotelIDSet[r.HotelID] = struct{}{}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	hotelIDs := make([]uint, 0, len(hotelIDSet))
	for id := range hotelIDSet {
		hotelIDs = append(hotelIDs, id)
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	usersByID := make(map[string]userModel.User, len(users))
	for _, u := range users {
		usersByID[u.ID.String()] = u
	}

	hotels, err := s.hotelRepo.GetByIDs(hotelIDs)
	if err != nil {
		return nil, err
	}

	hotelsByID := make(map[uint]hotelModel.Hotel, len(hotels))
	for _, h := range hotels {
		hotelsByID[h.ID] = h
	}

	summaries := make([]AdminReviewSummary, 0, len(reviews))

	for _, r := range reviews {
		summary := AdminReviewSummary{
			ID:        r.ID,
			Rating:    r.Rating,
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}

		if u, ok := usersByID[r.UserID]; ok {
			summary.User = &AdminReviewUser{
				ID:    u.ID.String(),
				Name:  strings.TrimSpace(u.FirstName + " " + u.LastName),
				Email: u.Email,
				Phone: u.Phone,
			}
		}

		if h, ok := hotelsByID[r.HotelID]; ok {
			summary.Hotel = &AdminReviewHotel{ID: h.ID, Name: h.Name}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
