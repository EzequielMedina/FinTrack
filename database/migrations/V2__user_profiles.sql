-- Add profile and tracking fields to users table
USE fintrack;

-- Add new columns for profile data and last login tracking
ALTER TABLE users 
ADD COLUMN profile_data JSON NULL COMMENT 'User profile information as JSON',
ADD COLUMN last_login_at DATETIME NULL COMMENT 'Last login timestamp';

-- Update role values to match the new constants (case-sensitive)
UPDATE users SET role = 'user' WHERE role = 'USER';
UPDATE users SET role = 'admin' WHERE role = 'ADMIN';
UPDATE users SET role = 'operator' WHERE role = 'OPERATOR';
UPDATE users SET role = 'treasurer' WHERE role = 'TREASURER';

-- Add index on last_login_at for performance
CREATE INDEX idx_users_last_login ON users(last_login_at);

-- Add index on role for filtering
CREATE INDEX idx_users_role ON users(role);

-- Add index on is_active for filtering active users
CREATE INDEX idx_users_is_active ON users(is_active);

-- Add composite index for common queries
CREATE INDEX idx_users_active_role ON users(is_active, role);