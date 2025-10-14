# Test EmailJS simulando browser headers
Write-Host "Testing EmailJS with browser simulation headers..." -ForegroundColor Yellow

$headers = @{
    "Content-Type" = "application/json"
    "User-Agent" = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
    "Accept" = "application/json, text/plain, */*"
    "Accept-Language" = "es-ES,es;q=0.9,en;q=0.8"
    "Origin" = "http://localhost:4200"
    "Referer" = "http://localhost:4200/"
    "Sec-Fetch-Dest" = "empty"
    "Sec-Fetch-Mode" = "cors"
    "Sec-Fetch-Site" = "cross-site"
}

$body = @{
    service_id = "service_ceg7xlp"
    template_id = "template_e43va39" 
    user_id = "MSBb87-PQcXWr1gWK"
    template_params = @{
        username = "Ezequiel"
        to_email = "eze.11.94.em@gmail.com"
        from_name = "FinTrack"
        message = "Test with browser headers"
    }
}

$jsonBody = $body | ConvertTo-Json -Depth 10

Write-Host "Request with browser headers:" -ForegroundColor Blue
Write-Host $jsonBody -ForegroundColor Gray

try {
    $response = Invoke-RestMethod -Uri "https://api.emailjs.com/api/v1.0/email/send" -Method POST -Body $jsonBody -Headers $headers -ErrorAction Stop
    Write-Host "SUCCESS: $response" -ForegroundColor Green
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    
    Write-Host "ERROR Details:" -ForegroundColor Red
    Write-Host "  Status Code: $statusCode" -ForegroundColor Red
    
    try {
        $errorStream = $_.Exception.Response.GetResponseStream()
        $reader = New-Object System.IO.StreamReader($errorStream)
        $errorBody = $reader.ReadToEnd()
        Write-Host "  Response Body: $errorBody" -ForegroundColor Red
    } catch {
        Write-Host "  Could not read response body" -ForegroundColor Red
    }
}