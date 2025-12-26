package support

import (
	"encoding/json"
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients mapped by session ID
	clients map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan *BroadcastMessage

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	SessionID string   `json:"session_id"`
	Message   *Message `json:"message"`
	Type      string   `json:"type"` // "message", "user_joined", "supporter_joined", "typing", "system"
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.SessionID] == nil {
				h.clients[client.SessionID] = make(map[*Client]bool)
			}
			h.clients[client.SessionID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered: SessionID=%s, Type=%s, Total clients in session=%d",
				client.SessionID, client.ClientType, len(h.clients[client.SessionID]))

		case client := <-h.Unregister:
			h.mu.Lock()
			if sessions, ok := h.clients[client.SessionID]; ok {
				if _, ok := sessions[client]; ok {
					delete(sessions, client)
					close(client.Send)
					log.Printf("Client unregistered: SessionID=%s, Type=%s, Remaining clients in session=%d",
						client.SessionID, client.ClientType, len(sessions))

					if len(sessions) == 0 {
						delete(h.clients, client.SessionID)
						log.Printf("No more clients in session %s, cleaned up", client.SessionID)
					}
				}
			}
			h.mu.Unlock()

		case broadcastMsg := <-h.broadcast:
			h.mu.RLock()
			sessionClients := h.clients[broadcastMsg.SessionID]
			h.mu.RUnlock()

			if sessionClients != nil {
				// Create the message payload
				payload := map[string]interface{}{
					"type":    broadcastMsg.Type,
					"message": broadcastMsg.Message,
				}

				messageBytes, err := json.Marshal(payload)
				if err != nil {
					log.Printf("Error marshaling broadcast message: %v", err)
					continue
				}

				// Broadcast to all clients in the session
				for client := range sessionClients {
					select {
					case client.Send <- messageBytes:
						log.Printf("Message broadcasted to client: SessionID=%s, Type=%s, MessageType=%s",
							client.SessionID, client.ClientType, broadcastMsg.Type)
					default:
						// Client's send channel is full, close and unregister
						h.mu.Lock()
						close(client.Send)
						delete(sessionClients, client)
						if len(sessionClients) == 0 {
							delete(h.clients, client.SessionID)
						}
						h.mu.Unlock()
						log.Printf("Client send buffer full, unregistered: SessionID=%s, Type=%s",
							client.SessionID, client.ClientType)
					}
				}
			}
		}
	}
}

// BroadcastToSession sends a message to all clients in a session
func (h *Hub) BroadcastToSession(sessionID string, message *Message, msgType string) {
	h.broadcast <- &BroadcastMessage{
		SessionID: sessionID,
		Message:   message,
		Type:      msgType,
	}
}

// GetSessionClientCount returns the number of connected clients in a session
func (h *Hub) GetSessionClientCount(sessionID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[sessionID])
}

// GetAllSessions returns all active session IDs
func (h *Hub) GetAllSessions() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	sessions := make([]string, 0, len(h.clients))
	for sessionID := range h.clients {
		sessions = append(sessions, sessionID)
	}
	return sessions
}
