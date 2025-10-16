package middlewarefx

import "go.uber.org/fx"

var Module = fx.Module(
	"middlewarefx",
	fx.Provide(
		NewAuthMiddleware,
		NewRequestBodyValidator,
	),
)
