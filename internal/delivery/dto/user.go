package dto

type User struct {
	ID        string `json:"id" validate:"omitempty"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	CreatedAt string `json:"created_at" validate:"omitempty"`
}
