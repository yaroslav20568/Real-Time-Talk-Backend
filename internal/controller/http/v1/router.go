package v1

import (
	"gin-real-time-talk/internal/controller/http/v1/auth"
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

	api := router.Group("/api/v1")
	{
		auth.SetupAuthRoutes(api, db)
	}

	return router
}
