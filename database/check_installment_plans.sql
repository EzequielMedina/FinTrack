-- =====================================================
-- Ver planes de cuotas en la base de datos
-- Ejecutar contra la BD del account-service (MySQL)
-- =====================================================

-- Listar todos los planes de cuotas (todos los estados)
SELECT 
  id,
  card_id,
  user_id,
  total_amount,
  installments_count,
  paid_installments,
  remaining_amount,
  status,
  description,
  created_at,
  completed_at
FROM installment_plans
ORDER BY created_at DESC;

-- Contar por estado
SELECT status, COUNT(*) AS cantidad
FROM installment_plans
GROUP BY status;

-- Planes por tarjeta (reemplazar CARD_ID por el id de tu tarjeta)
-- SELECT * FROM installment_plans WHERE card_id = 'CARD_ID' ORDER BY created_at DESC;
