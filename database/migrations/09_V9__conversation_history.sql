-- Migration: Create conversation_history table
-- Description: Store chat conversation history for the chatbot
-- Date: 2025-10-27

USE fintrack;

-- Create conversation_history table
CREATE TABLE IF NOT EXISTS conversation_history (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    conversation_id VARCHAR(36) NOT NULL,
    role ENUM('user', 'assistant') NOT NULL,
    message TEXT NOT NULL,
    context_data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_conversation_id (conversation_id),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at),
    INDEX idx_user_conversation (user_id, conversation_id),
    
    CONSTRAINT fk_conversation_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add comment to table
ALTER TABLE conversation_history COMMENT = 'Stores chat conversation history between users and the financial assistant';

-- Sample query to verify
-- SELECT * FROM conversation_history WHERE user_id = 'xxx' AND conversation_id = 'yyy' ORDER BY created_at;
