package mysql

import (
	"time"

	"gorm.io/gorm"
)

// UuidModel is the base model for MySQL entities with a binary(16) UUID primary key.
// Uses mysql.UUID which stores UUIDs efficiently as binary(16).
type UuidModel struct {
	ID        UUID `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *UuidModel) BeforeCreate(_ *gorm.DB) error {
	if m.ID.IsZero() {
		m.ID = NewUUID()
	}
	return nil
}
