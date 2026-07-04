package sqlite

import (
	"log/slog"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewDB(log *slog.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("MINSTACK_DB_URL")
	if dsn == "" {
		log.Warn("MINSTACK_DB_URL is not set, using in-memory SQLite database")
		dsn = ":memory:"
	}

	return gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 newGormLogger(log, 200*time.Millisecond),
	})
}
