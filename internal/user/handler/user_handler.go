package handler

import (
	"net/http"

	"github.com/Mosteben/hotel-booking-system/internal/user/service"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(
	service service.UserService,
) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// ListUsers godoc
// @Summary List all users
// @Description Retrieve all registered users with safe, non-sensitive fields. Admin/Manager only.
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure 401 {object} map[string]interface{} "User not authenticated"
// @Failure 403 {object} map[string]interface{} "Access denied"
// @Failure 500 {object} map[string]interface{} "Failed to get users"
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to get users",
			response.SanitizedError("ListUsers", err),
		)
		return
	}

	response.OK(c, "users retrieved successfully", users)
}
