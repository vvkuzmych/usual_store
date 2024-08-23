-- Step 1: Remove the foreign key constraint
ALTER TABLE orders
DROP FOREIGN KEY fk_customer_id;

-- Step 2: Drop the column
ALTER TABLE orders
DROP COLUMN customer_id;