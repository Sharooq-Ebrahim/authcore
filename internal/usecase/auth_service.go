package usecase

import (
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
