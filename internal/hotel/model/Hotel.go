package model

import (
	"time"

	roomModel "github.com/Mosteben/hotel-booking-system/internal/room/model"
)

type Hotel struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Name        string           `gorm:"not null" json:"name"`
	Description string           `json:"description"`
	Address     string           `gorm:"not null" json:"address"`
	City        string           `gorm:"not null" json:"city"`
	Country     string           `gorm:"not null" json:"country"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	Stars       int              `gorm:"default:1" json:"stars"`
	Rooms       []roomModel.Room `gorm:"foreignKey:HotelID" json:"rooms"`
	Images      []HotelImage     `gorm:"foreignKey:HotelID" json:"images"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`

	// Not persisted - computed from Images after every fetch so every hotel
	// endpoint returns the same image_url field, whether or not the hotel
	// has any images yet.
	ImageURL *string `gorm:"-" json:"image_url"`
}

// ResolveImageURL fills ImageURL from Images: the image explicitly marked
// as main if one exists, otherwise the first uploaded image, otherwise nil
// (a hotel with no images yet must not break existing frontend rendering).
func (h *Hotel) ResolveImageURL() {
	for _, img := range h.Images {
		if img.IsMain {
			url := img.URL
			h.ImageURL = &url
			return
		}
	}

	if len(h.Images) > 0 {
		url := h.Images[0].URL
		h.ImageURL = &url
	}
}
