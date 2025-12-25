-- Create AI feedback table to track response quality
CREATE TABLE IF NOT EXISTS ai_feedback (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES ai_messages(id) ON DELETE CASCADE,
    conversation_id INTEGER NOT NULL REFERENCES ai_conversations(id) ON DELETE CASCADE,
    helpful BOOLEAN,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    feedback_text TEXT,
    feedback_type VARCHAR(50) CHECK (feedback_type IN ('helpful', 'not_helpful', 'incorrect', 'inappropriate', 'other')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ai_feedback_message_id ON ai_feedback(message_id);
CREATE INDEX idx_ai_feedback_conversation_id ON ai_feedback(conversation_id);
CREATE INDEX idx_ai_feedback_helpful ON ai_feedback(helpful);
CREATE INDEX idx_ai_feedback_rating ON ai_feedback(rating);

COMMENT ON TABLE ai_feedback IS 'User feedback on AI assistant responses';
COMMENT ON COLUMN ai_feedback.helpful IS 'Was this response helpful?';
COMMENT ON COLUMN ai_feedback.rating IS 'User rating from 1-5 stars';
COMMENT ON COLUMN ai_feedback.feedback_text IS 'Optional user comment';

