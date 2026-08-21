package repository

import (
	profileModel "github.com/Mosteben/hotel-booking-system/internal/profile/model"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) Create(profile *profileModel.Profile) error {
	return r.db.Create(profile).Error
}

func (r *ProfileRepository) Update(profile *profileModel.Profile) error {
	return r.db.Save(profile).Error
}

func (r *ProfileRepository) FindByUserID(userID string) (*profileModel.Profile, error) {

	var profile profileModel.Profile

	err := r.db.
		Where("user_id = ?", userID).
		First(&profile).Error

	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfileRepository) CreateTx(tx *gorm.DB, profile *profileModel.Profile,) error {
	return tx.Create(profile).Error
}

func (r *ProfileRepository) UpdateTx(tx *gorm.DB,profile *profileModel.Profile,) error {
	return tx.Save(profile).Error
}