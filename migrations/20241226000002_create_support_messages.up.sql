-- Create support_messages table
CREATE TABLE IF NOT EXISTS support_messages (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    sender_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    sender_type VARCHAR(20) NOT NULL, -- user, supporter, system
    sender_name VARCHAR(255),
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_support_messages_ticket_id ON support_messages(ticket_id);
CREATE INDEX idx_support_messages_sender_id ON support_messages(sender_id);
CREATE INDEX idx_support_messages_created_at ON support_messages(created_at);
CREATE INDEX idx_support_messages_is_read ON support_messages(is_read);

