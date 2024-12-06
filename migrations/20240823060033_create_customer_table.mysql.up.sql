CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Alter the orders table to add customer_id column
ALTER TABLE orders ADD COLUMN customer_id INT;

-- Add foreign key constraint to the orders table referencing the customers table
ALTER TABLE orders
    ADD CONSTRAINT fk_customer_id
        FOREIGN KEY (customer_id) REFERENCES customers(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;