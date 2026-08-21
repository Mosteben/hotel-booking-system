package handler

import (
	authModel "github.com/Mosteben/hotel-booking-system/internal/auth/model"
	authService "github.com/Mosteben/hotel-booking-system/internal/auth/service"
	 "github.com/Mosteben/hotel-booking-system/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *authService.AuthService
}

func NewAuthHandler(service *authService.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req authModel.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.BadRequest(
			c,
			"Validation failed",
			err.Error(),
		)

		return
	}

	err := h.service.Register(req)

	if err != nil {

		response.BadRequest(
			c,
			err.Error(),
			nil,
		)

		return
	}

	response.Created(
		c,
		"User registered successfully",
		nil,
	)
}
func (h *AuthHandler) Login(c *gin.Context) {

	var req authModel.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.BadRequest(
			c,
			"Validation failed",
			err.Error(),
		)

		return
	}

	token, err := h.service.Login(req)

	if err != nil {

		response.Unauthorized(
			c,
			err.Error(),
		)

		return
	}

	response.OK(
		c,
		"Login successful",
		gin.H{
			"token": token,
		},
	)
}
func (h *AuthHandler) Me(c *gin.Context) {

	userIDValue, exists := c.Get("userID")

	if !exists {
		response.Unauthorized(
			c,
			"User not authenticated",
		)

		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		response.Unauthorized(
			c,
			"Invalid user ID",
		)

		return
	}

	user, err := h.service.GetCurrentUser(userID)

	if err != nil {
		response.NotFound(
			c,
			"User not found",
		)

		return
	}

	response.OK(
		c,
		"Current user",
		gin.H{
			"id":            user.ID,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
			"email":         user.Email,
			"phone":         user.Phone,
			"role":          user.Role,
			"is_active":     user.IsActive,
			"is_verified":   user.IsVerified,
			"profile":       user.Profile,
		},
	)
}
func (h *AuthHandler) UpdateProfile(c *gin.Context) {

	userIDValue, exists := c.Get("userID")

	if !exists {
		response.Unauthorized(
			c,
			"User not authenticated",
		)

		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		response.Unauthorized(
			c,
			"Invalid user ID",
		)

		return
	}

	var req authModel.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.BadRequest(
			c,
			"Validation failed",
			err.Error(),
		)

		return
	}

	err := h.service.UpdateProfile(userID, req)

	if err != nil {

		response.BadRequest(
			c,
			err.Error(),
			nil,
		)

		return
	}

	response.OK(
		c,
		"Profile updated successfully",
		nil,
	)
}
func (h *AuthHandler) ChangePassword(c *gin.Context) {

	userIDValue, exists := c.Get("userID")

	if !exists {
		response.Unauthorized(
			c,
			"User not authenticated",
		)
		return
	}

	userID, ok := userIDValue.(string)

	if !ok || userID == "" {
		response.Unauthorized(
			c,
			"Invalid user ID",
		)
		return
	}

	var req authModel.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.BadRequest(
			c,
			"Validation failed",
			err.Error(),
		)

		return
	}

	err := h.service.ChangePassword(userID, req)

	if err != nil {

		response.BadRequest(
			c,
			err.Error(),
			nil,
		)

		return
	}

	response.OK(
		c,
		"Password changed successfully",
		nil,
	)
}