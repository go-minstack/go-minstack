# MinStack

MinStack is a modular backend foundation for Go, now maintained in a single repository with focused sibling packages.

## Packages

- `core`
- `auth`
- `gin`
- `logger`
- `mysql`
- `postgres`
- `sqlite`
- `repository`
- `cli`
- `web`
- `migration`

## Import Style

Consumers import only the packages they need, for example:

```go
import (
    "github.com/go-minstack/go-minstack/auth"
    "github.com/go-minstack/go-minstack/core"
)
```

Package boundaries remain explicit: importing `auth` or `core` does not import database packages unless those packages are directly referenced in code.

## New Project AI Context

If you want a fresh app repository to come with built-in AI guidance for go-minstack usage, start from `AGENTS.example.md` and copy it into the new project's root as `AGENTS.md`.

That template tells an AI assistant:

- which import paths to use
- how to compose apps with `core.New(...)`
- how to structure `internal/` domains
- how to register routes, services, and repositories the MinStack way
