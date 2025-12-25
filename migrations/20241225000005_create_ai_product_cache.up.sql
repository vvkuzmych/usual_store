-- Create AI product cache for faster product context retrieval
CREATE TABLE IF NOT EXISTS ai_product_cache (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES widgets(id) ON DELETE CASCADE,
    description_text TEXT NOT NULL,
    search_keywords TEXT[],
    category VARCHAR(100),
    price_tier VARCHAR(20) CHECK (price_tier IN ('budget', 'mid', 'premium')),
    popularity_score INTEGER DEFAULT 0,
    last_mentioned_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id)
);

CREATE INDEX idx_ai_product_cache_product_id ON ai_product_cache(product_id);
CREATE INDEX idx_ai_product_cache_search_keywords ON ai_product_cache USING GIN (search_keywords);
CREATE INDEX idx_ai_product_cache_popularity_score ON ai_product_cache(popularity_score DESC);

COMMENT ON TABLE ai_product_cache IS 'Cached product information optimized for AI retrieval';
COMMENT ON COLUMN ai_product_cache.description_text IS 'Product description formatted for AI context';
COMMENT ON COLUMN ai_product_cache.search_keywords IS 'Keywords for semantic search';
COMMENT ON COLUMN ai_product_cache.popularity_score IS 'How often product is mentioned in chats';

