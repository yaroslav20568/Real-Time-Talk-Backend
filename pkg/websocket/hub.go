package websocket

import (
	"sync"

	"gin-real-time-talk/internal/entity"
)

type Hub struct {
	clients    map[uint]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Message struct {
	Type        string          `json:"type"`
	Data        interface{}     `json:"data"`
	RecipientID uint            `json:"recipientId,omitempty"`
	Message     *entity.Message `json:"message,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			var clientsToRemove []*Client
			if message.RecipientID > 0 {
				if clients, ok := h.clients[message.RecipientID]; ok {
					for client := range clients {
						select {
						case client.send <- message:
						default:
							clientsToRemove = append(clientsToRemove, client)
						}
					}
				}
			} else {
				for _, clients := range h.clients {
					for client := range clients {
						select {
						case client.send <- message:
						default:
							clientsToRemove = append(clientsToRemove, client)
						}
					}
				}
			}
			h.mu.RUnlock()

			if len(clientsToRemove) > 0 {
				h.mu.Lock()
				for _, client := range clientsToRemove {
					if clients, ok := h.clients[client.userID]; ok {
						if _, exists := clients[client]; exists {
							delete(clients, client)
							close(client.send)
							if len(clients) == 0 {
								delete(h.clients, client.userID)
							}
						}
					}
				}
				h.mu.Unlock()
			}
		}
	}
}

func (h *Hub) BroadcastToUser(userID uint, message *Message) {
	message.RecipientID = userID
	h.broadcast <- message
}

func (h *Hub) BroadcastToAll(message *Message) {
	message.RecipientID = 0
	h.broadcast <- message
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}
