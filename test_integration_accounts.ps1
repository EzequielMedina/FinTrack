# Test de Integración Frontend-Backend: Gestión de Cuentas
# Este script valida el flujo completo de operaciones de cuentas

param(
    [string]$BackendUrl = "http://localhost:8080",
    [string]$TestUserId = "test-user-123",
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"

# Colores para output
$Green = "Green"
$Red = "Red"
$Yellow = "Yellow"
$Blue = "Blue"

function Write-TestResult {
    param(
        [string]$TestName,
        [bool]$Success,
        [string]$Message = ""
    )
    
    if ($Success) {
        Write-Host "✅ $TestName" -ForegroundColor $Green
        if ($Verbose -and $Message) {
            Write-Host "   → $Message" -ForegroundColor Gray
        }
    } else {
        Write-Host "❌ $TestName" -ForegroundColor $Red
        if ($Message) {
            Write-Host "   → $Message" -ForegroundColor $Red
        }
    }
}

function Write-TestSection {
    param([string]$Section)
    Write-Host "`n🔍 $Section" -ForegroundColor $Blue
    Write-Host "=" * ($Section.Length + 3) -ForegroundColor $Blue
}

function Test-BackendEndpoint {
    param(
        [string]$Url,
        [string]$Method = "GET",
        [hashtable]$Headers = @{},
        [string]$Body = $null
    )
    
    try {
        $params = @{
            Uri = $Url
            Method = $Method
            Headers = $Headers
            UseBasicParsing = $true
        }
        
        if ($Body) {
            $params.Body = $Body
            $params.ContentType = "application/json"
        }
        
        $response = Invoke-WebRequest @params
        return @{
            Success = $true
            StatusCode = $response.StatusCode
            Content = $response.Content
            Response = $response
        }
    }
    catch {
        return @{
            Success = $false
            Error = $_.Exception.Message
            StatusCode = if ($_.Exception.Response) { $_.Exception.Response.StatusCode.value__ } else { 0 }
        }
    }
}

# Variables globales para el test
$script:CreatedAccountId = $null
$script:AuthToken = $null
$script:TestResults = @{
    Total = 0
    Passed = 0
    Failed = 0
}

function Add-TestResult {
    param([bool]$Success)
    $script:TestResults.Total++
    if ($Success) {
        $script:TestResults.Passed++
    } else {
        $script:TestResults.Failed++
    }
}

Write-Host "🚀 Iniciando Tests de Integración Frontend-Backend" -ForegroundColor $Blue
Write-Host "Backend URL: $BackendUrl" -ForegroundColor Gray
Write-Host "Test User ID: $TestUserId" -ForegroundColor Gray

# Test 1: Verificar que el backend esté disponible
Write-TestSection "1. Conectividad del Backend"

$healthCheck = Test-BackendEndpoint -Url "$BackendUrl/health"
$healthSuccess = $healthCheck.Success -and $healthCheck.StatusCode -eq 200
Write-TestResult "Health Check del Backend" $healthSuccess $healthCheck.Error
Add-TestResult $healthSuccess

# Test 2: Verificar endpoints de cuentas
Write-TestSection "2. Endpoints de API de Cuentas"

$accountsEndpoint = Test-BackendEndpoint -Url "$BackendUrl/api/accounts"
$accountsSuccess = $accountsEndpoint.Success
Write-TestResult "Endpoint GET /api/accounts" $accountsSuccess $accountsEndpoint.Error
Add-TestResult $accountsSuccess

# Test 3: Crear cuenta de prueba
Write-TestSection "3. Operaciones CRUD de Cuentas"

$createAccountPayload = @{
    userId = $TestUserId
    accountType = "SAVINGS"
    name = "Cuenta de Prueba Integration Test"
    currency = "ARS"
    dni = "12345678"
    initialBalance = 1000.50
    description = "Cuenta creada por test de integración"
    isActive = $true
} | ConvertTo-Json

$createResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts" -Method "POST" -Body $createAccountPayload
$createSuccess = $createResult.Success -and $createResult.StatusCode -eq 201

if ($createSuccess) {
    try {
        $createdAccount = $createResult.Content | ConvertFrom-Json
        $script:CreatedAccountId = $createdAccount.id
        Write-TestResult "Crear cuenta (POST)" $true "Cuenta creada con ID: $($script:CreatedAccountId)"
    }
    catch {
        Write-TestResult "Crear cuenta (POST)" $false "Error parseando respuesta: $($_.Exception.Message)"
        $createSuccess = $false
    }
} else {
    Write-TestResult "Crear cuenta (POST)" $false $createResult.Error
}
Add-TestResult $createSuccess

# Test 4: Leer cuenta creada
if ($script:CreatedAccountId) {
    $getAccountResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$($script:CreatedAccountId)"
    $getSuccess = $getAccountResult.Success -and $getAccountResult.StatusCode -eq 200
    
    if ($getSuccess) {
        try {
            $account = $getAccountResult.Content | ConvertFrom-Json
            $dataValid = $account.name -eq "Cuenta de Prueba Integration Test" -and $account.balance -eq 1000.50
            Write-TestResult "Leer cuenta (GET)" $dataValid "Balance: $($account.balance), Nombre: $($account.name)"
        }
        catch {
            Write-TestResult "Leer cuenta (GET)" $false "Error parseando respuesta"
            $getSuccess = $false
        }
    } else {
        Write-TestResult "Leer cuenta (GET)" $false $getAccountResult.Error
    }
    Add-TestResult $getSuccess
}

# Test 5: Actualizar cuenta
if ($script:CreatedAccountId) {
    $updatePayload = @{
        name = "Cuenta Actualizada Integration Test"
        description = "Descripción actualizada por test"
    } | ConvertTo-Json
    
    $updateResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$($script:CreatedAccountId)" -Method "PUT" -Body $updatePayload
    $updateSuccess = $updateResult.Success -and $updateResult.StatusCode -eq 200
    Write-TestResult "Actualizar cuenta (PUT)" $updateSuccess $updateResult.Error
    Add-TestResult $updateSuccess
}

# Test 6: Operaciones de Wallet
Write-TestSection "4. Operaciones de Wallet"

if ($script:CreatedAccountId) {
    # Test agregar fondos
    $addFundsPayload = @{
        amount = 500.25
        description = "Depósito de prueba"
        reference = "TEST-REF-001"
    } | ConvertTo-Json
    
    $addFundsResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$($script:CreatedAccountId)/add-funds" -Method "POST" -Body $addFundsPayload
    $addFundsSuccess = $addFundsResult.Success -and $addFundsResult.StatusCode -eq 200
    Write-TestResult "Agregar fondos" $addFundsSuccess $addFundsResult.Error
    Add-TestResult $addFundsSuccess
    
    # Test retirar fondos
    $withdrawFundsPayload = @{
        amount = 200.00
        description = "Retiro de prueba"
        reference = "TEST-REF-002"
    } | ConvertTo-Json
    
    $withdrawResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$($script:CreatedAccountId)/withdraw-funds" -Method "POST" -Body $withdrawFundsPayload
    $withdrawSuccess = $withdrawResult.Success -and $withdrawResult.StatusCode -eq 200
    Write-TestResult "Retirar fondos" $withdrawSuccess $withdrawResult.Error
    Add-TestResult $withdrawSuccess
}

# Test 7: Crear tarjeta de crédito
Write-TestSection "5. Operaciones de Tarjeta de Crédito"

$creditCardPayload = @{
    userId = $TestUserId
    accountType = "CREDIT_CARD"
    name = "Tarjeta de Crédito Test"
    currency = "ARS"
    dni = "12345678"
    creditLimit = 50000.00
    closingDate = "2025-01-15"
    dueDate = "2025-02-05"
    isActive = $true
} | ConvertTo-Json

$creditCardResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts" -Method "POST" -Body $creditCardPayload
$creditCardSuccess = $creditCardResult.Success -and $creditCardResult.StatusCode -eq 201

if ($creditCardSuccess) {
    try {
        $creditCard = $creditCardResult.Content | ConvertFrom-Json
        $creditCardId = $creditCard.id
        Write-TestResult "Crear tarjeta de crédito" $true "Tarjeta creada con ID: $creditCardId"
        
        # Test actualizar límite de crédito
        $updateCreditLimitPayload = @{
            creditLimit = 75000.00
        } | ConvertTo-Json
        
        $updateCreditResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$creditCardId/credit-limit" -Method "PUT" -Body $updateCreditLimitPayload
        $updateCreditSuccess = $updateCreditResult.Success -and $updateCreditResult.StatusCode -eq 200
        Write-TestResult "Actualizar límite de crédito" $updateCreditSuccess $updateCreditResult.Error
        Add-TestResult $updateCreditSuccess
        
        # Test obtener crédito disponible
        $availableCreditResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$creditCardId/available-credit"
        $availableCreditSuccess = $availableCreditResult.Success -and $availableCreditResult.StatusCode -eq 200
        Write-TestResult "Obtener crédito disponible" $availableCreditSuccess $availableCreditResult.Error
        Add-TestResult $availableCreditSuccess
        
    }
    catch {
        Write-TestResult "Crear tarjeta de crédito" $false "Error parseando respuesta: $($_.Exception.Message)"
        $creditCardSuccess = $false
    }
} else {
    Write-TestResult "Crear tarjeta de crédito" $false $creditCardResult.Error
}
Add-TestResult $creditCardSuccess

# Test 8: Listar cuentas por usuario
Write-TestSection "6. Consultas de Usuario"

$userAccountsResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/user/$TestUserId"
$userAccountsSuccess = $userAccountsResult.Success -and $userAccountsResult.StatusCode -eq 200

if ($userAccountsSuccess) {
    try {
        $userAccounts = $userAccountsResult.Content | ConvertFrom-Json
        $accountCount = if ($userAccounts -is [array]) { $userAccounts.Length } else { 1 }
        Write-TestResult "Listar cuentas por usuario" $true "Encontradas $accountCount cuentas"
    }
    catch {
        Write-TestResult "Listar cuentas por usuario" $false "Error parseando respuesta"
        $userAccountsSuccess = $false
    }
} else {
    Write-TestResult "Listar cuentas por usuario" $false $userAccountsResult.Error
}
Add-TestResult $userAccountsSuccess

# Test 9: Paginación
$paginationResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts?page=1`&pageSize=10"
$paginationSuccess = $paginationResult.Success -and $paginationResult.StatusCode -eq 200
Write-TestResult "Paginación de cuentas" $paginationSuccess $paginationResult.Error
Add-TestResult $paginationSuccess

# Cleanup: Eliminar cuentas de prueba
Write-TestSection "7. Limpieza (Cleanup)"

if ($script:CreatedAccountId) {
    $deleteResult = Test-BackendEndpoint -Url "$BackendUrl/api/accounts/$($script:CreatedAccountId)" -Method "DELETE"
    $deleteSuccess = $deleteResult.Success -and ($deleteResult.StatusCode -eq 200 -or $deleteResult.StatusCode -eq 204)
    Write-TestResult "Eliminar cuenta de prueba" $deleteSuccess $deleteResult.Error
    Add-TestResult $deleteSuccess
}

# Resumen final
Write-TestSection "📊 Resumen de Resultados"

$successRate = [math]::Round(($script:TestResults.Passed / $script:TestResults.Total) * 100, 2)

Write-Host "Total de pruebas: $($script:TestResults.Total)" -ForegroundColor Gray
Write-Host "Exitosas: $($script:TestResults.Passed)" -ForegroundColor $Green
Write-Host "Fallidas: $($script:TestResults.Failed)" -ForegroundColor $Red
Write-Host "Tasa de éxito: $successRate%" -ForegroundColor $(if ($successRate -ge 90) { $Green } elseif ($successRate -ge 70) { $Yellow } else { $Red })

if ($script:TestResults.Failed -eq 0) {
    Write-Host "`n🎉 ¡Todos los tests de integración pasaron exitosamente!" -ForegroundColor $Green
    exit 0
} else {
    Write-Host "`n⚠️  Algunos tests fallaron. Revisa la configuración del backend." -ForegroundColor $Yellow
    exit 1
}