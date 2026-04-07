package repository

import "authcore/internal/domain/entity"

type UserRepository interface {
	CreateUser(email, password, role string) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
	UpdateUserRole(id, role string) error
}
