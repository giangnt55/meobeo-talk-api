package main

import (
	"context"
	"log"
	"meobeo-talk-api/internal/config"
	"meobeo-talk-api/internal/interfaces/http/routes"
	"meobeo-talk-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.App.Env)
	defer appLogger.Sync()

	// Initialize database
	db, err := config.NewDatabase(cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect database", "error", err)
	}
	defer db.Close()

	// Migrate database
	if err := config.RunMigrations(db, cfg); err != nil {
		appLogger.Fatal("Failed to run migrations", "error", err)
	}

	// Setup router with all dependencies
	router := routes.SetupRouter(cfg, db, appLogger)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		appLogger.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Fatal("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exited")
}
