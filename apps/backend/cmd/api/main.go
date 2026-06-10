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
	"github.com/devrapture/pod-events/internal/routes"
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

	deps := routes.HandlerDependencies{}

	addr := fmt.Sprintf(":%s", cfg.Port)
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
