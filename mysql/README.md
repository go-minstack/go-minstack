# go-minstack/mysql

MySQL module for MinStack. Provides a GORM `*gorm.DB` and a `binary(16)` UUID type optimized for MySQL.

## Installation

```sh
go get github.com/go-minstack/go-minstack/mysql
```

## Usage

Set `DB_URL` to your MySQL DSN, then pass `mysql.Module()` to `core.New`.

```go
func main() {
    app := core.New(cli.Module(), mysql.Module())
    app.Provide(NewApp)
    app.Run()
}
```

```sh
DB_URL="user:pass@tcp(localhost:3306)/dbname?parseTime=True" ./myapp
```

## API

### `mysql.Module() fx.Option`
Registers `*gorm.DB` into the fx container. Reads `DB_URL` from the environment.

### `mysql.UUID`
A GORM-compatible UUID type stored as `binary(16)` in MySQL.

```go
type User struct {
    ID   mysql.UUID `gorm:"primaryKey"`
    Name string
}
```

| Function | Description |
|----------|-------------|
| `mysql.NewUUID() UUID` | Generate a new random UUID |
| `mysql.ParseUUID(s string) (UUID, error)` | Parse a UUID string |
| `mysql.MustParseUUID(s string) UUID` | Parse or panic |
| `(UUID).String() string` | Format as standard UUID string |
| `(UUID).IsZero() bool` | Check if UUID is the zero value |

## Example

See [examples/hello](examples/hello/main.go).

## Constraints

- Requires `DB_URL` environment variable to be set
- UUID type is MySQL-specific — PostgreSQL users should use `uuid.UUID` directly
- No HTTP server, no CLI runner
