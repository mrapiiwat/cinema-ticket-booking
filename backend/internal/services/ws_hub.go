package services

import (
	"encoding/json"
	"log"
	"sync"
)

type WSMessage struct {
	Type       string      `json:"type"`
	ShowtimeID string      `json:"showtime_id"`
	Payload    interface{} `json:"payload"`
}

type WSHub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewWSHub() *WSHub {
	return &WSHub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *WSHub) BroadcastMessage(msg WSMessage) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling WS message:", err)
		return
	}
	h.broadcast <- bytes
}

func (h *WSHub) Register(c *Client) {
	h.register <- c
}

func (h *WSHub) Unregister(c *Client) {
	h.unregister <- c
}
