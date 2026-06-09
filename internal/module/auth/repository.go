package auth

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Transaction(ctx context.Context, fn func(repo Repository) error) error
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	ExistsUserByEmail(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *User) error
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

func (r *repository) FindUserByID(ctx context.Context, id string) (*User, error) {
	user := new(User)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Limit(1).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}
