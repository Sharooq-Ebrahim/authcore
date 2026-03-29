package security

import (
	"authcore/internal/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey                  []byte
	expirationHours            int
	refreshTokenExpirationDays int
}

func NewJWTService(secretKey string, expirationHours int, refreshTokenExpirationDays int) *JWTService {
	return &JWTService{
		secretKey:                  []byte(secretKey),
		expirationHours:            expirationHours,
		refreshTokenExpirationDays: refreshTokenExpirationDays,
	}
}

func (j *JWTService) GenerateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * time.Duration(j.expirationHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) GenerateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * time.Duration(j.refreshTokenExpirationDays)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}
