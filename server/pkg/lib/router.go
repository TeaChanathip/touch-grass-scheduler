package libfx

import (
	"fmt"
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
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
	cyan   = "\033[36m"
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
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Authorization",
		"Accept",
		"User-Agent",
		"Cache-Control",
		"Pragma",
		"Referer",
		"Referrer-Policy",
	}
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

	params.Logger.Info("Router initialization succeeded")

	return router
}

// Custom Gin Middleware for Logger
func ginLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)
		statusCode := ctx.Writer.Status()
		method := ctx.Request.Method

		// Determine status code color
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = green
		case statusCode >= 300 && statusCode < 400:
			statusColor = cyan
		case statusCode >= 400 && statusCode < 500:
			statusColor = yellow
		default: // 5xx or other
			statusColor = red
		}

		// Construct the log message with ANSI colors for console output
		logMessage := fmt.Sprintf("%s |%s %3d %s| %13v | %15s |%s %-7s %s %s %s",
			end.Format("2006/01/02 - 15:04:05"), // Timestamp
			statusColor, statusCode, reset,      // Status code with color
			latency,
			ctx.ClientIP(),
			blue, method, reset, // Method with color
			path,
			query,
		)

		if len(ctx.Errors) > 0 {
			// Log error with the formatted message
			for _, e := range ctx.Errors.Errors() {
				logger.Error(logMessage + " " + e)
			}
		} else {
			// Log info with the formatted message
			logger.Info(logMessage)
		}
	}
}
