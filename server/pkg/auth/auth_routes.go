package authfx

import (
	middlewarefx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthRoutesParams struct {
	fx.In
	Logger               *zap.Logger
	Router               *gin.Engine
	AuthController       *AuthController
	RequestBodyValidator *middlewarefx.RequestBodyValidator
}

type AuthRoutes struct {
	Logger               *zap.Logger
	Router               *gin.Engine
	AuthController       *AuthController
	RequestBodyValidator *middlewarefx.RequestBodyValidator
}

func NewAuthRoutes(params AuthRoutesParams) *AuthRoutes {
	return &AuthRoutes{
		Logger:               params.Logger,
		Router:               params.Router,
		AuthController:       params.AuthController,
		RequestBodyValidator: params.RequestBodyValidator,
	}
}

func (routes *AuthRoutes) Setup() {
	routes.Logger.Info("Setting up [Auth] routes.")
	authGroup := routes.Router.Group("api/v1/auth")

	authGroup.POST("/register",
		routes.RequestBodyValidator.Handler("register", RegisterBody{}),
		routes.AuthController.Register)

	authGroup.POST("/login",
		routes.RequestBodyValidator.Handler("login", LoginBody{}),
		routes.AuthController.Login)
}
