package auth

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Transaction(ctx context.Context, fn func(repo Repository) error) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Transaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		repo := NewRepository(tx)
		if err := fn(repo); err != nil {
			return err
		}
		return nil
	})
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	user := new(User)
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
