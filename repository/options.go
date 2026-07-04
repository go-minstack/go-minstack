package repository

import "gorm.io/gorm"

// QueryOption is a composable function applied to a GORM query.
// Stack multiple options to build complex queries without verbose structs.
type QueryOption func(*gorm.DB) *gorm.DB

// Where adds a WHERE clause.
//
//	Where("name = ?", "MinStack")
func Where(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(query, args...)
	}
}

// Order adds an ORDER BY clause. Pass desc=true for descending order.
//
//	Order("created_at", true)  // ORDER BY created_at DESC
func Order(column string, desc ...bool) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(desc) > 0 && desc[0] {
			return db.Order(column + " DESC")
		}
		return db.Order(column)
	}
}

// Preload eagerly loads an association.
//
//	Preload("Posts")
//	Preload("Posts", "published = ?", true)
func Preload(query string, args ...any) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query, args...)
	}
}

// Limit caps the number of results.
func Limit(n int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(n)
	}
}

// Scope applies a raw GORM scope function. Useful for reusing complex query logic.
func Scope(fn func(*gorm.DB) *gorm.DB) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(fn)
	}
}

// Paginate applies a Pagination to the query.
func Paginate(p *Pagination) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		return p.GetScope(db)
	}
}

func applyOptions(db *gorm.DB, opts []QueryOption) *gorm.DB {
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}
