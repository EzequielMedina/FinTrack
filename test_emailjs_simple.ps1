# Test simple de EmailJS
Write-Host "Testing EmailJS API..." -ForegroundColor Yellow

$body = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39"
    user_id = "MSBb87-PQcXWr1gWK"
    accessToken = "sXLmpEZ8y2EYtCDtN5gZv"
    template_params = @{
        from_name = "FinTrack"
        to_email = "eze.11.94.em@gmail.com"
        subject = "Test EmailJS"
        message = "Test message"
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $body -ContentType "application/json" -ErrorAction Stop
    Write-Host "Success: $response" -ForegroundColor Green
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Details: $($_.ErrorDetails.Message)" -ForegroundColor Red
}