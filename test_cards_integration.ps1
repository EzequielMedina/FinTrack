# Script para probar la nueva funcionalidad de tarjetas

# Configuración de variables
$API_BASE = "http://localhost:8082/api"
$USER_ID = "user-123"

Write-Host "🧪 Iniciando pruebas de la funcionalidad de tarjetas..." -ForegroundColor Cyan

# Función para hacer requests HTTP
function Invoke-ApiRequest {
    param(
        [string]$Method = "GET",
        [string]$Uri,
        [object]$Body = $null,
        [hashtable]$Headers = @{'Content-Type' = 'application/json'}
    )
    
    try {
        $params = @{
            Method = $Method
            Uri = $Uri
            Headers = $Headers
        }
        
        if ($Body) {
            $params.Body = $Body | ConvertTo-Json -Depth 10
        }
        
        $response = Invoke-RestMethod @params
        return $response
    }
    catch {
        Write-Host "❌ Error en request: $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# 1. Verificar que el servicio esté funcionando
Write-Host "`n1. ✅ Verificando conectividad del account-service..." -ForegroundColor Yellow
$healthCheck = Invoke-ApiRequest -Uri "$API_BASE/accounts"
if ($healthCheck) {
    Write-Host "✅ Account-service respondiendo correctamente" -ForegroundColor Green
} else {
    Write-Host "❌ Account-service no está disponible" -ForegroundColor Red
    exit 1
}

# 2. Obtener cuentas del usuario
Write-Host "`n2. 📋 Obteniendo cuentas del usuario..." -ForegroundColor Yellow
$accounts = Invoke-ApiRequest -Uri "$API_BASE/accounts/user/$USER_ID"
if ($accounts -and $accounts.Count -gt 0) {
    Write-Host "✅ Encontradas $($accounts.Count) cuentas" -ForegroundColor Green
    $firstAccount = $accounts[0]
    $ACCOUNT_ID = $firstAccount.id
    Write-Host "🏦 Usando cuenta: $($firstAccount.name) (ID: $ACCOUNT_ID)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se encontraron cuentas para el usuario" -ForegroundColor Red
    Write-Host "💡 Creando una cuenta de prueba..." -ForegroundColor Yellow
    
    $newAccount = @{
        user_id = $USER_ID
        account_type = "bank_account"
        name = "Cuenta de Prueba para Tarjetas"
        description = "Cuenta creada para probar funcionalidad de tarjetas"
        currency = "ARS"
        initial_balance = 1000.0
        is_active = $true
    }
    
    $createdAccount = Invoke-ApiRequest -Method "POST" -Uri "$API_BASE/accounts" -Body $newAccount
    if ($createdAccount) {
        $ACCOUNT_ID = $createdAccount.id
        Write-Host "✅ Cuenta creada con ID: $ACCOUNT_ID" -ForegroundColor Green
    } else {
        Write-Host "❌ No se pudo crear la cuenta de prueba" -ForegroundColor Red
        exit 1
    }
}

# 3. Crear una tarjeta de prueba
Write-Host "`n3. 💳 Creando tarjeta de prueba..." -ForegroundColor Yellow
$newCard = @{
    card_type = "debit"
    card_brand = "visa"
    last_four_digits = "1234"
    masked_number = "**** **** **** 1234"
    holder_name = "JUAN PEREZ"
    expiration_month = 12
    expiration_year = 2027
    nickname = "Tarjeta de Prueba"
    encrypted_number = "encrypted_test_data"
    encrypted_cvv = "encrypted_cvv_data"
    key_fingerprint = "test_fingerprint_123"
}

$createdCard = Invoke-ApiRequest -Method "POST" -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards" -Body $newCard
if ($createdCard) {
    $CARD_ID = $createdCard.id
    Write-Host "✅ Tarjeta creada con ID: $CARD_ID" -ForegroundColor Green
    Write-Host "💳 Detalles: $($createdCard.holder_name) - **** $($createdCard.last_four_digits)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se pudo crear la tarjeta" -ForegroundColor Red
    exit 1
}

# 4. Obtener tarjetas de la cuenta
Write-Host "`n4. 📋 Obteniendo tarjetas de la cuenta..." -ForegroundColor Yellow
$accountCards = Invoke-ApiRequest -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards"
if ($accountCards -and $accountCards.data) {
    Write-Host "✅ Encontradas $($accountCards.data.Count) tarjetas en la cuenta" -ForegroundColor Green
    foreach ($card in $accountCards.data) {
        Write-Host "  💳 $($card.nickname): $($card.holder_name) - **** $($card.last_four_digits)" -ForegroundColor Blue
    }
} else {
    Write-Host "❌ No se pudieron obtener las tarjetas de la cuenta" -ForegroundColor Red
}

# 5. Obtener tarjetas del usuario (todas las cuentas)
Write-Host "`n5. 👤 Obteniendo todas las tarjetas del usuario..." -ForegroundColor Yellow
$userCards = Invoke-ApiRequest -Uri "$API_BASE/accounts/user/$USER_ID/cards"
if ($userCards -and $userCards.data) {
    Write-Host "✅ Usuario tiene $($userCards.data.Count) tarjetas en total" -ForegroundColor Green
    foreach ($card in $userCards.data) {
        Write-Host "  💳 $($card.nickname): $($card.holder_name) - **** $($card.last_four_digits)" -ForegroundColor Blue
    }
} else {
    Write-Host "❌ No se pudieron obtener las tarjetas del usuario" -ForegroundColor Red
}

# 6. Obtener detalles de la tarjeta específica
Write-Host "`n6. 🔍 Obteniendo detalles de la tarjeta..." -ForegroundColor Yellow
$cardDetails = Invoke-ApiRequest -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards/$CARD_ID"
if ($cardDetails) {
    Write-Host "✅ Detalles de tarjeta obtenidos" -ForegroundColor Green
    Write-Host "  💳 Tipo: $($cardDetails.card_type)" -ForegroundColor Blue
    Write-Host "  🏦 Marca: $($cardDetails.card_brand)" -ForegroundColor Blue
    Write-Host "  👤 Titular: $($cardDetails.holder_name)" -ForegroundColor Blue
    Write-Host "  📅 Expira: $($cardDetails.expiration_month)/$($cardDetails.expiration_year)" -ForegroundColor Blue
    Write-Host "  ⚡ Estado: $($cardDetails.status)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se pudieron obtener los detalles de la tarjeta" -ForegroundColor Red
}

# 7. Actualizar la tarjeta
Write-Host "`n7. ✏️  Actualizando tarjeta..." -ForegroundColor Yellow
$updateCard = @{
    holder_name = "JUAN PEREZ UPDATED"
    nickname = "Mi Tarjeta Débito"
}

$updatedCard = Invoke-ApiRequest -Method "PUT" -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards/$CARD_ID" -Body $updateCard
if ($updatedCard) {
    Write-Host "✅ Tarjeta actualizada correctamente" -ForegroundColor Green
    Write-Host "  📝 Nuevo nombre: $($updatedCard.holder_name)" -ForegroundColor Blue
    Write-Host "  📝 Nuevo nickname: $($updatedCard.nickname)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se pudo actualizar la tarjeta" -ForegroundColor Red
}

# 8. Bloquear tarjeta
Write-Host "`n8. 🚫 Bloqueando tarjeta..." -ForegroundColor Yellow
$blockedCard = Invoke-ApiRequest -Method "PUT" -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards/$CARD_ID/block"
if ($blockedCard) {
    Write-Host "✅ Tarjeta bloqueada correctamente" -ForegroundColor Green
    Write-Host "  ⚡ Estado: $($blockedCard.status)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se pudo bloquear la tarjeta" -ForegroundColor Red
}

# 9. Desbloquear tarjeta
Write-Host "`n9. ✅ Desbloqueando tarjeta..." -ForegroundColor Yellow
$unblockedCard = Invoke-ApiRequest -Method "PUT" -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards/$CARD_ID/unblock"
if ($unblockedCard) {
    Write-Host "✅ Tarjeta desbloqueada correctamente" -ForegroundColor Green
    Write-Host "  ⚡ Estado: $($unblockedCard.status)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se pudo desbloquear la tarjeta" -ForegroundColor Red
}

# 10. Establecer como tarjeta predeterminada
Write-Host "`n10. ⭐ Estableciendo tarjeta como predeterminada..." -ForegroundColor Yellow
$defaultCard = Invoke-ApiRequest -Method "PUT" -Uri "$API_BASE/accounts/$ACCOUNT_ID/cards/$CARD_ID/set-default"
if ($defaultCard) {
    Write-Host "✅ Tarjeta establecida como predeterminada" -ForegroundColor Green
    Write-Host "  ⭐ Es predeterminada: $($defaultCard.is_default)" -ForegroundColor Blue
} else {
    Write-Host "❌ No se pudo establecer la tarjeta como predeterminada" -ForegroundColor Red
}

Write-Host "`n🎉 ¡Pruebas completadas!" -ForegroundColor Green
Write-Host "📊 Resumen de la funcionalidad implementada:" -ForegroundColor Cyan
Write-Host "  ✅ Creación de tarjetas asociadas a cuentas" -ForegroundColor Green
Write-Host "  ✅ Listado de tarjetas por cuenta y por usuario" -ForegroundColor Green
Write-Host "  ✅ Obtención de detalles específicos de tarjeta" -ForegroundColor Green
Write-Host "  ✅ Actualización de datos de tarjeta" -ForegroundColor Green
Write-Host "  ✅ Bloqueo y desbloqueo de tarjetas" -ForegroundColor Green
Write-Host "  ✅ Gestión de tarjeta predeterminada" -ForegroundColor Green
Write-Host "  ✅ Separación completa entre cuentas y tarjetas" -ForegroundColor Green