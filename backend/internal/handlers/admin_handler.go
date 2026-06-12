package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type AdminHandler struct {
	bookingService ports.BookingServicePort
	auditService   ports.AuditServicePort
}

func NewAdminHandler(bookingService ports.BookingServicePort, auditService ports.AuditServicePort) *AdminHandler {
	return &AdminHandler{
		bookingService: bookingService,
		auditService:   auditService,
	}
}

func (h *AdminHandler) GetBookings(c echo.Context) error {
	filter := model.BookingFilter{
		UserID:     c.QueryParam("user_id"),
		ShowtimeID: c.QueryParam("showtime_id"),
		Status:     strings.ToUpper(c.QueryParam("status")),
	}

	bookings, err := h.bookingService.GetAllBookings(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bookings)
}

func (h *AdminHandler) GetAuditLogs(c echo.Context) error {
	logs, err := h.auditService.GetLogs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, logs)
}
