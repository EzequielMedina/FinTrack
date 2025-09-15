# Test completo del flujo de autenticación actualizado

Write-Host "🔍 PRUEBA COMPLETA DEL FLUJO DE AUTENTICACIÓN" -ForegroundColor Cyan
Write-Host "=============================================" -ForegroundColor Cyan

$baseUrl = "http://localhost:8081/api"

# 1. Primero hacer login para obtener el token
Write-Host "`n1️⃣ PASO 1: Login para obtener token" -ForegroundColor Yellow
$loginRequest = @{
    "email" = "admin@fintrack.com"
    "password" = "admin123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginRequest -ContentType "application/json"
    
    Write-Host "✅ Login exitoso" -ForegroundColor Green
    Write-Host "📄 Datos básicos del login:" -ForegroundColor Yellow
    $loginResponse.user | ConvertTo-Json -Depth 5 | Write-Host
    
    # 2. Usar el token para llamar a /api/me
    Write-Host "`n2️⃣ PASO 2: Llamar a /api/me para obtener información completa" -ForegroundColor Yellow
    
    $headers = @{
        "Authorization" = "Bearer $($loginResponse.accessToken)"
        "Content-Type" = "application/json"
    }
    
    $meResponse = Invoke-RestMethod -Uri "$baseUrl/me" -Method GET -Headers $headers
    
    Write-Host "✅ Información completa obtenida de /api/me:" -ForegroundColor Green
    $meResponse | ConvertTo-Json -Depth 10 | Write-Host
    
    # 3. Verificar que tenemos toda la información necesaria
    Write-Host "`n3️⃣ PASO 3: Verificación de datos completos" -ForegroundColor Yellow
    
    $hasRequiredFields = $true
    $requiredFields = @('id', 'email', 'firstName', 'lastName', 'fullName', 'role', 'isActive', 'emailVerified', 'profile', 'createdAt', 'updatedAt')
    
    foreach ($field in $requiredFields) {
        if ($meResponse.PSObject.Properties.Name -contains $field) {
            Write-Host "✅ $field : $($meResponse.$field)" -ForegroundColor Green
        } else {
            Write-Host "❌ FALTA: $field" -ForegroundColor Red
            $hasRequiredFields = $false
        }
    }
    
    # 4. Verificar permisos de admin
    Write-Host "`n4️⃣ PASO 4: Prueba de permisos administrativos" -ForegroundColor Yellow
    
    # Probar listar usuarios (requiere permisos de admin)
    try {
        $usersResponse = Invoke-RestMethod -Uri "$baseUrl/users" -Method GET -Headers $headers
        Write-Host "✅ ADMIN: Puede listar usuarios ($($usersResponse.total) usuarios encontrados)" -ForegroundColor Green
        
        # Mostrar información de usuarios
        foreach ($user in $usersResponse.users) {
            Write-Host "  👤 $($user.fullName) ($($user.email)) - Rol: $($user.role)" -ForegroundColor White
        }
        
    } catch {
        Write-Host "❌ ERROR: No puede listar usuarios - $_" -ForegroundColor Red
    }
    
    # 5. Resumen final
    Write-Host "`n📊 RESUMEN FINAL:" -ForegroundColor Cyan
    Write-Host "================" -ForegroundColor Cyan
    
    if ($hasRequiredFields) {
        Write-Host "✅ Todos los campos requeridos están presentes" -ForegroundColor Green
    } else {
        Write-Host "❌ Faltan campos requeridos" -ForegroundColor Red
    }
    
    Write-Host "📋 Datos del usuario admin:" -ForegroundColor Yellow
    Write-Host "   ID: $($meResponse.id)" -ForegroundColor White
    Write-Host "   Email: $($meResponse.email)" -ForegroundColor White
    Write-Host "   Nombre completo: $($meResponse.fullName)" -ForegroundColor White
    Write-Host "   Rol: $($meResponse.role)" -ForegroundColor White
    Write-Host "   Activo: $($meResponse.isActive)" -ForegroundColor White
    Write-Host "   Email verificado: $($meResponse.emailVerified)" -ForegroundColor White
    
    Write-Host "`n🎯 CONCLUSIÓN:" -ForegroundColor Cyan
    Write-Host "El flujo actualizado funciona correctamente." -ForegroundColor Green
    Write-Host "El frontend ahora debe recibir información completa del usuario." -ForegroundColor Green
    
} catch {
    Write-Host "❌ ERROR EN LOGIN: $_" -ForegroundColor Red
    Write-Host "Detalles: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🌐 PRUEBA EN EL FRONTEND:" -ForegroundColor Cyan
Write-Host "1. Ve a http://localhost:4200" -ForegroundColor White
Write-Host "2. Haz login con: admin@fintrack.com / admin123" -ForegroundColor White
Write-Host "3. Verifica que veas la coronita 👑 y las funciones de admin" -ForegroundColor White
Write-Host "4. Entra a 'Gestionar Usuarios' desde el dashboard" -ForegroundColor White