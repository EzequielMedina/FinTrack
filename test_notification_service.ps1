# Script de testing para notification-service
param(
    [string]$Mode = "manual"  # "manual", "cron", "debug"
)

Write-Host "🧪 FinTrack Notification Service - Script de Testing" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Green

$servicePath = "c:\Facultad\Alumno\PS\backend\services\notification-service"
$projectRoot = "c:\Facultad\Alumno\PS"

switch ($Mode) {
    "manual" {
        Write-Host "🔧 Modo: Trigger Manual" -ForegroundColor Yellow
        
        # Verificar si el servicio está corriendo
        Write-Host "📡 Verificando si el servicio está activo..."
        try {
            $response = Invoke-RestMethod -Uri "http://localhost:8088/health" -Method GET -TimeoutSec 5
            Write-Host "✅ Servicio activo: $($response.status)" -ForegroundColor Green
        }
        catch {
            Write-Host "❌ Servicio no disponible. Iniciando..." -ForegroundColor Red
            Write-Host "💡 Ejecuta: docker-compose up --build mysql notification-service" -ForegroundColor Cyan
            return
        }
        
        # Trigger manual del job
        Write-Host "🚀 Ejecutando job manualmente..."
        try {
            $triggerResponse = Invoke-RestMethod -Uri "http://localhost:8088/api/notifications/trigger-card-due-job" -Method POST -TimeoutSec 10
            Write-Host "✅ Job disparado: $($triggerResponse.message)" -ForegroundColor Green
        }
        catch {
            Write-Host "❌ Error disparando job: $($_.Exception.Message)" -ForegroundColor Red
            return
        }
        
        # Esperar un poco y verificar resultados
        Write-Host "⏳ Esperando resultados (5 segundos)..."
        Start-Sleep -Seconds 5
        
        # Ver historial
        Write-Host "📊 Obteniendo historial de jobs..."
        try {
            $history = Invoke-RestMethod -Uri "http://localhost:8088/api/notifications/job-history?limit=3" -Method GET
            foreach ($job in $history.data) {
                Write-Host "   📋 Job ID: $($job.run_id)" -ForegroundColor Cyan
                Write-Host "   📅 Iniciado: $($job.started_at)" -ForegroundColor White
                Write-Host "   📈 Estado: $($job.status)" -ForegroundColor $(if ($job.status -eq "completed") { "Green" } else { "Red" })
                Write-Host "   💳 Tarjetas encontradas: $($job.cards_found)" -ForegroundColor White
                Write-Host "   📧 Emails enviados: $($job.emails_sent)" -ForegroundColor White
                Write-Host "   ❌ Errores: $($job.errors)" -ForegroundColor $(if ($job.errors -gt 0) { "Red" } else { "Green" })
                Write-Host "   ⏱️  Duración: $($job.duration)" -ForegroundColor White
                Write-Host ""
            }
        }
        catch {
            Write-Host "❌ Error obteniendo historial: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    "cron" {
        Write-Host "⏰ Modo: Testing con Cron cada minuto" -ForegroundColor Yellow
        Write-Host "📝 Modificando configuración para ejecutar cada minuto..." -ForegroundColor Cyan
        
        # Copiar configuración de testing
        Copy-Item "$servicePath\.env.testing" "$servicePath\.env" -Force
        Write-Host "✅ Configuración actualizada (ejecutará cada minuto)" -ForegroundColor Green
        
        Write-Host "🔄 Reinicia el servicio con:" -ForegroundColor Cyan
        Write-Host "   docker-compose down notification-service" -ForegroundColor White
        Write-Host "   docker-compose up --build notification-service" -ForegroundColor White
        Write-Host ""
        Write-Host "📊 Luego monitorea con:" -ForegroundColor Cyan
        Write-Host "   .\test_notification_service.ps1 monitor" -ForegroundColor White
    }
    
    "debug" {
        Write-Host "🐛 Modo: Debug con logs detallados" -ForegroundColor Yellow
        
        # Verificar logs del contenedor
        Write-Host "📋 Últimos logs del contenedor..."
        docker logs fintrack-notification-service --tail 20
        
        # Verificar estado del scheduler
        Write-Host "⏰ Estado del scheduler..."
        try {
            $scheduler = Invoke-RestMethod -Uri "http://localhost:8088/api/notifications/scheduler/status" -Method GET
            Write-Host "   🔄 Scheduler corriendo: $($scheduler.scheduler_running)" -ForegroundColor $(if ($scheduler.scheduler_running) { "Green" } else { "Red" })
            Write-Host "   📅 Próxima ejecución: $($scheduler.next_scheduled_run)" -ForegroundColor White
            Write-Host "   ⏳ Tiempo hasta próxima: $($scheduler.time_to_next_run)" -ForegroundColor White
        }
        catch {
            Write-Host "❌ Error obteniendo estado del scheduler: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    
    "monitor" {
        Write-Host "📊 Modo: Monitoreo continuo" -ForegroundColor Yellow
        Write-Host "⌚ Presiona Ctrl+C para detener" -ForegroundColor Cyan
        
        while ($true) {
            Clear-Host
            Write-Host "🔄 Monitoreo Notification Service - $(Get-Date -Format 'HH:mm:ss')" -ForegroundColor Green
            Write-Host "================================================" -ForegroundColor Green
            
            # Health check
            try {
                $health = Invoke-RestMethod -Uri "http://localhost:8088/health" -Method GET -TimeoutSec 3
                Write-Host "✅ Servicio: $($health.status) | Scheduler: $($health.job_scheduler_running)" -ForegroundColor Green
                if ($health.next_scheduled_run) {
                    Write-Host "⏰ Próxima ejecución: $($health.next_scheduled_run)" -ForegroundColor Cyan
                }
            }
            catch {
                Write-Host "❌ Servicio no disponible" -ForegroundColor Red
            }
            
            # Último job
            try {
                $history = Invoke-RestMethod -Uri "http://localhost:8088/api/notifications/job-history?limit=1" -Method GET -TimeoutSec 3
                if ($history.data -and $history.data.Count -gt 0) {
                    $lastJob = $history.data[0]
                    Write-Host ""
                    Write-Host "📋 Último Job:" -ForegroundColor Yellow
                    Write-Host "   ID: $($lastJob.run_id)" -ForegroundColor White
                    Write-Host "   Estado: $($lastJob.status)" -ForegroundColor $(if ($lastJob.status -eq "completed") { "Green" } else { "Red" })
                    Write-Host "   Tarjetas: $($lastJob.cards_found) | Emails: $($lastJob.emails_sent) | Errores: $($lastJob.errors)" -ForegroundColor White
                    Write-Host "   Iniciado: $($lastJob.started_at)" -ForegroundColor Gray
                }
            }
            catch {
                Write-Host "❌ Error obteniendo historial" -ForegroundColor Red
            }
            
            Start-Sleep -Seconds 5
        }
    }
    
    default {
        Write-Host "❓ Uso:" -ForegroundColor Yellow
        Write-Host "   .\test_notification_service.ps1 manual    # Trigger manual" -ForegroundColor White
        Write-Host "   .\test_notification_service.ps1 cron      # Config para cada minuto" -ForegroundColor White
        Write-Host "   .\test_notification_service.ps1 debug     # Logs y debugging" -ForegroundColor White
        Write-Host "   .\test_notification_service.ps1 monitor   # Monitoreo continuo" -ForegroundColor White
    }
}

Write-Host ""
Write-Host "🔗 URLs útiles:" -ForegroundColor Cyan
Write-Host "   Health: http://localhost:8088/health" -ForegroundColor White
Write-Host "   API: http://localhost:8088/api/notifications/" -ForegroundColor White
Write-Host "   Trigger: POST http://localhost:8088/api/notifications/trigger-card-due-job" -ForegroundColor White