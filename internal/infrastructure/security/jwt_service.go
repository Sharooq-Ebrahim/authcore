package security

import (
	"authcore/internal/domain/entity"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey       []byte
	expirationHours int
}

func NewJWTService(secretKey string, expirationHours int) *JWTService {
	return &JWTService{
		secretKey:       []byte(secretKey),
		expirationHours: expirationHours,
	}
}

func (j *JWTService) GenerateToken(user *entity.User) (string, error) {
	cliams := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * time.Duration(j.expirationHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cliams)
	return token.SignedString(j.secretKey)
}
