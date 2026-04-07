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

	http.HandleFunc("/register", userhandler.Register)
	http.HandleFunc("/login", userhandler.Login)
	http.HandleFunc("/refresh", userhandler.RefreshToken)
	http.HandleFunc("/verify", userhandler.VerifyToken)
	http.HandleFunc("/profile", userhandler.GetUserProfile)
	http.HandleFunc("/assign-role", userhandler.AssignRole)

	err = http.ListenAndServe(":"+cfg.PORT, nil)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
