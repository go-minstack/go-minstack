package repository

import "gorm.io/gorm"

// Pagination holds page/limit values for paginated queries.
// Build it from whatever source you have (gin query params, gRPC, CLI flags, etc).
type Pagination struct {
	Page  int // 1-based, default 1
	Limit int // default 10, max 100
}

// NewPagination creates a Pagination with safe defaults applied.
func NewPagination(page, limit int) *Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	return &Pagination{Page: page, Limit: limit}
}

// GetScope returns a GORM scope function that applies offset and limit.
func (p *Pagination) GetScope(db *gorm.DB) *gorm.DB {
	offset := (p.Page - 1) * p.Limit
	return db.Offset(offset).Limit(p.Limit)
}
