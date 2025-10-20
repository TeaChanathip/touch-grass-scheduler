package usersfx

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/endpoints"
	middlewarefx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/middlewares"
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type UsersRoutesParams struct {
	fx.In
	Logger          *zap.Logger
	Router          *gin.Engine
	AuthMiddleware  *middlewarefx.AuthMiddleware
	UsersController *UsersController
}

type UsersRoutes struct {
	Logger          *zap.Logger
	Router          *gin.Engine
	UsersController *UsersController
	AuthMiddleware  *middlewarefx.AuthMiddleware
}

func NewUsersRoutes(params UsersRoutesParams) *UsersRoutes {
	return &UsersRoutes{
		Logger:          params.Logger,
		Router:          params.Router,
		UsersController: params.UsersController,
		AuthMiddleware:  params.AuthMiddleware,
	}
}

func (routes *UsersRoutes) Setup() {
	routes.Logger.Info("Setting up [Users] routes.")

	routes.Router.GET(string(endpoints.GetMeV1),
		routes.AuthMiddleware.Handler(),
		routes.UsersController.GetMe)

	routes.Router.GET(string(endpoints.GetUserWithIDV1)+"/:id",
		routes.AuthMiddleware.HandlerWithRole(types.UserRoleAdmin),
		routes.UsersController.GetUserByID)

	// usersGroup.PUT("users/:id")
	// usersGroup.DELETE("users/:id")
}
