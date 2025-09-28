package libfx

import (
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RouterParam struct {
	fx.In
	Flag   *configfx.FlagConfig
	Logger *zap.Logger
}

func NewRouter(param RouterParam) *gin.Engine {
	// Switch from debug mode to release mode in production
	if param.Flag.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	param.Logger.Info("Router initialized successfully")

	return router
}

func RunRouter(router *gin.Engine) {
	router.Run()
}
