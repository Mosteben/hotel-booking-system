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

	// Snapshot of the room/hotel as they existed at booking creation (or
	// last room change via UpdateBooking) time - set once server-side from
	// the real Room/Hotel rows, never re-fetched afterward. This is what
	// keeps a booking's historical meaning intact even if the room or hotel
	// is later deleted: RoomID/room lookups can go stale, but these columns
	// don't. Deliberately plain scalar columns with no GORM association/FK
	// to Room or Hotel - adding one here would make AutoMigrate attempt a
	// foreign key against the current database, which already has orphaned
	// room_id/user_id rows and would break on that.
	//
	// `default` avoids these ever being SQL NULL - including on the
	// existing bookings this ADD COLUMN backfills - so scanning a legacy
	// row (created before this field existed) always yields a plain empty/
	// zero value rather than a NULL-scan error. An empty value on an old
	// booking means "no snapshot was captured for this legacy record", not
	// "the room is free" - it is never treated as if it means the latter.
	HotelID       uint    `gorm:"default:0" json:"hotel_id"`
	HotelName     string  `gorm:"default:''" json:"hotel_name"`
	RoomNumber    string  `gorm:"default:''" json:"room_number"`
	RoomType      string  `gorm:"default:''" json:"room_type"`
	PricePerNight float64 `gorm:"default:0" json:"price_per_night"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
