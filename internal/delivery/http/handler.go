package http

import (
	"authcore/internal/usecase"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService *usecase.AuthService
}

func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	err := h.authService.Register(req.Email, req.Password)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, true, "User registered successfully", nil, nil)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, true, "Login successful", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)

}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(req.RefreshToken)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, true, "Token refreshed successfully", map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}, nil)

}

