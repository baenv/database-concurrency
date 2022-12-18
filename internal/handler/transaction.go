package handler

import (
	"context"

	"github.com/labstack/echo/v4"
)

func (h handler) Transaction(ctx echo.Context) error {
	hash := ctx.Param("hash")
	if len(hash) == 0 {
		return echo.NewHTTPError(400, "hash is required")
	}

	tx, err := h.ctrl.TransactionCtrl().OneByHash(context.Background(), hash)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.JSON(200, tx)
}
