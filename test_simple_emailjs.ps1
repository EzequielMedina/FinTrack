# Test EmailJS simple
Write-Host "Testing EmailJS minimal..." -ForegroundColor Yellow

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

Write-Host "Sending request..." -ForegroundColor Blue
Write-Host $body -ForegroundColor Gray

try {
    $response = Invoke-WebRequest -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $body -ContentType "application/json"
    Write-Host "SUCCESS Status: $($response.StatusCode)" -ForegroundColor Green
    Write-Host "Response: $($response.Content)" -ForegroundColor Green
} catch {
    Write-Host "ERROR StatusCode: $($_.Exception.Response.StatusCode.value__)" -ForegroundColor Red
    Write-Host "Error Message: $($_.Exception.Message)" -ForegroundColor Red
}