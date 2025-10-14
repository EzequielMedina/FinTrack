-- Script para crear datos de prueba para notification-service
-- EJECUTAR EN MySQL/Adminer (http://localhost:8080)

-- 1. Crear usuario de prueba
INSERT IGNORE INTO users (id, email, password_hash, first_name, last_name, role, is_active, email_verified, created_at, updated_at)
VALUES (
    'test-user-001', 
    'test@fintrack.com', 
    '$2a$10$example.hash', 
    'Juan', 
    'Pérez', 
    'USER', 
    1, 
    1, 
    NOW(), 
    NOW()
);

-- 2. Crear cuenta para el usuario
INSERT IGNORE INTO accounts (id, user_id, account_type, name, description, currency, balance, dni, is_active, created_at, updated_at)
VALUES (
    'test-account-001',
    'test-user-001',
    'bank_account',
    'Cuenta Principal Test',
    'Cuenta para testing',
    'ARS',
    50000.00,
    NULL,
    1,
    NOW(),
    NOW()
);

-- 3. Crear tarjeta que vence MAÑANA
INSERT IGNORE INTO cards (
    id, account_id, card_type, card_brand, last_four_digits, masked_number, 
    holder_name, expiration_month, expiration_year, status, is_default, 
    nickname, credit_limit, closing_date, due_date, encrypted_number, 
    key_fingerprint, created_at, updated_at
)
VALUES (
    'test-card-001',
    'test-account-001',
    'credit',
    'Visa',
    '1234',
    '**** **** **** 1234',
    'Juan Pérez',
    12,
    2025,
    'active',
    1,
    'Visa Principal',
    100000.00,
    DATE_SUB(CURDATE(), INTERVAL 5 DAY),
    DATE_ADD(CURDATE(), INTERVAL 1 DAY), -- VENCE MAÑANA
    'encrypted_data_example',
    'fingerprint_example',
    NOW(),
    NOW()
);

-- 4. Crear plan de cuotas
INSERT IGNORE INTO installment_plans (
    id, transaction_id, card_id, user_id, total_amount, installments_count,
    installment_amount, start_date, merchant_name, merchant_id, description,
    status, paid_installments, remaining_amount, created_at, updated_at
)
VALUES (
    'test-plan-001',
    'test-transaction-001',
    'test-card-001',
    'test-user-001',
    15000.00,
    6,
    2500.00,
    DATE_SUB(CURDATE(), INTERVAL 30 DAY),
    'Mercado Libre',
    'ML123',
    'Compra de electrónicos',
    'active',
    2,
    10000.00,
    NOW(),
    NOW()
);

-- 5. Crear cuotas pendientes (que vencen mañana o antes)
INSERT IGNORE INTO installments (id, plan_id, installment_number, amount, due_date, status, remaining_amount, created_at, updated_at)
VALUES 
    ('test-installment-001', 'test-plan-001', 3, 2500.00, DATE_ADD(CURDATE(), INTERVAL 1 DAY), 'pending', 2500.00, NOW(), NOW()),
    ('test-installment-002', 'test-plan-001', 4, 2500.00, DATE_ADD(CURDATE(), INTERVAL 31 DAY), 'pending', 2500.00, NOW(), NOW());

-- 6. Crear otro plan con cuota que ya venció
INSERT IGNORE INTO installment_plans (
    id, transaction_id, card_id, user_id, total_amount, installments_count,
    installment_amount, start_date, merchant_name, description,
    status, paid_installments, remaining_amount, created_at, updated_at
)
VALUES (
    'test-plan-002',
    'test-transaction-002', 
    'test-card-001',
    'test-user-001',
    8000.00,
    4,
    2000.00,
    DATE_SUB(CURDATE(), INTERVAL 60 DAY),
    'Supermercado Coto',
    'Compras del mes',
    'active',
    1,
    6000.00,
    NOW(),
    NOW()
);

-- 7. Cuota vencida (para incluir en el reporte)
INSERT IGNORE INTO installments (id, plan_id, installment_number, amount, due_date, status, remaining_amount, created_at, updated_at)
VALUES 
    ('test-installment-003', 'test-plan-002', 2, 2000.00, DATE_SUB(CURDATE(), INTERVAL 5 DAY), 'overdue', 2000.00, NOW(), NOW());

-- VERIFICAR LOS DATOS CREADOS
SELECT 'USUARIOS CREADOS' as tipo;
SELECT id, email, first_name, last_name FROM users WHERE id = 'test-user-001';

SELECT 'CUENTAS CREADAS' as tipo;
SELECT id, user_id, name FROM accounts WHERE id = 'test-account-001';

SELECT 'TARJETAS CREADAS (que vencen mañana)' as tipo;
SELECT id, card_brand, last_four_digits, due_date, 
       DATEDIFF(due_date, CURDATE()) as dias_hasta_vencimiento
FROM cards WHERE id = 'test-card-001';

SELECT 'CUOTAS PENDIENTES' as tipo;
SELECT i.id, ip.merchant_name, i.amount, i.due_date, i.status,
       DATEDIFF(i.due_date, CURDATE()) as dias_hasta_vencimiento
FROM installments i 
JOIN installment_plans ip ON i.plan_id = ip.id 
WHERE ip.card_id = 'test-card-001' 
AND i.status IN ('pending', 'overdue')
ORDER BY i.due_date;