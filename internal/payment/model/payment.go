package model

import "time"

type Payment struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	BookingID     uint      `gorm:"not null;uniqueIndex" json:"booking_id"`
	UserID        string    `gorm:"not null;index" json:"user_id"`
	Amount        float64   `gorm:"not null" json:"amount"`
	PaymentMethod string    `gorm:"not null" json:"payment_method"`
	Status        string    `gorm:"not null;default:pending" json:"status"`
	TransactionID *string   `gorm:"uniqueIndex" json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

