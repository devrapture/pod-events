package routes

import (
	"net/http"

	"github.com/devrapture/pod-events/internal/config"
	handlers "github.com/devrapture/pod-events/internal/handler"
	"github.com/devrapture/pod-events/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/devrapture/pod-events/docs"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HandlerDependencies struct {
	AuthHandler     *handlers.AuthHandler
	ShowHandler     *handlers.ShowHandler
	TelegramHandler *handlers.TelegramWebHookHandler
	ChannelHandler  *handlers.ChannelHandler
}

func Setup(db *gorm.DB, deps HandlerDependencies, cfg *config.Config, logger *zap.Logger) *gin.Engine {

	r := gin.New()
	r.Use(middleware.RequestLogger(logger))
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg.FrontendURL, logger))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	{
		v1.GET("/health", handlers.HealthHandler(db))

		// auth
		auth := v1.Group("/auth")

		auth.
			GET("/spotify/login", deps.AuthHandler.SpotifyLogin).
			GET("/spotify/callback", deps.AuthHandler.SpotifyCallback).
			POST("/exchange", deps.AuthHandler.ExchangeAuthCode)

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		protected.GET("/auth/me", deps.AuthHandler.Me)

		// shows
		shows := protected.Group("/shows")
		shows.
			GET("/saved", deps.ShowHandler.GetUserSavedShows).
			GET("/search", deps.ShowHandler.SearchShows).
			POST("/:spotifyShowId/subscribe", deps.ShowHandler.Subscribe)

		// subscriptions
		subscriptions := protected.Group("/subscriptions")
		subscriptions.
			GET("", deps.ShowHandler.GetSubscriptions).
			DELETE("/:id", deps.ShowHandler.Unsubscribe)

		// webhooks
		webhooks := v1.Group("/webhooks")
		webhooks.POST("/telegram", deps.TelegramHandler.Handle)

		// Telegram
		telegram := protected.Group("/telegram")
		telegram.
			POST("/generate-link", deps.TelegramHandler.CreateConnectLink)

		// channels
		channels := protected.Group("/channels")
		channels.
			POST("", deps.ChannelHandler.CreateChannel).
			GET("", deps.ChannelHandler.GetChannels).
			POST("/:channelID/toggle", deps.ChannelHandler.ToggleActive).
			DELETE("/:channelID", deps.ChannelHandler.Delete)

	}

	return r
}

func corsMiddleware(allowedOrigin string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		matched := origin == allowedOrigin

		logger.Debug("CORS check",
			zap.String("origin", origin),
			zap.String("allowed_origin", allowedOrigin),
			zap.Bool("matched", matched),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		if matched {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET,POST,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type,ngrok-skip-browser-warning")
			c.Header("Access-Control-Max-Age", "86400")
			c.Header("Vary", "Origin")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
