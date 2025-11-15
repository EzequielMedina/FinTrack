# Test del Sistema de Soporte FAQ
# Este script prueba el endpoint de soporte del microservicio de notificaciones

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   TEST: Sistema de Soporte FAQ" -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

# Configuración
$notificationServiceUrl = "http://localhost:8084/api/notifications/support"

# Datos de prueba
$testData = @{
    name = "Juan Pérez"
    email = "juan.perez@test.com"
    subject = "Consulta sobre tarjetas de crédito"
    message = "Hola equipo de FinTrack, tengo una consulta sobre cómo agregar una nueva tarjeta de crédito. He intentado desde la sección de tarjetas pero no encuentro el botón. ¿Me pueden ayudar? Gracias!"
} | ConvertTo-Json

Write-Host "🔍 Verificando que el microservicio esté corriendo..." -ForegroundColor White
try {
    $healthCheck = Invoke-RestMethod -Uri "http://localhost:8084/health" -Method Get -ErrorAction Stop
    Write-Host "✅ Notification-Service: " -ForegroundColor Green -NoNewline
    Write-Host "$($healthCheck.status) - v$($healthCheck.version)" -ForegroundColor White
    Write-Host ""
} catch {
    Write-Host "❌ Error: El notification-service no está corriendo" -ForegroundColor Red
    Write-Host "   Ejecuta: docker-compose up notification-service" -ForegroundColor Yellow
    exit 1
}

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   Enviando email de soporte de prueba..." -ForegroundColor Yellow
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "📧 Datos del email:" -ForegroundColor Cyan
Write-Host "   Nombre:  Juan Pérez" -ForegroundColor White
Write-Host "   Email:   juan.perez@test.com" -ForegroundColor White
Write-Host "   Asunto:  Consulta sobre tarjetas de crédito" -ForegroundColor White
Write-Host "   Destino: soporte@fintrack.com" -ForegroundColor White
Write-Host ""

Write-Host "⏳ Enviando solicitud a: $notificationServiceUrl" -ForegroundColor Gray
Write-Host ""

try {
    # Enviar solicitud
    $response = Invoke-RestMethod `
        -Uri $notificationServiceUrl `
        -Method Post `
        -Body $testData `
        -ContentType "application/json" `
        -ErrorAction Stop

    Write-Host "================================================" -ForegroundColor Green
    Write-Host "   ✅ EMAIL ENVIADO EXITOSAMENTE" -ForegroundColor Green
    Write-Host "================================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "📨 Respuesta del servidor:" -ForegroundColor Cyan
    Write-Host "   Mensaje:   $($response.message)" -ForegroundColor White
    Write-Host "   Timestamp: $($response.timestamp)" -ForegroundColor White
    Write-Host ""
    Write-Host "✨ Verifica tu bandeja de entrada en:" -ForegroundColor Yellow
    Write-Host "   soporte@fintrack.com" -ForegroundColor White
    Write-Host ""
    Write-Host "💡 El email debería tener:" -ForegroundColor Cyan
    Write-Host "   - From: Juan Pérez" -ForegroundColor Gray
    Write-Host "   - Reply-To: juan.perez@test.com" -ForegroundColor Gray
    Write-Host "   - Subject: Consulta sobre tarjetas de crédito" -ForegroundColor Gray
    Write-Host "   - Template: template_yst8bd2" -ForegroundColor Gray
    Write-Host ""

} catch {
    Write-Host "================================================" -ForegroundColor Red
    Write-Host "   ❌ ERROR AL ENVIAR EMAIL" -ForegroundColor Red
    Write-Host "================================================" -ForegroundColor Red
    Write-Host ""
    
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode.value__
        Write-Host "   Status Code: $statusCode" -ForegroundColor Yellow
        
        try {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $errorBody = $reader.ReadToEnd() | ConvertFrom-Json
            Write-Host "   Error: $($errorBody.error)" -ForegroundColor Red
            Write-Host "   Detalles: $($errorBody.details)" -ForegroundColor Yellow
        } catch {
            Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
        }
    } else {
        Write-Host "   Error: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "🔧 Posibles soluciones:" -ForegroundColor Cyan
    Write-Host "   1. Verifica que notification-service esté corriendo" -ForegroundColor White
    Write-Host "   2. Verifica las credenciales de EmailJS en el backend" -ForegroundColor White
    Write-Host "   3. Verifica que el template template_yst8bd2 exista" -ForegroundColor White
    Write-Host "   4. Revisa los logs del microservicio" -ForegroundColor White
    Write-Host ""
    
    exit 1
}

Write-Host "================================================" -ForegroundColor Cyan
Write-Host "   TEST COMPLETADO" -ForegroundColor Green
Write-Host "================================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "💡 Próximos pasos:" -ForegroundColor Yellow
Write-Host "   1. Accede a http://localhost:4200/faq" -ForegroundColor White
Write-Host "   2. Scroll hasta el final" -ForegroundColor White
Write-Host "   3. Click en 'Contacta con Soporte'" -ForegroundColor White
Write-Host "   4. Completa y envía el formulario" -ForegroundColor White
Write-Host ""
