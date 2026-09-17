package service

import (
	"strings"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
)

// AdminBookingSummary is what GET /bookings (admin/manager only) returns -
// the raw booking plus its user/hotel/room resolved server-side, so the
// Admin frontend never has to guess relationships from bare IDs or issue a
// lookup per row itself.
type AdminBookingSummary struct {
	ID         uint      `json:"id"`
	Status     string    `json:"status"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	Guests     int       `json:"guests"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Any of these can be nil - a booking's user_id/room_id might not
	// resolve to an existing record (e.g. stale/orphaned data from a
	// seed script, or a room/hotel deleted after the booking was made).
	// They are left nil rather than fabricated, so the Admin can tell
	// "no data" apart from "data we don't have".
	User  *AdminBookingUser  `json:"user"`
	Hotel *AdminBookingHotel `json:"hotel"`
	Room  *AdminBookingRoom  `json:"room"`

	// The room/hotel snapshot captured at booking time (see
	// applyBookingSnapshot). Unlike User/Hotel/Room above, this is never
	// nil/missing by nature - it's copied straight from the booking row -
	// but on bookings created before this snapshot existed (or otherwise
	// never captured), these are simply empty/zero, which the Admin UI
	// should treat as "no historical snapshot available", not as real data.
	HotelIDSnapshot       uint    `json:"hotel_id_snapshot"`
	HotelNameSnapshot     string  `json:"hotel_name_snapshot"`
	RoomNumberSnapshot    string  `json:"room_number_snapshot"`
	RoomTypeSnapshot      string  `json:"room_type_snapshot"`
	PricePerNightSnapshot float64 `json:"price_per_night_snapshot"`
}

type AdminBookingUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AdminBookingHotel struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AdminBookingRoom struct {
	ID         uint   `json:"id"`
	RoomNumber string `json:"room_number"`
	Type       string `json:"type"`
}

// resolveAdminBookingSummaries batch-fetches every user/room/hotel a set of
// bookings references (three queries total, regardless of booking count)
// instead of resolving relationships one booking at a time.
func (s *bookingService) resolveAdminBookingSummaries(
	bookings []model.Booking,
) ([]AdminBookingSummary, error) {

	userIDSet := make(map[string]struct{})
	roomIDSet := make(map[uint]struct{})

	for _, b := range bookings {
		userIDSet[b.UserID] = struct{}{}
		roomIDSet[b.RoomID] = struct{}{}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	roomIDs := make([]uint, 0, len(roomIDSet))
	for id := range roomIDSet {
		roomIDs = append(roomIDs, id)
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	usersByID := make(map[string]userModel.User, len(users))
	for _, u := range users {
		usersByID[u.ID.String()] = u
	}

	rooms, err := s.roomRepo.GetByIDs(roomIDs)
	if err != nil {
		return nil, err
	}

	roomsByID := make(map[uint]roomModel.Room, len(rooms))
	hotelIDSet := make(map[uint]struct{})
	for _, r := range rooms {
		roomsByID[r.ID] = r
		hotelIDSet[r.HotelID] = struct{}{}
	}

	hotelIDs := make([]uint, 0, len(hotelIDSet))
	for id := range hotelIDSet {
		hotelIDs = append(hotelIDs, id)
	}

	hotels, err := s.hotelRepo.GetByIDs(hotelIDs)
	if err != nil {
		return nil, err
	}

	hotelsByID := make(map[uint]hotelModel.Hotel, len(hotels))
	for _, h := range hotels {
		hotelsByID[h.ID] = h
	}

	summaries := make([]AdminBookingSummary, 0, len(bookings))

	for _, b := range bookings {
		summary := AdminBookingSummary{
			ID:         b.ID,
			Status:     b.Status,
			CheckIn:    b.CheckIn,
			CheckOut:   b.CheckOut,
			Guests:     b.Guests,
			TotalPrice: b.TotalPrice,
			CreatedAt:  b.CreatedAt,
			UpdatedAt:  b.UpdatedAt,

			HotelIDSnapshot:       b.HotelID,
			HotelNameSnapshot:     b.HotelName,
			RoomNumberSnapshot:    b.RoomNumber,
			RoomTypeSnapshot:      b.RoomType,
			PricePerNightSnapshot: b.PricePerNight,
		}

		if u, ok := usersByID[b.UserID]; ok {
			summary.User = &AdminBookingUser{
				ID:    u.ID.String(),
				Name:  strings.TrimSpace(u.FirstName + " " + u.LastName),
				Email: u.Email,
				Phone: u.Phone,
			}
		}

		if room, ok := roomsByID[b.RoomID]; ok {
			summary.Room = &AdminBookingRoom{
				ID:         room.ID,
				RoomNumber: room.RoomNumber,
				Type:       room.Type,
			}

			if hotel, ok := hotelsByID[room.HotelID]; ok {
				summary.Hotel = &AdminBookingHotel{
					ID:   hotel.ID,
					Name: hotel.Name,
				}
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
