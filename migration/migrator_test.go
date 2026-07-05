package migration_test

import (
	"context"
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"log/slog"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/go-minstack/go-minstack/migration"
	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed testdata/sqlonly
var sqlOnlyRoot embed.FS

//go:embed testdata/withduplicate
var withDuplicateRoot embed.FS

func migrationFS(t *testing.T, root embed.FS, dir string) fs.FS {
	t.Helper()
	sub, err := fs.Sub(root, dir)
	if err != nil {
		t.Fatalf("fs.Sub: %v", err)
	}
	return sub
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestMigrator_SQLOnly(t *testing.T) {
	db := openTestDB(t)
	m := migration.New(db, discardLogger(), migrationFS(t, sqlOnlyRoot, "testdata/sqlonly"))
	if err := m.Up(); err != nil {
		t.Fatalf("Up: %v", err)
	}

	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM items").Scan(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected empty items table, got %d rows", count)
	}
}

func TestMigrator_SQLAndGo(t *testing.T) {
	db := openTestDB(t)

	goMig := goose.NewGoMigration(2,
		&goose.GoFunc{
			RunTx: func(ctx context.Context, tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, `INSERT INTO items (name) VALUES ('from-go')`)
				return err
			},
		},
		&goose.GoFunc{
			RunTx: func(ctx context.Context, tx *sql.Tx) error {
				_, err := tx.ExecContext(ctx, `DELETE FROM items WHERE name = 'from-go'`)
				return err
			},
		},
	)

	m := migration.New(db, discardLogger(), migrationFS(t, sqlOnlyRoot, "testdata/sqlonly"), migration.WithGoMigrations(goMig))
	if err := m.Up(); err != nil {
		t.Fatalf("Up: %v", err)
	}

	var name string
	if err := db.Raw("SELECT name FROM items LIMIT 1").Scan(&name).Error; err != nil {
		t.Fatalf("select: %v", err)
	}
	if name != "from-go" {
		t.Fatalf("expected from-go, got %q", name)
	}
}

func TestMigrator_DuplicateVersionRejected(t *testing.T) {
	db := openTestDB(t)

	goMig := goose.NewGoMigration(2, &goose.GoFunc{
		RunTx: func(ctx context.Context, tx *sql.Tx) error { return nil },
	}, nil)

	m := migration.New(db, discardLogger(), migrationFS(t, withDuplicateRoot, "testdata/withduplicate"), migration.WithGoMigrations(goMig))
	if err := m.Up(); err == nil {
		t.Fatal("expected duplicate version error, got nil")
	}
}
