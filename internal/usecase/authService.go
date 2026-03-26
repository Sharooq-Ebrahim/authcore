package usecase

import (
	"authcore/internal/domain/repository"
	"errors"
	"log"
)

type AuthService struct {
	repo repository.UserRepository
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

	err = s.repo.CreateUser(email, password)

	if err != nil {
		return err
	}

	return nil

}

func (s AuthService) Login(email, password string) (string, error) {

	_, err := s.repo.GetUserByEmail(email)

	if err != nil {
		return "", err
	}

	return "Login successfull", nil

}
