package service

import (
	"authcore/internal/domain/entity"
)

type TokenService interface {
	GenerateToken(user *entity.User) (string, error)
	GenerateRefreshToken(user *entity.User) (string, error)
	ValidateRefreshToken(token string) (string, error)
	ValidateToken(token string) (map[string]interface{}, error)
}
