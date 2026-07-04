package mysql

import "go.uber.org/fx"

func Module() fx.Option {
	return fx.Module("mysql",
		fx.Provide(NewDB),
	)
}
