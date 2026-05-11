package repository

import (
	"authcore/internal/domain/entity"
	"context"
)

type OAuthRepository interface {
	FindByProvider(ctx context.Context, provider, providerID string) (*entity.OAuthAccount, error)
	LinkAccount(ctx context.Context, userID, provider, providerID, email string) (*entity.OAuthAccount, error)
}
