package libfx

import (
	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerParams struct {
	fx.In
	FlagConfig *configfx.FlagConfig
}

func NewLogger(params LoggerParams) *zap.Logger {
	var cfg zap.Config

	switch params.FlagConfig.Environment {
	case "production":
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
		cfg.OutputPaths = []string{"stdout"}
	case "test":
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		cfg.OutputPaths = []string{"stdout"}
	default: // development and others
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.OutputPaths = []string{"stdout"}
	}

	// Customize the time format to match Go's standard log
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05")

	logger, err := cfg.Build()
	if err != nil {
		// Fallback to a no-op logger if zap fails
		panic("Failed to initialize zap logger: " + err.Error())
	}

	logger.Info("Logger initialized successfully.",
		zap.String("min_level", cfg.Level.String()),
		zap.String("format", cfg.Encoding),
	)

	return logger
}
