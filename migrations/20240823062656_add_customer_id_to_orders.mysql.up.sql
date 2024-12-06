
ALTER TABLE orders ADD COLUMN customer_id INTEGER;

ALTER TABLE orders
    ADD CONSTRAINT fk_customer_id
        FOREIGN KEY (customer_id) REFERENCES customers(id)
            ON DELETE CASCADE
            ON UPDATE CASCADE;