package config

import (
	"os"

	"gin-real-time-talk/pkg/logger"

	"github.com/joho/godotenv"
)

type TConfig struct {
	Port        string
	Environment string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
}

var Env *TConfig

func init() {
	log := logger.New()
	err := godotenv.Load()

	if err != nil {
		log.Info("Warning: .env file not found, using environment variables")
	}

	Env = &TConfig{
		Port:        getEnv("PORT", "5000"),
		Environment: getEnv("ENVIRONMENT", "development"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "realtimetalk"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
