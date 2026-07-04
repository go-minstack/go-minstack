package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/go-minstack/go-minstack/core"
)

func run(log *slog.Logger) {
	log.Info("Hello from MinStack!", "app", "hello")
}

func main() {
	app := core.New()
	app.Invoke(run)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer app.Stop(ctx)
}
