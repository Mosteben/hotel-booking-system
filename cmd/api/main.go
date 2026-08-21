package main

import (
	"time"

	authHandler "github.com/Mosteben/hotel-booking-system/internal/auth/handler"
	authService "github.com/Mosteben/hotel-booking-system/internal/auth/service"

	hotelHandler "github.com/Mosteben/hotel-booking-system/internal/hotel/handler"
	hotelRepository "github.com/Mosteben/hotel-booking-system/internal/hotel/repository"
	hotelService "github.com/Mosteben/hotel-booking-system/internal/hotel/service"

	profileRepository "github.com/Mosteben/hotel-booking-system/internal/profile/repository"
	userRepository "github.com/Mosteben/hotel-booking-system/internal/user/repository"

	"github.com/Mosteben/hotel-booking-system/configs"
	"github.com/Mosteben/hotel-booking-system/pkg/database"
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

	// =========================
	// Handlers
	// =========================

	auth := authHandler.NewAuthHandler(
		authSrv,
	)

	hotel := hotelHandler.NewHotelHandler(
		hotelSrv,
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
	// Routes
	// =========================

	routes.RegisterRoutes(
		r,
		auth,
		hotel,
	)

	// =========================
	// Server
	// =========================

	r.Run(
		":" + configs.GetEnv("APP_PORT"),
	)
}
