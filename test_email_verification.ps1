# Test EmailJS con template simple para verificar recepción
Write-Host "Testing EmailJS with simple template..." -ForegroundColor Yellow

$body = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39" 
    user_id = "MSBb87-PQcXWr1gWK"
    template_params = @{
        to_email = "eze.11.94.em@gmail.com"
        from_name = "FinTrack Test"
        subject = "TEST - ¿Llega este email?"
        user_name = "Ezequiel"
        message = "Este es un email de prueba simple"
        html_content = "<h1>TEST EMAIL</h1><p>Este es un test para verificar que el email llega correctamente.</p>"
    }
}

$headers = @{
    "Content-Type" = "application/json"
    "User-Agent" = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
    "Origin" = "http://localhost:4200"
    "Referer" = "http://localhost:4200/"
}

$jsonBody = $body | ConvertTo-Json -Depth 10

Write-Host "Sending test email..." -ForegroundColor Blue

try {
    $response = Invoke-RestMethod -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $jsonBody -Headers $headers -ErrorAction Stop
    Write-Host "SUCCESS: $response" -ForegroundColor Green
    Write-Host "✅ Email enviado - revisa tu bandeja de entrada Y carpeta de SPAM" -ForegroundColor Green
} catch {
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
}