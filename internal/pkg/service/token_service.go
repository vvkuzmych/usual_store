package service

import (
	"context"
	"time"
	"usual_store/internal/models"
	"usual_store/internal/pkg/repository"
)

type TokenService struct {
	Repo repository.TokenRepository
}

func NewTokenService(repo repository.TokenRepository) *TokenService {
	return &TokenService{Repo: repo}
}

// CreateToken generates a token and stores it in the database.
func (s *TokenService) CreateToken(ctx context.Context, user models.User, ttl time.Duration, scope string) (*models.Token, error) {
	token, err := models.GenerateToken(user.ID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = s.Repo.InsertToken(ctx, token, user)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// GetUserForToken retrieves a user based on the provided token.
func (s *TokenService) GetUserForToken(ctx context.Context, token string) (*models.User, error) {
	return s.Repo.GetUserForToken(ctx, token)
}
