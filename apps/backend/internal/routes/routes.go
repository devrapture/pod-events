package routes

import (
	"github.com/devrapture/pod-events/internal/config"
	handlers "github.com/devrapture/pod-events/internal/handler"
	"github.com/devrapture/pod-events/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HandlerDependencies struct {
}

func Setup(db *gorm.DB, deps HandlerDependencies, cfg *config.Config, logger *zap.Logger) *gin.Engine {

	r := gin.New()
	r.Use(middleware.RequestLogger(logger))
	r.Use(gin.Recovery())
	v1 := r.Group("/api/v1")

	{
		v1.GET("/health", handlers.HealthHandler(db))
	}

	return r
}
