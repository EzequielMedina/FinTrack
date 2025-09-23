-- Migration initialization script
-- This file ensures all migrations are executed in the correct order
-- It will be executed first due to the 00_ prefix

-- Set character set and collation
SET NAMES utf8mb4;
SET character_set_client = utf8mb4;

-- Use the database
USE fintrack;

-- Log migration start
SELECT 'Starting database migrations...' AS message;

-- Create migration tracking table if it doesn't exist
CREATE TABLE IF NOT EXISTS migration_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    migration_file VARCHAR(255) NOT NULL UNIQUE,
    executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('success', 'failed') DEFAULT 'success'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Log that initialization is complete
INSERT INTO migration_history (migration_file) 
VALUES ('00_init_order.sql') 
ON DUPLICATE KEY UPDATE executed_at = CURRENT_TIMESTAMP;

SELECT 'Database initialization complete.' AS message;