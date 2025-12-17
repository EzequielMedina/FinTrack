# Script para probar el envío de emails desde el servicio de notificaciones

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Test de Notificaciones - FinTrack" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:3004"

# 1. Health Check
Write-Host "[1/4] Verificando salud del servicio..." -ForegroundColor Yellow
try {
    $health = Invoke-RestMethod -Uri "$baseUrl/health" -Method Get
    Write-Host "✓ Servicio: $($health.service)" -ForegroundColor Green
    Write-Host "✓ Status: $($health.status)" -ForegroundColor Green
    Write-Host "✓ Job Scheduler Running: $($health.job_scheduler_running)" -ForegroundColor Green
    if ($health.next_scheduled_run) {
        Write-Host "✓ Próxima ejecución programada: $($health.next_scheduled_run)" -ForegroundColor Green
    }
} catch {
    Write-Host "✗ Error al verificar salud del servicio: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
Write-Host ""

# 2. Enviar Email de Soporte (Este sí envía email)
Write-Host "[2/4] Enviando email de soporte de prueba..." -ForegroundColor Yellow
$supportPayload = @{
    name = "Usuario de Prueba"
    email = "test@example.com"
    subject = "Prueba de Email desde Notification Service"
    message = "Este es un mensaje de prueba enviado desde el script test_notification_email.ps1. Si recibes este email, significa que el servicio de notificaciones está funcionando correctamente."
} | ConvertTo-Json

try {
    $supportResponse = Invoke-RestMethod -Uri "$baseUrl/api/notifications/support" -Method Post -Body $supportPayload -ContentType "application/json"
    Write-Host "✓ Email de soporte enviado exitosamente!" -ForegroundColor Green
    Write-Host "  Mensaje: $($supportResponse.message)" -ForegroundColor Gray
    Write-Host "  Timestamp: $($supportResponse.timestamp)" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  📧 Revisa tu bandeja de entrada o spam!" -ForegroundColor Cyan
} catch {
    Write-Host "✗ Error al enviar email de soporte: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.ErrorDetails.Message) {
        $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
        Write-Host "  Detalles: $($errorDetails.details)" -ForegroundColor Red
    }
}
Write-Host ""

# 3. Trigger Card Due Job (Envía emails a usuarios con tarjetas por vencer)
Write-Host "[3/4] Disparando job de notificaciones de tarjetas por vencer..." -ForegroundColor Yellow
try {
    $cardDueResponse = Invoke-RestMethod -Uri "$baseUrl/api/notifications/trigger-card-due-job" -Method Post
    Write-Host "✓ Job disparado exitosamente!" -ForegroundColor Green
    Write-Host "  Mensaje: $($cardDueResponse.message)" -ForegroundColor Gray
    Write-Host "  Async: $($cardDueResponse.async)" -ForegroundColor Gray
    Write-Host "  Timestamp: $($cardDueResponse.timestamp)" -ForegroundColor Gray
    Write-Host ""
    Write-Host "  ℹ️  El job se ejecuta de forma asíncrona. Espera unos segundos..." -ForegroundColor Cyan
} catch {
    Write-Host "✗ Error al disparar job: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Esperar un poco para que el job se ejecute
Write-Host "  ⏳ Esperando 5 segundos para que el job se complete..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 4. Consultar historial de jobs
Write-Host "[4/4] Consultando historial de ejecuciones..." -ForegroundColor Yellow
try {
    $history = Invoke-RestMethod -Uri "$baseUrl/api/notifications/job-history?limit=5" -Method Get
    Write-Host "✓ Historial obtenido!" -ForegroundColor Green
    Write-Host "  Total de registros: $($history.count)" -ForegroundColor Gray
    Write-Host ""
    
    if ($history.data.Count -gt 0) {
        Write-Host "  Últimas ejecuciones:" -ForegroundColor Cyan
        foreach ($job in $history.data) {
            Write-Host "  ----------------------------------------" -ForegroundColor Gray
            Write-Host "  Run ID: $($job.run_id)" -ForegroundColor White
            Write-Host "  Status: $($job.status)" -ForegroundColor $(if ($job.status -eq "completed") { "Green" } else { "Yellow" })
            Write-Host "  Inicio: $($job.started_at)" -ForegroundColor Gray
            if ($job.completed_at) {
                Write-Host "  Fin: $($job.completed_at)" -ForegroundColor Gray
            }
            if ($job.duration) {
                Write-Host "  Duración: $($job.duration)" -ForegroundColor Gray
            }
            Write-Host "  Tarjetas encontradas: $($job.cards_found)" -ForegroundColor Gray
            Write-Host "  Emails enviados: $($job.emails_sent)" -ForegroundColor $(if ($job.emails_sent -gt 0) { "Green" } else { "Yellow" })
            Write-Host "  Errores: $($job.errors)" -ForegroundColor $(if ($job.errors -gt 0) { "Red" } else { "Green" })
            if ($job.error_message) {
                Write-Host "  Error: $($job.error_message)" -ForegroundColor Red
            }
        }
    } else {
        Write-Host "  No hay registros en el historial" -ForegroundColor Yellow
    }
} catch {
    Write-Host "✗ Error al consultar historial: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 5. Scheduler Status
Write-Host "[Bonus] Estado del scheduler..." -ForegroundColor Yellow
try {
    $schedulerStatus = Invoke-RestMethod -Uri "$baseUrl/api/notifications/scheduler/status" -Method Get
    Write-Host "✓ Scheduler Status:" -ForegroundColor Green
    Write-Host "  Running: $($schedulerStatus.scheduler_running)" -ForegroundColor Gray
    Write-Host "  Status: $($schedulerStatus.status)" -ForegroundColor Gray
    if ($schedulerStatus.next_scheduled_run) {
        Write-Host "  Próxima ejecución: $($schedulerStatus.next_scheduled_run)" -ForegroundColor Gray
        Write-Host "  Tiempo restante: $($schedulerStatus.time_to_next_run)" -ForegroundColor Gray
    }
} catch {
    Write-Host "✗ Error al consultar scheduler: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "           Test Completado!" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "📋 Resumen de endpoints disponibles:" -ForegroundColor Yellow
Write-Host ""
Write-Host "  POST /api/notifications/support" -ForegroundColor White
Write-Host "       → Envía email de soporte/contacto" -ForegroundColor Gray
Write-Host ""
Write-Host "  POST /api/notifications/trigger-card-due-job" -ForegroundColor White
Write-Host "       → Dispara job de notificaciones de tarjetas por vencer" -ForegroundColor Gray
Write-Host ""
Write-Host "  POST /api/notifications/trigger-update-due-dates-job" -ForegroundColor White
Write-Host "       → Actualiza fechas de vencimiento vencidas" -ForegroundColor Gray
Write-Host ""
Write-Host "  GET /api/notifications/job-history?limit=10" -ForegroundColor White
Write-Host "       → Consulta historial de ejecuciones de jobs" -ForegroundColor Gray
Write-Host ""
Write-Host "  GET /api/notifications/logs?job_run_id=xxx&limit=50" -ForegroundColor White
Write-Host "       → Consulta logs de notificaciones de un job específico" -ForegroundColor Gray
Write-Host ""
Write-Host "  GET /api/notifications/scheduler/status" -ForegroundColor White
Write-Host "       → Estado del job scheduler" -ForegroundColor Gray
Write-Host ""
Write-Host "  GET /health" -ForegroundColor White
Write-Host "       → Health check del servicio" -ForegroundColor Gray
Write-Host ""
