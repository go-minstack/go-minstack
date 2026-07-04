package core

import "go.uber.org/fx"

// Module returns the core fx.Option. It is included automatically by New()
// and does not need to be added manually.
func Module() fx.Option {
	return fx.Module("core")
}
