package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository is a Repo with uint primary key — matches the gorm.Model default.
type Repository[T any] = Repo[T, uint]

// UuidRepository is a Repo with uuid.UUID primary key.
type UuidRepository[T any] = Repo[T, uuid.UUID]

// NewRepository creates a Repository[T] (uint primary key).
func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return New[T, uint](db)
}

// NewUuidRepository creates a UuidRepository[T] (uuid.UUID primary key).
func NewUuidRepository[T any](db *gorm.DB) *UuidRepository[T] {
	return New[T, uuid.UUID](db)
}
