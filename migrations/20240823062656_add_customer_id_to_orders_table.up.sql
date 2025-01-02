DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'orders' AND column_name = 'customer_id') THEN
        ALTER TABLE orders ADD COLUMN customer_id INTEGER;
    END IF;
END $$;
