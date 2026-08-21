package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(".env file not found")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}