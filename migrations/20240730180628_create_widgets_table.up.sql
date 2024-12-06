-- Create the widgets table
CREATE TABLE widgets (
     id              SERIAL PRIMARY KEY,
     name            VARCHAR(255) NOT NULL,
     description     TEXT,
     inventory_level INTEGER,
     price           INTEGER,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
