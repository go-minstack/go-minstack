# go-minstack/postgres

PostgreSQL module for MinStack. Provides a GORM `*gorm.DB` connected to a PostgreSQL database.

## Installation

```sh
go get github.com/go-minstack/go-minstack/postgres
```

## Usage

Set `DB_URL` to your PostgreSQL DSN, then pass `postgres.Module()` to `core.New`.

```go
func main() {
    app := core.New(cli.Module(), postgres.Module())
    app.Provide(NewApp)
    app.Run()
}
```

```sh
DB_URL="host=localhost user=myuser password=mypass dbname=mydb port=5432 sslmode=disable" ./myapp
```

## API

### `postgres.Module() fx.Option`
Registers `*gorm.DB` into the fx container. Reads `DB_URL` from the environment.

## Example

See [examples/hello](examples/hello/main.go).

## Constraints

- Requires `DB_URL` environment variable to be set
- No HTTP server, no CLI runner
