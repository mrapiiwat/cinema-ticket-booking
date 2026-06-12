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
		"token":   user.Token,
		"user":    user.User,
	})
}

func (h *AuthHandler) GoogleConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, h.authService.GoogleConfig())
}

func (h *AuthHandler) GoogleCodeLogin(c echo.Context) error {
	var req struct {
		Code        string `json:"code"`
		RedirectURI string `json:"redirect_uri"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	auth, err := h.authService.ExchangeGoogleCode(c.Request().Context(), req.Code, req.RedirectURI)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"token":   auth.Token,
		"user":    auth.User,
	})
}

func (h *AuthHandler) DevLogin(c echo.Context) error {
	var req struct {
		Role string `json:"role"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	auth, err := h.authService.DevLogin(c.Request().Context(), req.Role)
	if err != nil {
		if err.Error() == "dev login is disabled" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Dev login successful",
		"token":   auth.Token,
		"user":    auth.User,
	})
}
