ALTER TABLE transactions ADD COLUMN payment_intent VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE transactions ADD COLUMN payment_method VARCHAR(255) NOT NULL DEFAULT '';