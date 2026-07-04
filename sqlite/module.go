package sqlite

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("sqlite",
		fx.Provide(NewDB),
	)
}
