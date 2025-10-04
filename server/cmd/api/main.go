package main

import (
	bootstrapfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/bootstrap"
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
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

		// All routes
		usersfx.Module,
		authfx.Module,

		// Forcing router to correctly initialize
		bootstrapfx.Module,
	).Run()
}
