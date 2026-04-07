package usecase

import (
	"authcore/internal/domain/entity"
	"authcore/internal/domain/repository"
	"authcore/internal/domain/service"

	"errors"
)

type AuthService struct {
	repo            repository.UserRepository
	passwordService service.PasswordService
	TokenService    service.TokenService
}

func NewAuthService(repo repository.UserRepository, passwordService service.PasswordService, tokenService service.TokenService) *AuthService {
	return &AuthService{repo: repo, passwordService: passwordService, TokenService: tokenService}
}

func (s *AuthService) Register(email, password, role string) error {

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if user != nil {
		if user.Email == email {
			return errors.New("Account already exists")
		}
	}

	hashedPassword, err := s.passwordService.HashPassword(password)

	if err != nil {
		return err
	}

	err = s.repo.CreateUser(email, hashedPassword, role)

	if err != nil {
		return err
	}

	return nil

}

func (s *AuthService) Login(email, password string) (string, string, error) {

	user, err := s.repo.GetUserByEmail(email)

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

	accessToken, err := s.TokenService.GenerateToken(user)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.TokenService.GenerateRefreshToken(user)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {

	userID, err := s.TokenService.ValidateRefreshToken(refreshToken)

	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserByID(userID)

	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", errors.New("user not found")
	}

	newAccessToken, err := s.TokenService.GenerateToken(user)

	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := s.TokenService.GenerateRefreshToken(user)

	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil

}

func (s *AuthService) VerifyToken(token string) (map[string]interface{}, error) {

	claims, err := s.TokenService.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	return claims, nil

}

func (s *AuthService) GetUserProfile(token string) (map[string]interface{}, error) {

	claims, err := s.TokenService.ValidateToken(token)

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

	user, err := s.repo.GetUserByID(userID)

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

func (s *AuthService) AssignRole(email, role string) error {

	user, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	if role != entity.RoleUser && role != entity.RoleAdmin {
		return errors.New("invalid role")
	}

	err = s.repo.UpdateUserRole(user.ID, role)

	if err != nil {
		return err
	}

	return nil

}
