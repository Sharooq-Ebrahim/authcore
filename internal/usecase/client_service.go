package usecase

import (
	"authcore/internal/domain/entity"
	"authcore/internal/domain/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type ClientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) CreateClient(ctx context.Context, name string) (*entity.Client, error) {
	if name == "" {
		return nil, errors.New("client name is required")
	}
	return s.repo.CreateClient(ctx, name)
}

func (s *ClientService) GetClient(ctx context.Context, id string) (*entity.Client, error) {
	if id == "" {
		return nil, errors.New("client id is required")
	}
	return s.repo.GetClientByID(ctx, id)
}

func (s *ClientService) GenerateCredential(ctx context.Context, clientID, name string) (*entity.ClientCredential, error) {
	if clientID == "" {
		return nil, errors.New("client id is required")
	}

	client, err := s.repo.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client not found")
	}

	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)

	return s.repo.CreateClientCredential(ctx, clientID, key, name)
}

func (s *ClientService) GetCredentials(ctx context.Context, clientID string) ([]*entity.ClientCredential, error) {
	if clientID == "" {
		return nil, errors.New("client id is required")
	}
	return s.repo.GetClientCredentialsByClientID(ctx, clientID)
}

func (s *ClientService) ValidateAPIKey(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", errors.New("api key is required")
	}

	cred, err := s.repo.GetClientCredentialByKey(ctx, key)
	if err != nil {
		return "", err
	}
	if cred == nil || !cred.IsActive {
		return "", errors.New("invalid or inactive api key")
	}

	return cred.ClientID, nil
}
