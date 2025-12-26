-- Create support_tickets table
CREATE TABLE IF NOT EXISTS support_tickets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    supporter_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    subject VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open', -- open, assigned, in_progress, resolved, closed
    priority VARCHAR(20) DEFAULT 'medium', -- low, medium, high, urgent
    session_id VARCHAR(255) UNIQUE NOT NULL,
    user_email VARCHAR(255),
    user_name VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_support_tickets_user_id ON support_tickets(user_id);
CREATE INDEX idx_support_tickets_supporter_id ON support_tickets(supporter_id);
CREATE INDEX idx_support_tickets_status ON support_tickets(status);
CREATE INDEX idx_support_tickets_session_id ON support_tickets(session_id);
CREATE INDEX idx_support_tickets_created_at ON support_tickets(created_at);

