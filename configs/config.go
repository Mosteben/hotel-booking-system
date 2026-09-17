package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// A missing .env file is expected and fine in most real deployments
	// (Docker/Railway/Render/Fly/etc. inject environment variables
	// directly rather than shipping a .env file) - the app must not crash
	// just because no physical .env exists on disk. godotenv.Load() never
	// overrides a variable that's already set in the process environment
	// (verified against its source: pkg/joho/godotenv - loadFile only
	// calls os.Setenv for keys not already present), so real
	// deployment-injected env vars always take precedence over .env when
	// both are present; this is already correct and unchanged.
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found - continuing with process environment variables only")
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
