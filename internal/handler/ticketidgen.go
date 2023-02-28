package handler

import (
	"database-concurrency/internal/handler/payload"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) GenTicketID(ctx echo.Context) error {
	var req payload.GenTicketIDRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	if req.UniqueID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "unique_id is required")
	}

	result, err := h.ctrl.GenCtrl().CreateTicketID(ctx.Request().Context(), req.UniqueID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.JSON(200, payload.GenTicketIDResponse{
		TicketID: result,
	})
}
