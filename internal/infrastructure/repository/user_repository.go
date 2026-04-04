package repository

import (
	"authcore/internal/domain/entity"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(email, password, role string) error {

	user := entity.User{}

	err := r.db.QueryRow(
		"INSERT INTO users (email, password, role, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		email, password, role,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {

	user := entity.User{}

	err := r.db.QueryRow(
		"SELECT id, email, password, role FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*entity.User, error) {

	user := entity.User{}

	err := r.db.QueryRow("SELECT id, email, password, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
