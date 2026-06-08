package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/cache"
	"github.com/AriaPutra01/go-commerce/internal/token"
	"github.com/AriaPutra01/go-commerce/internal/util"
)

type service struct {
	jwt        *token.JWTMaker
	repository Repository
	cache      cache.Cache
}

func NewService(jwt *token.JWTMaker, repository Repository, cache cache.Cache) *service {
	return &service{
		jwt:        jwt,
		repository: repository,
		cache:      cache,
	}
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repository.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("Invalid email or password")
	}

	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errors.New("Invalid email or password")
	}

	accessToken, err := s.jwt.GenerateToken(user.ID.String(), user.Email, user.Role, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken := fmt.Sprintf("refresh:%s", rand.Text())

	if err := s.cache.Save(ctx, refreshToken, user.ID.String(), 7*24*time.Hour); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
