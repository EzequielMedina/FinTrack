# Test EmailJS con template mínimo
Write-Host "🧪 Probando EmailJS con template mínimo..." -ForegroundColor Yellow

$body = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39" 
    user_id = "MSBb87-PQcXWr1gWK"
    accessToken = "sXLmpEZ8y2EYtCDtN5gZv"
    template_params = @{
        to_email = "eze.11.94.em@gmail.com"
        from_name = "FinTrack"
        message = "Test simple"
    }
} | ConvertTo-Json -Depth 10

Write-Host "📨 Enviando request..." -ForegroundColor Blue
Write-Host $body -ForegroundColor Gray

try {
    $response = Invoke-WebRequest -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $body -ContentType "application/json"
    Write-Host "✅ Status: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "✅ Response: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "❌ StatusCode: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $stream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($stream)
        $responseBody = $reader.ReadToEnd()
        Write-Host "❌ Response Body: $responseBody" -ForegroundColor Red
    }
}