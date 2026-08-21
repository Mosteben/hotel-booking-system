package model

import (
	"github.com/google/uuid"
	"time"
)

type Profile struct {

	ID uint `gorm:"primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex"`

	DateOfBirth time.Time

	Gender string `gorm:"size:20"`

	Nationality string `gorm:"size:100"`

	NationalID string `gorm:"size:30"`

	PassportNumber string `gorm:"size:30"`

	Address string

	City string

	State string

	Country string

	PostalCode string

	Avatar string

	Bio string

	CreatedAt time.Time

	UpdatedAt time.Time
}