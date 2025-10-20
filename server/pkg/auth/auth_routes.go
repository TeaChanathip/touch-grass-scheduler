package authfx

import (
	"github.com/TeaChanathip/touch-grass-scheduler/server/internal/endpoints"
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
	AuthMiddleware       *middlewarefx.AuthMiddleware
}

type AuthRoutes struct {
	Logger               *zap.Logger
	Router               *gin.Engine
	AuthController       *AuthController
	RequestBodyValidator *middlewarefx.RequestBodyValidator
	AuthMiddleware       *middlewarefx.AuthMiddleware
}

func NewAuthRoutes(params AuthRoutesParams) *AuthRoutes {
	return &AuthRoutes{
		Logger:               params.Logger,
		Router:               params.Router,
		AuthController:       params.AuthController,
		RequestBodyValidator: params.RequestBodyValidator,
		AuthMiddleware:       params.AuthMiddleware,
	}
}

func (routes *AuthRoutes) Setup() {
	// routes.Logger.Info("Setting up [Auth] routes.")

	routes.Router.GET(string(endpoints.GetRegistrationMailV1)+"/:email",
		routes.AuthController.GetRegistrationMail,
	)

	routes.Router.POST(string(endpoints.RegisterV1)+"/:registrationToken",
		routes.RequestBodyValidator.Handler("register", RegisterBody{}),
		routes.AuthController.Register)

	routes.Router.POST(string(endpoints.LoginV1),
		routes.RequestBodyValidator.Handler("login", LoginBody{}),
		routes.AuthController.Login)

	routes.Router.POST(string(endpoints.LogoutV1),
		routes.AuthMiddleware.Handler(),
		routes.AuthController.Logout)
}
