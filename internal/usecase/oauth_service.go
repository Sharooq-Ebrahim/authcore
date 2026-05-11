package usecase

import (
	"authcore/internal/domain/apperrors"
	"authcore/internal/domain/entity"
	"authcore/internal/domain/repository"
	"authcore/internal/domain/service"
	"authcore/internal/oauth"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type OAuthService struct {
	oauthRepo    repository.OAuthRepository
	userRepo     repository.UserRepository
	tokenService service.TokenService
	providers    map[string]oauth.Provider
}

func NewOAuthService(
	oauthRepo repository.OAuthRepository,
	userRepo repository.UserRepository,
	tokenService service.TokenService,
	providers []oauth.Provider,
) *OAuthService {
	pm := make(map[string]oauth.Provider, len(providers))
	for _, p := range providers {
		pm[p.Name()] = p
	}
	return &OAuthService{
		oauthRepo:    oauthRepo,
		userRepo:     userRepo,
		tokenService: tokenService,
		providers:    pm,
	}
}

func (s *OAuthService) GetProvider(name string) (oauth.Provider, error) {
	p, ok := s.providers[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", apperrors.ErrNotFound, name)
	}
	return p, nil
}

func (s *OAuthService) GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *OAuthService) HandleCallback(ctx context.Context, providerName, code, clientID string) (string, string, error) {
	p, err := s.GetProvider(providerName)
	if err != nil {
		return "", "", apperrors.ErrNotFound
	}

	info, err := p.Exchange(code)
	if err != nil {
		return "", "", fmt.Errorf("oauth exchange failed: %w", err)
	}

	existing, err := s.oauthRepo.FindByProvider(ctx, providerName, info.ProviderID)
	if err != nil {
		return "", "", err
	}

	var user *entity.User

	if existing != nil {
		user, err = s.userRepo.GetUserByID(ctx, existing.UserID)
		if err != nil {
			return "", "", err
		}
		if user == nil {
			return "", "", apperrors.ErrUserNotFound
		}
	} else {
		user, err = s.userRepo.GetUserByEmail(ctx, info.Email)
		if err != nil {
			return "", "", err
		}

		if user == nil {
			if err := s.userRepo.CreateUser(ctx, info.Email, "", entity.RoleUser, clientID); err != nil {
				return "", "", err
			}
			user, err = s.userRepo.GetUserByEmail(ctx, info.Email)
			if err != nil {
				return "", "", err
			}
			if user == nil {
				return "", "", apperrors.ErrInternalServer
			}
		}

		if _, err := s.oauthRepo.LinkAccount(ctx, user.ID, providerName, info.ProviderID, info.Email); err != nil {
			return "", "", err
		}
	}

	accessToken, err := s.tokenService.GenerateToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
