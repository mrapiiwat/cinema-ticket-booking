package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type ShowtimeHandler struct {
	showtimeService ports.ShowtimeServicePort
}

func NewShowtimeHandler(showtimeService ports.ShowtimeServicePort) *ShowtimeHandler {
	return &ShowtimeHandler{
		showtimeService: showtimeService,
	}
}

func (h *ShowtimeHandler) GetShowtimesByMovie(c echo.Context) error {
	movieID := c.Param("movieID")
	showtimes, err := h.showtimeService.GetShowtimesByMovie(c.Request().Context(), movieID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, showtimes)
}

func (h *ShowtimeHandler) GetShowtimeByID(c echo.Context) error {
	id := c.Param("id")
	showtime, err := h.showtimeService.GetShowtimeByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if showtime == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Showtime not found"})
	}
	return c.JSON(http.StatusOK, showtime)
}

func (h *ShowtimeHandler) CreateShowtime(c echo.Context) error {
	var showtime model.Showtime
	if err := c.Bind(&showtime); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.showtimeService.CreateShowtime(c.Request().Context(), &showtime); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, showtime)
}
