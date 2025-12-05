package bootstrapfx

import (
	"context"
	"net/http"
	"strconv"

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
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(appConfig.AppPort),
		Handler: router,
	}

	lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				// ======== SET UP COMPONENTS ========
				// Perform any necessary setup or initialization tasks for the routes
				routes.Setup()

				// Start the server by running it in a separate goroutine
				// If it don't run in goroutine, it will lock the Fx startup process indefinitely
				go func() {
					err := srv.ListenAndServe()
					if err != nil && err != http.ErrServerClosed {
						logger.Error("Failed to start server", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				logger.Info("Shutting down server...")

				// Gracefully shutdown the server
				if err := srv.Shutdown(ctx); err != nil {
					logger.Error("Failed to stop server", zap.Error(err))
					return err
				}

				logger.Info("Sever exiting")

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
