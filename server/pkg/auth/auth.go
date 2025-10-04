package authfx

import "go.uber.org/fx"

var Module = fx.Module(
	"authfx",
	fx.Provide(
		NewAuthRoutes,
		NewAuthController,
		NewAuthService,
	),
)
