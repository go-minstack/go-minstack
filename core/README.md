# go-minstack/core

The minimal foundation for MinStack. Every other module depends on this — nothing else does.

## Installation

```sh
go get github.com/go-minstack/go-minstack/core
```

## Usage

`New()` is for composing **modules**. `Provide` and `Invoke` are for your own **constructors and startup functions**.

```go
app := core.New()
app.Provide(NewGreeter)
app.Invoke(run)
app.Run()
```

## API

### `core.New(modules ...fx.Option) *App`
Creates the application. Includes `core.Module()` automatically.
Pass other MinStack modules as they become available.

### `app.Provide(constructors ...interface{})`
Registers your own constructors into the fx container.

### `app.Invoke(funcs ...interface{})`
Registers your own startup functions. fx resolves their dependencies automatically.

### `app.Run()`
Builds the fx app and blocks until a shutdown signal is received.

### `app.Start(ctx context.Context) error`
Builds and starts the fx app without blocking.

### `app.Stop(ctx context.Context) error`
Stops the running app.

### `core.Module() fx.Option`
Returns the core fx module. Included automatically by `New()` — no need to pass it manually.

## Example

See [examples/hello](examples/hello/main.go).

## Constraints

- No HTTP server, no database drivers, no concrete infrastructure
- Infrastructure-agnostic by design
- Interfaces live here, implementations belong in their own module
