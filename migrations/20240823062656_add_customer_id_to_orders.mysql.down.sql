-- Remove the foreign key constraint
ALTER TABLE orders
DROP FOREIGN KEY fk_customer_id;

-- Remove the column
ALTER TABLE orders
DROP COLUMN customer_id;