package model

import "time"

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	HotelID   uint      `gorm:"not null;index" json:"hotel_id"`
	CreatedAt time.Time `json:"created_at"`
}
