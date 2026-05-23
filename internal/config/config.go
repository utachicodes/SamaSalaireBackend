package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	MongoURI       string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("config: no .env file found; reading configuration from process environment")
		}

		expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
		if err != nil {
			expiryHours = 24
		}

		instance = &Config{
			Port:           getEnv("PORT", "8080"),
			MongoURI:       getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			DBName:         getEnv("DB_NAME", "samasalaire"),
			JWTSecret:      getEnv("JWT_SECRET", "change-me"),
			JWTExpiryHours: expiryHours,
		}
	})
	return instance
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
