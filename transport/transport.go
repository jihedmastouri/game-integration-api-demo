package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/jihedmastouri/game-integration-api-demo/transport/handlers"

	_ "github.com/jihedmastouri/game-integration-api-demo/docs"
)

func Web(address string, srv *service.Service, logger *slog.Logger) *echo.Echo {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			HandleError: true,
			LogRemoteIP: true,
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
	e.Validator = &CustomValidation{validator: validator.New()}
	e.Use(handlers.ErrorMiddlewareFactory())

	handlers.SetupRoutes(e, srv)

	return e
}

type CustomValidation struct {
	validator *validator.Validate
}

func (cv *CustomValidation) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fmt.Errorf("Request validation failed: %v", validationErrors[0].Field())
		}
		return fmt.Errorf("Request validation failed: %v", err)
	}
	return nil
}
