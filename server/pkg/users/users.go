package usersfx

import "go.uber.org/fx"

var Module = fx.Module(
	"usersfx",
	fx.Provide(
		NewUsersRoutes,
		NewUsersController,
		NewUserService,
	),
)
