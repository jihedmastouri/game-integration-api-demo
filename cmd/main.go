package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jihedmastouri/game-integration-api-demo/internal"
	"github.com/jihedmastouri/game-integration-api-demo/repository"
	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/jihedmastouri/game-integration-api-demo/transport"

	_ "github.com/jihedmastouri/game-integration-api-demo/repository/migrations"
)

func main() {
	logLevel := slog.LevelInfo
	if internal.Config.MODE == internal.ModeDevelopment {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}))
	slog.SetDefault(logger)

	repo, err := repository.Connect(internal.Config.DATABASE_URL)
	if err != nil {
		slog.Error("Failed to connect to db", "error", err)
		os.Exit(1) // Exit if database connection fails
	}

	srv := service.NewService(repo)

	// Start pending transaction worker
	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go srv.StartPendingTransactionWorker(workerCtx)

	server := transport.Web(internal.Config.APP_URL, srv, logger)

	done := make(chan struct{})
	go gracefulShutdown(server, workerCancel, done)

	// Start server
	if err := server.Start(internal.Config.APP_URL); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	<-done
	slog.Info("Application shutdown complete")
}

type ServerWithShutdown interface {
	Shutdown(context.Context) error
}

func gracefulShutdown(apiServer ServerWithShutdown, workerCancel context.CancelFunc, done chan struct{}) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	slog.Info("Shutdown signal received, shutting down gracefully. Press Ctrl+C again to force shutdown.")

	// Cancel the worker context first
	workerCancel()

	// Stop listening for new signals to allow force shutdown
	stop()

	// Create a timeout context for the shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown with error", "error", err)
	} else {
		slog.Info("Server shutdown completed successfully")
	}

	// Notify the main goroutine that the shutdown is complete
	done <- struct{}{}
}
