package libfx

import (
	"log"
	"time"

	configfx "github.com/TeaChanathip/touch-grass-scheduler/server/internal/config"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseParams struct {
	fx.In
	AppConfig *configfx.AppConfig
}

func NewDatabase(param DatabaseParams) *gorm.DB {
	dsn := param.AppConfig.GetDBConfig()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to DB: %+v", err)
	}

	// Get underlying SQL DB to check connection and get stats
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Warning: Could not get underlying SQL DB: %+v", err)
	} else {
		// Configure connection pool
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Test the connection
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Database ping failed: %+v", err)
		}

		// Log connection statistics
		stats := sqlDB.Stats()
		log.Printf("✅ Database connected successfully!")
		log.Printf("Connection Stats:")
		log.Printf("  - Open Connections: %d", stats.OpenConnections)
		log.Printf("  - Max Open Connections: %d", stats.MaxOpenConnections)
		log.Printf("  - In Use: %d", stats.InUse)
		log.Printf("  - Idle: %d", stats.Idle)
	}

	return db
}
