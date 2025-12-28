-- Tenant Database Schema
-- This schema is applied to each tenant database

-- Customers table
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Widgets table (products)
CREATE TABLE IF NOT EXISTS widgets (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    inventory_level INT NOT NULL DEFAULT 0,
    price INT NOT NULL,
    image VARCHAR(255),
    is_recurring BOOLEAN DEFAULT FALSE,
    plan_id VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Transaction statuses
CREATE TABLE IF NOT EXISTS transaction_statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    amount INT NOT NULL,
    currency VARCHAR(255) NOT NULL,
    last_four VARCHAR(4) NOT NULL,
    bank_return_code VARCHAR(255) NOT NULL DEFAULT '',
    expiry_month INT NOT NULL DEFAULT 0,
    expiry_year INT NOT NULL DEFAULT 0,
    payment_intent VARCHAR(255) NOT NULL,
    payment_method VARCHAR(255) NOT NULL,
    transaction_status_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (transaction_status_id) REFERENCES transaction_statuses(id) ON DELETE CASCADE
);

-- Statuses table
CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    widget_id INT NOT NULL,
    transaction_id INT NOT NULL,
    customer_id INT NOT NULL,
    status_id INT NOT NULL,
    quantity INT NOT NULL,
    amount INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (widget_id) REFERENCES widgets(id) ON DELETE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (status_id) REFERENCES statuses(id) ON DELETE CASCADE
);

-- Insert default statuses
INSERT INTO statuses (name, created_at, updated_at) 
VALUES 
    ('Cleared', NOW(), NOW()),
    ('Refunded', NOW(), NOW()),
    ('Cancelled', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Insert default transaction statuses
INSERT INTO transaction_statuses (name, created_at, updated_at)
VALUES
    ('Pending', NOW(), NOW()),
    ('Cleared', NOW(), NOW()),
    ('Declined', NOW(), NOW()),
    ('Refunded', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON orders(customer_id);
CREATE INDEX IF NOT EXISTS idx_orders_widget_id ON orders(widget_id);
CREATE INDEX IF NOT EXISTS idx_transactions_payment_intent ON transactions(payment_intent);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add triggers for updated_at
DROP TRIGGER IF EXISTS update_customers_updated_at ON customers;
CREATE TRIGGER update_customers_updated_at 
    BEFORE UPDATE ON customers
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_widgets_updated_at ON widgets;
CREATE TRIGGER update_widgets_updated_at 
    BEFORE UPDATE ON widgets
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_orders_updated_at ON orders;
CREATE TRIGGER update_orders_updated_at 
    BEFORE UPDATE ON orders
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
CREATE TRIGGER update_transactions_updated_at 
    BEFORE UPDATE ON transactions
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Grant permissions on sequences (needed for SERIAL columns)
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO PUBLIC;

-- Set default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public 
GRANT SELECT, INSERT, UPDATE ON TABLES TO PUBLIC;

ALTER DEFAULT PRIVILEGES IN SCHEMA public 
GRANT USAGE, SELECT ON SEQUENCES TO PUBLIC;

