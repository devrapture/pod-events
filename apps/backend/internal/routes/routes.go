package routes

import (
	"net/http"
	"strings"

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
	AuthHandler      *handlers.AuthHandler
	ShowHandler      *handlers.ShowHandler
	TelegramHandler  *handlers.TelegramWebHookHandler
	ChannelHandler   *handlers.ChannelHandler
	DashboardHandler *handlers.DashboardShowHandler
}

func Setup(db *gorm.DB, deps HandlerDependencies, cfg *config.Config, logger *zap.Logger) *gin.Engine {

	r := gin.New()
	r.Use(middleware.RequestLogger(logger))
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg.FrontendURL))

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

		// dashboard
		dashboard := protected.Group("/dashboard")
		dashboard.
			GET("/summary", deps.DashboardHandler.GetDashboardSummary)

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

func corsMiddleware(allowedOrigin string) gin.HandlerFunc {
	allowedOrigin = strings.TrimRight(strings.TrimSpace(allowedOrigin), "/")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		normalizedOrigin := strings.TrimRight(strings.TrimSpace(origin), "/")
		matched := normalizedOrigin == allowedOrigin

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
