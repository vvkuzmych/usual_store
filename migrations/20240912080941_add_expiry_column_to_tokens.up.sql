ALTER TABLE tokens
    ADD COLUMN expiry TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP;