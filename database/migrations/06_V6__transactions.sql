-- =====================================================
-- FinTrack Transaction Service - Database Migration
-- Version: V6__transactions.sql
-- Description: Create transactions table and related indexes
-- =====================================================

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    -- Core identity
    id VARCHAR(36) PRIMARY KEY,
    reference_id VARCHAR(100),
    external_id VARCHAR(100),
    
    -- Transaction details
    type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    
    -- Source and destination
    from_account_id VARCHAR(36),
    to_account_id VARCHAR(36),
    from_card_id VARCHAR(36),
    to_card_id VARCHAR(36),
    
    -- User information
    user_id VARCHAR(36) NOT NULL,
    initiated_by VARCHAR(36) NOT NULL,
    
    -- Transaction metadata
    description TEXT,
    payment_method VARCHAR(30),
    merchant_name VARCHAR(255),
    merchant_id VARCHAR(100),
    
    -- Balance tracking
    previous_balance DECIMAL(15,2) DEFAULT 0.00,
    new_balance DECIMAL(15,2) DEFAULT 0.00,
    
    -- Processing details
    processed_at TIMESTAMP NULL,
    failed_at TIMESTAMP NULL,
    failure_reason TEXT,
    
    -- Additional metadata (JSON)
    metadata JSON,
    tags JSON,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Foreign key constraints (will be added when other services are integrated)
    -- FOREIGN KEY (from_account_id) REFERENCES accounts(id),
    -- FOREIGN KEY (to_account_id) REFERENCES accounts(id),
    -- FOREIGN KEY (from_card_id) REFERENCES cards(id),
    -- FOREIGN KEY (to_card_id) REFERENCES cards(id),
    -- FOREIGN KEY (user_id) REFERENCES users(id),
    -- FOREIGN KEY (initiated_by) REFERENCES users(id),
    
    -- Constraints
    CONSTRAINT chk_amount_positive CHECK (amount > 0),
    CONSTRAINT chk_valid_type CHECK (type IN (
        'wallet_deposit', 'wallet_withdrawal', 'wallet_transfer',
        'credit_charge', 'credit_payment', 'credit_refund',
        'debit_purchase', 'debit_withdrawal', 'debit_refund',
        'account_transfer', 'account_deposit', 'account_withdraw'
    )),
    CONSTRAINT chk_valid_status CHECK (status IN (
        'pending', 'completed', 'failed', 'canceled', 'reversed'
    )),
    CONSTRAINT chk_valid_payment_method CHECK (
        payment_method IS NULL OR payment_method IN (
            'cash', 'bank_transfer', 'credit_card', 'debit_card', 'wallet'
        )
    ),
    CONSTRAINT chk_valid_currency CHECK (currency IN ('ARS', 'USD', 'EUR'))
);

-- Create indexes for performance optimization

-- Primary query indexes
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);

-- Account and card relationship indexes
CREATE INDEX idx_transactions_from_account_id ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account_id ON transactions(to_account_id);
CREATE INDEX idx_transactions_from_card_id ON transactions(from_card_id);
CREATE INDEX idx_transactions_to_card_id ON transactions(to_card_id);

-- Reference and external ID indexes
CREATE INDEX idx_transactions_reference_id ON transactions(reference_id);
CREATE INDEX idx_transactions_external_id ON transactions(external_id);

-- Date range query indexes
CREATE INDEX idx_transactions_processed_at ON transactions(processed_at);
CREATE INDEX idx_transactions_failed_at ON transactions(failed_at);

-- Composite indexes for common queries
CREATE INDEX idx_transactions_user_status_created ON transactions(user_id, status, created_at);
CREATE INDEX idx_transactions_user_type_created ON transactions(user_id, type, created_at);
CREATE INDEX idx_transactions_account_created ON transactions(from_account_id, created_at);
CREATE INDEX idx_transactions_card_created ON transactions(from_card_id, created_at);

-- Create transaction_rules table for business rules
CREATE TABLE IF NOT EXISTS transaction_rules (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36),
    account_id VARCHAR(36),
    card_id VARCHAR(36),
    transaction_type VARCHAR(50),
    
    -- Rule parameters
    max_daily_amount DECIMAL(15,2),
    max_single_amount DECIMAL(15,2),
    min_amount DECIMAL(15,2),
    requires_approval BOOLEAN DEFAULT FALSE,
    allowed_hours VARCHAR(50), -- JSON array of allowed hours
    allowed_days VARCHAR(20),  -- JSON array of allowed days
    
    -- Rule status
    is_active BOOLEAN DEFAULT TRUE,
    effective_from TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    effective_until TIMESTAMP NULL,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by VARCHAR(36) NOT NULL,
    
    -- Constraints
    CONSTRAINT chk_rule_amounts CHECK (
        (max_daily_amount IS NULL OR max_daily_amount > 0) AND
        (max_single_amount IS NULL OR max_single_amount > 0) AND
        (min_amount IS NULL OR min_amount >= 0)
    ),
    
    -- Indexes
    INDEX idx_transaction_rules_user_id (user_id),
    INDEX idx_transaction_rules_account_id (account_id),
    INDEX idx_transaction_rules_card_id (card_id),
    INDEX idx_transaction_rules_type (transaction_type),
    INDEX idx_transaction_rules_active (is_active),
    INDEX idx_transaction_rules_effective (effective_from, effective_until)
);

-- Create transaction_limits table for tracking daily/monthly limits
CREATE TABLE IF NOT EXISTS transaction_limits (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    account_id VARCHAR(36),
    card_id VARCHAR(36),
    transaction_type VARCHAR(50),
    
    -- Limit tracking
    period_type ENUM('daily', 'weekly', 'monthly') NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    
    -- Usage tracking
    transaction_count INT DEFAULT 0,
    total_amount DECIMAL(15,2) DEFAULT 0.00,
    
    -- Limits
    max_transactions INT,
    max_amount DECIMAL(15,2),
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Unique constraint for period tracking
    UNIQUE KEY unique_period (user_id, account_id, card_id, transaction_type, period_type, period_start),
    
    -- Indexes
    INDEX idx_transaction_limits_user_period (user_id, period_type, period_start, period_end),
    INDEX idx_transaction_limits_account_period (account_id, period_type, period_start, period_end),
    INDEX idx_transaction_limits_card_period (card_id, period_type, period_start, period_end)
);

-- Create transaction_audit table for audit trail
CREATE TABLE IF NOT EXISTS transaction_audit (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL,
    action VARCHAR(50) NOT NULL, -- 'created', 'updated', 'completed', 'failed', 'canceled', 'reversed'
    
    -- Change tracking
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    changed_fields JSON,
    
    -- Audit details
    changed_by VARCHAR(36) NOT NULL,
    change_reason TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_transaction_audit_transaction_id (transaction_id),
    INDEX idx_transaction_audit_action (action),
    INDEX idx_transaction_audit_changed_by (changed_by),
    INDEX idx_transaction_audit_created_at (created_at),
    
    -- Foreign key constraint
    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);

-- Insert default transaction rules (these can be customized per user/account later)
INSERT IGNORE INTO transaction_rules (
    id, 
    transaction_type, 
    max_daily_amount, 
    max_single_amount, 
    created_by
) VALUES 
    ('default-wallet-deposit', 'wallet_deposit', 50000.00, 10000.00, 'system'),
    ('default-wallet-withdrawal', 'wallet_withdrawal', 20000.00, 5000.00, 'system'),
    ('default-credit-charge', 'credit_charge', 100000.00, 50000.00, 'system'),
    ('default-debit-purchase', 'debit_purchase', 30000.00, 10000.00, 'system'),
    ('default-account-transfer', 'account_transfer', 100000.00, 25000.00, 'system');

-- Add comments for documentation
ALTER TABLE transactions COMMENT = 'Main transactions table storing all financial transactions';
ALTER TABLE transaction_rules COMMENT = 'Business rules and limits for different transaction types';
ALTER TABLE transaction_limits COMMENT = 'Period-based tracking of transaction limits and usage';
ALTER TABLE transaction_audit COMMENT = 'Audit trail for all changes to transactions';