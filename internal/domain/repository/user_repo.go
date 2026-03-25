package repository

import "authcore/internal/domain/entity"

type UserRepository interface {
	CreateUser(email, password string) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id string) (*entity.User, error)
}



