# Test directo de EmailJS REST API
$ErrorActionPreference = "Stop"

Write-Host "🧪 Probando EmailJS REST API directamente..." -ForegroundColor Yellow

# Credenciales desde el notification-service
$emailJSData = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39"
    user_id = "MSBb87-PQcXWr1gWK"
    accessToken = "sXLmpEZ8y2EYtCDtN5gZv"
    template_params = @{
        from_name = "FinTrack"
        to_email = "eze.11.94.em@gmail.com"
        subject = "Test EmailJS desde PowerShell"
        user_name = "Ezequiel"
        card_name = "TEST"
        bank_name = "TEST BANK"
        last_four = "1234"
        due_date = "15/10/2025"
        total_amount = "`$100.00"
        installment_count = "2"
        html_content = "<h1>Test email desde PowerShell</h1><p>Este es un test directo.</p>"
    }
}

try {
    $jsonBody = $emailJSData | ConvertTo-Json -Depth 10
    Write-Host "📨 Enviando request a EmailJS..." -ForegroundColor Blue
    Write-Host "Body: $jsonBody" -ForegroundColor Gray
    
    $response = Invoke-RestMethod -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $jsonBody -ContentType "application/json"
    
    Write-Host "✅ Respuesta de EmailJS: $response" -ForegroundColor Green
    Write-Host "✅ Email enviado exitosamente!" -ForegroundColor Green
}
catch {
    Write-Host "❌ Error al enviar email:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody" -ForegroundColor Red
    }
}