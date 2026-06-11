package product

import (
	"time"

	"github.com/AriaPutra01/go-commerce/internal/module/auth"
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID       `gorm:"column:id"`
	Name        string          `gorm:"column:name"`
	Description *string         `gorm:"column:description"`
	Price       int             `gorm:"column:price"`
	Stock       int             `gorm:"column:stock"`
	CreatedBy   uuid.UUID       `gorm:"column:created_by"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
	DeletedAt   *time.Time      `gorm:"column:deleted_at"`
	User        auth.User       `gorm:"foreignKey:CreatedBy"`
	Images      []*ProductImage `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	ID        uuid.UUID  `gorm:"column:id"`
	ProductID *uuid.UUID `gorm:"column:product_id"`
	UrlPath   string     `gorm:"column:url_path"`
	IsPrimary bool       `gorm:"column:is_primary"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}
