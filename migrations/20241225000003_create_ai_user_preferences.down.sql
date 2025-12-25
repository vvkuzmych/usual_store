-- Rollback AI user preferences table
DROP INDEX IF EXISTS idx_ai_user_preferences_session_id;
DROP INDEX IF EXISTS idx_ai_user_preferences_user_id;
DROP TABLE IF EXISTS ai_user_preferences;

