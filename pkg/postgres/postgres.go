package postgres

import (
	"fmt"

	"gin-real-time-talk/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func New() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Env.DBHost,
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBName,
		config.Env.DBPort,
		config.Env.DBSSLMode,
	)

	var logLevel logger.LogLevel
	if config.Env.Environment == "production" {
		logLevel = logger.Silent
	} else {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	return db, nil
}

func GetDB() *gorm.DB {
	return DB
}
