package cli

import (
	"context"
	"os"

	"go.uber.org/fx"
)

type ConsoleApp interface {
	Run(ctx context.Context) error
}

func newConsoleRunner(lc fx.Lifecycle, app ConsoleApp, shutdowner fx.Shutdowner) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := app.Run(context.Background()); err != nil {
					os.Exit(1)
				}
				shutdowner.Shutdown()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}
