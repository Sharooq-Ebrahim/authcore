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

	config := config.LoadEnv()

	log.Println("Database URL:", config.DatabaseURL)

	db, err := db.ConnectDB(config.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	userRepo := repository.NewUserRepository(db)

	jwtService := security.NewJWTService(config.JWTSecret, 24, 7)

	passwordService := security.NewBcryptService()

	userService := usecase.NewAuthService(userRepo, passwordService, jwtService)

	userhandler := handler.NewAuthHandler(userService)

	http.HandleFunc("/register", userhandler.Register)
	http.HandleFunc("/login", userhandler.Login)
	http.HandleFunc("/refresh", userhandler.RefreshToken)


	http.ListenAndServe(":"+config.PORT, nil)

}
