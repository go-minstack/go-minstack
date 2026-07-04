package main

import (
	"fmt"

	"github.com/go-minstack/go-minstack/core"
	"github.com/go-minstack/go-minstack/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func run(db *gorm.DB) {
	db.Create(&User{Name: "MinStack"})

	var user User
	db.First(&user)
	fmt.Printf("Hello, %s!\n", user.Name)
}

func main() {
	app := core.New(sqlite.Module())
	app.Invoke(migrate)
	app.Invoke(run)
	app.Run()
}
