package v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) Authenticate(c echo.Context) error {
	var req service.AuthRequest

	if err := c.Bind(&req); err != nil {
		return err
	}

	token, err := h.srv.AuthenticatePlayer(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}
