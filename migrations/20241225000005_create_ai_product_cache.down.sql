-- Rollback AI product cache table
DROP INDEX IF EXISTS idx_ai_product_cache_popularity_score;
DROP INDEX IF EXISTS idx_ai_product_cache_search_keywords;
DROP INDEX IF EXISTS idx_ai_product_cache_product_id;
DROP TABLE IF EXISTS ai_product_cache;

