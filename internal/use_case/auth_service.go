package usecase

import "authcore/internal/domain/repository"

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s AuthService) Register(email, password string) error {

	err := s.repo.CreateUser(email, password)

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
