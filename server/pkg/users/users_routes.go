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
	Logger               *zap.Logger
	Router               *gin.Engine
	AuthMiddleware       *middlewarefx.AuthMiddleware
	UsersController      *UsersController
	RequestBodyValidator *middlewarefx.RequestBodyValidator
}

type UsersRoutes struct {
	Logger               *zap.Logger
	Router               *gin.Engine
	UsersController      *UsersController
	AuthMiddleware       *middlewarefx.AuthMiddleware
	RequestBodyValidator *middlewarefx.RequestBodyValidator
}

func NewUsersRoutes(params UsersRoutesParams) *UsersRoutes {
	return &UsersRoutes{
		Logger:               params.Logger,
		Router:               params.Router,
		UsersController:      params.UsersController,
		AuthMiddleware:       params.AuthMiddleware,
		RequestBodyValidator: params.RequestBodyValidator,
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

	routes.Router.PUT(string(endpoints.UpdateUserV1),
		routes.AuthMiddleware.HandlerWithRole(types.UserRoleStudent,
			types.UserRoleTeacher,
			types.UserRoleGuardian),
		routes.RequestBodyValidator.Handler("update-user", UpdateUserBody{}),
		routes.UsersController.UpdateUser)

	routes.Router.GET(string(endpoints.GetUploadAvartarSignedURL),
		routes.AuthMiddleware.HandlerWithRole(types.UserRoleStudent,
			types.UserRoleTeacher,
			types.UserRoleGuardian),
		routes.UsersController.GetUploadAvartarSignedURL)

	// usersGroup.DELETE("users/:id")
}
