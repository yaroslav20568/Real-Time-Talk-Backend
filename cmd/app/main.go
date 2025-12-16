package main

import (
	"fmt"

	"gin-real-time-talk/config"
	"gin-real-time-talk/internal/app"
	"gin-real-time-talk/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	if config.Env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	log := logger.New()

	if err := app.Run(); err != nil {
		log.Fatal(fmt.Sprintf("Application failed to start: %v", err))
	}
}
