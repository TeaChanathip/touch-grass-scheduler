package libfx

import (
	"reflect"
	"strings"
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

func NewRouter(params RouterParam) *gin.Engine {
	// Set mode accroding to environment
	switch params.FlagConfig.Environment {
	case "production":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	}

	router := gin.New()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Specify allowed origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma", "Referer", "Referrer-Policy"}
	config.ExposeHeaders = []string{"Content-Length"} // Headers exposed to the client
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour // Cache preflight requests for 12 hours

	// Apply middleware
	router.Use(cors.New(config))
	router.Use(gin.Recovery())
	router.Use(ginLoggerMiddleware(params.Logger))

	// Use field name specified for JSON in validation
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	params.Logger.Info("Router initialized successfully.")

	return router
}

// Custom Gin Middleware for Logger
func ginLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
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
