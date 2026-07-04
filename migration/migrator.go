package migration

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Module provides a *Migrator into the FX container.
// Opt-in to running migrations by invoking migration.Run:
//
//	app := core.New(postgres.Module, migration.Module(migrationsFS))
//	app.Invoke(migration.Run)
func Module(fsys fs.FS) fx.Option {
	return fx.Module("migration",
		fx.Provide(func(db *gorm.DB) *Migrator {
			return New(db, slog.Default(), fsys)
		}),
	)
}

type Migrator struct {
	db  *gorm.DB
	log *slog.Logger
	fs  fs.FS
}

// New creates a Migrator with a custom logger. Use when wiring manually via Register.
func New(db *gorm.DB, log *slog.Logger, fsys fs.FS) *Migrator {
	if log == nil {
		log = slog.Default()
	}
	return &Migrator{db: db, log: log, fs: fsys}
}

// Run is the FX invoke target for manual wiring: app.Invoke(migration.Run).
func Run(m *Migrator) error {
	return m.Up()
}

// Up applies all pending migrations.
func (m *Migrator) Up() error {
	dialect, err := dialectOf(m.db)
	if err != nil {
		return err
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider(dialect, sqlDB, m.fs,
		goose.WithLogger(&gooseLogger{log: m.log}),
	)
	if err != nil {
		return err
	}

	_, err = provider.Up(context.Background())
	return err
}

func dialectOf(db *gorm.DB) (goose.Dialect, error) {
	switch db.Dialector.Name() {
	case "postgres":
		return goose.DialectPostgres, nil
	case "mysql":
		return goose.DialectMySQL, nil
	case "sqlite":
		return goose.DialectSQLite3, nil
	default:
		return "", fmt.Errorf("migration: unsupported dialect %q", db.Dialector.Name())
	}
}

type gooseLogger struct{ log *slog.Logger }

func (g *gooseLogger) Printf(format string, v ...any) {
	g.log.Info(fmt.Sprintf(format, v...))
}

func (g *gooseLogger) Fatalf(format string, v ...any) {
	g.log.Error(fmt.Sprintf(format, v...))
	os.Exit(1)
}
