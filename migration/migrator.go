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

type migratorConfig struct {
	goMigrations []*goose.Migration
}

// Option configures a Migrator or Module.
type Option func(*migratorConfig)

// WithGoMigrations registers programmatic goose migrations alongside embedded SQL.
// SQL and Go migrations share one version sequence — do not use the same version
// number in both (e.g. 00002_foo.sql and NewGoMigration(2, ...) together).
func WithGoMigrations(m ...*goose.Migration) Option {
	return func(c *migratorConfig) {
		c.goMigrations = append(c.goMigrations, m...)
	}
}

func applyOptions(opts []Option) migratorConfig {
	var cfg migratorConfig
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// Module provides a *Migrator into the FX container.
// Opt-in to running migrations by invoking migration.Run:
//
//	app := core.New(postgres.Module, migration.Module(migrationsFS))
//	app.Invoke(migration.Run)
func Module(fsys fs.FS, opts ...Option) fx.Option {
	return fx.Module("migration",
		fx.Provide(func(db *gorm.DB) *Migrator {
			return New(db, slog.Default(), fsys, opts...)
		}),
	)
}

type Migrator struct {
	db           *gorm.DB
	log          *slog.Logger
	fs           fs.FS
	goMigrations []*goose.Migration
}

// New creates a Migrator with a custom logger. Use when wiring manually via Register.
func New(db *gorm.DB, log *slog.Logger, fsys fs.FS, opts ...Option) *Migrator {
	if log == nil {
		log = slog.Default()
	}
	cfg := applyOptions(opts)
	return &Migrator{
		db:           db,
		log:          log,
		fs:           fsys,
		goMigrations: cfg.goMigrations,
	}
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

	providerOpts := []goose.ProviderOption{
		goose.WithLogger(&gooseLogger{log: m.log}),
		goose.WithDisableGlobalRegistry(true),
	}
	if len(m.goMigrations) > 0 {
		providerOpts = append(providerOpts, goose.WithGoMigrations(m.goMigrations...))
	}

	provider, err := goose.NewProvider(dialect, sqlDB, m.fs, providerOpts...)
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
