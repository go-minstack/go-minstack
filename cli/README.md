# go-minstack/cli

Run a Go program as a one-shot CLI process using the MinStack module system. The app exits automatically when your logic finishes — no manual shutdown needed.

## Installation

```sh
go get github.com/go-minstack/go-minstack/cli
```

## Usage

Implement `cli.ConsoleApp` and register it with `app.Provide`. The app exits when `Run` returns.

```go
type App struct{}

func NewApp() cli.ConsoleApp { return &App{} }

func (a *App) Run(ctx context.Context) error {
    fmt.Println("Hello!")
    return nil
}

func main() {
    app := core.New(cli.Module())
    app.Provide(NewApp)
    app.Run()
}
```

## API

### `cli.ConsoleApp`
Interface your app must implement:
```go
type ConsoleApp interface {
    Run(ctx context.Context) error
}
```

### `cli.Module() fx.Option`
Wires the console runner into the fx lifecycle. Call inside `core.New(...)`.

## Example

See [examples/hello](examples/hello/main.go).

## Constraints

- No HTTP server, no database drivers
- One `ConsoleApp` per process
- If `Run` returns a non-nil error, the process exits with code 1
