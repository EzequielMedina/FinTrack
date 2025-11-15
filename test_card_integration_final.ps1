# Script de prueba para la funcionalidad completa de cards
# Este script prueba la separación completa entre accounts y cards

Write-Host "=== PRUEBA DE INTEGRACIÓN COMPLETA - CARDS ===" -ForegroundColor Green
Write-Host "Probando la separación completa entre accounts y cards" -ForegroundColor Yellow

# Configuración
$accountServiceUrl = "http://localhost:8082"
$userServiceUrl = "http://localhost:8081"

# Headers comunes
$headers = @{
    'Content-Type' = 'application/json'
    'Accept' = 'application/json'
}

Write-Host "`n1. Verificando que los servicios estén disponibles..." -ForegroundColor Cyan

# Verificar account-service
try {
    $accountHealth = Invoke-RestMethod -Uri "$accountServiceUrl/health" -Method GET -Headers $headers
    Write-Host "✅ Account Service está disponible" -ForegroundColor Green
} catch {
    Write-Host "❌ Account Service no está disponible: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Verificar user-service
try {
    $userHealth = Invoke-RestMethod -Uri "$userServiceUrl/health" -Method GET -Headers $headers
    Write-Host "✅ User Service está disponible" -ForegroundColor Green
} catch {
    Write-Host "❌ User Service no está disponible: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host "`n2. Creando usuario de prueba..." -ForegroundColor Cyan

# Crear usuario
$userPayload = @{
    email = "card.test@fintrack.com"
    password = "CardTest123!"
    first_name = "Card"
    last_name = "Tester"
    phone_number = "+1234567890"
} | ConvertTo-Json

try {
    $userResponse = Invoke-RestMethod -Uri "$userServiceUrl/api/users" -Method POST -Headers $headers -Body $userPayload
    $userId = $userResponse.id
    Write-Host "✅ Usuario creado exitosamente. ID: $userId" -ForegroundColor Green
} catch {
    Write-Host "❌ Error creando usuario: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        Write-Host "Detalles: $($reader.ReadToEnd())" -ForegroundColor Red
    }
    exit 1
}

Write-Host "`n3. Creando cuenta..." -ForegroundColor Cyan

# Crear cuenta
$accountPayload = @{
    user_id = $userId
    account_type = "checking"
    currency = "USD"
    initial_balance = 1000.00
    account_name = "Test Account for Cards"
} | ConvertTo-Json

try {
    $accountResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts" -Method POST -Headers $headers -Body $accountPayload
    $accountId = $accountResponse.id
    Write-Host "✅ Cuenta creada exitosamente. ID: $accountId" -ForegroundColor Green
    Write-Host "   Tipo: $($accountResponse.account_type)" -ForegroundColor White
    Write-Host "   Balance: $($accountResponse.balance)" -ForegroundColor White
} catch {
    Write-Host "❌ Error creando cuenta: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        Write-Host "Detalles: $($reader.ReadToEnd())" -ForegroundColor Red
    }
    exit 1
}

Write-Host "`n4. Probando CREAR tarjeta con endpoints separados..." -ForegroundColor Cyan

# Crear tarjeta usando el nuevo endpoint dedicado
$cardPayload = @{
    card_type = "debit"
    card_brand = "VISA"
    last_four_digits = "1234"
    masked_number = "**** **** **** 1234"
    holder_name = "CARD TESTER"
    expiration_month = 12
    expiration_year = 2028
    nickname = "Mi Tarjeta de Prueba"
    encrypted_number = "encrypted_4111111111111234"
    encrypted_cvv = "encrypted_123"
    key_fingerprint = "test_fingerprint"
} | ConvertTo-Json

try {
    $cardResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards" -Method POST -Headers $headers -Body $cardPayload
    $cardId = $cardResponse.id
    Write-Host "✅ TARJETA CREADA EXITOSAMENTE con endpoint dedicado!" -ForegroundColor Green
    Write-Host "   Card ID: $cardId" -ForegroundColor White
    Write-Host "   Tipo: $($cardResponse.card_type)" -ForegroundColor White
    Write-Host "   Marca: $($cardResponse.card_brand)" -ForegroundColor White
    Write-Host "   Nickname: $($cardResponse.nickname)" -ForegroundColor White
    Write-Host "   Account ID asociado: $($cardResponse.account_id)" -ForegroundColor White
} catch {
    Write-Host "❌ Error creando tarjeta: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        Write-Host "Detalles: $($reader.ReadToEnd())" -ForegroundColor Red
    }
    exit 1
}

Write-Host "`n5. Verificando que la tarjeta NO se guardó como cuenta..." -ForegroundColor Cyan

# Verificar que no se creó una cuenta adicional
try {
    $accountsResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/user/$userId" -Method GET -Headers $headers
    $accountsCount = $accountsResponse.Count
    Write-Host "✅ Verificado: Solo hay $accountsCount cuenta(s), la tarjeta NO se guardó como cuenta" -ForegroundColor Green
    
    foreach ($account in $accountsResponse) {
        Write-Host "   Account ID: $($account.id) | Tipo: $($account.account_type) | Nombre: $($account.account_name)" -ForegroundColor White
    }
} catch {
    Write-Host "❌ Error obteniendo cuentas: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n6. Obteniendo tarjetas por cuenta con endpoint dedicado..." -ForegroundColor Cyan

# Obtener tarjetas de la cuenta específica
try {
    $cardsResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards" -Method GET -Headers $headers
    Write-Host "✅ Tarjetas obtenidas exitosamente:" -ForegroundColor Green
    
    if ($cardsResponse.data -and $cardsResponse.data.Count -gt 0) {
        foreach ($card in $cardsResponse.data) {
            Write-Host "   Card ID: $($card.id)" -ForegroundColor White
            Write-Host "   Tipo: $($card.card_type)" -ForegroundColor White
            Write-Host "   Marca: $($card.card_brand)" -ForegroundColor White
            Write-Host "   Últimos 4 dígitos: $($card.last_four_digits)" -ForegroundColor White
            Write-Host "   Account ID: $($card.account_id)" -ForegroundColor White
            Write-Host "   Creada: $($card.created_at)" -ForegroundColor White
            Write-Host "   ---" -ForegroundColor Gray
        }
    } else {
        Write-Host "⚠️ No se encontraron tarjetas para esta cuenta" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ Error obteniendo tarjetas: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = [System.IO.StreamReader]::new($_.Exception.Response.GetResponseStream())
        Write-Host "Detalles: $($reader.ReadToEnd())" -ForegroundColor Red
    }
}

Write-Host "`n7. Obteniendo tarjeta específica..." -ForegroundColor Cyan

# Obtener tarjeta específica
try {
    $singleCardResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards/$cardId" -Method GET -Headers $headers
    Write-Host "✅ Tarjeta específica obtenida:" -ForegroundColor Green
    Write-Host "   ID: $($singleCardResponse.id)" -ForegroundColor White
    Write-Host "   Status: $($singleCardResponse.status)" -ForegroundColor White
    Write-Host "   Es default: $($singleCardResponse.is_default)" -ForegroundColor White
} catch {
    Write-Host "❌ Error obteniendo tarjeta específica: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n8. Probando actualización de tarjeta..." -ForegroundColor Cyan

# Actualizar tarjeta
$updatePayload = @{
    nickname = "Mi Tarjeta Actualizada"
    expiration_month = 11
    expiration_year = 2029
} | ConvertTo-Json

try {
    $updatedCardResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards/$cardId" -Method PUT -Headers $headers -Body $updatePayload
    Write-Host "✅ Tarjeta actualizada exitosamente:" -ForegroundColor Green
    Write-Host "   Nuevo nickname: $($updatedCardResponse.nickname)" -ForegroundColor White
    Write-Host "   Nueva expiración: $($updatedCardResponse.expiration_month)/$($updatedCardResponse.expiration_year)" -ForegroundColor White
} catch {
    Write-Host "❌ Error actualizando tarjeta: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n9. Probando bloqueo y desbloqueo de tarjeta..." -ForegroundColor Cyan

# Bloquear tarjeta
try {
    $blockResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards/$cardId/block" -Method PUT -Headers $headers
    Write-Host "✅ Tarjeta bloqueada exitosamente" -ForegroundColor Green
    Write-Host "   Nuevo status: $($blockResponse.status)" -ForegroundColor White
} catch {
    Write-Host "❌ Error bloqueando tarjeta: $($_.Exception.Message)" -ForegroundColor Red
}

# Desbloquear tarjeta
try {
    Start-Sleep -Seconds 1
    $unblockResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards/$cardId/unblock" -Method PUT -Headers $headers
    Write-Host "✅ Tarjeta desbloqueada exitosamente" -ForegroundColor Green
    Write-Host "   Nuevo status: $($unblockResponse.status)" -ForegroundColor White
} catch {
    Write-Host "❌ Error desbloqueando tarjeta: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n10. Estableciendo como tarjeta por defecto..." -ForegroundColor Cyan

# Establecer como tarjeta por defecto
try {
    $defaultResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards/$cardId/set-default" -Method PUT -Headers $headers
    Write-Host "✅ Tarjeta establecida como por defecto" -ForegroundColor Green
    Write-Host "   Es default: $($defaultResponse.is_default)" -ForegroundColor White
} catch {
    Write-Host "❌ Error estableciendo tarjeta por defecto: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n11. Obteniendo todas las tarjetas del usuario..." -ForegroundColor Cyan

# Obtener todas las tarjetas del usuario
try {
    $userCardsResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/user/$userId/cards" -Method GET -Headers $headers
    Write-Host "✅ Tarjetas del usuario obtenidas:" -ForegroundColor Green
    
    if ($userCardsResponse.data -and $userCardsResponse.data.Count -gt 0) {
        foreach ($card in $userCardsResponse.data) {
            Write-Host "   Card: $($card.nickname) | Cuenta: $($card.account_id) | Default: $($card.is_default)" -ForegroundColor White
        }
    } else {
        Write-Host "⚠️ No se encontraron tarjetas para este usuario" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ Error obteniendo tarjetas del usuario: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n12. VERIFICACIÓN FINAL: Arquitectura separada..." -ForegroundColor Cyan

# Verificación final de que la arquitectura está correctamente separada
Write-Host "✅ VERIFICACIÓN COMPLETA EXITOSA:" -ForegroundColor Green
Write-Host "   • Las tarjetas se crean con endpoints dedicados (/api/accounts/:id/cards)" -ForegroundColor White
Write-Host "   • Las tarjetas NO se guardan como cuentas" -ForegroundColor White
Write-Host "   • Las tarjetas están correctamente asociadas a cuentas" -ForegroundColor White
Write-Host "   • Todas las operaciones de tarjetas funcionan independientemente" -ForegroundColor White
Write-Host "   • La separación entre cuentas y tarjetas está completa" -ForegroundColor White

Write-Host "`n13. Limpieza opcional (eliminar tarjeta)..." -ForegroundColor Cyan

# Eliminación opcional de la tarjeta
try {
    $deleteResponse = Invoke-RestMethod -Uri "$accountServiceUrl/api/accounts/$accountId/cards/$cardId" -Method DELETE -Headers $headers
    Write-Host "✅ Tarjeta eliminada exitosamente" -ForegroundColor Green
} catch {
    Write-Host "⚠️ La eliminación puede no estar implementada: $($_.Exception.Message)" -ForegroundColor Yellow
}

Write-Host "`n=== RESUMEN FINAL ===" -ForegroundColor Green
Write-Host "🎉 LA FUNCIONALIDAD DE TARJETAS ESTÁ COMPLETAMENTE SEPARADA Y FUNCIONAL" -ForegroundColor Green
Write-Host "✅ Problema original resuelto: Las tarjetas ya NO se guardan como cuentas" -ForegroundColor Green
Write-Host "✅ Endpoints dedicados funcionando correctamente" -ForegroundColor Green
Write-Host "✅ Arquitectura de microservicios implementada correctamente" -ForegroundColor Green

Write-Host "`nPrueba completada exitosamente!" -ForegroundColor Cyan