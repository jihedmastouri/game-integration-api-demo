package rest_v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
	"github.com/labstack/echo/v4"
)

// PlayerInfo godoc
// @Summary Get player information
// @Description Retrieves essential player details including user ID, balance, and currency
// @Tags Player
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Success 200 {object} shared.PlayerInfoResponse "Player information"
// @Failure 401 {object} shared.ErrorResponse "Unauthorized"
// @Failure 500 {object} shared.ErrorResponse "Internal server error"
// @Router /api/v1/player-info [get]
// @Security BearerAuth
func (h *Handlers) PlayerInfo(c echo.Context) error {
	player, ok := c.Get("player").(models.Player)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, shared.ErrorResponse{
			Code: shared.Unauthorized,
			Msg:  "player not found",
		})
	}

	// Get player info from service
	walletInfo, err := h.srv.WalletClient.GetBalance(player.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, shared.ErrorResponse{
			Code: shared.ServiceUnAvailable,
			Msg:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, shared.PlayerInfoResponse{
		PlayerID: player.ID,
		Balance:  walletInfo.Balance,
		Currency: walletInfo.Currency,
	})
}
