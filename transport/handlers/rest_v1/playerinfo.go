package rest_v1

import (
	"net/http"

	"github.com/jihedmastouri/game-integration-api-demo/models"
	"github.com/jihedmastouri/game-integration-api-demo/transport/shared"
	"github.com/labstack/echo/v4"
)

func (h *Handlers) PlayerInfo(c echo.Context) error {
	player, ok := c.Get("player").(models.Player)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid player context")
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
