package authfx

import (
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
	routes.Router.POST("/register", routes.AuthController.Register)
}
