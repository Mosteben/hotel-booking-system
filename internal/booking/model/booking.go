package model

import "time"

type Booking struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"not null;index" json:"user_id"`
	RoomID     uint      `gorm:"not null;index" json:"room_id"`
	CheckIn    time.Time `gorm:"not null" json:"check_in"`
	CheckOut   time.Time `gorm:"not null" json:"check_out"`
	Guests     int       `gorm:"not null" json:"guests"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
	Status     string    `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
