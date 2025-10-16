package main

import (
	bootstrapfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/bootstrap"
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	middlewarefx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/middlewares"
	authfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/auth"
	libfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/lib"
	usersfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/users"
	"go.uber.org/fx"
)

func main() {
	// Register every module and run with FX
	fx.New(
		// Prerequisite
		configfx.Module,
		libfx.Module,

		// Routes
		usersfx.Module,
		authfx.Module,

		// Middlewares
		middlewarefx.Module,

		// Forcing router to correctly initialize
		bootstrapfx.Module,
	).Run()
}
