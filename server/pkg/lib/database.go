package libfx

import (
	"os"
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseParams struct {
	fx.In
	AppConfig *configfx.AppConfig
	Logger    *zap.Logger
}

func NewDatabase(param DatabaseParams) *gorm.DB {
	logger := param.Logger
	dsn := param.AppConfig.GetDBConfig()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Error("Error connecting to DB", zap.Error(err))
		os.Exit(1)
	}

	// Get underlying SQL DB to check connection and get stats
	sqlDB, err := db.DB()
	if err != nil {
		logger.Warn("Could not get underlying SQL DB: %+v", zap.Error(err))
	} else {
		// Configure connection pool
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Test the connection
		if err := sqlDB.Ping(); err != nil {
			logger.Error("Database ping failed", zap.Error(err))
		}

		// Log connection statistics
		stats := sqlDB.Stats()
		logger.Info("✅ Database connected successfully",
			zap.Int("Open Connections", stats.OpenConnections),
			zap.Int("In Use", stats.InUse),
			zap.Int("Idle", stats.Idle))
	}

	return db
}
