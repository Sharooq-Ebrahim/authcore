package dto

type Client struct {
	ID        string `json:"id" validate:"omitempty"`
	Name      string `json:"name" validate:"required"`
	CreatedAt string `json:"created_at" validate:"omitempty"`
}

type ClientCredential struct {
	ID        string `json:"id" validate:"omitempty"`
	ClientID  string `json:"client_id" validate:"required"`
	Key       string `json:"key" validate:"required"`
	Name      string `json:"name" validate:"omitempty"`
	IsActive  bool   `json:"is_active" validate:"omitempty"`
	CreatedAt string `json:"created_at" validate:"omitempty"`
}
