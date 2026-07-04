# go-minstack/sqlite

SQLite module for MinStack. Provides a GORM `*gorm.DB` backed by SQLite — no CGO required.

## Installation

```sh
go get github.com/go-minstack/go-minstack/sqlite
```

## Usage

Set `DB_URL` to a file path or `:memory:`, then pass `sqlite.Module()` to `core.New`.

```go
func main() {
    app := core.New(cli.Module(), sqlite.Module())
    app.Provide(NewApp)
    app.Run()
}
```

```sh
DB_URL=./data.db ./myapp
# or in-memory:
DB_URL=:memory: ./myapp
```

## API

### `sqlite.Module() fx.Option`
Registers `*gorm.DB` into the fx container. Reads `DB_URL` from the environment.

## Example

See [examples/hello](examples/hello/main.go).

## Constraints

- Requires `DB_URL` environment variable to be set
- Pure Go — no CGO, no system SQLite library needed
- No HTTP server, no CLI runner
