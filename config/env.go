package config

import (
	"os"

	"gin-real-time-talk/pkg/logger"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port        string
	Environment string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWConfig struct {
	Secret        string
	AccessExpiry  string
	RefreshExpiry string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type Config struct {
	App  AppConfig
	DB   DBConfig
	JWT  JWConfig
	SMTP SMTPConfig
}

var Env *Config

func init() {
	log := logger.New()
	err := godotenv.Load()

	if err != nil {
		log.Info("Warning: .env file not found, using environment variables")
	}

	Env = &Config{
		App: AppConfig{
			Port:        getEnv("PORT", "5000"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "realtimetalk"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWConfig{
			Secret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
			RefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "7d"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
