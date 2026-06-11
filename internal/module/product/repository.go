package product

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	Transaction(ctx context.Context, fn func(repo Repository) error) error
	CreateProductImage(ctx context.Context, image *ProductImage) error
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

func (r *repository) CreateProductImage(ctx context.Context, image *ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}
