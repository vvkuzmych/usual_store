-- Add role column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(50) NOT NULL DEFAULT 'user';

-- Create index for faster role-based queries
CREATE INDEX idx_users_role ON users(role);

-- Update existing admin user
UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';

-- Add comment to explain the role column
COMMENT ON COLUMN users.role IS 'User role: admin, supporter, or user';

