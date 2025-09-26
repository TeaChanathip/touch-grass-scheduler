package configs

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"db"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBSSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

var AppConfig *Config

func LoadConfig(environment string) *Config {
	config := &Config{}

	// Load ENV by the set environment (relative to server root directory)
	if err := godotenv.Load("env/.env." + environment); err != nil {
		log.Fatalf("Error loading .env.%s file: %+v", environment, err)
	}

	// Parse ENV to Config
	if err := env.Parse(config); err != nil {
		log.Fatalf("Error parsing environment variables: %+v", err)
	}

	// Set global config
	AppConfig = config

	return config
}

func (c *Config) GetDBConfig() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode)
}
