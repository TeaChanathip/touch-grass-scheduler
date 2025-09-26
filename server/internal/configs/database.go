package configs

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDatabase() *gorm.DB {
	dsn := AppConfig.GetDBConfig()

	// Configure GORM logger for SQL query logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Include params in SQL log
			Colorful:                  true,        // Enable color
		},
	)

	// Log connection attempt
	log.Printf("Attempting to connect to database...")
	log.Printf("Database Host: %s", AppConfig.DBHost)
	log.Printf("Database Name: %s", AppConfig.DBName)
	log.Printf("Database User: %s", AppConfig.DBUser)
	log.Printf("Database Port: %d", AppConfig.DBPort)
	log.Printf("SSL Mode: %s", AppConfig.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // Use custom logger
	})

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
