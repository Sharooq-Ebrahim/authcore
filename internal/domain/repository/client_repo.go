package repository

import (
	"authcore/internal/domain/entity"
	"context"
)

type ClientRepository interface {
	CreateClient(ctx context.Context, name string) (*entity.Client, error)
	GetClientByID(ctx context.Context, id string) (*entity.Client, error)
	CreateClientCredential(ctx context.Context, clientID, key, name string) (*entity.ClientCredential, error)
	GetClientCredentialByKey(ctx context.Context, key string) (*entity.ClientCredential, error)
	GetClientCredentialsByClientID(ctx context.Context, clientID string) ([]*entity.ClientCredential, error)
}
