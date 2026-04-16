package repository

import (
	"authcore/internal/domain/entity"
	"context"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, email, password, role string) error {

	user := entity.User{}

	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO users (email, password, role, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		email, password, role,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {

	user := entity.User{}

	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, email, password, role FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*entity.User, error) {

	user := entity.User{}

	err := r.db.QueryRowContext(ctx, "SELECT id, email, password, role FROM users WHERE id = $1", id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateUserRole(ctx context.Context, id, role string) error {

	_, err := r.db.ExecContext(ctx, "UPDATE users SET role = $1 WHERE id = $2", role, id)

	if err != nil {
		return err
	}

	return nil

}
