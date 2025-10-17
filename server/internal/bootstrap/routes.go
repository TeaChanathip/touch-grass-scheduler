package bootstrapfx

import (
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	usersfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/users"
	"go.uber.org/fx"
)

// Every route must have Setup method
type Route interface {
	Setup()
}

type RoutesParams struct {
	fx.In
	AuthRoutes  *authfx.AuthRoutes
	UsersRoutes *usersfx.UsersRoutes
}

type Routes []Route

func NewRoutes(params RoutesParams) Routes {
	return Routes{
		params.AuthRoutes,
		params.UsersRoutes,
	}
}

func (routes Routes) Setup() {
	// Trigger Setup method of every route
	for _, route := range routes {
		route.Setup()
	}
}
