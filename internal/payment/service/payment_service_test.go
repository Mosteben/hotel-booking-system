package service

import (
	"errors"
	"testing"
	"time"

	bookingModel "github.com/Mosteben/hotel-booking-system/internal/booking/model"
	bookingRepository "github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	paymentModel "github.com/Mosteben/hotel-booking-system/internal/payment/model"
	paymentRepository "github.com/Mosteben/hotel-booking-system/internal/payment/repository"

	"gorm.io/gorm"
)

type mockPaymentRepository struct {
	payment         *paymentModel.Payment
	getByBookingErr error
	createCalled    bool
}

func (m *mockPaymentRepository) Create(
	payment *paymentModel.Payment,
) error {
	m.createCalled = true
	m.payment = payment
	return nil
}

func (m *mockPaymentRepository) GetByID(
	id uint,
) (*paymentModel.Payment, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m *mockPaymentRepository) GetByBookingID(
	bookingID uint,
) (*paymentModel.Payment, error) {
	return m.payment, m.getByBookingErr
}

func (m *mockPaymentRepository) GetByUserID(
	userID string,
) ([]paymentModel.Payment, error) {
	return nil, nil
}

func (m *mockPaymentRepository) GetAll() (
	[]paymentModel.Payment,
	error,
) {
	return nil, nil
}

func (m *mockPaymentRepository) Update(
	payment *paymentModel.Payment,
) error {
	return nil
}

func (m *mockPaymentRepository) WithTx(
	tx *gorm.DB,
) paymentRepository.PaymentRepository {
	return m
}

type mockBookingRepository struct {
	booking    *bookingModel.Booking
	getByIDErr error
}

func (m *mockBookingRepository) Create(
	booking *bookingModel.Booking,
) error {
	return nil
}

func (m *mockBookingRepository) GetAll() (
	[]bookingModel.Booking,
	error,
) {
	return nil, nil
}

func (m *mockBookingRepository) GetByID(
	id uint,
) (*bookingModel.Booking, error) {
	return m.booking, m.getByIDErr
}

func (m *mockBookingRepository) GetByUserID(
	userID string,
) ([]bookingModel.Booking, error) {
	return nil, nil
}

func (m *mockBookingRepository) Update(
	booking *bookingModel.Booking,
) error {
	return nil
}

func (m *mockBookingRepository) Delete(
	id uint,
) error {
	return nil
}

func (m *mockBookingRepository) IsRoomAvailable(
	roomID uint,
	checkIn time.Time,
	checkOut time.Time,
	excludeBookingID uint,
) (bool, error) {
	return true, nil
}

func (m *mockBookingRepository) WithTx(
	tx *gorm.DB,
) bookingRepository.BookingRepository {
	return m
}

// ========================================
// Test 1: Booking Not Found
// ========================================

func TestCreatePayment_BookingNotFound(t *testing.T) {
	paymentRepo := &mockPaymentRepository{
		getByBookingErr: gorm.ErrRecordNotFound,
	}

	bookingRepo := &mockBookingRepository{
		getByIDErr: gorm.ErrRecordNotFound,
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     999999,
		UserID:        "test-user",
		PaymentMethod: "card",
	}

	err := service.CreatePayment(payment)

	if !errors.Is(err, ErrBookingNotFound) {
		t.Fatalf(
			"expected ErrBookingNotFound, got %v",
			err,
		)
	}

	if paymentRepo.createCalled {
		t.Fatal("payment should not be created")
	}
}

// ========================================
// Test 2: Not Booking Owner
// ========================================

func TestCreatePayment_NotBookingOwner(t *testing.T) {
	paymentRepo := &mockPaymentRepository{
		getByBookingErr: gorm.ErrRecordNotFound,
	}

	bookingRepo := &mockBookingRepository{
		booking: &bookingModel.Booking{
			ID:     100,
			UserID: "real-owner",
			Status: "pending",
		},
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     100,
		UserID:        "different-user",
		PaymentMethod: "card",
	}

	err := service.CreatePayment(payment)

	if !errors.Is(err, ErrNotBookingOwner) {
		t.Fatalf(
			"expected ErrNotBookingOwner, got %v",
			err,
		)
	}

	if paymentRepo.createCalled {
		t.Fatal("payment should not be created")
	}
}
func TestCreatePayment_InvalidPaymentMethod(t *testing.T) {
	paymentRepo := &mockPaymentRepository{
		getByBookingErr: gorm.ErrRecordNotFound,
	}

	bookingRepo := &mockBookingRepository{
		booking: &bookingModel.Booking{
			ID:     101,
			UserID: "test-user",
			Status: "pending",
		},
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     101,
		UserID:        "test-user",
		PaymentMethod: "bitcoin",
	}

	err := service.CreatePayment(payment)

	if !errors.Is(err, ErrInvalidPaymentMethod) {
		t.Fatalf(
			"expected ErrInvalidPaymentMethod, got %v",
			err,
		)
	}

	if paymentRepo.createCalled {
		t.Fatal("payment should not be created")
	}
}
func TestCreatePayment_BookingNotPending(t *testing.T) {
	paymentRepo := &mockPaymentRepository{
		getByBookingErr: gorm.ErrRecordNotFound,
	}

	bookingRepo := &mockBookingRepository{
		booking: &bookingModel.Booking{
			ID:     102,
			UserID: "test-user",
			Status: "confirmed",
		},
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     102,
		UserID:        "test-user",
		PaymentMethod: "card",
	}

	err := service.CreatePayment(payment)

	if !errors.Is(err, ErrBookingNotPending) {
		t.Fatalf(
			"expected ErrBookingNotPending, got %v",
			err,
		)
	}

	if paymentRepo.createCalled {
		t.Fatal("payment should not be created")
	}
}
func TestCreatePayment_UsesBookingTotalPrice(t *testing.T) {
	paymentRepo := &mockPaymentRepository{
		getByBookingErr: gorm.ErrRecordNotFound,
	}

	bookingRepo := &mockBookingRepository{
		booking: &bookingModel.Booking{
			ID:         103,
			UserID:     "test-user",
			Status:     "pending",
			TotalPrice: 9000,
		},
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     103,
		UserID:        "test-user",
		Amount:        1,
		PaymentMethod: "card",
	}

	err := service.CreatePayment(payment)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !paymentRepo.createCalled {
		t.Fatal("expected payment to be created")
	}

	if paymentRepo.payment.Amount != 9000 {
		t.Fatalf(
			"expected amount 9000, got %v",
			paymentRepo.payment.Amount,
		)
	}
}
func TestCreatePayment_AlreadyExists(t *testing.T) {
	existingPayment := &paymentModel.Payment{
		ID:        10,
		BookingID: 104,
		UserID:    "test-user",
		Status:    "pending",
	}

	paymentRepo := &mockPaymentRepository{
		payment:         existingPayment,
		getByBookingErr: nil,
	}

	bookingRepo := &mockBookingRepository{
		booking: &bookingModel.Booking{
			ID:         104,
			UserID:     "test-user",
			Status:     "pending",
			TotalPrice: 5000,
		},
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     104,
		UserID:        "test-user",
		PaymentMethod: "card",
	}

	err := service.CreatePayment(payment)

	if !errors.Is(err, ErrPaymentAlreadyExists) {
		t.Fatalf(
			"expected ErrPaymentAlreadyExists, got %v",
			err,
		)
	}

	if paymentRepo.createCalled {
		t.Fatal("payment should not be created")
	}
}
func TestCreatePayment_Success(t *testing.T) {
	paymentRepo := &mockPaymentRepository{
		getByBookingErr: gorm.ErrRecordNotFound,
	}

	bookingRepo := &mockBookingRepository{
		booking: &bookingModel.Booking{
			ID:         105,
			UserID:     "test-user",
			Status:     "pending",
			TotalPrice: 7000,
		},
	}

	service := &paymentService{
		repo:        paymentRepo,
		bookingRepo: bookingRepo,
	}

	payment := &paymentModel.Payment{
		BookingID:     105,
		UserID:        "test-user",
		PaymentMethod: "card",
	}

	err := service.CreatePayment(payment)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !paymentRepo.createCalled {
		t.Fatal("payment was not created")
	}

	if payment.Amount != 7000 {
		t.Fatalf(
			"expected amount 7000, got %v",
			payment.Amount,
		)
	}

	if payment.Status != "pending" {
		t.Fatalf(
			"expected status pending, got %s",
			payment.Status,
		)
	}

	if payment.TransactionID != nil {
		t.Fatal("transaction id should be nil")
	}
}
