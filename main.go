package main

import (
	"gin-real-time-talk/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if config.Env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Test",
		})
	})
	router.Run(":" + config.Env.Port)
}
