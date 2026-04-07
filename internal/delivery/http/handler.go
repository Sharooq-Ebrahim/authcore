package http

import (
	"authcore/internal/usecase"
	"encoding/json"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService *usecase.AuthService
}

func NewAuthHandler(authService *usecase.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, false, "", nil, "Method not allowed")
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	err := h.authService.Register(req.Email, req.Password, role)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, err.Error())
		return
	}

	writeResponse(w, http.StatusCreated, true, "User registered successfully", nil, nil)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, false, "", nil, "Method not allowed")
		return
	}

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

	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, false, "", nil, "Method not allowed")
		return
	}

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

func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, false, "", nil, "Method not allowed")
		return
	}

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	claims, err := h.authService.VerifyToken(req.Token)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, true, "Token verified successfully", claims, nil)

}

func (h *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		writeResponse(w, http.StatusMethodNotAllowed, false, "", nil, "Method not allowed")
		return
	}

	authHeader := r.Header.Get("Authorization")
	token := authHeader
	if parts := strings.Split(authHeader, " "); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		token = parts[1]
	}

	if token == "" {
		writeResponse(w, http.StatusBadRequest, false, "", nil, "Token not provided")
		return
	}

	claims, err := h.authService.GetUserProfile(token)

	if err != nil {
		writeResponse(w, http.StatusBadRequest, false, "", nil, err.Error())
		return
	}

	writeResponse(w, http.StatusOK, true, "User profile retrieved successfully", claims, nil)

}
