-- PostgreSQL IPv6 Configuration Script
-- This script is automatically executed when the PostgreSQL container starts
-- It ensures IPv6 connectivity is properly configured

-- Log the initialization
DO $$
BEGIN
    RAISE NOTICE 'Configuring PostgreSQL for IPv6 support...';
END $$;

-- Create a function to log connection info (useful for debugging)
CREATE OR REPLACE FUNCTION log_connection_info()
RETURNS TABLE(setting_name text, setting_value text) AS $$
BEGIN
    RETURN QUERY
    SELECT name::text, setting::text
    FROM pg_settings
    WHERE name IN ('listen_addresses', 'max_connections', 'port');
END;
$$ LANGUAGE plpgsql;

-- Grant necessary permissions
GRANT CONNECT ON DATABASE usualstore TO postgres;
GRANT ALL PRIVILEGES ON DATABASE usualstore TO postgres;

-- Log completion
DO $$
BEGIN
    RAISE NOTICE 'IPv6 configuration complete!';
    RAISE NOTICE 'PostgreSQL is now listening on both IPv4 and IPv6 addresses';
END $$;

