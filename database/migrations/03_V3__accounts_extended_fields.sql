-- Add extended fields for credit cards and virtual wallets
USE fintrack;

-- First, create the accounts table if it doesn't exist (baseline)
CREATE TABLE IF NOT EXISTS accounts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_accounts_user_id (user_id),
    INDEX idx_accounts_account_type (account_type),
    INDEX idx_accounts_currency (currency),
    INDEX idx_accounts_is_active (is_active),
    INDEX idx_accounts_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add new columns for credit cards and virtual wallets
-- Credit card specific fields
ALTER TABLE accounts 
ADD COLUMN credit_limit DECIMAL(15,2) NULL COMMENT 'Credit limit for credit card accounts',
ADD COLUMN closing_date DATE NULL COMMENT 'Monthly closing date for credit cards',
ADD COLUMN due_date DATE NULL COMMENT 'Payment due date for credit cards';

-- Personal identification field for virtual wallets
ALTER TABLE accounts
ADD COLUMN dni VARCHAR(20) NULL COMMENT 'Document number for virtual wallet identification';

-- Add indexes for performance
CREATE INDEX idx_accounts_credit_limit ON accounts(credit_limit);
CREATE INDEX idx_accounts_closing_date ON accounts(closing_date);
CREATE INDEX idx_accounts_due_date ON accounts(due_date);
CREATE INDEX idx_accounts_dni ON accounts(dni);

-- Update account_type enum to include wallet if not already present
-- Note: MySQL doesn't support direct ALTER for ENUM, so we'll ensure the column accepts the new values
-- The application will validate the account_type values

-- Insert some sample data for testing (optional)
-- These will be credit card and wallet accounts
INSERT IGNORE INTO accounts (
    id, user_id, account_type, name, description, currency, balance, 
    credit_limit, closing_date, due_date, is_active, created_at, updated_at
) VALUES
(
    'sample-credit-001', 'sample-user-001', 'credit', 'Tarjeta de Crédito Visa',
    'Tarjeta de crédito principal', 'ARS', 0.00, 50000.00, '2024-01-25', '2024-02-15',
    TRUE, NOW(), NOW()
),
(
    'sample-wallet-001', 'sample-user-001', 'wallet', 'Billetera Virtual',
    'Cuenta de billetera virtual', 'ARS', 1500.00, NULL, NULL, NULL,
    TRUE, NOW(), NOW()
);

-- Add constraint to ensure credit_limit is only set for credit card accounts
-- This is enforced at application level for flexibility