# Test para obtener cuentas del usuario

$headers = @{
    "X-User-ID" = "3736c174-fa96-4259-8669-392b903ae24e"
}

Write-Host "Testing Account Service - Get Accounts..." -ForegroundColor Cyan
Write-Host "URL: http://localhost:8082/api/accounts/user/3736c174-fa96-4259-8669-392b903ae24e" -ForegroundColor Yellow
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8082/api/accounts/user/3736c174-fa96-4259-8669-392b903ae24e" `
        -Method GET `
        -Headers $headers
    
    Write-Host "Success! Found $($response.accounts.Count) accounts" -ForegroundColor Green
    Write-Host ""
    
    if ($response.accounts.Count -gt 0) {
        Write-Host "Accounts:" -ForegroundColor Cyan
        $response.accounts | ForEach-Object {
            Write-Host "  - $($_.name) ($($_.account_type)): $($_.currency) $($_.balance)" -ForegroundColor White
        }
    } else {
        Write-Host "No accounts found. You need to create accounts first." -ForegroundColor Yellow
    }
    
    Write-Host ""
    $response | ConvertTo-Json -Depth 10
}
catch {
    Write-Host "Error!" -ForegroundColor Red
    Write-Host "Status Code: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Error: $($_.Exception.Message)"
}
