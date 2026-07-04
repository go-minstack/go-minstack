# go-minstack/repository

Generic, type-safe GORM repository for MinStack. Eliminates data access boilerplate while keeping queries explicit and composable.

## Installation

```sh
go get github.com/go-minstack/go-minstack/repository
```

## Usage

Embed `*repository.Repository[T, ID]` inside your domain entity file. The base UUID model lives in each database module — pick the one that matches your driver.

### PostgreSQL

Uses `postgres.UuidModel` — native `uuid` column type.

```go
import (
    "github.com/go-minstack/go-minstack/postgres"
    "github.com/go-minstack/go-minstack/repository"
    "gorm.io/gorm"
)

type User struct {
    postgres.UuidModel
    Name  string
    Email string
}

type UserRepository struct {
    *repository.Repository[User, uuid.UUID]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{repository.New[User, uuid.UUID](db)}
}
```

### MySQL

Uses `mysql.UuidModel` — UUID stored as `binary(16)`.

```go
import (
    "github.com/go-minstack/go-minstack/mysql"
    "github.com/go-minstack/go-minstack/repository"
    "gorm.io/gorm"
)

type User struct {
    mysql.UuidModel
    Name  string
    Email string
}

type UserRepository struct {
    *repository.Repository[User, mysql.UUID]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{repository.New[User, mysql.UUID](db)}
}
```

### SQLite

SQLite has no native UUID type. Use `uint` with `gorm.Model` for simplicity, or store UUIDs as `text`.

```go
import (
    "github.com/go-minstack/go-minstack/repository"
    "gorm.io/gorm"
)

// Option A — uint primary key (simplest)
type User struct {
    gorm.Model
    Name  string
    Email string
}

type UserRepository struct {
    *repository.Repository[User, uint]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{repository.New[User, uint](db)}
}
```

### Registering with FX

Regardless of driver, register with `app.Provide` so FX injects `*gorm.DB` automatically:

```go
app.Provide(user_entities.NewUserRepository)
```

## Querying

```go
// FindAll with options
users, err := repo.FindAll(
    repository.Where("active = ?", true),
    repository.Order("created_at", true), // DESC
    repository.Paginate(repository.NewPagination(1, 20)),
)

// FindOne
user, err := repo.FindOne(repository.Where("email = ?", email))

// Count
total, err := repo.Count(repository.Where("active = ?", true))
```

## API

### `repository.Repository[T, ID]`

| Method | Description |
|--------|-------------|
| `FindByID(id)` | Single record by primary key |
| `FindAll(opts...)` | All records matching options |
| `FindOne(opts...)` | First record matching options |
| `Count(opts...)` | Count records matching options |
| `Save(entity)` | Create a new record |
| `Update(entity)` | Full update (all fields) |
| `Updates(entity, columns)` | Partial update via map |
| `UpdatesByID(id, columns)` | Partial update by ID |
| `UpdateAll(values, opts...)` | Bulk update — requires at least one option |
| `Delete(entity)` | Delete a record |
| `DeleteByID(id)` | Delete by primary key |
| `DeleteAll(opts...)` | Bulk delete — requires at least one option |
| `WithTransaction(fn)` | Run fn inside a transaction |
| `DB()` | Direct `*gorm.DB` access for complex queries |

### Base models

Two options — pick based on whether you need UUID or `uint` primary keys:

**`gorm.Model` — uint primary key, works with all drivers:**
```go
type UserRepository struct {
    *repository.Repository[User, uint]
}
```

**`UuidModel` — UUID primary key, driver-specific:**

| Driver | Embed | ID type | Import |
|--------|-------|---------|--------|
| PostgreSQL | `postgres.UuidModel` | `uuid.UUID` | `github.com/go-minstack/go-minstack/postgres` |
| MySQL | `mysql.UuidModel` | `mysql.UUID` | `github.com/go-minstack/go-minstack/mysql` |

UUID variants auto-generate the ID in Go via `BeforeCreate` — no database function required.

### QueryOptions

| Function | Description |
|----------|-------------|
| `Where(query, args...)` | Add a WHERE clause |
| `Order(column, desc...)` | Add ORDER BY |
| `Preload(query, args...)` | Eager load associations |
| `Limit(n)` | Cap number of results |
| `Scope(fn)` | Apply a raw GORM scope |
| `Paginate(p)` | Apply a Pagination |

### `repository.Pagination`
```go
p := repository.NewPagination(page, limit) // defaults: page=1, limit=10, max=100
```

## Custom domain queries

Use `DB()` for joins, aggregations, or anything beyond the standard API:

```go
func (r *UserRepository) FindByEmailDomain(domain string) ([]User, error) {
    var users []User
    err := r.DB().Where("email LIKE ?", "%@"+domain).Find(&users).Error
    return users, err
}
```

## Constraints

- Requires a `*gorm.DB` — pair with `mysql`, `postgres`, or `sqlite` modules
- No FX `Module()` — it's a utility package, not an infrastructure provider
- `UpdateAll` and `DeleteAll` require at least one `QueryOption` as a safety guard
