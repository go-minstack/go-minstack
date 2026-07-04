package postgres

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UuidModel is the base model for PostgreSQL entities with a UUID primary key.
// The ID is generated in Go via BeforeCreate using the native uuid type.
type UuidModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *UuidModel) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
