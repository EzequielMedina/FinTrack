# Script de Prueba Final - Chatbot FinTrack Optimizado
param(
    [string]$RealUserId = "6a67040e-79fe-4b98-8980-1929f2b5b8bb"
)

Write-Host "PRUEBAS FINALES DEL CHATBOT FINTRACK" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan

# Verificar que el chatbot este funcionando
Write-Host "Verificando estado del chatbot..." -ForegroundColor Blue
try {
    $healthCheck = Invoke-WebRequest -Uri "http://localhost:8090/health" -UseBasicParsing -TimeoutSec 5
    Write-Host "Chatbot activo y funcionando" -ForegroundColor Green
} catch {
    Write-Host "Chatbot no disponible. Verifica que este ejecutandose." -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "TEST 1: Gastos del dia" -ForegroundColor Yellow
$payload1 = @{
    userId = $RealUserId
    message = "cuales son los gastos del dia?"
    period = @{
        from = "2025-10-14"
        to = "2025-10-14"
    }
} | ConvertTo-Json -Depth 3

try {
    $response1 = Invoke-WebRequest -Uri "http://localhost:8090/api/chat/query" `
        -Method POST `
        -ContentType "application/json" `
        -Body $payload1 `
        -UseBasicParsing `
        -TimeoutSec 60
    
    $result1 = ($response1.Content | ConvertFrom-Json).reply
    Write-Host "Respuesta:" -ForegroundColor Green
    Write-Host $result1 -ForegroundColor White
    $test1Success = $true
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $test1Success = $false
}

Write-Host ""
Write-Host "TEST 2: Planes de cuotas" -ForegroundColor Yellow
$payload2 = @{
    userId = $RealUserId
    message = "consumo de plan de cuotas"
    period = @{
        from = "2025-10-14"
        to = "2025-10-14"
    }
} | ConvertTo-Json -Depth 3

try {
    $response2 = Invoke-WebRequest -Uri "http://localhost:8090/api/chat/query" `
        -Method POST `
        -ContentType "application/json" `
        -Body $payload2 `
        -UseBasicParsing `
        -TimeoutSec 60
    
    $result2 = ($response2.Content | ConvertFrom-Json).reply
    Write-Host "Respuesta:" -ForegroundColor Green
    Write-Host $result2 -ForegroundColor White
    $test2Success = $true
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $test2Success = $false
}

Write-Host ""
Write-Host "TEST 3: Resumen mensual" -ForegroundColor Yellow
$payload3 = @{
    userId = $RealUserId
    message = "cuanto gaste este mes y cuales son los principales comercios?"
    period = @{
        from = "2025-10-01"
        to = "2025-10-31"
    }
} | ConvertTo-Json -Depth 3

try {
    $response3 = Invoke-WebRequest -Uri "http://localhost:8090/api/chat/query" `
        -Method POST `
        -ContentType "application/json" `
        -Body $payload3 `
        -UseBasicParsing `
        -TimeoutSec 60
    
    $result3 = ($response3.Content | ConvertFrom-Json).reply
    Write-Host "Respuesta:" -ForegroundColor Green
    Write-Host $result3 -ForegroundColor White
    $test3Success = $true
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    $test3Success = $false
}

# Resumen
Write-Host ""
Write-Host "RESUMEN DE PRUEBAS" -ForegroundColor Magenta
Write-Host "=" * 60 -ForegroundColor Magenta

$totalTests = 3
$passedTests = 0
if ($test1Success) { $passedTests++; Write-Host "TEST 1: PASADO" -ForegroundColor Green } else { Write-Host "TEST 1: FALLIDO" -ForegroundColor Red }
if ($test2Success) { $passedTests++; Write-Host "TEST 2: PASADO" -ForegroundColor Green } else { Write-Host "TEST 2: FALLIDO" -ForegroundColor Red }
if ($test3Success) { $passedTests++; Write-Host "TEST 3: PASADO" -ForegroundColor Green } else { Write-Host "TEST 3: FALLIDO" -ForegroundColor Red }

Write-Host ""
Write-Host "Tests exitosos: $passedTests de $totalTests" -ForegroundColor Yellow

if ($passedTests -eq $totalTests) {
    Write-Host ""
    Write-Host "TODOS LOS TESTS PASARON!" -ForegroundColor Green
    Write-Host "El chatbot esta funcionando perfectamente" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "Algunos tests fallaron. Revisa los logs." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "ID de usuario real: $RealUserId" -ForegroundColor Blue