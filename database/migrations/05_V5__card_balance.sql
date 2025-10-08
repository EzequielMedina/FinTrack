-- Migration: 05_V5__card_balance.sql
-- Description: Add balance field to cards table for credit card debt management

USE fintrack;

-- Add balance field to cards table
ALTER TABLE cards 
ADD COLUMN balance DECIMAL(15,2) NOT NULL DEFAULT 0.00 COMMENT 'Credit card balance/debt (negative = debt, positive = overpayment)';

-- Add current_balance for debit cards (references account balance)
-- This is computed, not stored, but we add a virtual column for clarity
-- ALTER TABLE cards 
-- ADD COLUMN available_balance DECIMAL(15,2) AS (
--     CASE 
--         WHEN card_type = 'debit' THEN (SELECT balance FROM accounts WHERE id = account_id)
--         WHEN card_type = 'credit' THEN GREATEST(0, IFNULL(credit_limit, 0) - ABS(balance))
--         ELSE 0
--     END
-- ) VIRTUAL COMMENT 'Available balance: account balance for debit, available credit for credit cards';

-- Add index for balance queries
CREATE INDEX idx_cards_balance ON cards(balance);

-- Add comments to clarify the balance usage
ALTER TABLE cards MODIFY COLUMN balance DECIMAL(15,2) NOT NULL DEFAULT 0.00 
COMMENT 'For credit cards: debt amount (positive = owed to bank). For debit cards: should remain 0 (uses account balance)';

-- Update existing cards to have 0 balance
UPDATE cards SET balance = 0.00 WHERE balance IS NULL;

-- Add constraint to ensure debit cards have 0 balance
-- This will be enforced at application level for flexibility
-- ALTER TABLE cards ADD CONSTRAINT chk_debit_card_balance 
-- CHECK (card_type != 'debit' OR balance = 0.00);

COMMIT;