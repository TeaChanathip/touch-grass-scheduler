package usersfx

import "go.uber.org/fx"

var Module = fx.Module(
	"userfx",
	fx.Provide(
		NewUserService,
	),
)
