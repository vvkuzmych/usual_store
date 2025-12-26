package support

import (
	"database/sql"
	"time"
)

// Ticket represents a support ticket
type Ticket struct {
	ID          int        `json:"id"`
	UserID      *int       `json:"user_id,omitempty"`
	SupporterID *int       `json:"supporter_id,omitempty"`
	Subject     string     `json:"subject"`
	Status      string     `json:"status"`   // open, assigned, in_progress, resolved, closed
	Priority    string     `json:"priority"` // low, medium, high, urgent
	SessionID   string     `json:"session_id"`
	UserEmail   string     `json:"user_email"`
	UserName    string     `json:"user_name"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
}

// Message represents a support message
type Message struct {
	ID         int       `json:"id"`
	TicketID   int       `json:"ticket_id"`
	SenderID   *int      `json:"sender_id,omitempty"`
	SenderType string    `json:"sender_type"` // user, supporter, system
	SenderName string    `json:"sender_name"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// Session represents an active support session
type Session struct {
	ID                 int        `json:"id"`
	TicketID           int        `json:"ticket_id"`
	UserID             *int       `json:"user_id,omitempty"`
	SupporterID        *int       `json:"supporter_id,omitempty"`
	IsActive           bool       `json:"is_active"`
	StartedAt          time.Time  `json:"started_at"`
	EndedAt            *time.Time `json:"ended_at,omitempty"`
	UserConnected      bool       `json:"user_connected"`
	SupporterConnected bool       `json:"supporter_connected"`
}

// CreateTicket creates a new support ticket
func CreateTicket(db *sql.DB, ticket *Ticket) error {
	query := `
		INSERT INTO support_tickets (user_id, subject, status, priority, session_id, user_email, user_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	err := db.QueryRow(
		query,
		ticket.UserID,
		ticket.Subject,
		ticket.Status,
		ticket.Priority,
		ticket.SessionID,
		ticket.UserEmail,
		ticket.UserName,
		time.Now(),
		time.Now(),
	).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.UpdatedAt)
	return err
}

// GetTicketBySessionID retrieves a ticket by session ID
func GetTicketBySessionID(db *sql.DB, sessionID string) (*Ticket, error) {
	query := `
		SELECT id, user_id, supporter_id, subject, status, priority, session_id, user_email, user_name, created_at, updated_at, closed_at
		FROM support_tickets
		WHERE session_id = $1
	`
	ticket := &Ticket{}
	err := db.QueryRow(query, sessionID).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.SupporterID,
		&ticket.Subject,
		&ticket.Status,
		&ticket.Priority,
		&ticket.SessionID,
		&ticket.UserEmail,
		&ticket.UserName,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.ClosedAt,
	)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// GetTicketByID retrieves a ticket by ID
func GetTicketByID(db *sql.DB, ticketID int) (*Ticket, error) {
	query := `
		SELECT id, user_id, supporter_id, subject, status, priority, session_id, user_email, user_name, created_at, updated_at, closed_at
		FROM support_tickets
		WHERE id = $1
	`
	ticket := &Ticket{}
	err := db.QueryRow(query, ticketID).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.SupporterID,
		&ticket.Subject,
		&ticket.Status,
		&ticket.Priority,
		&ticket.SessionID,
		&ticket.UserEmail,
		&ticket.UserName,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
		&ticket.ClosedAt,
	)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

// UpdateTicketStatus updates the status of a ticket
func UpdateTicketStatus(db *sql.DB, ticketID int, status string, supporterID *int) error {
	query := `
		UPDATE support_tickets
		SET status = $1, supporter_id = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := db.Exec(query, status, supporterID, time.Now(), ticketID)
	return err
}

// GetAllOpenTickets retrieves all open tickets
func GetAllOpenTickets(db *sql.DB) ([]*Ticket, error) {
	query := `
		SELECT id, user_id, supporter_id, subject, status, priority, session_id, user_email, user_name, created_at, updated_at, closed_at
		FROM support_tickets
		WHERE status IN ('open', 'assigned', 'in_progress')
		ORDER BY priority DESC, created_at ASC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*Ticket
	for rows.Next() {
		ticket := &Ticket{}
		err := rows.Scan(
			&ticket.ID,
			&ticket.UserID,
			&ticket.SupporterID,
			&ticket.Subject,
			&ticket.Status,
			&ticket.Priority,
			&ticket.SessionID,
			&ticket.UserEmail,
			&ticket.UserName,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
			&ticket.ClosedAt,
		)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

// CreateMessage creates a new support message
func CreateMessage(db *sql.DB, message *Message) error {
	query := `
		INSERT INTO support_messages (ticket_id, sender_id, sender_type, sender_name, message, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`
	err := db.QueryRow(
		query,
		message.TicketID,
		message.SenderID,
		message.SenderType,
		message.SenderName,
		message.Message,
		message.IsRead,
		time.Now(),
	).Scan(&message.ID, &message.CreatedAt)
	return err
}

// GetMessagesByTicketID retrieves all messages for a ticket
func GetMessagesByTicketID(db *sql.DB, ticketID int) ([]*Message, error) {
	query := `
		SELECT id, ticket_id, sender_id, sender_type, sender_name, message, is_read, created_at
		FROM support_messages
		WHERE ticket_id = $1
		ORDER BY created_at ASC
	`
	rows, err := db.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.TicketID,
			&message.SenderID,
			&message.SenderType,
			&message.SenderName,
			&message.Message,
			&message.IsRead,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// MarkMessagesAsRead marks all messages in a ticket as read
func MarkMessagesAsRead(db *sql.DB, ticketID int, senderType string) error {
	query := `
		UPDATE support_messages
		SET is_read = TRUE
		WHERE ticket_id = $1 AND sender_type != $2 AND is_read = FALSE
	`
	_, err := db.Exec(query, ticketID, senderType)
	return err
}
