package http

import (
	"authcore/internal/usecase"
	"context"
	"net/http"
)

type contextKey string

const ClientIDKey contextKey = "client_id"

func APIKeyMiddleware(clientService *usecase.ClientService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("x-api-key")
			if apiKey == "" {
				WriteResponse(w, http.StatusUnauthorized, false, "Missing x-api-key header", nil, nil)
				return
			}

			clientID, err := clientService.ValidateAPIKey(r.Context(), apiKey)
			if err != nil {
				WriteResponse(w, http.StatusUnauthorized, false, err.Error(), nil, nil)
				return
			}

			ctx := context.WithValue(r.Context(), ClientIDKey, clientID)
			next(w, r.WithContext(ctx))
		}
	}
}
