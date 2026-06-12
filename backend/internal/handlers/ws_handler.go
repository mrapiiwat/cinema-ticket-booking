package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/mrapiiwat/cinema-ticket-booking/backend/internal/services"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WSHandler struct {
	hub *services.WSHub
}

func NewWSHandler(hub *services.WSHub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

func (h *WSHandler) ServeWS(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println("WS Upgrade error:", err)
		return err
	}

	client := &services.Client{Hub: h.hub, Conn: conn, Send: make(chan []byte, 256)}
	h.hub.Register(client)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()

	return nil
}
