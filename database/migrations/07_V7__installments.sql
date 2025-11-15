-- =====================================================
-- FinTrack Account Service - Database Migration
-- Version: V7__installments.sql
-- Description: Create installment plans and installments tables for credit card installment payments
-- =====================================================

-- Create installment_plans table
CREATE TABLE IF NOT EXISTS installment_plans (
    -- Core identity
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL,
    card_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    
    -- Plan details
    total_amount DECIMAL(15,2) NOT NULL,
    installments_count INT NOT NULL,
    installment_amount DECIMAL(15,2) NOT NULL,
    start_date DATE NOT NULL,
    
    -- Merchant information (optional)
    merchant_name VARCHAR(255),
    merchant_id VARCHAR(100),
    description TEXT,
    
    -- Status tracking
    status ENUM('active', 'completed', 'cancelled', 'suspended') NOT NULL DEFAULT 'active',
    paid_installments INT NOT NULL DEFAULT 0,
    remaining_amount DECIMAL(15,2) NOT NULL,
    
    -- Interest and fees (for future enhancement)
    interest_rate DECIMAL(5,2) DEFAULT 0.00,
    total_interest DECIMAL(15,2) DEFAULT 0.00,
    admin_fee DECIMAL(15,2) DEFAULT 0.00,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    cancelled_at TIMESTAMP NULL,
    
    -- Constraints
    CONSTRAINT chk_installment_plan_amount_positive CHECK (total_amount > 0),
    CONSTRAINT chk_installment_plan_count_valid CHECK (installments_count >= 1 AND installments_count <= 24),
    CONSTRAINT chk_installment_plan_installment_amount_positive CHECK (installment_amount > 0),
    CONSTRAINT chk_installment_plan_paid_valid CHECK (paid_installments >= 0 AND paid_installments <= installments_count),
    CONSTRAINT chk_installment_plan_remaining_valid CHECK (remaining_amount >= 0),
    CONSTRAINT chk_installment_plan_status_valid CHECK (status IN ('active', 'completed', 'cancelled', 'suspended')),
    
    -- Foreign key constraints (will be enabled when services are fully integrated)
    -- FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    -- FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE,
    -- FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Indexes for performance optimization
    INDEX idx_installment_plans_transaction_id (transaction_id),
    INDEX idx_installment_plans_card_id (card_id),
    INDEX idx_installment_plans_user_id (user_id),
    INDEX idx_installment_plans_status (status),
    INDEX idx_installment_plans_start_date (start_date),
    INDEX idx_installment_plans_created_at (created_at),
    INDEX idx_installment_plans_user_status (user_id, status),
    INDEX idx_installment_plans_card_status (card_id, status)
);

-- Create installments table
CREATE TABLE IF NOT EXISTS installments (
    -- Core identity
    id VARCHAR(36) PRIMARY KEY,
    plan_id VARCHAR(36) NOT NULL,
    installment_number INT NOT NULL,
    
    -- Payment details
    amount DECIMAL(15,2) NOT NULL,
    due_date DATE NOT NULL,
    paid_date TIMESTAMP NULL,
    status ENUM('pending', 'paid', 'overdue', 'cancelled', 'partial') NOT NULL DEFAULT 'pending',
    
    -- Payment information
    paid_amount DECIMAL(15,2) DEFAULT 0.00,
    remaining_amount DECIMAL(15,2) NOT NULL,
    payment_method VARCHAR(30),
    payment_reference VARCHAR(100),
    
    -- Transaction references
    payment_transaction_id VARCHAR(36) NULL,
    
    -- Late fees and penalties (for future enhancement)
    late_fee DECIMAL(15,2) DEFAULT 0.00,
    penalty_amount DECIMAL(15,2) DEFAULT 0.00,
    grace_period_days INT DEFAULT 0,
    
    -- Audit fields
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_installment_amount_positive CHECK (amount > 0),
    CONSTRAINT chk_installment_number_positive CHECK (installment_number > 0),
    CONSTRAINT chk_installment_paid_amount_valid CHECK (paid_amount >= 0 AND paid_amount <= amount),
    CONSTRAINT chk_installment_remaining_valid CHECK (remaining_amount >= 0 AND remaining_amount <= amount),
    CONSTRAINT chk_installment_status_valid CHECK (status IN ('pending', 'paid', 'overdue', 'cancelled', 'partial')),
    CONSTRAINT chk_installment_payment_method_valid CHECK (
        payment_method IS NULL OR payment_method IN (
            'automatic', 'manual', 'bank_transfer', 'credit_card', 'debit_card', 'cash', 'wallet'
        )
    ),
    
    -- Foreign key constraints
    FOREIGN KEY (plan_id) REFERENCES installment_plans(id) ON DELETE CASCADE,
    -- FOREIGN KEY (payment_transaction_id) REFERENCES transactions(id) ON DELETE SET NULL,
    
    -- Unique constraint for plan and installment number
    UNIQUE KEY unique_plan_installment (plan_id, installment_number),
    
    -- Indexes for performance optimization
    INDEX idx_installments_plan_id (plan_id),
    INDEX idx_installments_due_date (due_date),
    INDEX idx_installments_status (status),
    INDEX idx_installments_payment_transaction_id (payment_transaction_id),
    INDEX idx_installments_plan_status (plan_id, status),
    INDEX idx_installments_due_status (due_date, status),
    INDEX idx_installments_created_at (created_at)
);

-- Create installment_plan_audit table for change tracking
CREATE TABLE IF NOT EXISTS installment_plan_audit (
    id VARCHAR(36) PRIMARY KEY,
    plan_id VARCHAR(36) NOT NULL,
    action VARCHAR(50) NOT NULL, -- 'created', 'payment_applied', 'status_changed', 'cancelled', 'suspended'
    
    -- Change tracking
    old_status VARCHAR(20),
    new_status VARCHAR(20),
    old_paid_installments INT,
    new_paid_installments INT,
    old_remaining_amount DECIMAL(15,2),
    new_remaining_amount DECIMAL(15,2),
    
    -- Additional context
    installment_id VARCHAR(36) NULL, -- Related installment if applicable
    payment_amount DECIMAL(15,2) NULL,
    changed_by VARCHAR(36) NOT NULL,
    change_reason TEXT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    
    -- Metadata (JSON for flexible tracking)
    metadata JSON,
    
    -- Timestamp
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_installment_plan_audit_plan_id (plan_id),
    INDEX idx_installment_plan_audit_action (action),
    INDEX idx_installment_plan_audit_changed_by (changed_by),
    INDEX idx_installment_plan_audit_created_at (created_at),
    INDEX idx_installment_plan_audit_installment_id (installment_id),
    
    -- Foreign key constraint
    FOREIGN KEY (plan_id) REFERENCES installment_plans(id) ON DELETE CASCADE
);

-- Create view for installment plan summary (for easier querying)
CREATE OR REPLACE VIEW installment_plan_summary AS
SELECT 
    ip.id,
    ip.card_id,
    ip.user_id,
    ip.total_amount,
    ip.installments_count,
    ip.installment_amount,
    ip.start_date,
    ip.status,
    ip.paid_installments,
    ip.remaining_amount,
    ip.description,
    ip.merchant_name,
    ip.created_at,
    
    -- Next installment information
    (SELECT MIN(due_date) FROM installments i WHERE i.plan_id = ip.id AND i.status = 'pending') as next_due_date,
    (SELECT amount FROM installments i WHERE i.plan_id = ip.id AND i.status = 'pending' ORDER BY due_date LIMIT 1) as next_installment_amount,
    
    -- Overdue information
    (SELECT COUNT(*) FROM installments i WHERE i.plan_id = ip.id AND i.status = 'overdue') as overdue_count,
    (SELECT SUM(remaining_amount) FROM installments i WHERE i.plan_id = ip.id AND i.status = 'overdue') as overdue_amount,
    
    -- Progress calculation
    ROUND((ip.paid_installments / ip.installments_count) * 100, 2) as completion_percentage,
    (ip.installments_count - ip.paid_installments) as remaining_installments
    
FROM installment_plans ip
WHERE ip.status != 'cancelled';

-- Insert default configuration for installment limits and rules
INSERT IGNORE INTO transaction_rules (
    id, 
    transaction_type, 
    max_daily_amount, 
    max_single_amount,
    min_amount,
    created_by,
    created_at
) VALUES 
    ('installment-credit-charge-1', 'credit_charge', NULL, NULL, 1000.00, 'system', NOW()),
    ('installment-credit-charge-2', 'credit_charge', NULL, 50000.00, NULL, 'system', NOW());

-- Create trigger to automatically update installment plan status
DELIMITER $$

CREATE TRIGGER update_installment_plan_status 
AFTER UPDATE ON installments
FOR EACH ROW
BEGIN
    DECLARE plan_total_installments INT;
    DECLARE plan_paid_installments INT;
    DECLARE plan_remaining_amount DECIMAL(15,2);
    
    -- Get plan information
    SELECT installments_count INTO plan_total_installments
    FROM installment_plans 
    WHERE id = NEW.plan_id;
    
    -- Count paid installments
    SELECT COUNT(*) INTO plan_paid_installments
    FROM installments 
    WHERE plan_id = NEW.plan_id AND status = 'paid';
    
    -- Calculate remaining amount
    SELECT SUM(remaining_amount) INTO plan_remaining_amount
    FROM installments 
    WHERE plan_id = NEW.plan_id AND status != 'paid';
    
    -- Set remaining amount to 0 if NULL
    IF plan_remaining_amount IS NULL THEN
        SET plan_remaining_amount = 0;
    END IF;
    
    -- Update installment plan
    UPDATE installment_plans 
    SET 
        paid_installments = plan_paid_installments,
        remaining_amount = plan_remaining_amount,
        status = CASE 
            WHEN plan_paid_installments = plan_total_installments THEN 'completed'
            WHEN plan_paid_installments > 0 THEN 'active'
            ELSE status
        END,
        completed_at = CASE 
            WHEN plan_paid_installments = plan_total_installments THEN NOW()
            ELSE completed_at
        END,
        updated_at = NOW()
    WHERE id = NEW.plan_id;
    
    -- Log audit trail if status changed
    IF OLD.status != NEW.status THEN
        INSERT INTO installment_plan_audit (
            id, plan_id, action, old_status, new_status,
            installment_id, payment_amount, changed_by, change_reason
        ) VALUES (
            UUID(), NEW.plan_id, 'payment_applied', OLD.status, NEW.status,
            NEW.id, NEW.paid_amount, 'system', 'Installment payment processed'
        );
    END IF;
END$$

DELIMITER ;

-- Create trigger to mark overdue installments
DELIMITER $$

CREATE EVENT IF NOT EXISTS mark_overdue_installments
ON SCHEDULE EVERY 1 DAY
STARTS CURRENT_DATE + INTERVAL 1 DAY
DO
BEGIN
    -- Mark installments as overdue if past due date and still pending
    UPDATE installments 
    SET 
        status = 'overdue',
        updated_at = NOW()
    WHERE 
        status = 'pending' 
        AND due_date < CURDATE();
        
    -- Log audit for overdue installments
    INSERT INTO installment_plan_audit (
        id, plan_id, action, installment_id, changed_by, change_reason
    )
    SELECT 
        UUID(), i.plan_id, 'status_changed', i.id, 'system', 'Marked as overdue automatically'
    FROM installments i
    WHERE i.status = 'overdue' 
        AND DATE(i.updated_at) = CURDATE();
END$$

DELIMITER ;

-- Add comments for documentation
ALTER TABLE installment_plans COMMENT = 'Credit card installment plans for dividing purchases into monthly payments';
ALTER TABLE installments COMMENT = 'Individual installment payments within an installment plan';
ALTER TABLE installment_plan_audit COMMENT = 'Audit trail for all changes to installment plans and payments';

-- Create indexes for common query patterns
CREATE INDEX idx_installment_plans_active_user ON installment_plans(user_id, status);
CREATE INDEX idx_installments_pending_due ON installments(due_date, status);

-- Performance optimization: Create composite indexes for frequent joins
CREATE INDEX idx_installment_plans_card_user_status ON installment_plans(card_id, user_id, status);
CREATE INDEX idx_installments_plan_due_status ON installments(plan_id, due_date, status);