package http

import (
	"authcore/internal/usecase"
	"net/http"
)

type OAuthHandler struct {
	oauthService *usecase.OAuthService
}

func NewOAuthHandler(oauthService *usecase.OAuthService) *OAuthHandler {
	return &OAuthHandler{oauthService: oauthService}
}

func (h *OAuthHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")

	provider, err := h.oauthService.GetProvider(providerName)
	if err != nil {
		WriteResponse(w, http.StatusNotFound, false, "", nil, "OAuth provider not supported")
		return
	}

	state, err := h.oauthService.GenerateState()
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, false, "", nil, "Failed to generate state")
		return
	}

	http.Redirect(w, r, provider.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (h *OAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	providerName := r.PathValue("provider")

	code := r.URL.Query().Get("code")
	if code == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Authorization code not provided")
		return
	}

	clientID, ok := r.Context().Value(ClientIDKey).(string)
	if !ok || clientID == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Client ID not provided")
		return
	}

	accessToken, refreshToken, err := h.oauthService.HandleCallback(r.Context(), providerName, code, clientID)
	if err != nil {
		WriteResponse(w, http.StatusUnauthorized, false, "", nil, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, true, "OAuth login successful", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}
