package authfx

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type AuthRoutesParams struct {
	fx.In
	Logger         *zap.Logger
	Router         *gin.Engine
	AuthController *AuthController
}

type AuthRoutes struct {
	Logger         *zap.Logger
	Router         *gin.Engine
	AuthController *AuthController
}

func NewAuthRoutes(params AuthRoutesParams) *AuthRoutes {
	return &AuthRoutes{
		Logger:         params.Logger,
		Router:         params.Router,
		AuthController: params.AuthController,
	}
}

func (routes *AuthRoutes) Setup() {
	routes.Logger.Info("Setting up [Auth] routes.")
	authGroup := routes.Router.Group("api/v1/auth")

	authGroup.POST("/register",
		middlewares.RequestBodyValidator(routes.Logger, "register", RegisterBody{}),
		routes.AuthController.Register)

	authGroup.POST("/login",
		middlewares.RequestBodyValidator(routes.Logger, "login", LoginBody{}),
		routes.AuthController.Login)
}
