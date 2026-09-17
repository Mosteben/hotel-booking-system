package service

import (
	"strings"
	"time"

	bookingModel "github.com/Mosteben/hotel-booking-system/internal/booking/model"
	hotelModel "github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	paymentModel "github.com/Mosteben/hotel-booking-system/internal/payment/model"
	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
)

// AdminPaymentSummary is what GET /payments (admin/manager only) returns -
// the raw payment plus the payer, booking, hotel, and room resolved
// server-side.
type AdminPaymentSummary struct {
	ID            uint      `json:"id"`
	BookingID     uint      `json:"booking_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	TransactionID *string   `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Any of these can be nil - a payment's user_id/booking_id, or the
	// booking's room_id, might not resolve to an existing record. Left
	// nil rather than fabricated.
	User    *AdminPaymentUser    `json:"user"`
	Booking *AdminPaymentBooking `json:"booking"`
	Hotel   *AdminPaymentHotel   `json:"hotel"`
	Room    *AdminPaymentRoom    `json:"room"`
}

type AdminPaymentUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type AdminPaymentBooking struct {
	ID       uint      `json:"id"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
	Status   string    `json:"status"`
}

type AdminPaymentHotel struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AdminPaymentRoom struct {
	ID         uint   `json:"id"`
	RoomNumber string `json:"room_number"`
	Type       string `json:"type"`
}

func (s *paymentService) resolveAdminPaymentSummaries(
	payments []paymentModel.Payment,
) ([]AdminPaymentSummary, error) {

	userIDSet := make(map[string]struct{})
	bookingIDSet := make(map[uint]struct{})

	for _, p := range payments {
		userIDSet[p.UserID] = struct{}{}
		bookingIDSet[p.BookingID] = struct{}{}
	}

	userIDs := make([]string, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	bookingIDs := make([]uint, 0, len(bookingIDSet))
	for id := range bookingIDSet {
		bookingIDs = append(bookingIDs, id)
	}

	users, err := s.userRepo.GetByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	usersByID := make(map[string]userModel.User, len(users))
	for _, u := range users {
		usersByID[u.ID.String()] = u
	}

	bookings, err := s.bookingRepo.GetByIDs(bookingIDs)
	if err != nil {
		return nil, err
	}

	bookingsByID := make(map[uint]bookingModel.Booking, len(bookings))
	roomIDSet := make(map[uint]struct{})
	for _, b := range bookings {
		bookingsByID[b.ID] = b
		roomIDSet[b.RoomID] = struct{}{}
	}

	roomIDs := make([]uint, 0, len(roomIDSet))
	for id := range roomIDSet {
		roomIDs = append(roomIDs, id)
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

	summaries := make([]AdminPaymentSummary, 0, len(payments))

	for _, p := range payments {
		summary := AdminPaymentSummary{
			ID:            p.ID,
			BookingID:     p.BookingID,
			Amount:        p.Amount,
			PaymentMethod: p.PaymentMethod,
			Status:        p.Status,
			TransactionID: p.TransactionID,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}

		if u, ok := usersByID[p.UserID]; ok {
			summary.User = &AdminPaymentUser{
				ID:    u.ID.String(),
				Name:  strings.TrimSpace(u.FirstName + " " + u.LastName),
				Email: u.Email,
				Phone: u.Phone,
			}
		}

		if b, ok := bookingsByID[p.BookingID]; ok {
			summary.Booking = &AdminPaymentBooking{
				ID:       b.ID,
				CheckIn:  b.CheckIn,
				CheckOut: b.CheckOut,
				Status:   b.Status,
			}

			if room, ok := roomsByID[b.RoomID]; ok {
				summary.Room = &AdminPaymentRoom{
					ID:         room.ID,
					RoomNumber: room.RoomNumber,
					Type:       room.Type,
				}

				if hotel, ok := hotelsByID[room.HotelID]; ok {
					summary.Hotel = &AdminPaymentHotel{
						ID:   hotel.ID,
						Name: hotel.Name,
					}
				}
			}
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}
