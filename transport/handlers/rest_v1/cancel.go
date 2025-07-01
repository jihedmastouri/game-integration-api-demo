package rest_v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
	"github.com/labstack/echo/v4"
)

// Cancel godoc
// @Summary Cancel a transaction
// @Description Reverts a previously processed transaction by reversing its financial impact
// @Tags Betting
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param request body shared.CancelRequest true "Cancel request details"
// @Success 200 {object} shared.BetOperationResponse "Transaction cancelled successfully"
// @Failure 400 {object} shared.ErrorResponse "Bad request"
// @Failure 401 {object} shared.ErrorResponse "Unauthorized"
// @Failure 404 {object} shared.ErrorResponse "Transaction not found"
// @Failure 500 {object} shared.ErrorResponse "Internal server error"
// @Router /api/v1/cancel [post]
// @Security BearerAuth
func (h *Handlers) Cancel(c echo.Context) error {
	// Get player from auth middleware
	player, ok := c.Get("player").(models.Player)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid player context")
	}

	// Bind request
	var req shared.CancelRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
		})
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	// Process cancel through service
	cancelResponse, err := h.srv.ProcessCancel(c.Request().Context(), &player, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cancelResponse)
}
