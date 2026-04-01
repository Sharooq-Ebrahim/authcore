package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL               string
	PORT                      string
	JWTSecret                 string
	JWTExpirationMinutes      int
	JWTRefreshExpirationHours int
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Not found .env file")
	}

	expMin, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		log.Fatal("Invalid JWT_EXPIRATION_MINUTES")
	}

	refreshHours, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRATION_HOURS"))
	if err != nil {
		log.Fatal("Invalid JWT_REFRESH_EXPIRATION_HOURS")
	}

	return &Config{
		DatabaseURL:               os.Getenv("DATABASE_URL"),
		PORT:                      os.Getenv("PORT"),
		JWTSecret:                 os.Getenv("JWT_SECRET"),
		JWTExpirationMinutes:      expMin,
		JWTRefreshExpirationHours: refreshHours,
	}
}
