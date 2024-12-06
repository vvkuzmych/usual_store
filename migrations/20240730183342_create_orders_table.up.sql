-- Create the orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    widget_id INTEGER,
    transaction_id INTEGER,
    status_id INTEGER,
    quantity INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (widget_id) REFERENCES widgets(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (status_id) REFERENCES transaction_statuses(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
