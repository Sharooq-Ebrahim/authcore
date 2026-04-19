package main

import (
	"authcore/internal/config"
	handler "authcore/internal/delivery/http"
	"authcore/internal/infrastructure/db"
	"authcore/internal/infrastructure/repository"
	"authcore/internal/infrastructure/security"
	"authcore/internal/usecase"
	"log"
	"net/http"
)

func main() {

	cfg := config.LoadEnv()

	dbConn, err := db.ConnectDB(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repository.NewUserRepository(dbConn)

	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationMinutes, cfg.JWTRefreshExpirationHours)

	passwordService := security.NewBcryptService()

	userService := usecase.NewAuthService(userRepo, passwordService, jwtService)

	userhandler := handler.NewAuthHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", userhandler.Register)
	mux.HandleFunc("/login", userhandler.Login)
	mux.HandleFunc("/refresh", userhandler.RefreshToken)
	mux.HandleFunc("/verify", userhandler.VerifyToken)
	mux.HandleFunc("/profile", userhandler.GetUserProfile)
	mux.HandleFunc("/assign-role", userhandler.AssignRole)

	err = http.ListenAndServe(":"+cfg.PORT, mux)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
