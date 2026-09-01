package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Mosteben/hotel-booking-system/configs"
	"github.com/Mosteben/hotel-booking-system/pkg/database"
	"github.com/Mosteben/hotel-booking-system/pkg/middleware"
	"github.com/Mosteben/hotel-booking-system/routes"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	authService "github.com/Mosteben/hotel-booking-system/internal/auth/service"

	hotelHandler "github.com/Mosteben/hotel-booking-system/internal/hotel/handler"
	hotelRepository "github.com/Mosteben/hotel-booking-system/internal/hotel/repository"
	hotelService "github.com/Mosteben/hotel-booking-system/internal/hotel/service"

	profileRepository "github.com/Mosteben/hotel-booking-system/internal/profile/repository"
	userRepository "github.com/Mosteben/hotel-booking-system/internal/user/repository"

	roomHandler "github.com/Mosteben/hotel-booking-system/internal/room/handler"
	roomRepository "github.com/Mosteben/hotel-booking-system/internal/room/repository"
	roomService "github.com/Mosteben/hotel-booking-system/internal/room/service"
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
	)

	// =========================
	// Room Routes
	// =========================

	// Get all rooms for a specific hotel
	r.GET(
		"/rooms/hotel/:hotel_id",
		room.GetRoomsByHotelID,
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
	// Server
	// =========================

	r.Run(
		":" + configs.GetEnv("APP_PORT"),
	)
}
