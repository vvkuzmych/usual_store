-- Rollback AI feedback table
DROP INDEX IF EXISTS idx_ai_feedback_rating;
DROP INDEX IF EXISTS idx_ai_feedback_helpful;
DROP INDEX IF EXISTS idx_ai_feedback_conversation_id;
DROP INDEX IF EXISTS idx_ai_feedback_message_id;
DROP TABLE IF EXISTS ai_feedback;

