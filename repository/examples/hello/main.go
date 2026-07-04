package main

import (
	"context"
	"fmt"

	"github.com/go-minstack/go-minstack/cli"
	"github.com/go-minstack/go-minstack/core"
	"github.com/go-minstack/go-minstack/repository"
	"github.com/go-minstack/go-minstack/sqlite"
	"gorm.io/gorm"
)

// ── Entity ───────────────────────────────────────────────────────────────────

// User uses gorm.Model (uint primary key) — compatible with all drivers.
type User struct {
	gorm.Model
	Name  string
	Email string
}

// UserRepository wraps the generic repository with domain-specific queries.
type UserRepository struct {
	*repository.Repository[User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{repository.NewRepository[User](db)}
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	return r.FindOne(repository.Where("email = ?", email))
}

// ── App ──────────────────────────────────────────────────────────────────────

type App struct {
	users *UserRepository
}

func NewApp(users *UserRepository) cli.ConsoleApp {
	return &App{users: users}
}

func (a *App) Run(_ context.Context) error {
	// Seed
	seeds := []User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Carol", Email: "carol@example.com"},
	}
	for i := range seeds {
		if err := a.users.Create(&seeds[i]); err != nil {
			return err
		}
	}

	// FindAll — ordered by name
	all, err := a.users.FindAll(repository.Order("name"))
	if err != nil {
		return err
	}
	fmt.Println("All users:")
	for _, u := range all {
		fmt.Printf("  [%d] %s <%s>\n", u.ID, u.Name, u.Email)
	}

	// FindByEmail — custom domain query
	alice, err := a.users.FindByEmail("alice@example.com")
	if err != nil {
		return err
	}
	fmt.Printf("\nFound by email: %s\n", alice.Name)

	// UpdatesByID — partial update
	if err := a.users.UpdatesByID(alice.ID, map[string]interface{}{"name": "Alice Updated"}); err != nil {
		return err
	}

	updated, _ := a.users.FindByID(alice.ID)
	fmt.Printf("After update:   %s\n", updated.Name)

	// Count
	total, _ := a.users.Count()
	fmt.Printf("\nTotal users: %d\n", total)

	// Paginate — page 1, limit 2
	page, err := a.users.FindAll(
		repository.Order("name"),
		repository.Paginate(repository.NewPagination(1, 2)),
	)
	if err != nil {
		return err
	}
	fmt.Println("\nPage 1 (limit 2):")
	for _, u := range page {
		fmt.Printf("  %s\n", u.Name)
	}

	// DeleteByID
	if err := a.users.DeleteByID(alice.ID); err != nil {
		return err
	}
	remaining, _ := a.users.Count()
	fmt.Printf("\nAfter delete: %d users remaining\n", remaining)

	return nil
}

// ── Bootstrap ─────────────────────────────────────────────────────────────────

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func main() {
	app := core.New(cli.Module(), sqlite.Module())
	app.Provide(NewUserRepository)
	app.Provide(NewApp)
	app.Invoke(migrate)
	app.Run()
}
