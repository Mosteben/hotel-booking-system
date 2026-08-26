package model

import "time"

type Room struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	HotelID       uint      `gorm:"not null;index" json:"hotel_id"`
	RoomNumber    string    `gorm:"not null" json:"room_number"`
	Type          string    `gorm:"not null" json:"type"`
	Description   string    `json:"description"`
	PricePerNight float64   `gorm:"not null" json:"price_per_night"`
	Capacity      int       `gorm:"not null" json:"capacity"`
	Status        string    `gorm:"not null;default:'available'" json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
