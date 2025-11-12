# Test Chatbot Conversacional - Fase 1
# Este script prueba las nuevas funcionalidades conversacionales del chatbot

$baseUrl = "http://localhost:8090/api/chat"
$userId = "6a67040e-79fe-4b98-8980-1929f2b5b8bb"
$headers = @{
    "Content-Type" = "application/json"
    "X-User-ID" = $userId
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "PRUEBA 1: Primera consulta (sin conversationId)" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$body1 = @{
    userId = $userId
    message = "¿cuánto gasté hoy?"
} | ConvertTo-Json

Write-Host "Request:" -ForegroundColor Yellow
Write-Host $body1 -ForegroundColor Gray

try {
    $response1 = Invoke-RestMethod -Uri "$baseUrl/query" -Method Post -Headers $headers -Body $body1
    Write-Host "`nResponse:" -ForegroundColor Green
    $response1 | ConvertTo-Json -Depth 10 | Write-Host
    
    $conversationId = $response1.conversationId
    Write-Host "`n--- ConversationID generado: $conversationId" -ForegroundColor Green
    Write-Host "--- Periodo inferido: $($response1.inferredPeriod)" -ForegroundColor Green
    Write-Host "--- Contexto inferido: $($response1.inferredContext)" -ForegroundColor Green
    
    if ($response1.quickSuggestions) {
        Write-Host "`nSugerencias rápidas:" -ForegroundColor Magenta
        $response1.quickSuggestions | ForEach-Object { Write-Host "  - $_" -ForegroundColor Magenta }
    }

    Start-Sleep -Seconds 2

    # ============================================
    Write-Host "`n`n========================================" -ForegroundColor Cyan
    Write-Host "PRUEBA 2: Pregunta de seguimiento (con conversationId)" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan

    $body2 = @{
        userId = $userId
        message = "¿y mis tarjetas?"
        conversationId = $conversationId
    } | ConvertTo-Json

    Write-Host "Request:" -ForegroundColor Yellow
    Write-Host $body2 -ForegroundColor Gray

    $response2 = Invoke-RestMethod -Uri "$baseUrl/query" -Method Post -Headers $headers -Body $body2
    Write-Host "`nResponse:" -ForegroundColor Green
    $response2 | ConvertTo-Json -Depth 10 | Write-Host
    
    Write-Host "`n--- Mismo ConversationID: $($response2.conversationId)" -ForegroundColor Green
    Write-Host "--- Periodo inferido: $($response2.inferredPeriod)" -ForegroundColor Green
    Write-Host "--- Contexto inferido: $($response2.inferredContext)" -ForegroundColor Green

    if ($response2.quickSuggestions) {
        Write-Host "`nSugerencias rápidas:" -ForegroundColor Magenta
        $response2.quickSuggestions | ForEach-Object { Write-Host "  - $_" -ForegroundColor Magenta }
    }

    Start-Sleep -Seconds 2

    # ============================================
    Write-Host "`n`n========================================" -ForegroundColor Cyan
    Write-Host "PRUEBA 3: Recuperar historial de conversación" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan

    Write-Host "Request: GET /history/$conversationId" -ForegroundColor Yellow

    $historyResponse = Invoke-RestMethod -Uri "$baseUrl/history/$conversationId" -Method Get -Headers $headers
    Write-Host "`nResponse:" -ForegroundColor Green
    $historyResponse | ConvertTo-Json -Depth 10 | Write-Host

    Write-Host "`n--- Total de mensajes: $($historyResponse.total)" -ForegroundColor Green
    Write-Host "`nHistorial de conversacion:" -ForegroundColor Cyan
    
    foreach ($msg in $historyResponse.messages) {
        if ($msg.role -eq "user") {
            $role = "Usuario"
            $color = "Blue"
        } else {
            $role = "Asistente"
            $color = "Green"
        }
        
        $time = ([datetime]$msg.createdAt).ToString("HH:mm:ss")
        Write-Host "`n[$time] $role" -ForegroundColor $color
        Write-Host "  $($msg.message)" -ForegroundColor Gray
        
        if ($msg.contextData) {
            Write-Host "  Contexto: $($msg.contextData.inferredContext) | Periodo: $($msg.contextData.inferredPeriod)" -ForegroundColor DarkGray
        }
    }

    # ============================================
    Write-Host "`n`n========================================" -ForegroundColor Cyan
    Write-Host "PRUEBA 4: Inferencia de períodos variados" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan

    $testMessages = @(
        "¿cuánto gasté ayer?",
        "muéstrame gastos de esta semana",
        "estado de cuotas del mes pasado",
        "ingresos de los últimos 30 días"
    )

    foreach ($testMsg in $testMessages) {
        Write-Host "`nProbando: '$testMsg'" -ForegroundColor Yellow
        
        $testBody = @{
            userId = $userId
            message = $testMsg
        } | ConvertTo-Json

        $testResp = Invoke-RestMethod -Uri "$baseUrl/query" -Method Post -Headers $headers -Body $testBody
        Write-Host "  --- Periodo inferido: $($testResp.inferredPeriod)" -ForegroundColor Green
        Write-Host "  --- Contexto inferido: $($testResp.inferredContext)" -ForegroundColor Green
        
        Start-Sleep -Milliseconds 500
    }

    # ============================================
    Write-Host "`n`n========================================" -ForegroundColor Cyan
    Write-Host "PRUEBA 5: Inferencia de contextos variados" -ForegroundColor Cyan
    Write-Host "========================================`n" -ForegroundColor Cyan

    $contextMessages = @(
        @{ msg = "estado de mis tarjetas"; expectedContext = "cards" },
        @{ msg = "cuánto debo en cuotas"; expectedContext = "installments" },
        @{ msg = "mis gastos"; expectedContext = "expenses" },
        @{ msg = "ingresos"; expectedContext = "income" },
        @{ msg = "principales comercios"; expectedContext = "merchants" }
    )

    foreach ($test in $contextMessages) {
        Write-Host "`nProbando: '$($test.msg)'" -ForegroundColor Yellow
        
        $testBody = @{
            userId = $userId
            message = $test.msg
        } | ConvertTo-Json

        $testResp = Invoke-RestMethod -Uri "$baseUrl/query" -Method Post -Headers $headers -Body $testBody
        
        $isCorrect = $testResp.inferredContext -eq $test.expectedContext
        if ($isCorrect) {
            $status = "OK"
            $color = "Green"
        } else {
            $status = "ERROR"
            $color = "Red"
        }
        
        Write-Host "  $status Contexto inferido: $($testResp.inferredContext) (esperado: $($test.expectedContext))" -ForegroundColor $color
        
        Start-Sleep -Milliseconds 500
    }

    # ============================================
    Write-Host "`n`n========================================" -ForegroundColor Green
    Write-Host "TODAS LAS PRUEBAS COMPLETADAS" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    
    Write-Host "`nResumen:" -ForegroundColor Cyan
    Write-Host "  - Generacion automatica de conversationId" -ForegroundColor Green
    Write-Host "  - Inferencia de periodos temporales" -ForegroundColor Green
    Write-Host "  - Inferencia de contextos" -ForegroundColor Green
    Write-Host "  - Continuidad conversacional" -ForegroundColor Green
    Write-Host "  - Historial de conversacion" -ForegroundColor Green
    Write-Host "  - Sugerencias rapidas contextuales" -ForegroundColor Green

} catch {
    Write-Host "`nERROR:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host "`nDetalles:" -ForegroundColor Yellow
    Write-Host $_.ErrorDetails.Message -ForegroundColor Gray
}

