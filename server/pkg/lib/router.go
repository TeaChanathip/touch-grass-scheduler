package libfx

import (
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RouterParam struct {
	fx.In
	AppConfig  *configfx.AppConfig
	FlagConfig *configfx.FlagConfig
	Logger     *zap.Logger
}

func NewRouter(param RouterParam) *gin.Engine {
	// Switch from debug mode to release mode in production
	if param.FlagConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Use Logger as Gin middleware
	router.Use(GinLoggerMiddleware(param.Logger))

	param.Logger.Info("Router initialized successfully")

	return router
}

func RunRouter(router *gin.Engine, appConfig *configfx.AppConfig) {
	router.Run(":" + appConfig.ServerPort)
}

// Custom Gin Middleware for Logger
func GinLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		fields := []zapcore.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", latency),
		}

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info(path, fields...)
		}
	}
}
