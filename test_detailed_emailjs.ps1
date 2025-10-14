# Test detallado de EmailJS con manejo de errores completo
Write-Host "Testing EmailJS with detailed error handling..." -ForegroundColor Yellow

$headers = @{
    "Content-Type" = "application/json"
    "User-Agent" = "FinTrack-NotificationService/1.0"
}

$body = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39" 
    user_id = "MSBb87-PQcXWr1gWK"
    accessToken = "sXLmpEZ8y2EYtCDtN5gZv"
    template_params = @{
        to_email = "eze.11.94.em@gmail.com"
        from_name = "FinTrack"
        message = "Test detailed"
    }
}

$jsonBody = $body | ConvertTo-Json -Depth 10

Write-Host "Request Body:" -ForegroundColor Blue
Write-Host $jsonBody -ForegroundColor Gray

try {
    $response = Invoke-RestMethod -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $jsonBody -Headers $headers -ErrorAction Stop
    Write-Host "SUCCESS: $response" -ForegroundColor Green
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    $statusDescription = $_.Exception.Response.StatusDescription
    
    Write-Host "ERROR Details:" -ForegroundColor Red
    Write-Host "  Status Code: $statusCode" -ForegroundColor Red
    Write-Host "  Status Description: $statusDescription" -ForegroundColor Red
    Write-Host "  Exception Message: $($_.Exception.Message)" -ForegroundColor Red
    
    try {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Response Body: $errorBody" -ForegroundColor Red
    } catch {
        Write-Host "  Could not read response body" -ForegroundColor Red
    }
}