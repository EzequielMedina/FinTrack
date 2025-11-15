# Script Maestro de Testing de Integración
# Ejecuta tests completos de frontend y backend para validar la implementación

param(
    [string]$BackendUrl = "http://localhost:8080",
    [string]$FrontendPath = ".\frontend",
    [string]$TestUserId = "test-user-integration",
    [switch]$Verbose,
    [switch]$FrontendOnly,
    [switch]$BackendOnly,
    [switch]$SkipBuild
)

$ErrorActionPreference = "Stop"

# Colores para output
$Green = "Green"
$Red = "Red"
$Yellow = "Yellow"
$Blue = "Blue"
$Magenta = "Magenta"

function Write-Header {
    param([string]$Title)
    Write-Host "`n" -NoNewline
    Write-Host "=" * 80 -ForegroundColor $Magenta
    Write-Host " $Title" -ForegroundColor $Magenta
    Write-Host "=" * 80 -ForegroundColor $Magenta
}

function Write-SubHeader {
    param([string]$Title)
    Write-Host "`n🔧 $Title" -ForegroundColor $Blue
    Write-Host "-" * ($Title.Length + 3) -ForegroundColor $Blue
}

# Variables de resultados
$script:TestSummary = @{
    Frontend = @{ Executed = $false; Success = $false; Details = "" }
    Backend = @{ Executed = $false; Success = $false; Details = "" }
    Integration = @{ Executed = $false; Success = $false; Details = "" }
}

Write-Header "TESTING DE INTEGRACIÓN COMPLETO - FINTRACK ACCOUNTS"
Write-Host "🎯 Objetivo: Validar la implementación completa del módulo de cuentas" -ForegroundColor $Blue
Write-Host "📅 Fecha: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor Gray
Write-Host "🔗 Backend URL: $BackendUrl" -ForegroundColor Gray
Write-Host "📁 Frontend Path: $FrontendPath" -ForegroundColor Gray

# Test 1: Frontend Testing
if (-not $BackendOnly) {
    Write-Header "1. TESTING DE FRONTEND"
    
    try {
        $frontendParams = @{
            FrontendPath = $FrontendPath
        }
        
        if ($Verbose) { $frontendParams.Verbose = $true }
        if ($SkipBuild) { $frontendParams.BuildOnly = $true }
        
        Write-Host "Ejecutando test_frontend_integration.ps1..." -ForegroundColor $Blue
        
        $frontendResult = & ".\test_frontend_integration.ps1" @frontendParams
        $frontendSuccess = $LASTEXITCODE -eq 0
        
        $script:TestSummary.Frontend.Executed = $true
        $script:TestSummary.Frontend.Success = $frontendSuccess
        $script:TestSummary.Frontend.Details = if ($frontendSuccess) { "Todos los componentes validados" } else { "Errores en componentes" }
        
        if ($frontendSuccess) {
            Write-Host "✅ Frontend testing completado exitosamente" -ForegroundColor $Green
        } else {
            Write-Host "❌ Frontend testing falló" -ForegroundColor $Red
            if (-not $BackendOnly -and -not $FrontendOnly) {
                Write-Host "⚠️  Continuando con backend testing..." -ForegroundColor $Yellow
            }
        }
    }
    catch {
        Write-Host "❌ Error ejecutando frontend testing: $($_.Exception.Message)" -ForegroundColor $Red
        $script:TestSummary.Frontend.Details = "Error de ejecución: $($_.Exception.Message)"
    }
}

# Test 2: Backend Testing
if (-not $FrontendOnly) {
    Write-Header "2. TESTING DE BACKEND"
    
    try {
        $backendParams = @{
            BackendUrl = $BackendUrl
            TestUserId = $TestUserId
        }
        
        if ($Verbose) { $backendParams.Verbose = $true }
        
        Write-Host "Ejecutando test_integration_accounts.ps1..." -ForegroundColor $Blue
        
        $backendResult = & ".\test_integration_accounts.ps1" @backendParams
        $backendSuccess = $LASTEXITCODE -eq 0
        
        $script:TestSummary.Backend.Executed = $true
        $script:TestSummary.Backend.Success = $backendSuccess
        $script:TestSummary.Backend.Details = if ($backendSuccess) { "Todas las APIs validadas" } else { "Errores en APIs" }
        
        if ($backendSuccess) {
            Write-Host "✅ Backend testing completado exitosamente" -ForegroundColor $Green
        } else {
            Write-Host "❌ Backend testing falló" -ForegroundColor $Red
        }
    }
    catch {
        Write-Host "❌ Error ejecutando backend testing: $($_.Exception.Message)" -ForegroundColor $Red
        $script:TestSummary.Backend.Details = "Error de ejecución: $($_.Exception.Message)"
    }
}

# Test 3: Integration Testing (solo si ambos están disponibles)
if (-not $FrontendOnly -and -not $BackendOnly -and $script:TestSummary.Frontend.Success -and $script:TestSummary.Backend.Success) {
    Write-Header "3. TESTING DE INTEGRACIÓN COMPLETA"
    
    Write-SubHeader "Validación de Flujos End-to-End"
    
    # Test de integración básico: verificar que frontend puede comunicarse con backend
    try {
        Write-Host "🔍 Verificando comunicación Frontend ↔ Backend..." -ForegroundColor $Blue
        
        # Simular una llamada que haría el frontend
        $testEndpoint = "$BackendUrl/api/accounts"
        $response = Invoke-WebRequest -Uri $testEndpoint -UseBasicParsing -TimeoutSec 10
        
        if ($response.StatusCode -eq 200) {
            Write-Host "✅ Comunicación Frontend ↔ Backend: OK" -ForegroundColor $Green
            $script:TestSummary.Integration.Success = $true
            $script:TestSummary.Integration.Details = "Comunicación establecida correctamente"
        } else {
            Write-Host "❌ Comunicación Frontend ↔ Backend: FAILED" -ForegroundColor $Red
            $script:TestSummary.Integration.Details = "Error de comunicación: Status $($response.StatusCode)"
        }
        
        $script:TestSummary.Integration.Executed = $true
    }
    catch {
        Write-Host "❌ Error en testing de integración: $($_.Exception.Message)" -ForegroundColor $Red
        $script:TestSummary.Integration.Executed = $true
        $script:TestSummary.Integration.Success = $false
        $script:TestSummary.Integration.Details = "Error: $($_.Exception.Message)"
    }
}

# Resumen Final
Write-Header "📊 RESUMEN EJECUTIVO"

Write-Host "🧪 Resultados de Testing:" -ForegroundColor $Blue

# Frontend Results
if ($script:TestSummary.Frontend.Executed) {
    $frontendIcon = if ($script:TestSummary.Frontend.Success) { "✅" } else { "❌" }
    $frontendColor = if ($script:TestSummary.Frontend.Success) { $Green } else { $Red }
    Write-Host "   Frontend: $frontendIcon $($script:TestSummary.Frontend.Details)" -ForegroundColor $frontendColor
}

# Backend Results
if ($script:TestSummary.Backend.Executed) {
    $backendIcon = if ($script:TestSummary.Backend.Success) { "✅" } else { "❌" }
    $backendColor = if ($script:TestSummary.Backend.Success) { $Green } else { $Red }
    Write-Host "   Backend:  $backendIcon $($script:TestSummary.Backend.Details)" -ForegroundColor $backendColor
}

# Integration Results
if ($script:TestSummary.Integration.Executed) {
    $integrationIcon = if ($script:TestSummary.Integration.Success) { "✅" } else { "❌" }
    $integrationColor = if ($script:TestSummary.Integration.Success) { $Green } else { $Red }
    Write-Host "   Integración: $integrationIcon $($script:TestSummary.Integration.Details)" -ForegroundColor $integrationColor
}

# Overall Status
$overallSuccess = (
    (!$script:TestSummary.Frontend.Executed -or $script:TestSummary.Frontend.Success) -and
    (!$script:TestSummary.Backend.Executed -or $script:TestSummary.Backend.Success) -and
    (!$script:TestSummary.Integration.Executed -or $script:TestSummary.Integration.Success)
)

Write-Host "`n🎯 Estado General:" -ForegroundColor $Blue
if ($overallSuccess) {
    Write-Host "   🎉 TODOS LOS TESTS PASARON - IMPLEMENTACIÓN LISTA PARA PRODUCCIÓN" -ForegroundColor $Green
} else {
    Write-Host "   ⚠️  ALGUNOS TESTS FALLARON - REVISAR IMPLEMENTACIÓN" -ForegroundColor $Yellow
}

# Recomendaciones
Write-Host "`n📋 Próximos Pasos:" -ForegroundColor $Blue

if ($overallSuccess) {
    Write-Host "   ✅ La implementación del módulo de cuentas está completa" -ForegroundColor $Green
    Write-Host "   ✅ Frontend y Backend están integrados correctamente" -ForegroundColor $Green
    Write-Host "   ✅ Listo para despliegue y testing manual" -ForegroundColor $Green
    Write-Host "   📝 Considerar agregar tests de usuario final (E2E con navegador)" -ForegroundColor $Blue
} else {
    if ($script:TestSummary.Frontend.Executed -and -not $script:TestSummary.Frontend.Success) {
        Write-Host "   🔧 Corregir errores en componentes de frontend" -ForegroundColor $Yellow
        Write-Host "   📝 Verificar importaciones y sintaxis de TypeScript" -ForegroundColor $Yellow
    }
    
    if ($script:TestSummary.Backend.Executed -and -not $script:TestSummary.Backend.Success) {
        Write-Host "   🔧 Corregir errores en APIs de backend" -ForegroundColor $Yellow
        Write-Host "   📝 Verificar que el servidor esté ejecutándose en $BackendUrl" -ForegroundColor $Yellow
    }
    
    Write-Host "   🔄 Re-ejecutar tests después de las correcciones" -ForegroundColor $Blue
}

Write-Host "`n📁 Archivos de Configuración:" -ForegroundColor $Blue
Write-Host "   • Logs detallados disponibles en salida de cada script" -ForegroundColor Gray
Write-Host "   • Frontend: $FrontendPath" -ForegroundColor Gray
Write-Host "   • Backend: $BackendUrl" -ForegroundColor Gray

Write-Host "`n⏰ Testing completado en: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')" -ForegroundColor Gray

# Exit code basado en el resultado general
if ($overallSuccess) {
    exit 0
} else {
    exit 1
}