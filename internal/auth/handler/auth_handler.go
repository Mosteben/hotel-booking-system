package handler

import (
	"github.com/Mosteben/hotel-booking-system/internal/auth/model"
	authService "github.com/Mosteben/hotel-booking-system/internal/auth/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
	validatorPkg "github.com/Mosteben/hotel-booking-system/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *authService.AuthService
}

func NewAuthHandler(
	service *authService.AuthService,
) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with profile information.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.RegisterRequest true "Registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]interface{} "Validation failed"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"Validation failed",
			response.InvalidRequestMessage,
		)
		return
	}

	err := h.service.Register(req)

	if err != nil {
		if fields := validatorPkg.FieldErrors(err); fields != nil {
			response.BadRequest(c, "validation failed", fields)
			return
		}
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

// Login godoc
// @Summary Login user
// @Description Authenticate a user and return a JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 400 {object} map[string]interface{} "Validation failed"
// @Failure 401 {object} map[string]interface{} "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"Validation failed",
			response.InvalidRequestMessage,
		)
		return
	}

	token, err := h.service.Login(req)

	if err != nil {
		if validatorPkg.FieldErrors(err) != nil {
			response.Unauthorized(c, "validation failed")
			return
		}
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

// Me godoc
// @Summary Get current user
// @Description Return the authenticated user's information and profile.
// @Tags Authentication
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Current user"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /auth/me [get]
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
			"id":          user.ID,
			"first_name":  user.FirstName,
			"last_name":   user.LastName,
			"email":       user.Email,
			"phone":       user.Phone,
			"role":        user.Role,
			"is_active":   user.IsActive,
			"is_verified": user.IsVerified,
			"profile":     user.Profile,
		},
	)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information.
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation failed"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Router /auth/profile [put]
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

	var req model.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"Validation failed",
			response.InvalidRequestMessage,
		)
		return
	}

	err := h.service.UpdateProfile(
		userID,
		req,
	)

	if err != nil {
		if fields := validatorPkg.FieldErrors(err); fields != nil {
			response.BadRequest(c, "validation failed", fields)
			return
		}
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

// ChangePassword godoc
// @Summary Change password
// @Description Change the authenticated user's password.
// @Tags Authentication
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body model.ChangePasswordRequest true "Password change data"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Validation failed"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Router /auth/change-password [put]
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

	var req model.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(
			c,
			"Validation failed",
			response.InvalidRequestMessage,
		)
		return
	}

	err := h.service.ChangePassword(
		userID,
		req,
	)

	if err != nil {
		if fields := validatorPkg.FieldErrors(err); fields != nil {
			response.BadRequest(c, "validation failed", fields)
			return
		}
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
