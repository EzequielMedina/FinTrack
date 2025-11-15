# Test directo del endpoint de Transaction Service
# Simula exactamente la llamada que hace Account Service

$headers = @{
    "Content-Type" = "application/json"
    "X-User-ID" = "018d9fa4-ba1d-703b-bd46-38e867b04e05"  # Tu userID
}

$body = @{
    type = "installment_payment"
    amount = 13333.33
    currency = "ARS"
    fromAccountId = "019447f2-95e1-7cba-b7b2-a087d0dc8b93"  # Bank account
    description = "Installment payment via bank transfer"
    paymentMethod = "bank_transfer"
    referenceId = "test-installment-payment"
    metadata = @{
        cardId = "019447f3-4e2c-7be2-9ee4-55da7c3ec86c"
        installmentId = "6c0ad793-5819-4784-b48c-1ce5201d10f2"
        installmentPlanId = "ef3f9bd6-4f4e-4fa6-ad61-f5c41b14c28b"
        installmentNumber = "1"
        category = "installment_payment"
        recordOnly = $true  # Changed to true for testing - Account Service already updated the balance
    }
} | ConvertTo-Json -Depth 10

Write-Host "Testing Transaction Service directly..."
Write-Host "URL: http://localhost:8083/api/v1/transactions"
Write-Host ""
Write-Host "Headers:" -ForegroundColor Cyan
$headers | ConvertTo-Json
Write-Host ""
Write-Host "Body:" -ForegroundColor Cyan
Write-Host $body
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8083/api/v1/transactions" `
        -Method POST `
        -Headers $headers `
        -Body $body `
        -ContentType "application/json"
    
    Write-Host "Success!" -ForegroundColor Green
    Write-Host ""
    $response | ConvertTo-Json -Depth 10
}
catch {
    Write-Host "Error!" -ForegroundColor Red
    Write-Host "Status Code: $($_.Exception.Response.StatusCode.value__)"
    Write-Host "Status Description: $($_.Exception.Response.StatusDescription)"
    Write-Host ""
    
    # Try to read error response
    $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
    $errorBody = $reader.ReadToEnd()
    Write-Host "Error Response:" -ForegroundColor Yellow
    Write-Host $errorBody
}
