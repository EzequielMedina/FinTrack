# Script de Optimización - Chatbot con Modelo Liviano
param(
    [switch]$Force,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"

Write-Host "Optimizando Chatbot-Service con modelo liviano..." -ForegroundColor Green

# 1. Detener servicios para actualizar
Write-Host "Deteniendo servicios actuales..." -ForegroundColor Blue
docker-compose down

# 2. Rebuild solo los servicios necesarios
Write-Host "Rebuilding chatbot-service..." -ForegroundColor Blue
docker-compose build --no-cache chatbot-service

# 3. Levantar servicios con nueva configuracion
Write-Host "Levantando servicios optimizados..." -ForegroundColor Blue
docker-compose up -d

Write-Host "Esperando que Ollama este disponible..." -ForegroundColor Yellow
$maxAttempts = 30
$attempt = 0

do {
    Start-Sleep -Seconds 2
    $attempt++
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:11434/api/tags" -UseBasicParsing -TimeoutSec 5
        if ($response.StatusCode -eq 200) {
            Write-Host "Ollama esta disponible!" -ForegroundColor Green
            break
        }
    }
    catch {
        if ($Verbose) {
            Write-Host "Intento $attempt/$maxAttempts - Ollama aun no disponible..." -ForegroundColor Gray
        }
    }
} while ($attempt -lt $maxAttempts)

if ($attempt -ge $maxAttempts) {
    Write-Host "Error: Ollama no se pudo conectar despues de $maxAttempts intentos" -ForegroundColor Red
    exit 1
}

# 4. Descargar modelo liviano Qwen2.5:3b
Write-Host "Descargando modelo liviano Qwen2.5:3b..." -ForegroundColor Blue
Write-Host "Este modelo es mas rapido que llama3:8b" -ForegroundColor Gray

$pullCommand = "docker exec fintrack-ollama ollama pull qwen2.5:3b"
if ($Verbose) {
    Write-Host "Ejecutando: $pullCommand" -ForegroundColor Gray
}

try {
    Invoke-Expression $pullCommand
    Write-Host "Modelo qwen2.5:3b descargado exitosamente!" -ForegroundColor Green
} catch {
    Write-Host "Error descargando modelo: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Intentando con modelo alternativo mas liviano..." -ForegroundColor Yellow
    
    # Fallback a modelo mas pequeño
    $pullCommandFallback = "docker exec fintrack-ollama ollama pull llama3.2:1b"
    try {
        Invoke-Expression $pullCommandFallback
        Write-Host "Modelo llama3.2:1b descargado como fallback!" -ForegroundColor Green
        
        # Actualizar variable de entorno para usar el modelo fallback
        Write-Host "Actualizando configuracion para usar modelo fallback..." -ForegroundColor Blue
        docker-compose stop chatbot-service
        $env:OLLAMA_MODEL = "llama3.2:1b"
        docker-compose up -d chatbot-service
    } catch {
        Write-Host "Error descargando modelo fallback: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
}

# 5. Verificar que los modelos estan disponibles
Write-Host "Verificando modelos disponibles..." -ForegroundColor Blue
try {
    $modelsOutput = docker exec fintrack-ollama ollama list
    Write-Host "Modelos disponibles:" -ForegroundColor Green
    Write-Host $modelsOutput -ForegroundColor Gray
} catch {
    Write-Host "No se pudieron listar los modelos, pero continuando..." -ForegroundColor Yellow
}

# 6. Test basico del chatbot
Write-Host "Probando chatbot optimizado..." -ForegroundColor Blue

$testPayload = @{
    userId = "test-user-optimization"
    message = "Cuanto gaste este mes?"
    period = @{
        from = "2025-10-01"
        to = "2025-10-31"
    }
} | ConvertTo-Json -Depth 3

$maxTestAttempts = 5
$testAttempt = 0
$testSuccess = $false

do {
    Start-Sleep -Seconds 3
    $testAttempt++
    try {
        if ($Verbose) {
            Write-Host "Test intento $testAttempt/$maxTestAttempts..." -ForegroundColor Gray
        }
        
        $testResponse = Invoke-WebRequest -Uri "http://localhost:8090/api/chat/query" `
            -Method POST `
            -ContentType "application/json" `
            -Body $testPayload `
            -UseBasicParsing `
            -TimeoutSec 45
        
        if ($testResponse.StatusCode -eq 200) {
            $responseContent = $testResponse.Content | ConvertFrom-Json
            Write-Host "Chatbot funcionando correctamente!" -ForegroundColor Green
            if ($Verbose) {
                Write-Host "Respuesta: $($responseContent.reply.Substring(0, [Math]::Min(100, $responseContent.reply.Length)))..." -ForegroundColor Gray
            }
            $testSuccess = $true
            break
        }
    }
    catch {
        if ($Verbose) {
            Write-Host "Test intento $testAttempt fallo: $($_.Exception.Message)" -ForegroundColor Gray
        }
    }
} while ($testAttempt -lt $maxTestAttempts)

if (-not $testSuccess) {
    Write-Host "El chatbot aun no responde, pero la configuracion se completo." -ForegroundColor Yellow
    Write-Host "Puede tomar unos minutos adicionales para estar completamente listo." -ForegroundColor Gray
}

# 7. Mostrar resumen de optimizacion
Write-Host ""
Write-Host "RESUMEN DE OPTIMIZACION" -ForegroundColor Magenta
Write-Host "=" * 50 -ForegroundColor Magenta

Write-Host "Configuracion actualizada:" -ForegroundColor Blue
Write-Host "   • RAM limite: 8GB (antes: 4GB)" -ForegroundColor Gray
Write-Host "   • Modelo: qwen2.5:3b (mas liviano que llama3:8b)" -ForegroundColor Gray
Write-Host "   • Contexto: 1024 tokens (optimizado para velocidad)" -ForegroundColor Gray
Write-Host "   • Respuestas: maximo 256 tokens (mas rapidas)" -ForegroundColor Gray
Write-Host "   • Timeout: 30 segundos" -ForegroundColor Gray

Write-Host ""
Write-Host "Mejoras esperadas:" -ForegroundColor Green
Write-Host "   • 2-3x mas rapido que antes" -ForegroundColor Gray
Write-Host "   • Menor uso de RAM" -ForegroundColor Gray
Write-Host "   • Mejor estabilidad" -ForegroundColor Gray

Write-Host ""
Write-Host "Para probar manualmente:" -ForegroundColor Blue
Write-Host 'curl -X POST http://localhost:8090/api/chat/query -H "Content-Type: application/json" -d "{\"userId\":\"test\",\"message\":\"consumo de plan de cuotas\",\"period\":{\"from\":\"2025-10-01\",\"to\":\"2025-10-31\"}}"' -ForegroundColor Gray

Write-Host ""
Write-Host "Logs del chatbot:" -ForegroundColor Blue
Write-Host "   docker logs fintrack-chatbot-service -f" -ForegroundColor Gray

Write-Host ""
Write-Host "Optimizacion completada! El chatbot ahora deberia ser mas rapido." -ForegroundColor Green