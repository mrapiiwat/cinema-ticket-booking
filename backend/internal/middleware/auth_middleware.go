package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/model"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/security"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized: missing bearer token"})
		}

		claims, err := security.ParseToken(strings.TrimPrefix(authHeader, "Bearer "))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get("role").(string)
		if role != model.RoleAdmin {
			return c.JSON(http.StatusForbidden, echo.Map{"error": "Forbidden: requires ADMIN role"})
		}

		return next(c)
	}
}
