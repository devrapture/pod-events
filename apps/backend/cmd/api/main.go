// @title           PodEvents API
// @version         1.0
// @description     API for PodEvents - podcast notification platform
// @termsOfService  https://podevents.app/terms
//
// @contact.name   API Support
// @contact.email  devrapture@proton.me
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
//
// @tag.name Auth
// @tag.description Authentication endpoints
// @tag.name Shows
// @tag.description Podcast show operations
// @tag.name Subscriptions
// @tag.description Subscription management
// @tag.name Channels
// @tag.description Notification channel management
// @tag.name Telegram
// @tag.description Telegram integration
// @tag.name Health
// @tag.description Service health check
// @tag.name Dashboard
// @tag.description Dashboard overview
// @tag.name Cron
// @tag.description Cron job endpoints (protected by secret header)
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/devrapture/pod-events/internal/config"
	"github.com/devrapture/pod-events/internal/cron"
	"github.com/devrapture/pod-events/internal/database"
	handlers "github.com/devrapture/pod-events/internal/handler"
	"github.com/devrapture/pod-events/internal/notifications/telegram"
	"github.com/devrapture/pod-events/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"github.com/devrapture/pod-events/internal/routes"
	"github.com/devrapture/pod-events/internal/services"
	"github.com/devrapture/pod-events/internal/spotify"
	"github.com/devrapture/pod-events/pkg/logger"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration %v", err)
	}

	logger, err := logger.NewLogger(!cfg.IsProduction())
	if err != nil {
		log.Fatalf("Failed to create logger %v", err)
	}
	defer logger.Sync()
	logger.Info("Starting PodEvents server", zap.String("env", cfg.AppEnv), zap.String("port", cfg.Port))
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:                   cfg.SentryDSN,
		Environment:           cfg.AppEnv,
		Release:               cfg.SentryRelease,
		AttachStacktrace:      true,
		EnableTracing:         cfg.SentryTracesSampleRate > 0,
		TracesSampleRate:      cfg.SentryTracesSampleRate,
		DisableLogs:           !cfg.SentryEnableLogs,
		BeforeSend:            scrubSentryEvent,
		BeforeSendTransaction: scrubSentryEvent,
	}); err != nil {
		log.Fatalf("Failed to initialize Sentry: %v", err)
	}
	if cfg.SentryDSN != "" {
		defer func() {
			if !sentry.Flush(2 * time.Second) {
				logger.Warn("Timed out flushing Sentry events")
			}
		}()
		logger.Info(
			"Sentry initialized",
			zap.Float64("traces_sample_rate", cfg.SentryTracesSampleRate),
			zap.Bool("logs_enabled", cfg.SentryEnableLogs),
		)
	}

	db, err := database.ConnectDb(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// ── Cache ────────────────────────────────────────────────
	appCache := cache.New(10*time.Minute, 15*time.Minute)

	// ── Clients ────────────────────────────────────────────────
	spotifyClient := spotify.NewSpotifyClient(cfg, logger)

	// ── Notifier ────────────────────────────────────────────────
	telegramNotifier := telegram.NewNotifier(cfg, 0, logger)

	// ── Repositories ────────────────────────────────────────────────
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db, cfg.TokenEncryptionKey)
	channelRepo := repositories.NewChannelRepository(db, cfg.TokenEncryptionKey)
	telegramConnectionRepo := repositories.NewTelegramConnectionRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	showRepository := repositories.NewShowRepository(db)
	dashboardSummaryRepo := repositories.NewDashboardSummaryRepository(db, tokenRepo)
	episodeRepo := repositories.NewEpisodeRepository(db)
	showRepo := repositories.NewShowRepository(db)
	notificationLogRepo := repositories.NewNotificationLogRepository(db)

	// ── Services ────────────────────────────────────────────────
	authService := services.NewAuthService(cfg, tokenRepo, userRepo, spotifyClient, appCache, logger)
	showService := services.NewShowServices(spotifyClient, authService, appCache, subscriptionRepo, showRepository)
	channelService := services.NewChannelServices(channelRepo)
	telegramConnectionService := services.NewTelegramConnectionService(telegramConnectionRepo, channelRepo, cfg)
	dashboardService := services.NewDashboardSummaryService(dashboardSummaryRepo)
	notifService := services.NewNotificationService(notificationLogRepo, logger, cfg, channelRepo)

	// ── Cron ────────────────────────────────────────────────
	episodeChecker := cron.NewEpisodeChecker(subscriptionRepo, episodeRepo, showRepo, authService, notifService, logger, spotifyClient)

	// ── Handlers ────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authService, logger, cfg, userRepo)
	showHandler := handlers.NewShowHandler(showService, logger)
	telegramHandler := handlers.NewTelegramWebHookHandler(cfg, telegramNotifier, telegramConnectionService, logger)
	channelHandler := handlers.NewChannelHandler(channelService, logger)
	dashboardHandler := handlers.NewDashboardShowHandler(dashboardService, logger)
	cronHandler := handlers.NewCronJobHandler(logger, episodeChecker)

	deps := routes.HandlerDependencies{
		AuthHandler:      authHandler,
		ShowHandler:      showHandler,
		TelegramHandler:  telegramHandler,
		ChannelHandler:   channelHandler,
		DashboardHandler: dashboardHandler,
		CronHandler:      cronHandler,
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}
	r := routes.Setup(db, deps, cfg, logger)

	srv := &http.Server{
		Addr:    addr,
		Handler: r.Handler(),
	}

	go func() {
		logger.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			logger.Fatal("Failed to close database", zap.Error(err))
		}
	}

	logger.Info("Server exited")
}

func scrubSentryEvent(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
	if event.Request == nil {
		return event
	}

	event.Request.Cookies = ""
	for header := range event.Request.Headers {
		switch strings.ToLower(header) {
		case "authorization", "cookie", "x-cron-secret", "x-telegram-bot-api-secret-token":
			delete(event.Request.Headers, header)
		}
	}
	return event
}
