package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	bookingRepository "github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	paymentModel "github.com/Mosteben/hotel-booking-system/internal/payment/model"
	"github.com/Mosteben/hotel-booking-system/internal/payment/repository"
)

var (
	ErrPaymentAlreadyExists = errors.New("payment already exists for this booking")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
	ErrBookingNotFound      = errors.New("booking not found")
	ErrNotBookingOwner      = errors.New("you are not allowed to pay for this booking")
	ErrBookingNotPending    = errors.New("only pending bookings can be paid")
)

type PaymentService interface {
	CreatePayment(
		payment *paymentModel.Payment,
	) error

	GetPaymentByID(
		id uint,
	) (*paymentModel.Payment, error)

	GetPaymentsByUserID(
		userID string,
	) ([]paymentModel.Payment, error)

	GetAllPayments() ([]paymentModel.Payment, error)

	UpdatePaymentStatus(
		id uint,
		status string,
	) error
}

type paymentService struct {
	repo        repository.PaymentRepository
	bookingRepo bookingRepository.BookingRepository
	db          *gorm.DB
}

func NewPaymentService(
	repo repository.PaymentRepository,
	bookingRepo bookingRepository.BookingRepository,
	db *gorm.DB,
) PaymentService {
	return &paymentService{
		repo:        repo,
		bookingRepo: bookingRepo,
		db:          db,
	}
}

// =========================
// Create Payment
// =========================

func (s *paymentService) CreatePayment(
	payment *paymentModel.Payment,
) error {

	if payment == nil {
		return errors.New("payment data is required")
	}

	if strings.TrimSpace(payment.UserID) == "" {
		return errors.New("user id is required")
	}

	if payment.BookingID == 0 {
		return errors.New("booking id is required")
	}

	paymentMethod := strings.ToLower(
		strings.TrimSpace(payment.PaymentMethod),
	)

	if paymentMethod != "cash" &&
		paymentMethod != "card" {
		return ErrInvalidPaymentMethod
	}

	booking, err := s.bookingRepo.GetByID(payment.BookingID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookingNotFound
		}

		return err
	}

	if booking.UserID != payment.UserID {
		return ErrNotBookingOwner
	}

	if booking.Status != "pending" {
		return ErrBookingNotPending
	}

	_, err = s.repo.GetByBookingID(payment.BookingID)

	if err == nil {
		return ErrPaymentAlreadyExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	payment.Amount = booking.TotalPrice
	payment.PaymentMethod = paymentMethod
	payment.Status = "pending"
	payment.TransactionID = nil

	return s.repo.Create(payment)
}

// =========================
// Get Payment By ID
// =========================

func (s *paymentService) GetPaymentByID(
	id uint,
) (*paymentModel.Payment, error) {

	if id == 0 {
		return nil, errors.New("invalid payment id")
	}

	return s.repo.GetByID(id)
}

// =========================
// Get My Payments
// =========================

func (s *paymentService) GetPaymentsByUserID(
	userID string,
) ([]paymentModel.Payment, error) {

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("invalid user id")
	}

	return s.repo.GetByUserID(userID)
}

// =========================
// Get All Payments
// =========================

func (s *paymentService) GetAllPayments() (
	[]paymentModel.Payment,
	error,
) {
	return s.repo.GetAll()
}

// =========================
// Update Payment Status
// =========================

func (s *paymentService) UpdatePaymentStatus(
	id uint,
	status string,
) error {

	if id == 0 {
		return errors.New("invalid payment id")
	}

	status = strings.ToLower(strings.TrimSpace(status))

	if status != "paid" &&
		status != "failed" &&
		status != "refunded" {
		return ErrInvalidPaymentStatus
	}

	// Start database transaction.
	return s.db.Transaction(func(tx *gorm.DB) error {

		// Create repositories using the same transaction.
		txPaymentRepo := s.repo.WithTx(tx)
		txBookingRepo := s.bookingRepo.WithTx(tx)

		// Get payment inside transaction.
		payment, err := txPaymentRepo.GetByID(id)

		if err != nil {
			return err
		}

		// Payment already has this status.
		if payment.Status == status {
			return errors.New(
				"payment already has this status",
			)
		}

		// Get related booking inside transaction.
		booking, err := txBookingRepo.GetByID(payment.BookingID)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrBookingNotFound
			}

			return err
		}

		// =========================
		// Pending
		// =========================

		if payment.Status == "pending" {

			// pending -> paid
			if status == "paid" {

				if booking.Status != "pending" {
					return errors.New(
						"only pending bookings can be confirmed by payment",
					)
				}

				// Generate transaction ID.
				transactionID := uuid.New().String()
				payment.TransactionID = &transactionID

				payment.Status = "paid"
				booking.Status = "confirmed"

				// Update booking first.
				if err := txBookingRepo.Update(booking); err != nil {
					return err
				}

				// Update payment second.
				if err := txPaymentRepo.Update(payment); err != nil {
					return err
				}

				return nil
			}

			// pending -> failed
			if status == "failed" {

				if booking.Status != "pending" {
					return errors.New(
						"failed payment requires a pending booking",
					)
				}

				payment.Status = "failed"

				return txPaymentRepo.Update(payment)
			}

			// pending -> refunded is not allowed.
			if status == "refunded" {
				return errors.New(
					"only paid payments can be refunded",
				)
			}
		}

		// =========================
		// Paid
		// =========================

		if payment.Status == "paid" {

			// paid -> refunded
			if status == "refunded" {

				if booking.Status != "confirmed" {
					return errors.New(
						"only confirmed bookings can be refunded",
					)
				}

				payment.Status = "refunded"
				booking.Status = "cancelled"

				// Update booking first.
				if err := txBookingRepo.Update(booking); err != nil {
					return err
				}

				// Update payment second.
				if err := txPaymentRepo.Update(payment); err != nil {
					return err
				}

				return nil
			}

			// paid -> failed is not allowed.
			if status == "failed" {
				return errors.New(
					"paid payment cannot be marked as failed",
				)
			}
		}

		// =========================
		// Failed
		// =========================

		if payment.Status == "failed" {

			// failed -> paid is not allowed.
			if status == "paid" {
				return errors.New(
					"failed payment cannot be marked as paid",
				)
			}

			// failed -> refunded is not allowed.
			if status == "refunded" {
				return errors.New(
					"failed payment cannot be refunded",
				)
			}
		}

		// =========================
		// Refunded
		// =========================

		if payment.Status == "refunded" {
			return errors.New(
				"refunded payment cannot be updated",
			)
		}

		return ErrInvalidPaymentStatus
	})
}

