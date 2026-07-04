package repository

import (
	"errors"

	"gorm.io/gorm"
)

// Repo is the underlying generic GORM-backed data access layer.
// T is the entity type, ID is its primary key type.
// For the common cases prefer the typed aliases:
//   - Repository[T]     — uint primary key (matches gorm.Model)
//   - UuidRepository[T] — uuid.UUID primary key
type Repo[T any, ID any] struct {
	db *gorm.DB
}

// New creates a Repo for entity T with primary key type ID.
// For the common cases prefer NewRepository (uint) or NewUuidRepository (uuid.UUID).
func New[T any, ID any](db *gorm.DB) *Repo[T, ID] {
	return &Repo[T, ID]{db: db}
}

// DB returns the underlying *gorm.DB for queries beyond the standard API.
// Use sparingly — prefer QueryOptions when possible.
func (r *Repo[T, ID]) DB() *gorm.DB {
	return r.db
}

// ── Reads ────────────────────────────────────────────────────────────────────

// FindByID returns a single entity by primary key, or an error if not found.
func (r *Repo[T, ID]) FindByID(id ID) (*T, error) {
	var entity T
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll returns all entities matching the given options.
func (r *Repo[T, ID]) FindAll(opts ...QueryOption) ([]T, error) {
	var entities []T
	db := applyOptions(r.db, opts)
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// FindOne returns the first entity matching the given options, or an error if none found.
func (r *Repo[T, ID]) FindOne(opts ...QueryOption) (*T, error) {
	var entity T
	db := applyOptions(r.db, opts)
	if err := db.First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// Count returns the number of entities matching the given options.
func (r *Repo[T, ID]) Count(opts ...QueryOption) (int64, error) {
	var entity T
	var count int64
	db := applyOptions(r.db, opts)
	if err := db.Model(&entity).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ── Writes ───────────────────────────────────────────────────────────────────

// Create inserts a new entity into the database.
func (r *Repo[T, ID]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update performs a full update of the entity (all fields).
func (r *Repo[T, ID]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Updates performs a partial update using the provided column map.
func (r *Repo[T, ID]) Updates(entity *T, columns map[string]interface{}) error {
	return r.db.Model(entity).Updates(columns).Error
}

// UpdatesByID performs a partial update by primary key.
func (r *Repo[T, ID]) UpdatesByID(id ID, columns map[string]interface{}) error {
	var entity T
	return r.db.Model(&entity).Where("id = ?", id).Updates(columns).Error
}

// UpdateAll performs a conditional bulk update. Requires at least one QueryOption
// to prevent accidental full-table updates. Use DB() for intentional full-table ops.
func (r *Repo[T, ID]) UpdateAll(values map[string]any, opts ...QueryOption) error {
	if len(opts) == 0 {
		return errors.New("repository: UpdateAll requires at least one QueryOption to prevent full-table updates; use DB() for intentional full-table operations")
	}
	var entity T
	db := applyOptions(r.db.Model(&entity), opts)
	return db.Updates(values).Error
}

// ── Deletes ──────────────────────────────────────────────────────────────────

// Delete removes a single entity.
func (r *Repo[T, ID]) Delete(entity *T) error {
	return r.db.Delete(entity).Error
}

// DeleteByID removes an entity by primary key.
func (r *Repo[T, ID]) DeleteByID(id ID) error {
	var entity T
	return r.db.Delete(&entity, "id = ?", id).Error
}

// DeleteAll performs a conditional bulk delete. Requires at least one QueryOption
// to prevent accidental full-table deletes. Use DB() for intentional full-table ops.
func (r *Repo[T, ID]) DeleteAll(opts ...QueryOption) error {
	if len(opts) == 0 {
		return errors.New("repository: DeleteAll requires at least one QueryOption to prevent full-table deletes; use DB() for intentional full-table operations")
	}
	var entity T
	db := applyOptions(r.db, opts)
	return db.Delete(&entity).Error
}

// ── Transactions ─────────────────────────────────────────────────────────────

// WithTransaction runs fn inside a database transaction.
// The transaction is automatically committed on success and rolled back on error.
func (r *Repo[T, ID]) WithTransaction(fn func(*Repo[T, ID]) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(&Repo[T, ID]{db: tx})
	})
}
