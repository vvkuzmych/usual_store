-- Create AI user preferences table to learn from user interactions
CREATE TABLE IF NOT EXISTS ai_user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255),
    preferred_categories TEXT[],
    budget_min DECIMAL(10, 2),
    budget_max DECIMAL(10, 2),
    interaction_count INTEGER DEFAULT 0,
    last_products_viewed INTEGER[],
    last_products_purchased INTEGER[],
    conversation_style VARCHAR(50) CHECK (conversation_style IN ('concise', 'detailed', 'friendly', 'professional')),
    preferred_language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_ai_user_preferences_user_id ON ai_user_preferences(user_id);
CREATE INDEX idx_ai_user_preferences_session_id ON ai_user_preferences(session_id);

COMMENT ON TABLE ai_user_preferences IS 'Stores learned preferences from user interactions with AI';
COMMENT ON COLUMN ai_user_preferences.preferred_categories IS 'Array of category names user is interested in';
COMMENT ON COLUMN ai_user_preferences.budget_min IS 'Minimum budget mentioned by user';
COMMENT ON COLUMN ai_user_preferences.budget_max IS 'Maximum budget mentioned by user';
COMMENT ON COLUMN ai_user_preferences.last_products_viewed IS 'Array of product IDs viewed in chat';
COMMENT ON COLUMN ai_user_preferences.conversation_style IS 'How the user prefers to communicate';

