package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anacardo89/page-visitor-counter/config"
	"github.com/Anacardo89/page-visitor-counter/internal/api"
	"github.com/Anacardo89/page-visitor-counter/internal/middleware"
	"github.com/Anacardo89/page-visitor-counter/internal/repo"
	"github.com/Anacardo89/page-visitor-counter/internal/server"
)

func main() {
	// Setup
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	visitorRepo := repo.NewVisitorRepo()
	visitorAPI := api.NewVisitorAPI(visitorRepo)
	mw := middleware.NewMiddlewareHandler(cfg.WriteTimeout)
	srv := server.NewServer(cfg, visitorAPI, mw)

	// Serve
	errChan := make(chan error, 1)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		errChan <- srv.Start()
	}()
	select {
	case sig := <-stopChan:
		slog.Info("Shutting down...", "signal", sig)
		if err := srv.Shutdown(); err != nil {
			slog.Error("Failed to shutdown server gracefully", "error", err)
			os.Exit(1)
		}
		slog.Info("Server stopped gracefully")
	case err := <-errChan:
		if err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}
}
