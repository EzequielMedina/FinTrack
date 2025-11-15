# Script para analizar la respuesta completa del backend
$baseUrl = "http://localhost:8081/api"

# Login como admin para analizar la respuesta
$loginRequest = @{
    "email" = "admin@fintrack.com"
    "password" = "admin123"
} | ConvertTo-Json

Write-Host "=== ANÁLISIS DE RESPUESTA DEL BACKEND ===" -ForegroundColor Cyan
Write-Host "Login como admin..." -ForegroundColor Yellow

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginRequest -ContentType "application/json"
    
    Write-Host "`n✅ LOGIN EXITOSO" -ForegroundColor Green
    Write-Host "📄 Respuesta completa:" -ForegroundColor Yellow
    $loginResponse | ConvertTo-Json -Depth 10 | Write-Host
    
    # Decodificar el JWT para ver qué contiene
    Write-Host "`n🔍 ANÁLISIS DEL JWT TOKEN:" -ForegroundColor Cyan
    $token = $loginResponse.accessToken
    $parts = $token.Split('.')
    
    if ($parts.Length -eq 3) {
        $payload = $parts[1]
        # Agregar padding si es necesario
        while ($payload.Length % 4 -ne 0) {
            $payload += "="
        }
        
        try {
            $decodedBytes = [System.Convert]::FromBase64String($payload)
            $decodedText = [System.Text.Encoding]::UTF8.GetString($decodedBytes)
            $jwtPayload = $decodedText | ConvertFrom-Json
            
            Write-Host "📋 Contenido del JWT:" -ForegroundColor Yellow
            $jwtPayload | ConvertTo-Json -Depth 5 | Write-Host
            
            Write-Host "`n🔑 Información clave del token:" -ForegroundColor Green
            Write-Host "- User ID (sub): $($jwtPayload.sub)" -ForegroundColor White
            Write-Host "- Email: $($jwtPayload.email)" -ForegroundColor White
            Write-Host "- Role: $($jwtPayload.role)" -ForegroundColor White
            Write-Host "- Issuer: $($jwtPayload.iss)" -ForegroundColor White
            Write-Host "- Expires: $(Get-Date -UnixTimeSeconds $jwtPayload.exp)" -ForegroundColor White
            
        } catch {
            Write-Host "❌ Error decodificando JWT payload: $_" -ForegroundColor Red
        }
    }
    
    # Probar endpoint /api/me para obtener más información del usuario
    Write-Host "`n🔍 PROBANDO ENDPOINT /api/me:" -ForegroundColor Cyan
    $headers = @{
        "Authorization" = "Bearer $($loginResponse.accessToken)"
        "Content-Type" = "application/json"
    }
    
    try {
        $meResponse = Invoke-RestMethod -Uri "$baseUrl/me" -Method GET -Headers $headers
        Write-Host "📄 Respuesta de /api/me:" -ForegroundColor Yellow
        $meResponse | ConvertTo-Json -Depth 10 | Write-Host
    } catch {
        Write-Host "❌ Error en /api/me: $_" -ForegroundColor Red
    }
    
    # Probar listar usuarios (función de admin)
    Write-Host "`n🔍 PROBANDO ENDPOINT /api/users (lista de usuarios):" -ForegroundColor Cyan
    try {
        $usersResponse = Invoke-RestMethod -Uri "$baseUrl/users" -Method GET -Headers $headers
        Write-Host "📄 Respuesta de /api/users:" -ForegroundColor Yellow
        $usersResponse | ConvertTo-Json -Depth 10 | Write-Host
    } catch {
        Write-Host "❌ Error en /api/users: $_" -ForegroundColor Red
        Write-Host "Detalles: $($_.Exception.Message)" -ForegroundColor Red
    }
    
} catch {
    Write-Host "❌ ERROR EN LOGIN: $_" -ForegroundColor Red
}

Write-Host "`n=== ANÁLISIS COMPLETO ===" -ForegroundColor Cyan