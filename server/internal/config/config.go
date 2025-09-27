package configfx

import "go.uber.org/fx"

var Module = fx.Module(
	"configfx",
	fx.Provide(
		NewFlagsConfig,
		NewAppConfig,
	),
)
