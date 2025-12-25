-- Create AI messages table to store individual chat messages
CREATE TABLE IF NOT EXISTS ai_messages (
    id SERIAL PRIMARY KEY,
    conversation_id INTEGER NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content TEXT NOT NULL,
    tokens_used INTEGER DEFAULT 0,
    response_time_ms INTEGER,
    model VARCHAR(50) DEFAULT 'gpt-3.5-turbo',
    temperature DECIMAL(3, 2),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ai_messages_conversation_id ON ai_messages(conversation_id);
CREATE INDEX idx_ai_messages_created_at ON ai_messages(created_at);
CREATE INDEX idx_ai_messages_role ON ai_messages(role);
CREATE INDEX idx_ai_messages_metadata ON ai_messages USING GIN (metadata);

COMMENT ON TABLE ai_messages IS 'Stores individual messages in AI conversations';
COMMENT ON COLUMN ai_messages.role IS 'Message sender: user, assistant, or system';
COMMENT ON COLUMN ai_messages.content IS 'The actual message text';
COMMENT ON COLUMN ai_messages.tokens_used IS 'OpenAI tokens consumed by this message';
COMMENT ON COLUMN ai_messages.response_time_ms IS 'API response time in milliseconds';
COMMENT ON COLUMN ai_messages.metadata IS 'Additional data like product IDs mentioned, intents, etc.';

