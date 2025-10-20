package libfx

import "go.uber.org/fx"

var Module = fx.Module(
	"libfx",
	fx.Provide(
		NewLogger,
		NewDatabase,
		NewRouter,
		NewMailClient,
	),
)
