package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/ports"
)

type AuthHandler struct {
	authService ports.AuthServicePort
}

func NewAuthHandler(authService ports.AuthServicePort) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	user, err := h.authService.VerifyGoogleToken(c.Request().Context(), req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}


	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user":    user,
	})
}