package repository

import (
	"authcore/internal/domain/entity"
	"context"
	"database/sql"
)

type oauthRepository struct {
	db *sql.DB
}

func NewOAuthRepository(db *sql.DB) *oauthRepository {
	return &oauthRepository{db: db}
}

func (r *oauthRepository) FindByProvider(ctx context.Context, provider, providerID string) (*entity.OAuthAccount, error) {
	acc := entity.OAuthAccount{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, provider, provider_id, email, created_at
		 FROM oauth_accounts
		 WHERE provider = $1 AND provider_id = $2`,
		provider, providerID,
	).Scan(&acc.ID, &acc.UserID, &acc.Provider, &acc.ProviderID, &acc.Email, &acc.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &acc, nil
}

func (r *oauthRepository) LinkAccount(ctx context.Context, userID, provider, providerID, email string) (*entity.OAuthAccount, error) {
	acc := entity.OAuthAccount{}

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO oauth_accounts (user_id, provider, provider_id, email)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, provider, provider_id, email, created_at`,
		userID, provider, providerID, email,
	).Scan(&acc.ID, &acc.UserID, &acc.Provider, &acc.ProviderID, &acc.Email, &acc.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &acc, nil
}
