package routes

import (
	"net/http"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	hotelHandler "github.com/Mosteben/hotel-booking-system/internal/hotel/handler"
	"github.com/Mosteben/hotel-booking-system/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	auth *authHandler.AuthHandler,
	hotel *hotelHandler.HotelHandler,
) {

	// =========================
	// Public Routes
	// =========================

	r.GET("/health", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})

	})

	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)

	// =========================
	// Hotel Routes
	// =========================

	r.GET("/hotels", hotel.GetAllHotels)
	r.GET("/hotels/:id", hotel.GetHotelByID)

	r.POST(
		"/hotels",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		hotel.CreateHotel,
	)

	r.PUT(
		"/hotels/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		hotel.UpdateHotel,
	)

	r.DELETE(
		"/hotels/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin"),
		hotel.DeleteHotel,
	)

	// =========================
	// Protected Routes
	// =========================

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	{
		authorized.GET("/auth/me", auth.Me)
		authorized.PUT("/auth/profile", auth.UpdateProfile)
		authorized.PUT("/auth/change-password", auth.ChangePassword)
	}
}
