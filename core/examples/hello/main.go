package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-minstack/go-minstack/core"
)

type Greeter struct {
	name string
}

func NewGreeter() *Greeter {
	return &Greeter{name: "MinStack"}
}

func (g *Greeter) Hello() string {
	return fmt.Sprintf("Hello from %s!", g.name)
}

func run(g *Greeter) {
	fmt.Println(g.Hello())
}

func main() {
	app := core.New()
	app.Provide(NewGreeter)
	app.Invoke(run)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
	defer app.Stop(ctx)
}
