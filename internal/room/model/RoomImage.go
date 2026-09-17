package model

import "time"

type RoomImage struct {
	ID uint `gorm:"primaryKey" json:"id"`

	RoomID uint `gorm:"not null;index" json:"room_id"`

	URL string `gorm:"not null" json:"url"`

	// The storage provider's identifier for this file, needed to delete it
	// from storage later. Internal detail - never exposed to the frontend.
	PublicID string `gorm:"not null" json:"-"`

	IsMain bool `gorm:"default:false" json:"is_main"`

	CreatedAt time.Time `json:"created_at"`
}
