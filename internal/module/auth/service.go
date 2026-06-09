package auth

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/cache"
	"github.com/AriaPutra01/go-commerce/internal/constant"
	"github.com/AriaPutra01/go-commerce/internal/token"
	"github.com/AriaPutra01/go-commerce/internal/util"
	"github.com/google/uuid"
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
		return nil, constant.ErrInvalidCredentials
	}

	if err := util.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, constant.ErrInvalidCredentials
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

func (s *service) Register(ctx context.Context, req *RegisterRequest) error {
	return s.repository.Transaction(ctx, func(repo Repository) error {
		isTaken, err := repo.ExistsUserByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		if isTaken {
			return constant.ErrEmailAlreadyExists
		}

		hashedPass, err := util.HashPassword(req.Password)
		if err != nil {
			return err
		}

		newUser := &User{
			ID:           uuid.New(),
			Email:        req.Email,
			PasswordHash: hashedPass,
			FullName:     req.FullName,
			Phone:        &req.Phone,
			Role:         string(constant.RoleUser),
		}

		return repo.CreateUser(ctx, newUser)
	})
}

func (s *service) Refresh(ctx context.Context, refreshCookie string) (*RefreshResponse, error) {
	userID, err := s.cache.Get(ctx, refreshCookie)
	if err != nil {
		return nil, constant.ErrInvalidRefreshToken
	}

	user, err := s.repository.FindUserByID(ctx, userID)
	if err != nil {
		return nil, constant.ErrInvalidRefreshToken
	}

	newAccessToken, err := s.jwt.GenerateToken(user.ID.String(), user.Email, user.Role, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	newRefreshToken := fmt.Sprintf("refresh:%s", rand.Text())

	if err := s.cache.Save(ctx, newRefreshToken, user.ID.String(), 7*24*time.Hour); err != nil {
		return nil, err
	}

	if err := s.cache.Delete(ctx, refreshCookie); err != nil {
		return nil, err
	}

	return &RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
