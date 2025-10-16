package bootstrapfx

import (
	"context"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// =============== Private Methods ===============
func registerHooks(
	lc fx.Lifecycle,
	appConfig *configfx.AppConfig,
	router *gin.Engine,
	logger *zap.Logger,
	routes Routes,
) {
	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				// ======== SET UP COMPONENTS ========
				// Perform any necessary setup or initialization tasks for the routes
				routes.Setup()

				// Start the router by running it in a separate goroutine
				// If it don't run in goroutine, it will lock the Fx startup process indefinitely
				go router.Run(":" + appConfig.AppPort)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				// Log the stop of the application and any associated error
				logger.Fatal("Stopping application. Error:", zap.Error(ctx.Err()))

				return nil
			},
		},
	)
}

// =============== Exports ===============
var Module = fx.Module(
	"bootstrapfx",
	fx.Provide(NewRoutes),
	fx.Invoke(registerHooks),
)
