-- Create the statuses table with UNSIGNED id
CREATE TABLE statuses (
    id INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert default statuses
INSERT INTO statuses (name) VALUES ('Cleared');
INSERT INTO statuses (name) VALUES ('Refunded');
INSERT INTO statuses (name) VALUES ('Cancelled');

-- Alter the orders table to add a foreign key constraint
ALTER TABLE orders
    ADD CONSTRAINT fk_status_id
        FOREIGN KEY (status_id) REFERENCES statuses(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;
