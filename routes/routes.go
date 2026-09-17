package routes

import (
	"net/http"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	favoriteHandler "github.com/Mosteben/hotel-booking-system/internal/favorite/handler"
	hotelHandler "github.com/Mosteben/hotel-booking-system/internal/hotel/handler"
	paymentHandler "github.com/Mosteben/hotel-booking-system/internal/payment/handler"
	reviewHandler "github.com/Mosteben/hotel-booking-system/internal/review/handler"
	userHandler "github.com/Mosteben/hotel-booking-system/internal/user/handler"
	"github.com/Mosteben/hotel-booking-system/pkg/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	auth *authHandler.AuthHandler,
	hotel *hotelHandler.HotelHandler,
	review *reviewHandler.ReviewHandler,
	favorite *favoriteHandler.FavoriteHandler,
	payment *paymentHandler.PaymentHandler,
	user *userHandler.UserHandler,
) {

	// =========================
	// Swagger
	// =========================

	r.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	// =========================
	// Public Routes
	// =========================

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Rate-limited per IP - these are the endpoints most worth protecting
	// against brute-force/credential-stuffing and registration spam.
	r.POST("/auth/register", middleware.AuthRateLimiter(), auth.Register)
	r.POST("/auth/login", middleware.AuthRateLimiter(), auth.Login)

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

	r.POST(
		"/hotels/:id/images",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		hotel.UploadHotelImages,
	)

	r.DELETE(
		"/hotels/:id/images/:imageId",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		hotel.DeleteHotelImage,
	)

	r.PATCH(
		"/hotels/:id/images/:imageId/main",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		hotel.SetMainHotelImage,
	)

	// =========================
	// Review Routes
	// =========================

	r.POST(
		"/hotels/:id/reviews",
		middleware.AuthMiddleware(),
		review.CreateReview,
	)

	r.GET(
		"/hotels/:id/reviews",
		review.GetHotelReviews,
	)

	r.GET(
		"/hotels/:id/rating",
		review.GetHotelAverageRating,
	)

	// Admin/Manager review moderation - all reviews across every hotel,
	// and deleting an inappropriate one. Registered before any /reviews/:id
	// wildcard would matter, though there is only this one here.
	r.GET(
		"/reviews",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		review.GetAllReviews,
	)

	r.DELETE(
		"/reviews/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		review.DeleteReview,
	)

	// =========================
	// User Routes (Admin/Manager)
	// =========================

	r.GET(
		"/users",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		user.ListUsers,
	)

	// =========================
	// Favorite Routes
	// =========================

	r.POST(
		"/hotels/:id/favorite",
		middleware.AuthMiddleware(),
		favorite.AddFavorite,
	)

	r.GET(
		"/favorites",
		middleware.AuthMiddleware(),
		favorite.GetMyFavorites,
	)

	r.DELETE(
		"/hotels/:id/favorite",
		middleware.AuthMiddleware(),
		favorite.RemoveFavorite,
	)

	// =========================
	// Payment Routes
	// =========================

	// Create payment for my booking
	r.POST(
		"/bookings/:id/payment",
		middleware.AuthMiddleware(),
		payment.CreatePayment,
	)

	// Get my payments
	//
	// IMPORTANT:
	// This route must come before /payments/:id
	// because /payments/:id is a wildcard route.
	r.GET(
		"/payments/my",
		middleware.AuthMiddleware(),
		payment.GetMyPayments,
	)

	// Get payment by ID
	r.GET(
		"/payments/:id",
		middleware.AuthMiddleware(),
		payment.GetPaymentByID,
	)

	// Get all payments
	// Admin and Manager only.
	r.GET(
		"/payments",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		payment.GetAllPayments,
	)

	// Update payment status
	// Admin and Manager only.
	r.PATCH(
		"/payments/:id/status",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		payment.UpdatePaymentStatus,
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
