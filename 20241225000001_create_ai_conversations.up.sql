-- Create AI conversations table to track chat sessions
CREATE TABLE IF NOT EXISTS ai_conversations (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    total_messages INTEGER DEFAULT 0,
    resulted_in_purchase BOOLEAN DEFAULT FALSE,
    total_tokens_used INTEGER DEFAULT 0,
    total_cost DECIMAL(10, 6) DEFAULT 0.00,
    user_agent TEXT,
    ip_address VARCHAR(45),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ai_conversations_session_id ON ai_conversations(session_id);
CREATE INDEX idx_ai_conversations_user_id ON ai_conversations(user_id);
CREATE INDEX idx_ai_conversations_started_at ON ai_conversations(started_at);
CREATE INDEX idx_ai_conversations_resulted_in_purchase ON ai_conversations(resulted_in_purchase);

COMMENT ON TABLE ai_conversations IS 'Tracks AI assistant chat sessions';
COMMENT ON COLUMN ai_conversations.session_id IS 'Unique identifier for the chat session';
COMMENT ON COLUMN ai_conversations.user_id IS 'User ID if authenticated, null for anonymous';
COMMENT ON COLUMN ai_conversations.total_tokens_used IS 'Total OpenAI tokens consumed';
COMMENT ON COLUMN ai_conversations.total_cost IS 'Estimated cost in USD';

