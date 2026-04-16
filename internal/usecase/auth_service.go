package usecase

import (
	"authcore/internal/domain/entity"
	"authcore/internal/domain/repository"
	"authcore/internal/domain/service"
	"context"

	"errors"
)

type AuthService struct {
	repo            repository.UserRepository
	passwordService service.PasswordService
	tokenService    service.TokenService
}

func NewAuthService(repo repository.UserRepository, passwordService service.PasswordService, tokenService service.TokenService) *AuthService {
	return &AuthService{repo: repo, passwordService: passwordService, tokenService: tokenService}
}

func (s *AuthService) Register(ctx context.Context, email, password, role string) error {

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user != nil {
		return errors.New("account already exists")
	}

	hashedPassword, err := s.passwordService.HashPassword(password)

	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, email, hashedPassword, role)

	if err != nil {
		return err
	}

	return nil

}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {

	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("Invalid Email or Password")
	}

	err = s.passwordService.CheckPassword(password, user.PasswordHash)

	if err != nil {
		return "", "", errors.New("Invalid Email or Password")
	}

	accessToken, err := s.tokenService.GenerateToken(user)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {

	userID, err := s.tokenService.ValidateRefreshToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserByID(ctx, userID)

	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("user not found")
	}

	newAccessToken, err := s.tokenService.GenerateToken(user)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.tokenService.GenerateRefreshToken(user)

	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil

}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (map[string]interface{}, error) {

	claims, err := s.tokenService.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	return claims, nil

}

func (s *AuthService) GetUserProfile(ctx context.Context, token string) (map[string]interface{}, error) {

	claims, err := s.tokenService.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	tokenType, typeOk := claims["type"].(string)
	if !typeOk || tokenType != "access" {
		return nil, errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)

	if !ok {
		return nil, errors.New("userID not found in token")
	}

	user, err := s.repo.GetUserByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	}, nil

}

func (s *AuthService) AssignRole(ctx context.Context, email, role string) error {

	user, err := s.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if role != entity.RoleUser && role != entity.RoleAdmin {
		return errors.New("invalid role")
	}

	err = s.repo.UpdateUserRole(ctx, user.ID, role)

	if err != nil {
		return err
	}

	return nil

}
