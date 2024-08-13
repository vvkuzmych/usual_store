-- Create the transactions table
CREATE TABLE transactions (
              id INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY,
              amount INTEGER NOT NULL,
              currency VARCHAR(255) NOT NULL,
              last_four VARCHAR(4) NOT NULL,
              bank_return_code VARCHAR(255),
              transaction_status_id INTEGER UNSIGNED,
              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
              FOREIGN KEY (transaction_status_id) REFERENCES transaction_statuses(id)
                  ON DELETE CASCADE
                  ON UPDATE CASCADE
);

-- Insert example data into transactions
-- You would typically insert data after ensuring your schema is in place
