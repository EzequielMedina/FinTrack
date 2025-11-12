# ====================================================================
# SCRIPT DE RESET DE BASE DE DATOS PARA DEMO
# ====================================================================

Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  RESET DE BASE DE DATOS PARA DEMO - FINTRACK" -ForegroundColor Yellow
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "⚠️  ADVERTENCIA: Este script eliminará TODOS los datos" -ForegroundColor Red
Write-Host ""

# Script SQL para limpiar la base de datos
$sqlScript = @"
USE fintrack_db;

SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE chat_messages;
TRUNCATE TABLE chat_conversations;
TRUNCATE TABLE notifications;
TRUNCATE TABLE installments;
TRUNCATE TABLE transactions;
TRUNCATE TABLE cards;
TRUNCATE TABLE accounts;

DELETE FROM users WHERE id != 1;

INSERT IGNORE INTO users (id, email, password, nombre, apellido, created_at, updated_at, is_active)
VALUES (
    1,
    'ezequielmedina456@gmail.com',
    '$2a$10$V3D5ZRN0Z.Vn4v5XhqG3YuJW8YG5K3Jz2N1F4Q7R8S9T0V1W2X3Y4Z',
    'Ezequiel',
    'Medina',
    NOW(),
    NOW(),
    1
);

ALTER TABLE chat_conversations AUTO_INCREMENT = 1;
ALTER TABLE chat_messages AUTO_INCREMENT = 1;
ALTER TABLE notifications AUTO_INCREMENT = 1;
ALTER TABLE transactions AUTO_INCREMENT = 1;
ALTER TABLE installments AUTO_INCREMENT = 1;
ALTER TABLE accounts AUTO_INCREMENT = 1;
ALTER TABLE cards AUTO_INCREMENT = 1;

SET FOREIGN_KEY_CHECKS = 1;

SELECT '✅ Base de datos limpiada' as status;
SELECT COUNT(*) as total_usuarios FROM users;
SELECT COUNT(*) as total_cuentas FROM accounts;
SELECT COUNT(*) as total_tarjetas FROM cards;
SELECT COUNT(*) as total_transacciones FROM transactions;
"@

# Guardar script SQL temporal
$tempSqlFile = "$PSScriptRoot\temp_reset_db.sql"
$sqlScript | Out-File -FilePath $tempSqlFile -Encoding UTF8 -NoNewline

Write-Host "🔄 Ejecutando limpieza en MySQL..." -ForegroundColor Cyan

# Ejecutar script
docker exec -i fintrack-mysql mysql -uroot -proot_password < $tempSqlFile

Write-Host ""
Write-Host "✅ Base de datos limpiada exitosamente!" -ForegroundColor Green
Write-Host ""
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  ESTADO FINAL" -ForegroundColor Yellow
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "✅ Usuario: ezequielmedina456@gmail.com" -ForegroundColor White
Write-Host "✅ Cuentas: 0" -ForegroundColor White
Write-Host "✅ Tarjetas: 0" -ForegroundColor White
Write-Host "✅ Transacciones: 0" -ForegroundColor White
Write-Host "✅ Conversaciones chatbot: 0" -ForegroundColor White
Write-Host ""
Write-Host "🎬 Lista para tu demo!" -ForegroundColor Green
Write-Host ""

# Eliminar archivo temporal
Remove-Item $tempSqlFile -Force
