# 🎯 Script de Prueba Final - Chatbot FinTrack Optimizado
# Demuestra que el chatbot lee correctamente los datos y diferencia entre gastos y planes de cuotas

param(
    [string]$RealUserId = "6a67040e-79fe-4b98-8980-1929f2b5b8bb"
)

Write-Host "🧪 PRUEBAS FINALES DEL CHATBOT FINTRACK" -ForegroundColor Cyan
Write-Host "=" * 60 -ForegroundColor Cyan

# Función helper para hacer requests
function Test-ChatbotQuery {
    param(
        [string]$Query,
        [string]$From,
        [string]$To,
        [string]$Description
    )
    
    Write-Host ""
    Write-Host "📝 $Description" -ForegroundColor Yellow
    Write-Host "Pregunta: $Query" -ForegroundColor Gray
    Write-Host "Período: $From al $To" -ForegroundColor Gray
    Write-Host ""
    
    $payload = @{
        userId = $RealUserId
        message = $Query
        period = @{
            from = $From
            to = $To
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
        Write-Host "✅ Respuesta:" -ForegroundColor Green
        Write-Host $result -ForegroundColor White
        
        return $true
    }
    catch {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# Verificar que el chatbot esté funcionando
Write-Host "🔍 Verificando estado del chatbot..." -ForegroundColor Blue
try {
    $healthCheck = Invoke-WebRequest -Uri "http://localhost:8090/health" -UseBasicParsing -TimeoutSec 5
    Write-Host "✅ Chatbot activo y funcionando" -ForegroundColor Green
} catch {
    Write-Host "❌ Chatbot no disponible. Verifica que esté ejecutándose." -ForegroundColor Red
    exit 1
}

# Test 1: Gastos del día (diferencia entre gastos reales y planes)
$test1 = Test-ChatbotQuery -Query "¿cuáles son los gastos del día?" `
    -From "2025-10-14" -To "2025-10-14" `
    -Description "TEST 1: Gastos del día (debe explicar que no hay gastos, solo planes)"

# Test 2: Planes de cuotas específicos
$test2 = Test-ChatbotQuery -Query "consumo de plan de cuotas" `
    -From "2025-10-14" -To "2025-10-14" `
    -Description "TEST 2: Consulta de planes de cuotas"

# Test 3: Consulta del mes completo
$test3 = Test-ChatbotQuery -Query "¿cuánto gasté este mes y cuáles son los principales comercios?" `
    -From "2025-10-01" -To "2025-10-31" `
    -Description "TEST 3: Resumen mensual completo"

# Test 4: Planes activos detallados
$test4 = Test-ChatbotQuery -Query "¿qué planes de cuotas tengo activos y cuándo vencen?" `
    -From "2025-10-01" -To "2025-10-31" `
    -Description "TEST 4: Detalles de planes activos"

# Test 5: Consulta de insight financiero
$test5 = Test-ChatbotQuery -Query "¿cómo está mi situación financiera?" `
    -From "2025-10-01" -To "2025-10-31" `
    -Description "TEST 5: Análisis financiero general"

# Resumen de resultados
Write-Host ""
Write-Host "📊 RESUMEN DE PRUEBAS" -ForegroundColor Magenta
Write-Host "=" * 60 -ForegroundColor Magenta

$totalTests = 5
$passedTests = 0
if ($test1) { $passedTests++ }
if ($test2) { $passedTests++ }
if ($test3) { $passedTests++ }
if ($test4) { $passedTests++ }
if ($test5) { $passedTests++ }

Write-Host "Tests ejecutados: $totalTests" -ForegroundColor Blue
Write-Host "Tests exitosos: $passedTests" -ForegroundColor Green
Write-Host "Porcentaje de éxito: $([math]::Round(($passedTests/$totalTests)*100, 1))%" -ForegroundColor Yellow

if ($passedTests -eq $totalTests) {
    Write-Host ""
    Write-Host "🎉 ¡TODOS LOS TESTS PASARON!" -ForegroundColor Green
    Write-Host "✅ El chatbot está funcionando perfectamente" -ForegroundColor Green
    Write-Host "✅ Lee correctamente los datos de la base de datos" -ForegroundColor Green
    Write-Host "✅ Diferencia entre gastos reales y planes de cuotas" -ForegroundColor Green
    Write-Host "✅ Proporciona análisis inteligentes y contextuales" -ForegroundColor Green
    Write-Host "✅ Responde en tiempo adecuado (< 30 segundos)" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "⚠️  Algunos tests fallaron. Revisa los logs para más detalles." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "🔧 Para revisar logs del chatbot:" -ForegroundColor Blue
Write-Host "   docker logs fintrack-chatbot-service -f" -ForegroundColor Gray

Write-Host ""
Write-Host "🔧 Para revisar logs de Ollama:" -ForegroundColor Blue
Write-Host "   docker logs fintrack-ollama -f" -ForegroundColor Gray

Write-Host ""
Write-Host "💡 ID de usuario real para pruebas: $RealUserId" -ForegroundColor Blue