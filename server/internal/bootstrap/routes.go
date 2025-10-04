package bootstrapfx

import (
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	"go.uber.org/fx"
)

// Every route must have Setup method
type Route interface {
	Setup()
}

type RoutesParams struct {
	fx.In
	AuthRoutes *authfx.AuthRoutes
}

type Routes []Route

func NewRoutes(params RoutesParams) Routes {
	return Routes{
		params.AuthRoutes,
	}
}

func (routes Routes) Setup() {
	// Trigger Setup method of every route
	for _, route := range routes {
		route.Setup()
	}
}
