package main

import (
	"log"
	"strings"
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

	paymentHandler "github.com/Mosteben/hotel-booking-system/internal/payment/handler"
	paymentRepository "github.com/Mosteben/hotel-booking-system/internal/payment/repository"
	paymentService "github.com/Mosteben/hotel-booking-system/internal/payment/service"

	profileRepository "github.com/Mosteben/hotel-booking-system/internal/profile/repository"

	reviewHandler "github.com/Mosteben/hotel-booking-system/internal/review/handler"
	reviewRepository "github.com/Mosteben/hotel-booking-system/internal/review/repository"
	reviewService "github.com/Mosteben/hotel-booking-system/internal/review/service"

	roomHandler "github.com/Mosteben/hotel-booking-system/internal/room/handler"
	roomRepository "github.com/Mosteben/hotel-booking-system/internal/room/repository"
	roomService "github.com/Mosteben/hotel-booking-system/internal/room/service"

	userHandler "github.com/Mosteben/hotel-booking-system/internal/user/handler"
	userRepository "github.com/Mosteben/hotel-booking-system/internal/user/repository"
	userService "github.com/Mosteben/hotel-booking-system/internal/user/service"

	"github.com/Mosteben/hotel-booking-system/configs"
	"github.com/Mosteben/hotel-booking-system/pkg/database"
	"github.com/Mosteben/hotel-booking-system/pkg/middleware"
	"github.com/Mosteben/hotel-booking-system/pkg/storage"
	"github.com/Mosteben/hotel-booking-system/routes"

	_ "github.com/Mosteben/hotel-booking-system/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title Hotel Booking System API
// @version 1.0
// @description RESTful API for a Hotel Booking System built with Go, Gin, GORM, PostgreSQL, and JWT authentication.
// @description
// @description Features include authentication, hotel management, room management, bookings, availability, reviews, favorites, and payments.
// @host localhost:8081
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// =========================
	// Configuration
	// =========================

	configs.LoadEnv()

	// JWT_SECRET signs and verifies every auth token issued by this app -
	// an empty secret would make tokens trivially forgeable (anyone could
	// sign their own with HS256 + ""). Fail fast at boot rather than
	// silently issuing insecure tokens at runtime.
	if configs.GetEnv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is not set - refusing to start")
	}

	// =========================
	// Database
	// =========================

	database.Connect()

	// =========================
	// Image Storage
	// =========================

	// Only hotel/room image upload/delete actually needs this - the server
	// still boots and every other endpoint still works normally without
	// it, so a missing Cloudinary config doesn't take down the whole app.
	// Hotel and room images share the same Cloudinary account/credentials
	// but land in separate folders, so the two galleries don't mix.
	var hotelImageStorage storage.Storage
	var roomImageStorage storage.Storage

	cloudinaryCloudName := configs.GetEnv("CLOUDINARY_CLOUD_NAME")
	cloudinaryAPIKey := configs.GetEnv("CLOUDINARY_API_KEY")
	cloudinaryAPISecret := configs.GetEnv("CLOUDINARY_API_SECRET")

	if cloudinaryCloudName == "" || cloudinaryAPIKey == "" || cloudinaryAPISecret == "" {
		log.Println(
			"CLOUDINARY_CLOUD_NAME, CLOUDINARY_API_KEY, and/or CLOUDINARY_API_SECRET are not set - hotel/room image upload/delete will be unavailable until they are configured",
		)
		hotelImageStorage = storage.NewUnconfiguredStorage()
		roomImageStorage = storage.NewUnconfiguredStorage()
	} else {
		hotelCldStorage, err := storage.NewCloudinaryStorage(
			cloudinaryCloudName,
			cloudinaryAPIKey,
			cloudinaryAPISecret,
			"nilestay/hotels",
		)

		if err != nil {
			log.Fatal("Failed to initialize hotel image storage: ", err)
		}

		roomCldStorage, err := storage.NewCloudinaryStorage(
			cloudinaryCloudName,
			cloudinaryAPIKey,
			cloudinaryAPISecret,
			"nilestay/rooms",
		)

		if err != nil {
			log.Fatal("Failed to initialize room image storage: ", err)
		}

		hotelImageStorage = hotelCldStorage
		roomImageStorage = roomCldStorage
	}

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

	paymentRepo := paymentRepository.NewPaymentRepository(
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
		hotelImageStorage,
		database.DB,
	)

	roomSrv := roomService.NewRoomService(
		roomRepo,
		roomImageStorage,
		database.DB,
	)

	bookingSrv := bookingService.NewBookingService(
		bookingRepo,
		roomRepo,
		userRepo,
		hotelRepo,
		database.DB,
	)

	reviewSrv := reviewService.NewReviewService(
		reviewRepo,
		userRepo,
		hotelRepo,
	)

	favoriteSrv := favoriteService.NewFavoriteService(
		favoriteRepo,
	)

	// Payment service uses the database
	// to manage atomic transactions.
	paymentSrv := paymentService.NewPaymentService(
		paymentRepo,
		bookingRepo,
		userRepo,
		roomRepo,
		hotelRepo,
		database.DB,
	)

	userSrv := userService.NewUserService(
		userRepo,
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

	payment := paymentHandler.NewPaymentHandler(
		paymentSrv,
	)

	user := userHandler.NewUserHandler(
		userSrv,
	)

	// =========================
	// Router
	// =========================

	r := gin.Default()

	// =========================
	// CORS
	// =========================

	// The allowed frontend origin(s) come from the environment so the same
	// binary works in dev and in a real deployment without a code change.
	// FRONTEND_URL accepts one origin or a comma-separated list (e.g. the
	// deployed frontend and localhost together) - this is what lets both
	// a production frontend and local dev work against the same backend
	// at once. AllowCredentials is true (JWT is sent from the browser),
	// so every entry must be an exact origin - never "*" - or browsers
	// will reject the credentialed request outright. gin-contrib/cors
	// itself never echoes back "*" here; with AllowCredentials it reflects
	// only whichever configured origin actually matches the request.
	frontendURL := configs.GetEnv("FRONTEND_URL")
	if frontendURL == "" {
		log.Println(
			"FRONTEND_URL is not set, defaulting CORS to http://localhost:5173 (local dev only)",
		)
		frontendURL = "http://localhost:5173"
	}

	var allowedOrigins []string
	for _, origin := range strings.Split(frontendURL, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins = append(allowedOrigins, origin)
		}
	}

	r.Use(middleware.SecurityHeaders())

	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,

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
		payment,
		user,
	)

	// =========================
	// Room Routes
	// =========================

	r.GET(
		"/rooms/hotel/:hotel_id",
		room.GetRoomsByHotelID,
	)

	r.GET(
		"/rooms/:id/availability",
		middleware.AuthMiddleware(),
		booking.CheckRoomAvailability,
	)

	r.GET(
		"/rooms/:id",
		room.GetRoomByID,
	)

	r.POST(
		"/rooms/hotel/:hotel_id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.CreateRoom,
	)

	r.PUT(
		"/rooms/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.UpdateRoom,
	)

	r.DELETE(
		"/rooms/:id",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin"),
		room.DeleteRoom,
	)

	r.POST(
		"/rooms/:id/images",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.UploadRoomImages,
	)

	r.DELETE(
		"/rooms/:id/images/:imageId",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.DeleteRoomImage,
	)

	r.PATCH(
		"/rooms/:id/images/:imageId/main",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		room.SetMainRoomImage,
	)

	// =========================
	// Booking Routes
	// =========================

	r.POST(
		"/bookings",
		middleware.AuthMiddleware(),
		booking.CreateBooking,
	)

	r.GET(
		"/bookings",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		booking.GetAllBookings,
	)

	r.GET(
		"/bookings/my",
		middleware.AuthMiddleware(),
		booking.GetMyBookings,
	)

	r.PATCH(
		"/bookings/:id/status",
		middleware.AuthMiddleware(),
		middleware.RequireRoles("admin", "manager"),
		booking.UpdateBookingStatus,
	)

	r.GET(
		"/bookings/:id",
		middleware.AuthMiddleware(),
		booking.GetBookingByID,
	)

	r.PUT(
		"/bookings/:id",
		middleware.AuthMiddleware(),
		booking.UpdateBooking,
	)

	r.DELETE(
		"/bookings/:id",
		middleware.AuthMiddleware(),
		booking.DeleteBooking,
	)

	// =========================
	// Server
	// =========================

	// PORT is the de facto standard most hosting platforms (Render,
	// Railway, Heroku, Koyeb, Fly, etc.) inject automatically and expect
	// the app to bind to - APP_PORT (this project's own local-dev
	// convention) is checked as a fallback so existing setups keep
	// working. Binding to neither would fall back to ":" (net.Listen
	// picks a random free port), which starts the server but leaves it
	// unreachable at the port the platform actually health-checks -
	// exactly the kind of gap that makes a working build look like a
	// failed deployment. 8080 is the last-resort default so the server
	// never silently binds to a random port either way.
	port := configs.GetEnv("PORT")
	if port == "" {
		port = configs.GetEnv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
