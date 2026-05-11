package main

import (
	"authcore/internal/config"
	handler "authcore/internal/delivery/http"
	"authcore/internal/infrastructure/db"
	"authcore/internal/infrastructure/repository"
	"authcore/internal/infrastructure/security"
	"authcore/internal/oauth"
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

	if err := db.RunMigration(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to apply migration:", err)
	}

	userRepo := repository.NewUserRepository(dbConn)
	clientRepo := repository.NewClientRepository(dbConn)
	oauthRepo := repository.NewOAuthRepository(dbConn)

	jwtService := security.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationMinutes, cfg.JWTRefreshExpirationHours)

	passwordService := security.NewBcryptService()

	userService := usecase.NewAuthService(userRepo, passwordService, jwtService)
	clientService := usecase.NewClientService(clientRepo)

	oauthProviders := []oauth.Provider{
		oauth.NewGoogleProvider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleCallbackURL),
	}
	oauthService := usecase.NewOAuthService(oauthRepo, userRepo, jwtService, oauthProviders)

	userhandler := handler.NewAuthHandler(userService)
	clientHandler := handler.NewClientHandler(clientService)
	oauthHandler := handler.NewOAuthHandler(oauthService)

	withAPIKey := handler.APIKeyMiddleware(clientService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", withAPIKey(userhandler.Register))
	mux.HandleFunc("POST /login", withAPIKey(userhandler.Login))
	mux.HandleFunc("POST /refresh", withAPIKey(userhandler.RefreshToken))
	mux.HandleFunc("GET /verify", withAPIKey(userhandler.VerifyToken))
	mux.HandleFunc("GET /profile", withAPIKey(userhandler.GetUserProfile))
	mux.HandleFunc("POST /assign-role", withAPIKey(userhandler.AssignRole))

	mux.HandleFunc("GET /oauth/{provider}", (oauthHandler.Redirect))
	mux.HandleFunc("GET /oauth/{provider}/callback", (oauthHandler.Callback))

	mux.HandleFunc("POST /clients", clientHandler.CreateClient)
	mux.HandleFunc("GET /clients/{id}", clientHandler.GetClient)
	mux.HandleFunc("POST /clients/{id}/credentials", clientHandler.CreateCredential)
	mux.HandleFunc("GET /clients/{id}/credentials", clientHandler.GetCredentials)

	err = http.ListenAndServe(":"+cfg.PORT, mux)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
