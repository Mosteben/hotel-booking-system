package repository

import (
	"errors"

	"github.com/Mosteben/hotel-booking-system/internal/user/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {

	var user model.User

	err := r.db.Preload("Profile").
		First(&user, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {

	var user model.User

	err := r.db.
		Where("email = ?", email).
		Preload("Profile").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ExistsByEmail(email string) (bool, error) {

	var user model.User

	err := r.db.
		Select("id").
		Where("email = ?", email).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}

func (r *UserRepository) ExistsByPhone(phone string) (bool, error) {

	var user model.User

	err := r.db.
		Select("id").
		Where("phone = ?", phone).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return err == nil, err
}
func (r *UserRepository) CreateTx(tx *gorm.DB, user *model.User) error {
	return tx.Create(user).Error
}
func (r *UserRepository) UpdateTx(tx *gorm.DB,user *model.User,) error {
	return tx.Save(user).Error
}
