package database

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Mosteben/hotel-booking-system/configs"
	bookingModel "github.com/Mosteben/hotel-booking-system/internal/booking/model"
	favoriteModel "github.com/Mosteben/hotel-booking-system/internal/favorite/model"
	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
	paymentModel "github.com/Mosteben/hotel-booking-system/internal/payment/model"
	profileModel "github.com/Mosteben/hotel-booking-system/internal/profile/model"
	reviewModel "github.com/Mosteben/hotel-booking-system/internal/review/model"
	room "github.com/Mosteben/hotel-booking-system/internal/room/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	// DB_SSLMODE is documented in .env / .env.example but was previously
	// hardcoded to "disable" here regardless of what was set. "disable" is
	// kept as the default so existing local dev setups keep working
	// unchanged; a real deployment against hosted Postgres sets
	// DB_SSLMODE=require (or similar) to actually get SSL.
	sslMode := configs.GetEnv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
	"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=public",
	configs.GetEnv("DB_HOST"),
	configs.GetEnv("DB_USER"),
	configs.GetEnv("DB_PASSWORD"),
	configs.GetEnv("DB_NAME"),
	configs.GetEnv("DB_PORT"),
	sslMode,
)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database Connection Failed: ", err)
	}

	DB = db

	configurePool(db)

	err = DB.AutoMigrate(
		&userModel.User{},
		&profileModel.Profile{},
		&model.Hotel{},
		&model.HotelImage{},
		&room.Room{},
		&room.RoomImage{},
		&bookingModel.Booking{},
		&reviewModel.Review{},
		&favoriteModel.Favorite{},
		&paymentModel.Payment{},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected Successfully")
}

// Default connection pool settings, used whenever the corresponding env var
// is unset, empty, or not a valid positive integer. These are reasonable
// starting points for a single-instance deployment against a typical
// managed Postgres plan, not a hard requirement - tune via env vars to
// match the actual deployment/DB plan.
const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 10
	defaultConnMaxLifeMins = 30
)

// configurePool applies connection pool limits to the underlying *sql.DB.
// Without this, database/sql defaults to unlimited open connections, which
// can exhaust a Postgres provider's connection limit under load.
func configurePool(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("failed to get underlying *sql.DB for pool configuration:", err)
		return
	}

	sqlDB.SetMaxOpenConns(envInt("DB_MAX_OPEN_CONNS", defaultMaxOpenConns))
	sqlDB.SetMaxIdleConns(envInt("DB_MAX_IDLE_CONNS", defaultMaxIdleConns))
	sqlDB.SetConnMaxLifetime(
		time.Duration(envInt("DB_CONN_MAX_LIFETIME_MINUTES", defaultConnMaxLifeMins)) * time.Minute,
	)
}

// envInt reads a positive integer from the environment, falling back to
// def if the variable is unset, empty, not a valid integer, or not positive.
func envInt(key string, def int) int {
	raw := configs.GetEnv(key)
	if raw == "" {
		return def
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		log.Printf("invalid %s=%q, falling back to default %d", key, raw, def)
		return def
	}

	return value
}
