-- Rollback AI messages table
DROP INDEX IF EXISTS idx_ai_messages_metadata;
DROP INDEX IF EXISTS idx_ai_messages_role;
DROP INDEX IF EXISTS idx_ai_messages_created_at;
DROP INDEX IF EXISTS idx_ai_messages_conversation_id;
DROP TABLE IF EXISTS ai_messages;

