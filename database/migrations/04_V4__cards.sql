-- Migration: 04_V4__cards.sql
-- Description: Add cards table to support multiple cards per account

-- Create cards table
CREATE TABLE IF NOT EXISTS cards (
    id VARCHAR(36) PRIMARY KEY,
    account_id VARCHAR(36) NOT NULL,
    card_type VARCHAR(10) NOT NULL,
    card_brand VARCHAR(20) NOT NULL,
    last_four_digits VARCHAR(4) NOT NULL,
    masked_number VARCHAR(19) NOT NULL,
    holder_name VARCHAR(100) NOT NULL,
    expiration_month INT NOT NULL,
    expiration_year INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    nickname VARCHAR(50),
    
    -- Credit card specific fields
    credit_limit DECIMAL(15,2) NULL,
    closing_date DATE NULL,
    due_date DATE NULL,
    
    -- Security fields (encrypted data stored separately)
    encrypted_number TEXT NOT NULL,
    key_fingerprint VARCHAR(64) NOT NULL,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_cards_account_id (account_id),
    INDEX idx_cards_status (status),
    INDEX idx_cards_deleted_at (deleted_at),
    INDEX idx_cards_is_default (is_default),
    INDEX idx_cards_card_type (card_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Now add the foreign key constraint separately
ALTER TABLE cards 
ADD CONSTRAINT fk_cards_account_id 
    FOREIGN KEY (account_id) 
    REFERENCES accounts(id) 
    ON DELETE CASCADE;

-- Add new account type for bank accounts
ALTER TABLE accounts 
MODIFY COLUMN account_type VARCHAR(20) NOT NULL;

-- Add comment to document the new account types
ALTER TABLE accounts 
COMMENT = 'Financial accounts table. Supports wallet (virtual) and bank_account (with cards) types';

ALTER TABLE cards 
COMMENT = 'Payment cards associated with bank accounts. Each card belongs to one account.';

-- Insert some sample data for testing (optional - remove in production)
-- This shows the new relationship structure
INSERT INTO accounts (id, user_id, account_type, name, description, currency, balance, dni, is_active, created_at, updated_at) 
VALUES 
(UUID(), 'user-123', 'wallet', 'Mi Billetera Virtual', 'Billetera para pagos digitales', 'ARS', 5000.00, '12345678', TRUE, NOW(), NOW()),
(UUID(), 'user-123', 'bank_account', 'Cuenta Principal', 'Cuenta bancaria con tarjetas', 'ARS', 15000.00, NULL, TRUE, NOW(), NOW());

-- Note: Cards will be added through the application API, not via SQL