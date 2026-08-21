package model

import "time"

type Hotel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Address     string    `gorm:"not null" json:"address"`
	City        string    `gorm:"not null" json:"city"`
	Country     string    `gorm:"not null" json:"country"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Stars       int       `gorm:"default:1" json:"stars"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}