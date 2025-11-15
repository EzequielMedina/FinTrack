# Test EmailJS siguiendo el ejemplo oficial de JS (SIN accessToken)
Write-Host "Testing EmailJS following official JS example (NO accessToken)..." -ForegroundColor Yellow

$body = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39" 
    user_id = "MSBb87-PQcXWr1gWK"
    # NO accessToken - siguiendo el ejemplo oficial
    template_params = @{
        username = "Ezequiel"
        to_email = "eze.11.94.em@gmail.com"
        from_name = "FinTrack"
        message = "Test following official example"
    }
}

$jsonBody = $body | ConvertTo-Json -Depth 10

Write-Host "Request Body (NO accessToken):" -ForegroundColor Blue
Write-Host $jsonBody -ForegroundColor Gray

try {
    $response = Invoke-RestMethod -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $jsonBody -ContentType "application/json" -ErrorAction Stop
    Write-Host "SUCCESS: $response" -ForegroundColor Green
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    $statusDescription = $_.Exception.Response.StatusDescription
    
    Write-Host "ERROR Details:" -ForegroundColor Red
    Write-Host "  Status Code: $statusCode" -ForegroundColor Red
    Write-Host "  Status Description: $statusDescription" -ForegroundColor Red
    
    try {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Response Body: $errorBody" -ForegroundColor Red
    } catch {
        Write-Host "  Could not read response body" -ForegroundColor Red
    }
}