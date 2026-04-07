package security

import (
	"authcore/internal/domain/entity"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey                   []byte
	expirationMinutes           int
	refreshTokenExpirationHours int
}

func NewJWTService(secretKey string, expirationMinutes int, refreshTokenExpirationHours int) *JWTService {
	return &JWTService{
		secretKey:                   []byte(secretKey),
		expirationMinutes:           expirationMinutes,
		refreshTokenExpirationHours: refreshTokenExpirationHours,
	}
}

func (j *JWTService) GenerateToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * time.Duration(j.expirationMinutes)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) GenerateRefreshToken(user *entity.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * time.Duration(j.refreshTokenExpirationHours)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *JWTService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, typeOk := claims["type"].(string)
		if !typeOk || tokenType != "refresh" {
			return "", errors.New("invalid token type")
		}

		userID, userOk := claims["user_id"].(string)

		if !userOk {
			return "", errors.New("userID not found in token")
		}

		return userID, nil

	}

	return "", errors.New("invalid token")

}

func (j *JWTService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
