package main

import (
	"fmt"

	"gin-real-time-talk/config"
	"gin-real-time-talk/internal/app"
	"gin-real-time-talk/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @title Real-Time Talk API
// @version 1.0
// @description API for Real-Time Talk application

// @host localhost:5000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token for authentication. Use format: Bearer {token}

func main() {
	if config.Env.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.New()

	if err := app.Run(); err != nil {
		log.Fatal(fmt.Sprintf("Application failed to start: %v", err))
	}
}
