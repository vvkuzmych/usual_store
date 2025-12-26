package support

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 8192
)

// Client represents a WebSocket client
type Client struct {
	Hub        *Hub
	Conn       *websocket.Conn
	Send       chan []byte
	SessionID  string
	ClientType string // "user" or "supporter"
	UserID     *int
	UserName   string
	DB         *sql.DB
}

// IncomingMessage represents a message received from the client
type IncomingMessage struct {
	Type    string                 `json:"type"` // "message", "typing", "read"
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("Error setting read deadline: %v", err)
		return
	}
	c.Conn.SetPongHandler(func(string) error {
		return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse incoming message
		var incomingMsg IncomingMessage
		if err := json.Unmarshal(messageBytes, &incomingMsg); err != nil {
			log.Printf("Error unmarshaling incoming message: %v", err)
			continue
		}

		// Handle different message types
		switch incomingMsg.Type {
		case "message":
			c.handleIncomingMessage(incomingMsg.Message)
		case "typing":
			c.handleTypingIndicator()
		case "read":
			c.handleReadReceipt()
		default:
			log.Printf("Unknown message type: %s", incomingMsg.Type)
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error setting write deadline: %v", err)
				return
			}
			if !ok {
				// The hub closed the channel
				if err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					log.Printf("Error writing close message: %v", err)
				}
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				log.Printf("Error writing message: %v", err)
			}

			// Add queued messages to the current WebSocket message
			n := len(c.Send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					log.Printf("Error writing newline: %v", err)
				}
				msg := <-c.Send
				if _, err := w.Write(msg); err != nil {
					log.Printf("Error writing queued message: %v", err)
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				log.Printf("Error setting write deadline for ping: %v", err)
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleIncomingMessage processes incoming chat messages
func (c *Client) handleIncomingMessage(messageText string) {
	// Get the ticket
	ticket, err := GetTicketBySessionID(c.DB, c.SessionID)
	if err != nil {
		log.Printf("Error getting ticket: %v", err)
		return
	}

	// Create the message in the database
	message := &Message{
		TicketID:   ticket.ID,
		SenderID:   c.UserID,
		SenderType: c.ClientType,
		SenderName: c.UserName,
		Message:    messageText,
		IsRead:     false,
	}

	if err := CreateMessage(c.DB, message); err != nil {
		log.Printf("Error creating message: %v", err)
		return
	}

	// Broadcast the message to all clients in the session
	c.Hub.BroadcastToSession(c.SessionID, message, "message")
}

// handleTypingIndicator processes typing indicators
func (c *Client) handleTypingIndicator() {
	// Create a system message for typing indicator
	message := &Message{
		SenderType: c.ClientType,
		SenderName: c.UserName,
		Message:    "is typing...",
	}

	// Broadcast typing indicator (doesn't save to DB)
	c.Hub.BroadcastToSession(c.SessionID, message, "typing")
}

// handleReadReceipt marks messages as read
func (c *Client) handleReadReceipt() {
	// Get the ticket
	ticket, err := GetTicketBySessionID(c.DB, c.SessionID)
	if err != nil {
		log.Printf("Error getting ticket: %v", err)
		return
	}

	// Mark messages as read
	if err := MarkMessagesAsRead(c.DB, ticket.ID, c.ClientType); err != nil {
		log.Printf("Error marking messages as read: %v", err)
	}
}
