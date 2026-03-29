package service

import (
	"authcore/internal/domain/entity"
)

type TokenService interface {
	GenerateToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	// ValidateToken(token string) (*entity.User, error)
}

