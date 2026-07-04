package sqlite

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UuidModel is the base model for SQLite entities with a UUID primary key.
// SQLite has no native UUID type — the ID is stored as text.
// The ID is auto-generated in Go via BeforeCreate.
type UuidModel struct {
	ID        uuid.UUID `gorm:"type:text;primaryKey"`
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
