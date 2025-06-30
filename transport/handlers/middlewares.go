package handlers

import (
	"net/http"
	"strings"

	"github.com/jihedmastouri/game-integration-api-demo/internal"
	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/labstack/echo/v4"
)

// AuthMiddlewareFactory creates an AuthMiddleware with the provided repository
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

			c.Set("player_id", player) // Also set as string for compatibility

			err = next(c)
			return err
		}
	}
}

func ErrorMiddlewareFactory() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
					"code": err.Error(),
				})
			}
			return err
		}
	}
}
