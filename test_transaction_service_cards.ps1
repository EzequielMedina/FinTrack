# Script para probar los endpoints de tarjetas en Transaction Service
Write-Host "=== Probando endpoints de tarjetas en Transaction Service ===" -ForegroundColor Green

$baseUrl = "http://localhost:8083"

# Test 1: Health check del Transaction Service
Write-Host "`n1. Verificando health check del Transaction Service..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET
    Write-Host "✓ Transaction Service está funcionando: $($healthResponse.status)" -ForegroundColor Green
} catch {
    Write-Host "✗ Error en health check: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Test 2: Cargo a tarjeta de crédito
Write-Host "`n2. Probando cargo a tarjeta de crédito..." -ForegroundColor Yellow
$chargeRequest = @{
    card_id = "card_123"
    amount = 100.50
    description = "Compra en tienda"
    merchant_name = "Tienda ABC"
} | ConvertTo-Json

try {
    $chargeResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/cards/credit/charge" -Method POST -Body $chargeRequest -ContentType "application/json"
    Write-Host "✓ Cargo exitoso - Transaction ID: $($chargeResponse.transaction_id)" -ForegroundColor Green
    Write-Host "  Status: $($chargeResponse.status)" -ForegroundColor Cyan
    Write-Host "  Message: $($chargeResponse.message)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Error en cargo: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error details: $errorBody" -ForegroundColor Red
    }
}

# Test 3: Pago a tarjeta de crédito
Write-Host "`n3. Probando pago a tarjeta de crédito..." -ForegroundColor Yellow
$paymentRequest = @{
    card_id = "card_123"
    amount = 50.25
    payment_method = "bank_transfer"
} | ConvertTo-Json

try {
    $paymentResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/cards/credit/payment" -Method POST -Body $paymentRequest -ContentType "application/json"
    Write-Host "✓ Pago exitoso - Transaction ID: $($paymentResponse.transaction_id)" -ForegroundColor Green
    Write-Host "  Status: $($paymentResponse.status)" -ForegroundColor Cyan
    Write-Host "  Message: $($paymentResponse.message)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Error en pago: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error details: $errorBody" -ForegroundColor Red
    }
}

# Test 4: Transacción con tarjeta de débito
Write-Host "`n4. Probando transacción con tarjeta de débito..." -ForegroundColor Yellow
$debitRequest = @{
    card_id = "card_456"
    amount = 75.00
    description = "Compra en supermercado"
    merchant_name = "Supermercado XYZ"
} | ConvertTo-Json

try {
    $debitResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/cards/debit/transaction" -Method POST -Body $debitRequest -ContentType "application/json"
    Write-Host "✓ Transacción débito exitosa - Transaction ID: $($debitResponse.transaction_id)" -ForegroundColor Green
    Write-Host "  Status: $($debitResponse.status)" -ForegroundColor Cyan
    Write-Host "  Message: $($debitResponse.message)" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Error en transacción débito: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Error details: $errorBody" -ForegroundColor Red
    }
}

Write-Host "`n=== Pruebas completadas ===" -ForegroundColor Green