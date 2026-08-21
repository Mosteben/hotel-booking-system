package middleware

import (
	"strings"

	"github.com/Mosteben/hotel-booking-system/pkg/jwt"
	"github.com/Mosteben/hotel-booking-system/pkg/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {

			response.Unauthorized(
				c,
				"Authorization header is required",
			)

			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {

			response.Unauthorized(
				c,
				"Invalid authorization header",
			)

			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwt.ParseToken(token)

		if err != nil {

			response.Unauthorized(
				c,
				"Invalid or expired token",
			)

			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}