package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	handlers "github.com/devrapture/pod-events/internal/handler"
	"github.com/devrapture/pod-events/internal/metrics"
	"github.com/devrapture/pod-events/internal/middleware"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"

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
	CronHandler      *handlers.CronJobHandler
}

func Setup(db *gorm.DB, deps HandlerDependencies, cfg *config.Config, logger *zap.Logger, recorder metrics.Recorder) *gin.Engine {
	r := gin.New()
	ipStore := middleware.NewRateLimiterStore(rate.Limit(5), 10)
	r.Use(middleware.IPRateLimiter(ipStore))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.MetricsRecorder(recorder))
	r.Use(gin.Recovery())
	if cfg.SentryDSN != "" {
		r.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
			Timeout: 2 * time.Second,
		}))
	}
	r.Use(corsMiddleware(cfg.FrontendURL))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	{
		v1.GET("/health", handlers.HealthHandler(db))

		// Cron endpoint — protected by secret header, NOT user auth
		// cron-job.org is a machine caller, not a user
		cronGroup := r.Group("/cron")
		cronGroup.Use(middleware.CronMiddleware(cfg.CronSecret, logger))

		cronGroup.POST("/check-episodes", deps.CronHandler.CheckEpisodes)

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
			POST("/subscribe", deps.ShowHandler.Subscribe)

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
