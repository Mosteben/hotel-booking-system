package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"github.com/Mosteben/hotel-booking-system/internal/booking/service"
	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	service service.BookingService
}

func NewBookingHandler(service service.BookingService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

// =========================
// Create Booking
// =========================

func (h *BookingHandler) CreateBooking(c *gin.Context) {
	var booking model.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	booking.UserID = userID

	if err := h.service.CreateBooking(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, booking)
}

// =========================
// Get All Bookings
// =========================

func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	bookings, err := h.service.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// =========================
// Get Booking By ID
// =========================

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid booking id",
		})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	role, ok := roleValue.(string)
	if !ok || strings.TrimSpace(role) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user role",
		})
		return
	}

	booking, err := h.service.GetBookingByID(
		uint(id),
		userID,
		role,
	)

	if err != nil {
		if service.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}

		if err.Error() == "you are not allowed to view this booking" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// =========================
// Get My Bookings
// =========================

func (h *BookingHandler) GetMyBookings(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	bookings, err := h.service.GetBookingsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// =========================
// Update Booking
// =========================

func (h *BookingHandler) UpdateBooking(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid booking id",
		})
		return
	}

	var booking model.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	role, ok := roleValue.(string)
	if !ok || strings.TrimSpace(role) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user role",
		})
		return
	}

	booking.UserID = userID

	if err := h.service.UpdateBooking(
		uint(id),
		userID,
		&booking,
	); err != nil {

		if service.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}

		if err.Error() == "you are not allowed to update this booking" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	updatedBooking, err := h.service.GetBookingByID(
		uint(id),
		userID,
		role,
	)

	if err != nil {
		if service.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedBooking)
}

// =========================
// Delete / Cancel Booking
// =========================

func (h *BookingHandler) DeleteBooking(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid booking id",
		})
		return
	}

	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	if err := h.service.DeleteBooking(
		uint(id),
		userID,
	); err != nil {

		if service.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}

		if err.Error() == "you are not allowed to cancel this booking" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "booking cancelled successfully",
	})
}

// =========================
// Update Booking Status
// =========================

func (h *BookingHandler) UpdateBookingStatus(c *gin.Context) {
	id, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid booking id",
		})
		return
	}

	var request struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	request.Status = strings.ToLower(
		strings.TrimSpace(request.Status),
	)

	if request.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status is required",
		})
		return
	}

	if err := h.service.UpdateBookingStatus(
		uint(id),
		request.Status,
	); err != nil {

		if service.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get current user information.
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		return
	}

	roleValue, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user role not found",
		})
		return
	}

	role, ok := roleValue.(string)
	if !ok || strings.TrimSpace(role) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user role",
		})
		return
	}

	booking, err := h.service.GetBookingByID(
		uint(id),
		userID,
		role,
	)

	if err != nil {
		if service.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "booking not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// =========================
// Check Room Availability
// =========================

func (h *BookingHandler) CheckRoomAvailability(c *gin.Context) {
	// Get room ID from /rooms/:id/availability
	roomID, err := strconv.ParseUint(
		c.Param("id"),
		10,
		64,
	)

	if err != nil || roomID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid room id",
		})
		return
	}

	// Get dates from query parameters
	checkInStr := c.Query("check_in")
	checkOutStr := c.Query("check_out")

	if checkInStr == "" || checkOutStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "check_in and check_out are required",
		})
		return
	}

	// Parse check-in date
	checkIn, err := time.Parse(
		"2006-01-02",
		checkInStr,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid check_in date format, use YYYY-MM-DD",
		})
		return
	}

	// Parse check-out date
	checkOut, err := time.Parse(
		"2006-01-02",
		checkOutStr,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid check_out date format, use YYYY-MM-DD",
		})
		return
	}

	// Check date order
	if !checkOut.After(checkIn) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "check_out must be after check_in",
		})
		return
	}

	// Check availability
	available, err := h.service.IsRoomAvailable(
		uint(roomID),
		checkIn,
		checkOut,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room_id":   roomID,
		"check_in":  checkInStr,
		"check_out": checkOutStr,
		"available": available,
	})
}
