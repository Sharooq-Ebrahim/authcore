package dto

type User struct {
	ID        string `json:"id" validate:"omitempty"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	Role      string `json:"role" validate:"omitempty"`
	CreatedAt string `json:"created_at" validate:"omitempty"`
}
