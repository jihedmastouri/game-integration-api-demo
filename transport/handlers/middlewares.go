package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/jihedmastouri/game-integration-api-demo/internal"
	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/labstack/echo/v4"
)

func AuthMiddlewareFactory(s *service.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, internal.ErrRestResponse{
					Code: internal.Unauthorized,
					Msg:  "Authorization header is empty",
				})
			}

			token = strings.TrimSpace(strings.TrimPrefix(
				strings.TrimSpace(token),
				"Bearer",
			))

			player, err := s.AuthorizePlayer(c.Request().Context(), token)
			if err != nil {
				c.Logger().Errorf("failed to validate token: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if player == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Failed to found player")
			}

			p := *player
			c.Set("player", p)
			return next(c)
		}
	}
}

func ErrorMiddlewareFactory() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				// Maybe send to sentry
				slog.Error(err.Error())
			}
			return err
		}
	}
}
