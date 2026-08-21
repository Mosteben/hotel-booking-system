package middleware

import (
	"github.com/Mosteben/hotel-booking-system/pkg/response"
	"github.com/gin-gonic/gin"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {

	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")

		if !exists {
			response.Forbidden(
				c,
				"User role not found",
			)

			c.Abort()
			return
		}

		role, ok := roleValue.(string)

		if !ok || role == "" {
			response.Forbidden(
				c,
				"Invalid user role",
			)

			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {

			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Forbidden(
			c,
			"You do not have permission to perform this action",
		)

		c.Abort()
	}
}