-- Create the transaction_statuses table
CREATE TABLE transaction_statuses (
    id  SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default transaction statuses
INSERT INTO transaction_statuses (name) VALUES ('Pending');
INSERT INTO transaction_statuses (name) VALUES ('Cleared');
INSERT INTO transaction_statuses (name) VALUES ('Declined');
INSERT INTO transaction_statuses (name) VALUES ('Refunded');
INSERT INTO transaction_statuses (name) VALUES ('Partially refunded');
