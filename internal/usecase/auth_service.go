package usecase

import (
	"authcore/internal/domain/repository"
	"authcore/internal/infrastructure/security"
	"errors"
	"log"
)

type AuthService struct {
	repo          repository.UserRepository
	BcryptService security.BcryptService
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s AuthService) Register(email, password string) error {

	user, err := s.repo.GetUserByEmail(email)

	if err != nil {

		log.Println("User not found", err)
	}

	if user != nil {
		if user.Email == email {
			return errors.New("Account already exists")
		}
	}

	hashedPassword, err := s.BcryptService.HashedPassword(password)

	if err != nil {
		return err
	}

	err = s.repo.CreateUser(email, hashedPassword)

	if err != nil {
		return err
	}

	return nil

}

func (s AuthService) Login(email, password string) (string, error) {

	user, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("User not found")
	}

	err = s.BcryptService.CheckPassword(password, user.PasswordHash)

	if err != nil {
		return "", errors.New("Password Mismatch")
	}

	return "Login successfull", nil

}
