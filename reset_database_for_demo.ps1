# ====================================================================
# SCRIPT DE RESET DE BASE DE DATOS PARA DEMO
# ====================================================================
# Limpia completamente la base de datos y la deja lista para una demo
# desde cero con el usuario de prueba
# ====================================================================

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  RESET DE BASE DE DATOS PARA DEMO - FINTRACK" -ForegroundColor Yellow
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

# Confirmación
Write-Host "⚠️  ADVERTENCIA: Este script eliminará TODOS los datos" -ForegroundColor Red
Write-Host "   de la base de datos y dejará solo el usuario de prueba." -ForegroundColor Red
Write-Host ""
$confirmation = Read-Host "¿Estás seguro de continuar? (escribe 'SI' para confirmar)"

if ($confirmation -ne "SI") {
    Write-Host ""
    Write-Host "❌ Operación cancelada" -ForegroundColor Yellow
    Write-Host ""
    exit
}

Write-Host ""
Write-Host "🔄 Iniciando limpieza de base de datos..." -ForegroundColor Cyan
Write-Host ""

# Script SQL para limpiar la base de datos
$sqlScript = @"
-- ====================================================================
-- LIMPIEZA COMPLETA DE BASE DE DATOS PARA DEMO
-- ====================================================================

USE fintrack_db;

-- Deshabilitar verificación de claves foráneas temporalmente
SET FOREIGN_KEY_CHECKS = 0;

-- ====================================================================
-- 1. ELIMINAR DATOS DE CHATBOT
-- ====================================================================
TRUNCATE TABLE chat_messages;
TRUNCATE TABLE chat_conversations;

-- ====================================================================
-- 2. ELIMINAR DATOS DE NOTIFICACIONES
-- ====================================================================
TRUNCATE TABLE notifications;

-- ====================================================================
-- 3. ELIMINAR DATOS DE TRANSACCIONES Y CUOTAS
-- ====================================================================
TRUNCATE TABLE installments;
TRUNCATE TABLE transactions;

-- ====================================================================
-- 4. ELIMINAR CUENTAS Y TARJETAS
-- ====================================================================
TRUNCATE TABLE cards;
TRUNCATE TABLE accounts;

-- ====================================================================
-- 5. MANTENER SOLO EL USUARIO DE PRUEBA (ID 1)
-- ====================================================================
DELETE FROM users WHERE id != 1;

-- Verificar que existe el usuario de prueba, si no, crearlo
INSERT IGNORE INTO users (id, email, password, nombre, apellido, created_at, updated_at, is_active)
VALUES (
    1,
    'ezequielmedina456@gmail.com',
    '\$2a\$10\$V3D5ZRN0Z.Vn4v5XhqG3YuJW8YG5K3Jz2N1F4Q7R8S9T0V1W2X3Y4Z',
    'Ezequiel',
    'Medina',
    NOW(),
    NOW(),
    1
);

-- ====================================================================
-- 6. RESETEAR AUTO_INCREMENT
-- ====================================================================
ALTER TABLE chat_conversations AUTO_INCREMENT = 1;
ALTER TABLE chat_messages AUTO_INCREMENT = 1;
ALTER TABLE notifications AUTO_INCREMENT = 1;
ALTER TABLE transactions AUTO_INCREMENT = 1;
ALTER TABLE installments AUTO_INCREMENT = 1;
ALTER TABLE accounts AUTO_INCREMENT = 1;
ALTER TABLE cards AUTO_INCREMENT = 1;

-- Habilitar verificación de claves foráneas
SET FOREIGN_KEY_CHECKS = 1;

-- ====================================================================
-- VERIFICACIÓN FINAL
-- ====================================================================
SELECT '✅ Base de datos limpiada exitosamente' as status;
SELECT COUNT(*) as total_usuarios FROM users;
SELECT COUNT(*) as total_cuentas FROM accounts;
SELECT COUNT(*) as total_tarjetas FROM cards;
SELECT COUNT(*) as total_transacciones FROM transactions;
SELECT COUNT(*) as total_conversaciones FROM chat_conversations;
SELECT COUNT(*) as total_mensajes FROM chat_messages;
SELECT COUNT(*) as total_notificaciones FROM notifications;
"@

# Guardar script SQL temporal
$tempSqlFile = "c:\Facultad\Alumno\PS\temp_reset_db.sql"
$sqlScript | Out-File -FilePath $tempSqlFile -Encoding UTF8

Write-Host "📝 Script SQL generado" -ForegroundColor Green

# Ejecutar script en MySQL
Write-Host "🔄 Ejecutando limpieza en MySQL..." -ForegroundColor Cyan
Write-Host ""

try {
    # Ejecutar el script SQL en el contenedor de MySQL
    docker exec -i fintrack-mysql mysql -uroot -proot_password < $tempSqlFile
    
    Write-Host ""
    Write-Host "✅ Base de datos limpiada exitosamente!" -ForegroundColor Green
    Write-Host ""
    
    # Eliminar archivo temporal
    Remove-Item $tempSqlFile -Force
    
    Write-Host "============================================================" -ForegroundColor Cyan
    Write-Host "  ESTADO FINAL DE LA BASE DE DATOS" -ForegroundColor Yellow
    Write-Host "============================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "✅ Usuario de prueba: ezequielmedina456@gmail.com" -ForegroundColor White
    Write-Host "✅ Todas las cuentas eliminadas" -ForegroundColor White
    Write-Host "✅ Todas las tarjetas eliminadas" -ForegroundColor White
    Write-Host "✅ Todas las transacciones eliminadas" -ForegroundColor White
    Write-Host "✅ Todas las cuotas eliminadas" -ForegroundColor White
    Write-Host "✅ Todas las conversaciones del chatbot eliminadas" -ForegroundColor White
    Write-Host "✅ Todas las notificaciones eliminadas" -ForegroundColor White
    Write-Host ""
    Write-Host "============================================================" -ForegroundColor Cyan
    Write-Host "  PRÓXIMOS PASOS PARA LA DEMO" -ForegroundColor Yellow
    Write-Host "============================================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "1. Inicia sesión con: ezequielmedina456@gmail.com" -ForegroundColor White
    Write-Host ""
    Write-Host "2. Crea tus primeras cuentas y tarjetas" -ForegroundColor White
    Write-Host ""
    Write-Host "3. Agrega algunas transacciones de ejemplo" -ForegroundColor White
    Write-Host ""
    Write-Host "4. Prueba el chatbot conversacional:" -ForegroundColor White
    Write-Host "   - '¿cuánto gasté hoy?'" -ForegroundColor Gray
    Write-Host "   - 'muéstrame mis tarjetas'" -ForegroundColor Gray
    Write-Host "   - 'estado de cuotas'" -ForegroundColor Gray
    Write-Host ""
    Write-Host "============================================================" -ForegroundColor Cyan
    Write-Host ""
    
} catch {
    Write-Host ""
    Write-Host "❌ ERROR al limpiar la base de datos:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Yellow
    Write-Host ""
    
    # Eliminar archivo temporal en caso de error
    if (Test-Path $tempSqlFile) {
        Remove-Item $tempSqlFile -Force
    }
    
    exit 1
}
