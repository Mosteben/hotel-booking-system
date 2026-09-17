package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Mosteben/hotel-booking-system/internal/booking/model"
	"github.com/Mosteben/hotel-booking-system/internal/booking/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
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

// CreateBooking godoc
// @Summary Create a booking
// @Description Create a new hotel room booking for the authenticated user.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param request body model.Booking true "Booking data"
// @Success 201 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /bookings [post]
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

// GetAllBookings godoc
// @Summary Get all bookings
// @Description Get all bookings. Admin and manager users only.
// @Tags Bookings
// @Produce json
// @Success 200 {array} model.Booking
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings [get]
func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	bookings, err := h.service.GetAllBookings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": response.SanitizedError("GetAllBookings", err),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// =========================
// Get Booking By ID
// =========================

// GetBookingByID godoc
// @Summary Get booking by ID
// @Description Get a booking by ID. Customers can only view their own bookings, while admins and managers can view any booking.
// @Tags Bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/{id} [get]
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
			"error": response.SanitizedError("GetBookingByID", err),
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// =========================
// Get My Bookings
// =========================

// GetMyBookings godoc
// @Summary Get my bookings
// @Description Get all bookings belonging to the authenticated user.
// @Tags Bookings
// @Produce json
// @Success 200 {array} model.Booking
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/my [get]
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
			"error": response.SanitizedError("GetMyBookings", err),
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// =========================
// Update Booking
// =========================

// UpdateBooking godoc
// @Summary Update a booking
// @Description Update an existing booking. Customers can only update their own bookings.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param request body model.Booking true "Updated booking data"
// @Success 200 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/{id} [put]
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
			"error": response.SanitizedError("UpdateBooking", err),
		})
		return
	}

	c.JSON(http.StatusOK, updatedBooking)
}

// =========================
// Delete / Cancel Booking
// =========================

// DeleteBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking belonging to the authenticated user. The booking is marked as cancelled instead of being deleted.
// @Tags Bookings
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/{id} [delete]
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

// UpdateBookingStatus godoc
// @Summary Update booking status
// @Description Confirm or cancel a booking. Admin and manager users only.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Param request body object{status=string} true "Booking status"
// @Success 200 {object} model.Booking
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /bookings/{id}/status [patch]
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
			"error": response.SanitizedError("UpdateBookingStatus", err),
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

// =========================
// Check Room Availability
// =========================

// CheckRoomAvailability godoc
// @Summary Check room availability
// @Description Check whether a room is available for the specified date range.
// @Tags Room Availability
// @Produce json
// @Param id path int true "Room ID"
// @Param check_in query string true "Check-in date in YYYY-MM-DD format"
// @Param check_out query string true "Check-out date in YYYY-MM-DD format"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /rooms/{id}/availability [get]
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
			"error": response.SanitizedError("CheckRoomAvailability", err),
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
