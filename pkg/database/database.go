package database

import (
	"fmt"
	"log"

	"github.com/Mosteben/hotel-booking-system/configs"
	profileModel "github.com/Mosteben/hotel-booking-system/internal/profile/model"
	userModel "github.com/Mosteben/hotel-booking-system/internal/user/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/Mosteben/hotel-booking-system/internal/hotel/model"
)

var DB *gorm.DB

func Connect() {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		configs.GetEnv("DB_HOST"),
		configs.GetEnv("DB_USER"),
		configs.GetEnv("DB_PASSWORD"),
		configs.GetEnv("DB_NAME"),
		configs.GetEnv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database Connection Failed: ", err)
	}

	DB = db

	err = DB.AutoMigrate(
		&userModel.User{},
		&profileModel.Profile{},
		&model.Hotel{},
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database Connected Successfully")
}
