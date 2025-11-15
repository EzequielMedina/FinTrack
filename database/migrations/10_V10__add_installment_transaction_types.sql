-- Migration: Add installment transaction types to transactions table
-- Description: Adds 'installment_payment', 'installment_refund', and 'installment_plan_completion' 
--              to the type constraint and 'installment_completion' to payment_method constraint
-- Date: 2025-11-05

-- Drop the existing type constraint
ALTER TABLE transactions
DROP CHECK chk_valid_type;

-- Add the new type constraint with installment types
ALTER TABLE transactions
ADD CONSTRAINT chk_valid_type CHECK (type IN (
    'wallet_deposit', 'wallet_withdrawal', 'wallet_transfer',
    'credit_charge', 'credit_payment', 'credit_refund',
    'debit_purchase', 'debit_withdrawal', 'debit_refund',
    'account_transfer', 'account_deposit', 'account_withdraw',
    'installment_payment', 'installment_refund', 'installment_plan_completion',
    'credit_purchase_installments'
));

-- Drop the existing payment_method constraint
ALTER TABLE transactions
DROP CHECK chk_valid_payment_method;

-- Add the new payment_method constraint with installment types
ALTER TABLE transactions
ADD CONSTRAINT chk_valid_payment_method CHECK (
    payment_method IS NULL OR payment_method IN (
        'cash', 'bank_transfer', 'credit_card', 'debit_card', 'wallet',
        'installment_payment', 'installment_completion', 'credit_card_installments'
    )
);
