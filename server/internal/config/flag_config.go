package configfx

import (
	"flag"
	"log"
	"slices"
)

type FlagConfig struct {
	Environment string
}

func NewFlagsConfig() *FlagConfig {
	flagConfig := &FlagConfig{}

	flag.StringVar(&flagConfig.Environment, "e", "development", "Environment to run in (shorthand)")
	flag.Parse()

	if !isValidEnvironment(flagConfig.Environment) {
		log.Printf("Warning: Invalid environment '%s', using 'development'", flagConfig.Environment)
		flagConfig.Environment = "development"
	}

	return flagConfig
}

func isValidEnvironment(env string) bool {
	validEnvs := []string{"development", "production", "test"}
	return slices.Contains(validEnvs, env)
}
