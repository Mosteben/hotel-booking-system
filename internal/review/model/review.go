package model

import "time"

type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HotelID   uint      `gorm:"not null;index" json:"hotel_id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Rating    int       `gorm:"not null" json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
