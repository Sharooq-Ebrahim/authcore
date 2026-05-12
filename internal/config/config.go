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

	GoogleClientID     string
	GoogleClientSecret string
	GoogleCallbackURL  string

	RateLimitRPS   float64
	RateLimitBurst int
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

	rps, err := strconv.ParseFloat(os.Getenv("RATE_LIMIT_RPS"), 64)
	if err != nil {
		rps = 10
	}

	burst, err := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST"))
	if err != nil {
		burst = 20
	}

	return &Config{
		DatabaseURL:               os.Getenv("DATABASE_URL"),
		PORT:                      os.Getenv("PORT"),
		JWTSecret:                 os.Getenv("JWT_SECRET"),
		JWTExpirationMinutes:      expMin,
		JWTRefreshExpirationHours: refreshHours,

		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleCallbackURL:  os.Getenv("GOOGLE_CALLBACK_URL"),

		RateLimitRPS:   rps,
		RateLimitBurst: burst,
	}
}
