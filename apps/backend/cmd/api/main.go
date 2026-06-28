package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devrapture/pod-events/internal/config"
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

	db, err := database.ConnectDb(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// ── Cache ────────────────────────────────────────────────
	appCache := cache.New(10*time.Minute, 15*time.Minute)

	// ── Clients ────────────────────────────────────────────────
	spotifyClient := spotify.NewSpotifyClient(cfg, logger)

	// ── Notifier ────────────────────────────────────────────────
	telegramNotifier := telegram.NewNotifier(cfg)

	// ── Repositories ────────────────────────────────────────────────
	userRepo := repositories.NewUserRepository(db)
	tokenRepo := repositories.NewTokenRepository(db, cfg.TokenEncryptionKey)
	channelRepo := repositories.NewChannelRepository(db)

	// ── Services ────────────────────────────────────────────────
	authService := services.NewAuthService(cfg, tokenRepo, userRepo, spotifyClient, logger)
	showService := services.NewShowServices(spotifyClient, authService, cfg, appCache)
	channelService := services.NewChannelServices(channelRepo)

	// ── Handlers ────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(authService, logger, cfg)
	showHandler := handlers.NewShowHandler(showService, logger)
	telegramHandler := handlers.NewTelegramWebHookHandler(cfg, telegramNotifier, logger)
	channelHandler := handlers.NewChannelHandler(channelService, logger)

	deps := routes.HandlerDependencies{
		AuthHandler:     authHandler,
		ShowHandler:     showHandler,
		TelegramHandler: telegramHandler,
		ChannelHandler:  channelHandler,
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
