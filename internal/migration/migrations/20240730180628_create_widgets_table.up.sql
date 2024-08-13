-- Create the widgets table
CREATE TABLE widgets (
    id INTEGER UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    inventory_level INTEGER,
    price INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Insert a record into the widgets table
INSERT INTO widgets (name, description, inventory_level, price, created_at, updated_at)
VALUES ('Widget', 'A very nice widget.', 10, 1000, NOW(), NOW());
