package http

import (
	"authcore/internal/usecase"
	"encoding/json"
	"net/http"
	"strings"
)

type ClientHandler struct {
	clientService *usecase.ClientService
}

func NewClientHandler(clientService *usecase.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {


	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Name is required")
		return
	}

	client, err := h.clientService.CreateClient(r.Context(), req.Name)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, false, "", nil, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, true, "Client created successfully", client, nil)
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Client ID is required")
		return
	}

	client, err := h.clientService.GetClient(r.Context(), id)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, false, "", nil, err.Error())
		return
	}
	if client == nil {
		WriteResponse(w, http.StatusNotFound, false, "", nil, "Client not found")
		return
	}

	WriteResponse(w, http.StatusOK, true, "Client retrieved successfully", client, nil)
}

func (h *ClientHandler) CreateCredential(w http.ResponseWriter, r *http.Request) {


	id := r.PathValue("id")
	if id == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Client ID is required")
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Invalid request body")
		return
	}

	cred, err := h.clientService.GenerateCredential(r.Context(), id, req.Name)
	if err != nil {
		if err.Error() == "client not found" {
			WriteResponse(w, http.StatusNotFound, false, "", nil, err.Error())
			return
		}
		WriteResponse(w, http.StatusInternalServerError, false, "", nil, err.Error())
		return
	}

	WriteResponse(w, http.StatusCreated, true, "Credential created successfully", cred, nil)
}

func (h *ClientHandler) GetCredentials(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")
	if id == "" {
		WriteResponse(w, http.StatusBadRequest, false, "", nil, "Client ID is required")
		return
	}

	creds, err := h.clientService.GetCredentials(r.Context(), id)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, false, "", nil, err.Error())
		return
	}

	WriteResponse(w, http.StatusOK, true, "Credentials retrieved successfully", creds, nil)
}
