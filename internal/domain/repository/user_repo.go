package repository

import (
	"authcore/internal/domain/entity"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, password, role string) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	UpdateUserRole(ctx context.Context, id, role string) error
}
