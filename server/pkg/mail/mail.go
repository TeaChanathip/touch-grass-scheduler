package mailfx

import "go.uber.org/fx"

var Module = fx.Module(
	"mailfx",
	fx.Provide(
		NewMailService,
	),
)
