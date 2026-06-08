package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"column:id"`
	Email        string     `gorm:"column:email"`
	PasswordHash string     `gorm:"column:password_hash"`
	FullName     string     `gorm:"column:full_name"`
	Phone        *string    `gorm:"column:phone"`
	Role         string     `gorm:"column:role"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}
