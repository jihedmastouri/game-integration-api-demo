package rest_v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) Authenticate(c echo.Context) error {
	var req service.AuthRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, shared.ErrorResponse{
			Code: shared.ValidationError,
			Msg:  err.Error(),
		})
	}

	token, err := h.srv.AuthenticatePlayer(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.ErrorResponse{
			Code: shared.ServiceUnAvailable,
			Msg:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.AuthResponse{Token: token})
}
