package main

import (
	"fmt"
	"log"

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
	logger.Info("Server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server %v", err)
	}
}
