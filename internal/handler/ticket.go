package handler

import (
	"database-concurrency/internal/handler/payload"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h handler) Book(ctx echo.Context) error {
	var req payload.BookRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	if req.TicketID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "ticket_id is required")
	}
	if req.UserID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "user_id is required")
	}

	result, err := h.ctrl.TicketCtrl().Book(ctx.Request().Context(), req.TicketID, req.UserID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.JSON(200, payload.BookResponse{
		Ticket: result,
	})
}

func (h handler) Reserve(ctx echo.Context) error {
	var req payload.ReserveRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	if req.TicketID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "ticket_id is required")
	}

	if req.UserID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "user_id is required")
	}

	result, err := h.ctrl.TicketCtrl().Reserve(ctx.Request().Context(), req.TicketID, req.UserID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.JSON(200, payload.ReserveResponse{
		Ticket: result,
	})
}

func (h handler) Cancel(ctx echo.Context) error {
	var req payload.CancelRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return echo.NewHTTPError(400, err.Error())
	}

	if req.TicketID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "ticket_id is required")
	}

	if req.UserID.String() == uuid.Nil.String() {
		return echo.NewHTTPError(400, "user_id is required")
	}

	result, err := h.ctrl.TicketCtrl().Cancel(ctx.Request().Context(), req.TicketID, req.UserID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.JSON(200, payload.CancelResponse{
		Ticket: result,
	})
}