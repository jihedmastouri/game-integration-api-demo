package rest_v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/service"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
	"github.com/labstack/echo/v4"
)

// Authenticate godoc
// @Summary Authenticate player
// @Description Authenticates a player using username and password, returns a JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body service.AuthRequest true "Authentication credentials"
// @Success 200 {object} map[string]string "Authentication successful"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} shared.ErrorResponse "Unauthorized"
// @Failure 500 {object} shared.ErrorResponse "Internal server error"
// @Router /api/v1/auth [post]
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
