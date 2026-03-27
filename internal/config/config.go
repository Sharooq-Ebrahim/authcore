package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	PORT               string
	JWTSecret          string
	JWTExpirationHours string
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		PORT:               os.Getenv("PORT"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpirationHours: os.Getenv("JWT_EXPIRATION_HOURS"),
	}
}
