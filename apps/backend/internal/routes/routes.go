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
	AuthHandler     *handlers.AuthHandler
	ShowHandler     *handlers.ShowHandler
	TelegramHandler *handlers.TelegramWebHookHandler
}

func Setup(db *gorm.DB, deps HandlerDependencies, cfg *config.Config, logger *zap.Logger) *gin.Engine {

	r := gin.New()
	r.Use(middleware.RequestLogger(logger))
	r.Use(gin.Recovery())
	v1 := r.Group("/api/v1")

	{
		v1.GET("/health", handlers.HealthHandler(db))

		// auth
		auth := v1.Group("/auth")

		auth.
			GET("/spotify/login", deps.AuthHandler.SpotifyLogin).
			GET("/spotify/callback", deps.AuthHandler.SpotifyCallback)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))

		// shows
		shows := protected.Group("/shows")
		shows.
			GET("/saved", deps.ShowHandler.GetUserSavedShows).
			GET("/search", deps.ShowHandler.SearchShows)

		// webhooks
		webhooks := v1.Group("/webhooks")
		webhooks.POST("/telegram", deps.TelegramHandler.Handle)

		// Telegram
		telegram := protected.Group("/telegram")
		telegram.
			POST("/send-test-message", deps.TelegramHandler.SendTestMessage)

	}

	return r
}
