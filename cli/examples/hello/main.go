package main

import (
	"context"
	"fmt"

	"github.com/go-minstack/go-minstack/cli"
	"github.com/go-minstack/go-minstack/core"
	"github.com/go-minstack/go-minstack/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

type App struct {
	db *gorm.DB
}

func NewApp(db *gorm.DB) cli.ConsoleApp {
	return &App{db: db}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func (a *App) Run(ctx context.Context) error {
	a.db.Create(&User{Name: "MinStack"})

	var user User
	a.db.First(&user)
	fmt.Printf("Hello, %s!\n", user.Name)

	return nil
}

func main() {
	app := core.New(cli.Module(), sqlite.Module())
	app.Provide(NewApp)
	app.Invoke(migrate)
	app.Run()
}
