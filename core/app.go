package core

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"github.com/go-minstack/go-minstack/logger"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type App struct {
	opts []fx.Option
	app  *fx.App
}

func New(opts ...fx.Option) *App {
	dotenvErr := godotenv.Load()

	baseOpts := []fx.Option{
		logger.Module(),
		coreModule(dotenvErr),
	}
	return &App{opts: append(baseOpts, opts...)}
}

func coreModule(dotenvErr error) fx.Option {
	return fx.Module("core", fx.Invoke(func(log *slog.Logger) {
		if dotenvErr == nil {
			return
		}
		if errors.Is(dotenvErr, os.ErrNotExist) {
			log.Debug("no .env file found, using environment variables")
			return
		}
		log.Warn("could not load .env file", "err", dotenvErr)
	}))
}

func (a *App) Provide(constructors ...interface{}) {
	a.opts = append(a.opts, fx.Provide(constructors...))
}

// Use appends raw fx.Options to the application. This is the correct way
// to register options returned by helpers like ProvideAs.
func (a *App) Use(opts ...fx.Option) {
	a.opts = append(a.opts, opts...)
}

// ProvideAs returns an fx.Option that provides a constructor's return value
// as interface type T. Works with both app.Provide and fx.Module.
//
// Usage:
//
//	app.Provide(core.ProvideAs[MyInterface](NewMyImplementation))
//
//	fx.Module("name", core.ProvideAs[MyInterface](NewMyImplementation))
func ProvideAs[T any](constructor any) fx.Option {
	return fx.Provide(fx.Annotate(
		constructor,
		fx.As(new(T)),
	))
}

func (a *App) Invoke(constructors ...interface{}) {
	a.opts = append(a.opts, fx.Invoke(constructors...))
}

func (a *App) Run() {
	a.app = fx.New(a.opts...)
	a.app.Run()
}

func (a *App) Start(ctx context.Context) error {
	a.app = fx.New(a.opts...)
	return a.app.Start(ctx)
}

func (a *App) Stop(ctx context.Context) error {
	return a.app.Stop(ctx)
}
