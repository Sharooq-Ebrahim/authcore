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

func (r *userRepository) CreateUser(email, password string) error {

	user := entity.User{}

	err := r.db.QueryRow(
		"INSERT INTO users (email, password, created_at) VALUES ($1, $2, NOW()) RETURNING id",
		email, password,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {

	user := entity.User{}

	err := r.db.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(id string) (*entity.User, error) {

	user := entity.User{}

	err := r.db.QueryRow("SELECT id,email, password FROM users WHERE id =$1", id).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
