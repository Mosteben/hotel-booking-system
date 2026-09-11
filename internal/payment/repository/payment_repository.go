package repository

import (
	"github.com/Mosteben/hotel-booking-system/internal/payment/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	GetByID(id uint) (*model.Payment, error)
	GetByBookingID(bookingID uint) (*model.Payment, error)
	GetByUserID(userID string) ([]model.Payment, error)
	GetAll() ([]model.Payment, error)
	Update(payment *model.Payment) error

	// Create a repository that uses the provided transaction.
	WithTx(tx *gorm.DB) PaymentRepository
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) WithTx(
	tx *gorm.DB,
) PaymentRepository {
	return &paymentRepository{
		db: tx,
	}
}

func (r *paymentRepository) Create(
	payment *model.Payment,
) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) GetByID(
	id uint,
) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.
		First(&payment, id).
		Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) GetByBookingID(
	bookingID uint,
) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.
		Where("booking_id = ?", bookingID).
		First(&payment).
		Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *paymentRepository) GetByUserID(
	userID string,
) ([]model.Payment, error) {

	var payments []model.Payment

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&payments).
		Error

	return payments, err
}

func (r *paymentRepository) GetAll() (
	[]model.Payment,
	error,
) {

	var payments []model.Payment

	err := r.db.
		Order("created_at DESC").
		Find(&payments).
		Error

	return payments, err
}

func (r *paymentRepository) Update(
	payment *model.Payment,
) error {
	return r.db.Save(payment).Error
}
