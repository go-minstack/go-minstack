package main

import (
	"fmt"

	"github.com/go-minstack/go-minstack/core"
	"github.com/go-minstack/go-minstack/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID   mysql.UUID `gorm:"primaryKey"`
	Name string
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}

func run(db *gorm.DB) {
	db.Create(&User{ID: mysql.NewUUID(), Name: "MinStack"})

	var user User
	db.First(&user)
	fmt.Printf("Hello, %s!\n", user.Name)
}

func main() {
	app := core.New(mysql.Module())
	app.Invoke(migrate)
	app.Invoke(run)
	app.Run()
}
