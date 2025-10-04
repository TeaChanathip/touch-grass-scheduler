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
	// Run with FX
	fx.New(
		configfx.Module,
		libfx.Module,
		usersfx.Module,
		authfx.Module,
		bootstrapfx.Module,
	).Run()
}
