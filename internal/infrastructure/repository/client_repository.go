package repository

import (
	"authcore/internal/domain/entity"
	"context"
	"database/sql"
)

type clientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) *clientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) CreateClient(ctx context.Context, name string) (*entity.Client, error) {
	client := entity.Client{}
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO clients (name, created_at) VALUES ($1, NOW()) RETURNING id, name, created_at",
		name,
	).Scan(&client.ID, &client.Name, &client.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) GetClientByID(ctx context.Context, id string) (*entity.Client, error) {
	client := entity.Client{}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, name, created_at FROM clients WHERE id = $1",
		id,
	).Scan(&client.ID, &client.Name, &client.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (r *clientRepository) CreateClientCredential(ctx context.Context, clientID, key, name string) (*entity.ClientCredential, error) {
	cred := entity.ClientCredential{}
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO client_credentials (client_id, key, name, is_active, created_at) VALUES ($1, $2, $3, true, NOW()) RETURNING id, client_id, key, name, is_active, created_at",
		clientID, key, name,
	).Scan(&cred.ID, &cred.ClientID, &cred.Key, &cred.Name, &cred.IsActive, &cred.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &cred, nil
}

func (r *clientRepository) GetClientCredentialByKey(ctx context.Context, key string) (*entity.ClientCredential, error) {
	cred := entity.ClientCredential{}
	var name sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		"SELECT id, client_id, key, name, is_active, created_at FROM client_credentials WHERE key = $1",
		key,
	).Scan(&cred.ID, &cred.ClientID, &cred.Key, &name, &cred.IsActive, &cred.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if name.Valid {
		cred.Name = name.String
	}
	return &cred, nil
}

func (r *clientRepository) GetClientCredentialsByClientID(ctx context.Context, clientID string) ([]*entity.ClientCredential, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, client_id, key, name, is_active, created_at FROM client_credentials WHERE client_id = $1",
		clientID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var creds []*entity.ClientCredential
	for rows.Next() {
		cred := entity.ClientCredential{}
		var name sql.NullString
		if err := rows.Scan(&cred.ID, &cred.ClientID, &cred.Key, &name, &cred.IsActive, &cred.CreatedAt); err != nil {
			return nil, err
		}
		if name.Valid {
			cred.Name = name.String
		}
		creds = append(creds, &cred)
	}
	return creds, nil
}
