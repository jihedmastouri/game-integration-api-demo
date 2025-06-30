package main

import (
	"context"
	"errors"
	"log"
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
	"github.com/joho/godotenv"
)

func main() {
	loadDotenv()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := repository.Connect(internal.Config.DATABASE_URL)
	srv := service.NewService(&repo)

	server := transport.Web(internal.Config.APP_URL, srv, logger)

	done := make(chan bool, 1)
	go gracefulShutdown(server, done)

	// Start server
	if err := server.Start(internal.Config.APP_URL); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

	<-done
}

func loadDotenv() {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("Warning: .env file not found")
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading the .env")
		}
	}
}

type ServerWithShutdown interface {
	Shutdown(context.Context) error
}

func gracefulShutdown(apiServer ServerWithShutdown, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // disables the signal NotifyContext and allowing Ctrl+C to force shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}
