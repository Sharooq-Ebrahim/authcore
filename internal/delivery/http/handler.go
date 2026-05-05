package http

import (
	"authcore/internal/domain/entity"
	"authcore/internal/usecase"
	"encoding/json"
	"log"
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

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Email and password are required")
		return
	}

	if len(req.Password) < 8 {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Password must be at least 8 characters")
		return
	}

	clientID := r.Context().Value(ClientIDKey).(string)

	log.Println("ClientID:", clientID)

	err := h.authService.Register(r.Context(), req.Email, req.Password, entity.RoleUser, clientID)

	if err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, err)
		return
	}

	WriteResponse(w, http.StatusCreated, true, "User registered successfully", nil, nil)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Email and password are required")
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)

	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, false, "", nil, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, true, "Login successful", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)

}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)

	if err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, err)
		return
	}

	WriteResponse(w, http.StatusOK, true, "Token refreshed successfully", map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	}, nil)

}

func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	claims, err := h.authService.VerifyToken(r.Context(), req.Token)

	if err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, err)
		return
	}

	WriteResponse(w, http.StatusOK, true, "Token verified successfully", claims, nil)

}

func (h *AuthHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	token := authHeader
	if parts := strings.Split(authHeader, " "); len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		token = parts[1]
	}

	if token == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Token not provided")
		return
	}

	claims, err := h.authService.GetUserProfile(r.Context(), token)

	if err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, err)
		return
	}

	WriteResponse(w, http.StatusOK, true, "User profile retrieved successfully", claims, nil)

}

func (h *AuthHandler) AssignRole(w http.ResponseWriter, r *http.Request) {

	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	err := h.authService.AssignRole(r.Context(), req.Email, req.Role)

	if err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, err)
		return
	}

	WriteResponse(w, http.StatusOK, true, "Role assigned successfully", nil, nil)

}
