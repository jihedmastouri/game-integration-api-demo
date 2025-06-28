package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	loadDotenv()

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		getDefaultEnv("PG_USER", "postgres"),
		getDefaultEnv("PG_PASS", "postgres"),
		getDefaultEnv("PG_URL", "localhost"),
		getDefaultEnv("PG_PORT", "5432"),
		getDefaultEnv("PG_DB", "postgres"),
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(databaseUrl)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// pool, err := sql.Open("postgres", databaseUrl)
	// if err != nil {
	// 	log.Fatal("unable to use data source", err)
	// }
	// defer pool.Close()
	//
	// pool.SetMaxIdleConns(3)
	// pool.SetConnMaxLifetime(0)
	// pool.SetMaxOpenConns(20)
	// pingTest(pool)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if v.URI == "/health" {
					return nil
				}
				if v.Error == nil {
					logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.String("remote_ip", v.RemoteIP),
					)
				} else {
					logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.String("remote_ip", v.RemoteIP),
						slog.String("error", v.Error.Error()),
					)
				}
				return nil
			},
		},
	))
	e.Use(middleware.Recover())

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

func getDefaultEnv(name, defaultValue string) string {
	if envValue := os.Getenv(name); envValue != "" {
		return envValue
	}
	return defaultValue
}

func pingTest(pool *sql.DB) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
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
