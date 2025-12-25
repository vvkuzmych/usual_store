-- Rollback AI conversations table
DROP INDEX IF EXISTS idx_ai_conversations_resulted_in_purchase;
DROP INDEX IF EXISTS idx_ai_conversations_started_at;
DROP INDEX IF EXISTS idx_ai_conversations_user_id;
DROP INDEX IF EXISTS idx_ai_conversations_session_id;
DROP TABLE IF EXISTS ai_conversations;

