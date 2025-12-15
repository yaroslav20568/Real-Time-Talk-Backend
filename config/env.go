package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type TConfig struct {
	Port        string
	Environment string
}

var Env *TConfig

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	Env = &TConfig{
		Port:        getEnv("PORT", "5000"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
