package main

import (
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	libfx "github.com/TeaChanathip/touch-grass-scheduler/server/pkg/lib"
	"go.uber.org/fx"
)

func main() {
	// Run each module with FX
	fx.New(
		configfx.Module,
		libfx.Module,
	).Run()
}
