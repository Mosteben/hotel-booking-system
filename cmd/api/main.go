package main

import (
	"time"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	authService "github.com/Mosteben/hotel-booking-system/internal/auth/service"

	bookingHandler "github.com/Mosteben/hotel-booking-system/internal/booking/handler"
	bookingRepository "github.com/Mosteben/hotel-booking-system/internal/booking/repository"
	bookingService "github.com/Mosteben/hotel-booking-system/internal/booking/service"

	favoriteHandler "github.com/Mosteben/hotel-booking-system/internal/favorite/handler"
	favoriteRepository "github.com/Mosteben/hotel-booking-system/internal/favorite/repository"
	favoriteService "github.com/Mosteben/hotel-booking-system/internal/favorite/service"

	hotelHandler "github.com/Mosteben/hotel-booking-system/internal/hotel/handler"
	hotelRepository "github.com/Mosteben/hotel-booking-system/internal/hotel/repository"
	hotelService "github.com/Mosteben/hotel-booking-system/internal/hotel/service"

	profileRepository "github.com/Mosteben/hotel-booking-system/internal/profile/repository"
	reviewHandler "github.com/Mosteben/hotel-booking-system/internal/review/handler"
	reviewRepository "github.com/Mosteben/hotel-booking-system/internal/review/repository"
	reviewService "github.com/Mosteben/hotel-booking-system/internal/review/service"
	roomHandler "github.com/Mosteben/hotel-booking-system/internal/room/handler"
	roomRepository "github.com/Mosteben/hotel-booking-system/internal/room/repository"
	roomService "github.com/Mosteben/hotel-booking-system/internal/room/service"
	userRepository "github.com/Mosteben/hotel-booking-system/internal/user/repository"

	"github.com/Mosteben/hotel-booking-system/configs"
	"github.com/Mosteben/hotel-booking-system/pkg/database"
	"github.com/Mosteben/hotel-booking-system/pkg/middleware"
	"github.com/Mosteben/hotel-booking-system/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// =========================
	// Configuration
	// =========================

	configs.LoadEnv()

	// =========================
	// Database
	// =========================

	database.Connect()

	// =========================
	// Repositories
	// =========================

	userRepo := userRepository.NewUserRepository(
		database.DB,
	)

	profileRepo := profileRepository.NewProfileRepository(
		database.DB,
	)

	hotelRepo := hotelRepository.NewHotelRepository(
		database.DB,
	)

	roomRepo := roomRepository.NewRoomRepository(
		database.DB,
	)

	bookingRepo := bookingRepository.NewBookingRepository(
		database.DB,
	)

	reviewRepo := reviewRepository.NewReviewRepository(
		database.DB,
	)

	favoriteRepo := favoriteRepository.NewFavoriteRepository(
		database.DB,
	)

	// =========================
	// Services
	// =========================

	authSrv := authService.NewAuthService(
		database.DB,
		userRepo,
		profileRepo,
	)

	hotelSrv := hotelService.NewHotelService(
		hotelRepo,
	)

	roomSrv := roomService.NewRoomService(
		roomRepo,
	)

	bookingSrv := bookingService.NewBookingService(
		bookingRepo,
		roomRepo,
	)

	reviewSrv := reviewService.NewReviewService(
		reviewRepo,
	)

	favoriteSrv := favoriteService.NewFavoriteService(
		favoriteRepo,
	)

	// =========================
	// Handlers
	// =========================

	auth := authHandler.NewAuthHandler(
		authSrv,
	)

	hotel := hotelHandler.NewHotelHandler(
		hotelSrv,
	)

	room := roomHandler.NewRoomHandler(
		roomSrv,
	)

	booking := bookingHandler.NewBookingHandler(
		bookingSrv,
	)

	review := reviewHandler.NewReviewHandler(
		reviewSrv,
	)

	favorite := favoriteHandler.NewFavoriteHandler(
		favoriteSrv,
	)

	// =========================
	// Router
	// =========================

	r := gin.Default()

	// =========================
	// CORS
	// =========================

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		ExposeHeaders: []string{
			"Content-Length",
		},

		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))

	// =========================
	// Existing Routes
	// =========================

	routes.RegisterRoutes(
		r,
		auth,
		hotel,
		review,
		favorite,
	)

	// =========================
	// Room Routes
	// =========================

	// Get all rooms for a specific hotel
	r.GET(
		"/rooms/hotel/:hotel_id",
		room.GetRoomsByHotelID,
	)

	// Check room availability
	//
	// IMPORTANT:
	// This route must come BEFORE /rooms/:id
	// because /rooms/:id is a wildcard route.
	r.GET(
		"/rooms/:id/availability",
		middleware.AuthMiddleware(),
		booking.CheckRoomAvailability,
	)

	// Get room by ID
	r.GET(
		"/rooms/:id",
		room.GetRoomByID,
	)

	// Create room
	r.POST(
		"/rooms/hotel/:hotel_id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.CreateRoom,
	)

	// Update room
	r.PUT(
		"/rooms/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.UpdateRoom,
	)

	// Delete room
	r.DELETE(
		"/rooms/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin"),
		room.DeleteRoom,
	)

	// =========================
	// Booking Routes
	// =========================

	// Create booking
	r.POST(
		"/bookings",
		middleware.AuthMiddleware(),
		booking.CreateBooking,
	)

	// Get all bookings
	r.GET(
		"/bookings",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		booking.GetAllBookings,
	)

	// Get my bookings
	r.GET(
		"/bookings/my",
		middleware.AuthMiddleware(),
		booking.GetMyBookings,
	)

	// Update booking status
	// Admin and manager only
	r.PATCH(
		"/bookings/:id/status",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		booking.UpdateBookingStatus,
	)

	// Get booking by ID
	r.GET(
		"/bookings/:id",
		middleware.AuthMiddleware(),
		booking.GetBookingByID,
	)

	// Update booking
	r.PUT(
		"/bookings/:id",
		middleware.AuthMiddleware(),
		booking.UpdateBooking,
	)

	// Delete / cancel booking
	r.DELETE(
		"/bookings/:id",
		middleware.AuthMiddleware(),
		booking.DeleteBooking,
	)

	// =========================
	// Server
	// =========================

	r.Run(
		":" + configs.GetEnv("APP_PORT"),
	)
}
