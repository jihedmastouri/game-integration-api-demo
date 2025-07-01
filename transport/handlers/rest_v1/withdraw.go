package rest_v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
	"github.com/labstack/echo/v4"
)

// Bet godoc
// @Summary Process a bet
// @Description Processes a withdrawal from a player's balance. Each request represents a bet placement action.
// @Tags Betting
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param request body shared.WithdrawRequest true "Bet request details"
// @Success 200 {object} shared.BetOperationResponse "Bet processed successfully"
// @Failure 400 {object} shared.ErrorResponse "Bad request"
// @Failure 401 {object} shared.ErrorResponse "Unauthorized"
// @Failure 500 {object} shared.ErrorResponse "Internal server error"
// @Router /api/v1/withdraw [post]
// @Security BearerAuth
func (h *Handlers) Withdraw(c echo.Context) error {
	// Get player from auth middleware
	player, ok := c.Get("player").(models.Player)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid player context")
	}

	// Bind request
	var req shared.WithdrawRequest
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

	// Process bet through service
	betResponse, err := h.srv.ProcessBet(c.Request().Context(), &player, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, betResponse)
}
