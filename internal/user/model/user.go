package model

import (
	profileModel "github.com/Mosteben/hotel-booking-system/internal/profile/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	FirstName string `gorm:"size:100;not null"`

	LastName string `gorm:"size:100;not null"`

	Email string `gorm:"size:255;uniqueIndex;not null"`

	Password string `gorm:"not null"`

	Phone string `gorm:"size:20;uniqueIndex;not null"`

	Role string `gorm:"type:varchar(20);default:'customer'"`

	IsVerified bool `gorm:"default:false"`

	IsActive bool `gorm:"default:true"`

	LastLoginAt *time.Time

	Profile profileModel.Profile `gorm:"foreignKey:UserID"`

	CreatedAt time.Time

	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}