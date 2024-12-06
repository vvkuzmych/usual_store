-- Add the `is_recurring` column with a boolean type and default value as false
ALTER TABLE widgets ADD COLUMN is_recurring BOOL NOT NULL DEFAULT false;

-- Add the `plan_id` column with a string type and default value as an empty string
ALTER TABLE widgets ADD COLUMN plan_id VARCHAR(255) NOT NULL DEFAULT '';
