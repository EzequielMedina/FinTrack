# Script para probar la integración backend-frontend
# Primero crear un usuario con rol admin

$baseUrl = "http://localhost:8081/api"

# Crear usuario administrador vía registro
$registerRequest = @{
    "email" = "admin@fintrack.com"
    "password" = "admin123"
    "firstName" = "Administrador"
    "lastName" = "Sistema"
} | ConvertTo-Json -Depth 3

Write-Host "Registrando usuario administrador..."
$response = Invoke-RestMethod -Uri "$baseUrl/auth/register" -Method POST -Body $registerRequest -ContentType "application/json" -ErrorAction SilentlyContinue

if ($response) {
    Write-Host "Usuario registrado exitosamente:"
    $response | ConvertTo-Json -Depth 5
    
    # Ahora intentar login
    $loginRequest = @{
        "email" = "admin@fintrack.com"
        "password" = "admin123"
    } | ConvertTo-Json
    
    Write-Host "`nIntentando login..."
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginRequest -ContentType "application/json" -ErrorAction SilentlyContinue
    
    if ($loginResponse) {
        Write-Host "Login exitoso:"
        $loginResponse | ConvertTo-Json -Depth 5
    } else {
        Write-Host "Error en login"
    }
} else {
    Write-Host "Error registrando usuario"
}

Write-Host "`nPuedes acceder al frontend en: http://localhost:4200"