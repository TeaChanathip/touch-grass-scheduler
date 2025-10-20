package configfx

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type AppConfigParams struct {
	fx.In
	Flag *FlagConfig
}

type AppConfig struct {
	// Database
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"db"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBSSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`

	// Client
	ClientURL  string `env:"CLIENT_URL" envDefault:"http://localhost"`
	ClientPort int    `env:"CLIENT_PORT" envDefault:"3000"`

	// Server
	AppDomain string `env:"APP_DOMAIN" envDefault:"localhost"`
	AppPort   int    `env:"APP_PORT" envDefault:"8080"`

	// JWT
	JWTSecret    string `env:"JWT_SECRET,required"`
	JWTExpiresIn int    `env:"JWT_EXPIRES_IN" envDefault:"24"`

	// Mail Service
	MailHost     string `env:"MAIL_HOST,required"`
	MailPort     int    `env:"MAIL_PORT,required"`
	MailUser     string `env:"MAIL_USER,required"`
	MailPassword string `env:"MAIL_PASSWORD,required"`
}

func NewAppConfig(params AppConfigParams) *AppConfig {
	config := &AppConfig{}

	// Load ENV by the set environment (relative to server root directory)
	if err := godotenv.Load("env/.env." + params.Flag.Environment); err != nil {
		log.Fatalf("Error loading .env.%s file: %+v", params.Flag.Environment, err)
	}

	/*
		Parse ENV to Config
		Fallback to default values if not defined
	*/
	if err := env.Parse(config); err != nil {
		log.Fatalf("Error parsing environment variables: %+v", err)
	}

	return config
}

func (c *AppConfig) GetDBConfig() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode)
}
