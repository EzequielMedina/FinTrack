# Test de Debug - Verificar contexto del chatbot
param(
    [string]$RealUserId = "6a67040e-79fe-4b98-8980-1929f2b5b8bb"
)

Write-Host "VERIFICANDO DATOS EN BASE DIRECTAMENTE" -ForegroundColor Cyan

# 1. Verificar planes creados hoy
Write-Host "1. Planes de cuotas creados hoy (2025-10-14):" -ForegroundColor Yellow
docker exec -it fintrack-mysql mysql -u fintrack_user -pfintrack_password fintrack -e "SELECT id, merchant_name, description, installments_count, remaining_amount, status, created_at FROM installment_plans WHERE user_id = '$RealUserId' AND DATE(created_at) = '2025-10-14' ORDER BY created_at DESC;"

Write-Host ""
Write-Host "2. Todas las transacciones del día:" -ForegroundColor Yellow
docker exec -it fintrack-mysql mysql -u fintrack_user -pfintrack_password fintrack -e "SELECT id, type, amount, merchant_name, status, created_at FROM transactions WHERE user_id = '$RealUserId' AND DATE(created_at) = '2025-10-14' ORDER BY created_at DESC;"

Write-Host ""
Write-Host "3. Información de cuentas:" -ForegroundColor Yellow
docker exec -it fintrack-mysql mysql -u fintrack_user -pfintrack_password fintrack -e "SELECT id, account_type, balance, currency, status FROM accounts WHERE user_id = '$RealUserId';"

Write-Host ""
Write-Host "4. Información de tarjetas:" -ForegroundColor Yellow
docker exec -it fintrack-mysql mysql -u fintrack_user -pfintrack_password fintrack -e "SELECT id, card_brand, last_four_digits, card_type, status FROM cards WHERE user_id = '$RealUserId';"

Write-Host ""
Write-Host "5. Ahora probando chatbot con pregunta simple:" -ForegroundColor Yellow

$payload = @{
    userId = $RealUserId
    message = "Lista todos los planes de cuotas que tengo"
    period = @{
        from = "2025-10-01"
        to = "2025-10-31"
    }
} | ConvertTo-Json -Depth 3

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8090/api/chat/query" `
        -Method POST `
        -ContentType "application/json" `
        -Body $payload `
        -UseBasicParsing `
        -TimeoutSec 60
    
    $result = ($response.Content | ConvertFrom-Json).reply
    Write-Host "Respuesta del chatbot:" -ForegroundColor Green
    Write-Host $result -ForegroundColor White
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "RESUMEN:" -ForegroundColor Magenta
Write-Host "- Compara los datos de la BD con la respuesta del chatbot" -ForegroundColor Gray
Write-Host "- Los planes creados hoy deberian aparecer en la respuesta" -ForegroundColor Gray