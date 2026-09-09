package routes

import (
	"net/http"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	hotelHandler "github.com/Mosteben/hotel-booking-system/internal/hotel/handler"
	reviewHandler "github.com/Mosteben/hotel-booking-system/internal/review/handler"
	"github.com/Mosteben/hotel-booking-system/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	auth *authHandler.AuthHandler,
	hotel *hotelHandler.HotelHandler,
	review *reviewHandler.ReviewHandler,
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

	// Search hotels
	//
	// IMPORTANT:
	// This route must come before /hotels/:id
	// because /hotels/:id is a wildcard route.
	r.GET("/hotels/search", hotel.SearchHotels)

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
	// Review Routes
	// =========================

	// Create review
	r.POST(
		"/hotels/:id/reviews",
		middleware.AuthMiddleware(),
		review.CreateReview,
	)

	// Get hotel reviews
	r.GET(
		"/hotels/:id/reviews",
		review.GetHotelReviews,
	)

	// Get hotel average rating
	r.GET(
		"/hotels/:id/rating",
		review.GetHotelAverageRating,
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