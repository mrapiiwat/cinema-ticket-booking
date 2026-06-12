package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/database"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/registry"
)

func main() {
	e := echo.New()

	database.ConnectDB()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fmt.Printf("[ECHO] %s | %d | %13v | %s | %s\n",
				v.StartTime.Format("2006-01-02 15:04:05"),
				v.Status,
				v.Latency,
				v.Method,
				v.URI,
			)
			return nil
		},
	}))

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins(),
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	container := registry.NewContainer(database.MongoDB)
	if err := database.SeedDefaults(e.Logger); err != nil {
		log.Fatalf("Failed to seed defaults: %v", err)
	}
	registry.SetupRoutes(e, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Echo Server is starting on port :%s \n", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to run Echo server: %v", err)
	}
}

func corsOrigins() []string {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		return []string{"http://localhost", "http://localhost:5173", "http://127.0.0.1:5173"}
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}
