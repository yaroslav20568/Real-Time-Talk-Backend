package v1

import (
	"gin-real-time-talk/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewRouter(db *gorm.DB, logger *logger.Logger) *gin.Engine {
	_ = &Router{
		db:     db,
		logger: logger,
	}

	router := gin.Default()

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Test",
		})
	})

	return router
}
