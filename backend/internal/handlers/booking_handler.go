package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type BookingHandler struct {
	bookingService ports.BookingServicePort
}

func NewBookingHandler(bookingService ports.BookingServicePort) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
	}
}

type LockSeatRequest struct {
	ShowtimeID string   `json:"showtime_id"`
	SeatNo     string   `json:"seat_no"`
	Seats      []string `json:"seats"`
}

func (h *BookingHandler) LockSeat(c echo.Context) error {
	var req LockSeatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	userID := c.Get("user_id").(string)
	seats := req.Seats
	if req.SeatNo != "" {
		seats = append(seats, req.SeatNo)
	}

	booking, err := h.bookingService.LockSeats(c.Request().Context(), req.ShowtimeID, seats, userID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "seats locked successfully", "booking": booking})
}

type ConfirmBookingRequest struct {
	BookingID string `json:"booking_id"`
}

func (h *BookingHandler) ConfirmBooking(c echo.Context) error {
	var req ConfirmBookingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	userID := c.Get("user_id").(string)
	booking, err := h.bookingService.ConfirmBooking(c.Request().Context(), req.BookingID, userID)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "booking confirmed successfully", "booking": booking})
}

func (h *BookingHandler) GetUserBookings(c echo.Context) error {
	userID := c.Get("user_id").(string)

	bookings, err := h.bookingService.GetUserBookings(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bookings)
}
