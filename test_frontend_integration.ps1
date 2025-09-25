# Test de Integración Frontend: Validación de Componentes y Servicios
# Este script valida que los componentes del frontend se compilen y funcionen correctamente

param(
    [string]$FrontendPath = ".\frontend",
    [switch]$Verbose,
    [switch]$BuildOnly
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

function Test-FileExists {
    param([string]$Path, [string]$Description)
    $exists = Test-Path $Path
    Write-TestResult $Description $exists $(if (-not $exists) { "Archivo no encontrado: $Path" })
    return $exists
}

function Test-AngularBuild {
    param([string]$ProjectPath)
    
    try {
        Push-Location $ProjectPath
        
        Write-Host "Ejecutando ng build..." -ForegroundColor Gray
        $buildOutput = & npm run build 2>&1
        $buildSuccess = $LASTEXITCODE -eq 0
        
        if ($Verbose -and $buildOutput) {
            Write-Host "Build output:" -ForegroundColor Gray
            $buildOutput | ForEach-Object { Write-Host "  $_" -ForegroundColor Gray }
        }
        
        return $buildSuccess
    }
    catch {
        return $false
    }
    finally {
        Pop-Location
    }
}

function Test-TypeScriptCompilation {
    param([string]$ProjectPath)
    
    try {
        Push-Location $ProjectPath
        
        Write-Host "Ejecutando verificación de TypeScript..." -ForegroundColor Gray
        $tscOutput = & npx tsc --noEmit 2>&1
        $tscSuccess = $LASTEXITCODE -eq 0
        
        if (-not $tscSuccess -and $Verbose) {
            Write-Host "TypeScript errors:" -ForegroundColor Red
            $tscOutput | ForEach-Object { Write-Host "  $_" -ForegroundColor Red }
        }
        
        return $tscSuccess
    }
    catch {
        return $false
    }
    finally {
        Pop-Location
    }
}

# Variables globales para el test
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

Write-Host "🚀 Iniciando Tests de Integración Frontend" -ForegroundColor $Blue
Write-Host "Frontend Path: $FrontendPath" -ForegroundColor Gray

# Verificar que existe el directorio frontend
if (-not (Test-Path $FrontendPath)) {
    Write-Host "❌ Directorio frontend no encontrado: $FrontendPath" -ForegroundColor $Red
    exit 1
}

# Test 1: Verificar estructura de archivos
Write-TestSection "1. Estructura de Archivos"

$requiredFiles = @(
    "$FrontendPath\src\app\pages\accounts\accounts.component.ts",
    "$FrontendPath\src\app\pages\accounts\accounts.component.html",
    "$FrontendPath\src\app\pages\accounts\accounts.component.css",
    "$FrontendPath\src\app\pages\accounts\account-form\account-form.component.ts",
    "$FrontendPath\src\app\pages\accounts\account-form\account-form.component.html",
    "$FrontendPath\src\app\pages\accounts\account-form\account-form.component.css",
    "$FrontendPath\src\app\pages\accounts\account-list\account-list.component.ts",
    "$FrontendPath\src\app\pages\accounts\wallet-dialog\wallet-dialog.component.ts",
    "$FrontendPath\src\app\pages\accounts\credit-dialog\credit-dialog.component.ts",
    "$FrontendPath\src\app\services\account.service.ts",
    "$FrontendPath\src\app\services\wallet.service.ts",
    "$FrontendPath\src\app\services\credit.service.ts",
    "$FrontendPath\src\app\services\account-validation.service.ts",
    "$FrontendPath\src\app\models\account.model.ts",
    "$FrontendPath\src\app\app.routes.ts"
)

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    $exists = Test-FileExists $file (Split-Path $file -Leaf)
    Add-TestResult $exists
    if (-not $exists) { $allFilesExist = $false }
}

# Test 2: Verificar contenido de archivos críticos
Write-TestSection "2. Validación de Contenido"

# Verificar que la ruta accounts existe en app.routes.ts
$routesContent = Get-Content "$FrontendPath\src\app\app.routes.ts" -Raw
$accountsRouteExists = $routesContent -like "*path: 'accounts'*"
Write-TestResult "Ruta 'accounts' en app.routes.ts" $accountsRouteExists
Add-TestResult $accountsRouteExists

# Verificar que AccountFormComponent existe y tiene las importaciones correctas
$accountFormContent = Get-Content "$FrontendPath\src\app\pages\accounts\account-form\account-form.component.ts" -Raw
$accountFormValid = $accountFormContent -like "*export class AccountFormComponent*" -and 
                   $accountFormContent -like "*MatDialogModule*" -and
                   $accountFormContent -like "*ReactiveFormsModule*"
Write-TestResult "AccountFormComponent válido" $accountFormValid
Add-TestResult $accountFormValid

# Verificar que AccountsComponent tiene las importaciones del nuevo componente
$accountsContent = Get-Content "$FrontendPath\src\app\pages\accounts\accounts.component.ts" -Raw
$accountsImportsValid = $accountsContent -like "*AccountFormComponent*" -and
                       $accountsContent -like "*onAddAccount*" -and
                       $accountsContent -like "*onEditAccount*"
Write-TestResult "AccountsComponent con importaciones válidas" $accountsImportsValid
Add-TestResult $accountsImportsValid

# Test 3: Verificar servicios
Write-TestSection "3. Validación de Servicios"

# Verificar AccountService
$accountServiceContent = Get-Content "$FrontendPath\src\app\services\account.service.ts" -Raw
$accountServiceValid = $accountServiceContent -like "*createAccount*" -and
                      $accountServiceContent -like "*updateAccount*" -and
                      $accountServiceContent -like "*deleteAccount*"
Write-TestResult "AccountService con métodos CRUD" $accountServiceValid
Add-TestResult $accountServiceValid

# Verificar AccountValidationService
$validationServiceContent = Get-Content "$FrontendPath\src\app\services\account-validation.service.ts" -Raw
$validationServiceValid = $validationServiceContent -like "*validateCreateAccount*" -and
                         $validationServiceContent -like "*validateUpdateAccount*"
Write-TestResult "AccountValidationService con métodos de validación" $validationServiceValid
Add-TestResult $validationServiceValid

# Test 4: Verificar modelos
Write-TestSection "4. Validación de Modelos"

$accountModelContent = Get-Content "$FrontendPath\src\app\models\account.model.ts" -Raw
$accountModelValid = $accountModelContent -like "*CreateAccountRequest*" -and
                    $accountModelContent -like "*UpdateAccountRequest*" -and
                    $accountModelContent -like "*AccountType*" -and
                    $accountModelContent -like "*Currency*"
Write-TestResult "Modelos de Account completos" $accountModelValid
Add-TestResult $accountModelValid

# Test 5: Compilación TypeScript
Write-TestSection "5. Compilación TypeScript"

$tscSuccess = Test-TypeScriptCompilation $FrontendPath
Write-TestResult "Compilación TypeScript sin errores" $tscSuccess
Add-TestResult $tscSuccess

# Test 6: Build de Angular (solo si se solicita)
if (-not $BuildOnly) {
    Write-TestSection "6. Build de Angular"
    
    $buildSuccess = Test-AngularBuild $FrontendPath
    Write-TestResult "Build de Angular exitoso" $buildSuccess
    Add-TestResult $buildSuccess
}

# Test 7: Verificar dependencias npm
Write-TestSection "7. Dependencias"

try {
    Push-Location $FrontendPath
    $packageJson = Get-Content "package.json" | ConvertFrom-Json
    
    $requiredDeps = @(
        "@angular/material",
        "@angular/cdk",
        "@angular/forms",
        "@angular/router"
    )
    
    $allDepsExist = $true
    foreach ($dep in $requiredDeps) {
        $exists = $packageJson.dependencies.PSObject.Properties.Name -contains $dep
        Write-TestResult "Dependencia $dep" $exists
        Add-TestResult $exists
        if (-not $exists) { $allDepsExist = $false }
    }
}
catch {
    Write-TestResult "Verificación de package.json" $false $_.Exception.Message
    Add-TestResult $false
}
finally {
    Pop-Location
}

# Test 8: Verificar archivos de configuración
Write-TestSection "8. Configuración"

$angularJsonExists = Test-FileExists "$FrontendPath\angular.json" "angular.json"
Add-TestResult $angularJsonExists

$tsconfigExists = Test-FileExists "$FrontendPath\tsconfig.json" "tsconfig.json"
Add-TestResult $tsconfigExists

$tsconfigAppExists = Test-FileExists "$FrontendPath\tsconfig.app.json" "tsconfig.app.json"
Add-TestResult $tsconfigAppExists

# Resumen final
Write-TestSection "📊 Resumen de Resultados"

$successRate = if ($script:TestResults.Total -gt 0) { 
    [math]::Round(($script:TestResults.Passed / $script:TestResults.Total) * 100, 2) 
} else { 0 }

Write-Host "Total de pruebas: $($script:TestResults.Total)" -ForegroundColor Gray
Write-Host "Exitosas: $($script:TestResults.Passed)" -ForegroundColor $Green
Write-Host "Fallidas: $($script:TestResults.Failed)" -ForegroundColor $Red
Write-Host "Tasa de éxito: $successRate%" -ForegroundColor $(if ($successRate -ge 90) { $Green } elseif ($successRate -ge 70) { $Yellow } else { $Red })

# Recomendaciones
if ($script:TestResults.Failed -gt 0) {
    Write-Host "`n📋 Recomendaciones:" -ForegroundColor $Yellow
    
    if (-not $accountsRouteExists) {
        Write-Host "  • Agregar ruta 'accounts' en app.routes.ts" -ForegroundColor $Yellow
    }
    
    if (-not $accountFormValid) {
        Write-Host "  • Verificar AccountFormComponent y sus importaciones" -ForegroundColor $Yellow
    }
    
    if (-not $tscSuccess) {
        Write-Host "  • Corregir errores de TypeScript antes de continuar" -ForegroundColor $Yellow
    }
    
    Write-Host "  • Ejecutar 'npm install' si faltan dependencias" -ForegroundColor $Yellow
    Write-Host "  • Verificar que todos los archivos estén correctamente creados" -ForegroundColor $Yellow
}

if ($script:TestResults.Failed -eq 0) {
    Write-Host "`n🎉 ¡Todos los tests de frontend pasaron exitosamente!" -ForegroundColor $Green
    Write-Host "   El frontend está listo para integración con el backend." -ForegroundColor $Green
    exit 0
} else {
    Write-Host "`n⚠️  Algunos tests fallaron. Revisa los componentes del frontend." -ForegroundColor $Yellow
    exit 1
}