-- Create support_sessions table to track active WebSocket connections
CREATE TABLE IF NOT EXISTS support_sessions (
    id SERIAL PRIMARY KEY,
    ticket_id INTEGER NOT NULL REFERENCES support_tickets(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    supporter_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMP,
    user_connected BOOLEAN DEFAULT FALSE,
    supporter_connected BOOLEAN DEFAULT FALSE
);

-- Create indexes
CREATE INDEX idx_support_sessions_ticket_id ON support_sessions(ticket_id);
CREATE INDEX idx_support_sessions_user_id ON support_sessions(user_id);
CREATE INDEX idx_support_sessions_supporter_id ON support_sessions(supporter_id);
CREATE INDEX idx_support_sessions_is_active ON support_sessions(is_active);

