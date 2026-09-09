package routes

import (
	"net/http"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	favoriteHandler "github.com/Mosteben/hotel-booking-system/internal/favorite/handler"
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
	favorite *favoriteHandler.FavoriteHandler,
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

	// Get hotel details with rooms
	//
	// IMPORTANT:
	// This route must come before /hotels/:id
	// because /hotels/:id is a wildcard route.
	r.GET("/hotels/:id/details", hotel.GetHotelDetails)

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
	// Favorite Routes
	// =========================

	// Add hotel to favorites
	r.POST(
		"/hotels/:id/favorite",
		middleware.AuthMiddleware(),
		favorite.AddFavorite,
	)

	// Get my favorites
	r.GET(
		"/favorites",
		middleware.AuthMiddleware(),
		favorite.GetMyFavorites,
	)

	// Remove hotel from favorites
	r.DELETE(
		"/hotels/:id/favorite",
		middleware.AuthMiddleware(),
		favorite.RemoveFavorite,
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
