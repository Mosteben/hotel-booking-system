package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Mosteben/hotel-booking-system/internal/payment/model"
	"github.com/Mosteben/hotel-booking-system/internal/payment/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

// =========================
// Create Payment
// =========================

// CreatePayment godoc
// @Summary Create payment
// @Description Creates a payment for the authenticated user's booking.
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Booking ID"
// @Param request body object{payment_method=string} true "Payment data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /bookings/{id}/payment [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {

	bookingID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || bookingID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid booking id",
		})
		return
	}

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
		})
		return
	}

	var request struct {
		PaymentMethod string `json:"payment_method"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid payment data",
			"error":   err.Error(),
		})
		return
	}

	payment := &model.Payment{
		BookingID:     uint(bookingID),
		UserID:        userID,
		PaymentMethod: request.PaymentMethod,
	}

	if err := h.service.CreatePayment(payment); err != nil {

		switch {
		case errors.Is(err, service.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})

		case errors.Is(err, service.ErrNotBookingOwner):
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": err.Error(),
			})

		case errors.Is(err, service.ErrPaymentAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": err.Error(),
			})

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "payment created successfully",
		"data":    payment,
	})
}

// =========================
// Get Payment By ID
// =========================

// GetPaymentByID godoc
// @Summary Get payment by ID
// @Description Retrieves a payment by ID. Customers can only access their own payments. Admins and managers can access any payment.
// @Tags Payments
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPaymentByID(c *gin.Context) {

	paymentID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || paymentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid payment id",
		})
		return
	}

	payment, err := h.service.GetPaymentByID(
		uint(paymentID),
	)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "payment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get payment",
			"error":   err.Error(),
		})
		return
	}

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
		})
		return
	}

	roleValue, _ := c.Get("role")
	role, _ := roleValue.(string)

	if role != "admin" &&
		role != "manager" &&
		payment.UserID != userID {

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "you are not allowed to view this payment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payment retrieved successfully",
		"data":    payment,
	})
}

// =========================
// Get My Payments
// =========================

// GetMyPayments godoc
// @Summary Get my payments
// @Description Retrieves all payments belonging to the authenticated user.
// @Tags Payments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments/my [get]
func (h *PaymentHandler) GetMyPayments(c *gin.Context) {

	userIDValue, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
		})
		return
	}

	payments, err := h.service.GetPaymentsByUserID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get payments",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payments retrieved successfully",
		"data":    payments,
	})
}

// =========================
// Get All Payments
// =========================

// GetAllPayments godoc
// @Summary Get all payments
// @Description Retrieves all payments. Admin and manager only.
// @Tags Payments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments [get]
func (h *PaymentHandler) GetAllPayments(c *gin.Context) {

	payments, err := h.service.GetAllPayments()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get payments",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "all payments retrieved successfully",
		"data":    payments,
	})
}

// =========================
// Update Payment Status
// =========================

// UpdatePaymentStatus godoc
// @Summary Update payment status
// @Description Updates a payment status. Admin and manager only. Supported statuses are paid, failed, and refunded.
// @Tags Payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Param request body object{status=string} true "Payment status"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /payments/{id}/status [patch]
func (h *PaymentHandler) UpdatePaymentStatus(c *gin.Context) {

	paymentID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || paymentID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid payment id",
		})
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid payment status data",
			"error":   err.Error(),
		})
		return
	}

	if err := h.service.UpdatePaymentStatus(
		uint(paymentID),
		request.Status,
	); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "payment not found",
			})
			return
		}

		if errors.Is(err, service.ErrBookingNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	payment, err := h.service.GetPaymentByID(
		uint(paymentID),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "payment status updated but failed to retrieve payment",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payment status updated successfully",
		"data":    payment,
	})
}
