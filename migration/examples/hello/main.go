package main

import (
	"github.com/go-minstack/go-minstack/core"
	"github.com/go-minstack/go-minstack/migration"
	"github.com/go-minstack/go-minstack/migration/examples/hello/migrations"
	"github.com/go-minstack/go-minstack/sqlite"
)

func main() {
	app := core.New(
		sqlite.Module(),
		migration.Module(migrations.FS),
	)
	app.Invoke(migration.Run)
	app.Run()
}
