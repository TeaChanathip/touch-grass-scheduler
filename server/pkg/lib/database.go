package libfx

import (
	"fmt"
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

func NewDatabase(params DatabaseParams) (*gorm.DB, error) {
	dsn := params.AppConfig.GetDBConfig()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database session: %w", err)
	}

	// Get underlying SQL DB to check connection and get stats
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying database: %w", err)
	} else {
		// Configure connection pool
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Test the connection
		if err := sqlDB.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping database: %w", err)
		}

		// Log connection statistics
		stats := sqlDB.Stats()
		params.Logger.Info("Database connection succeeded",
			zap.Int("open_connections", stats.OpenConnections),
			zap.Int("in_use", stats.InUse),
			zap.Int("idle", stats.Idle))
	}

	return db, nil
}
