package mysql

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(log *slog.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("MINSTACK_DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("MINSTACK_DB_URL is not set")
	}

	return gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 newGormLogger(log, 200*time.Millisecond),
	})
}
