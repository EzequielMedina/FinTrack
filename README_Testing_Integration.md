# Testing de Integración - Módulo de Cuentas FinTrack

## 📋 Descripción

Este conjunto de scripts permite validar la implementación completa del módulo de gestión de cuentas, incluyendo:

- ✅ **Frontend**: Componentes Angular, servicios, y routing
- ✅ **Backend**: APIs REST, operaciones CRUD, y lógica de negocio
- ✅ **Integración**: Comunicación end-to-end entre frontend y backend

## 🚀 Scripts Disponibles

### 1. `test_complete_integration.ps1` (Script Principal)

Script maestro que ejecuta todos los tests de integración.

**Uso básico:**
```powershell
.\test_complete_integration.ps1
```

**Parámetros disponibles:**
```powershell
.\test_complete_integration.ps1 `
    -BackendUrl "http://localhost:8080" `
    -FrontendPath ".\frontend" `
    -TestUserId "test-user-123" `
    -Verbose
```

**Opciones específicas:**
```powershell
# Solo frontend
.\test_complete_integration.ps1 -FrontendOnly

# Solo backend
.\test_complete_integration.ps1 -BackendOnly

# Sin build (más rápido)
.\test_complete_integration.ps1 -SkipBuild
```

### 2. `test_frontend_integration.ps1`

Valida específicamente los componentes del frontend.

**Funcionalidades:**
- ✅ Verifica estructura de archivos
- ✅ Valida importaciones y dependencias
- ✅ Compila TypeScript
- ✅ Build de Angular (opcional)
- ✅ Verifica configuración de routing

**Uso:**
```powershell
.\test_frontend_integration.ps1 -FrontendPath ".\frontend" -Verbose
```

### 3. `test_integration_accounts.ps1`

Valida las APIs del backend y operaciones de cuentas.

**Funcionalidades:**
- ✅ Health check del backend
- ✅ CRUD completo de cuentas
- ✅ Operaciones de wallet (agregar/retirar fondos)
- ✅ Gestión de tarjetas de crédito
- ✅ Validaciones de negocio
- ✅ Cleanup automático

**Uso:**
```powershell
.\test_integration_accounts.ps1 -BackendUrl "http://localhost:8080" -Verbose
```

## 🛠️ Requisitos Previos

### Backend
```powershell
# 1. Compilar el servicio de cuentas
cd backend\services\account-service
go build -o account-service.exe .\cmd\api

# 2. Ejecutar el servicio
.\account-service.exe
```

### Frontend
```powershell
# 1. Instalar dependencias
cd frontend
npm install

# 2. Verificar que compile
npm run build
```

### Base de Datos
```powershell
# Ejecutar MySQL con Docker
docker-compose up mysql
```

## 📊 Interpretación de Resultados

### Códigos de Salida
- **0**: Todos los tests pasaron ✅
- **1**: Algunos tests fallaron ❌

### Símbolos en la Salida
- ✅ **Verde**: Test exitoso
- ❌ **Rojo**: Test fallido
- ⚠️ **Amarillo**: Advertencia
- 🔍 **Azul**: Información/Sección

### Ejemplo de Salida Exitosa
```
🚀 Iniciando Tests de Integración Frontend-Backend
Backend URL: http://localhost:8080

🔍 1. Conectividad del Backend
===============================
✅ Health Check del Backend

🔍 2. Endpoints de API de Cuentas
=====================================
✅ Endpoint GET /api/accounts

📊 Resumen de Resultados
=========================
Total de pruebas: 15
Exitosas: 15
Fallidas: 0
Tasa de éxito: 100%

🎉 ¡Todos los tests de integración pasaron exitosamente!
```

## 🔧 Solución de Problemas

### Error: "Backend no disponible"
```powershell
# Verificar que el servicio esté ejecutándose
curl http://localhost:8080/health

# Si no responde, compilar y ejecutar:
cd backend\services\account-service
go run .\cmd\api
```

### Error: "Archivos frontend no encontrados"
```powershell
# Verificar la ruta
.\test_frontend_integration.ps1 -FrontendPath ".\tu-ruta-frontend"

# Verificar que los archivos existan
ls frontend\src\app\pages\accounts\
```

### Error: "Compilación TypeScript fallida"
```powershell
cd frontend
npx tsc --noEmit  # Ver errores específicos
npm run build     # Intentar build completo
```

### Error: "Dependencias faltantes"
```powershell
cd frontend
npm install
npm audit fix  # Si hay vulnerabilidades
```

## 📝 Customización

### Agregar Nuevos Tests

**Frontend:**
Editar `test_frontend_integration.ps1`, sección "Test X":
```powershell
# Test X: Tu nuevo test
Write-TestSection "X. Tu Nueva Validación"

$tuTest = Test-TuFuncionalidad
Write-TestResult "Tu test" $tuTest
Add-TestResult $tuTest
```

**Backend:**
Editar `test_integration_accounts.ps1`, agregar después del test 8:
```powershell
# Test 9: Tu nuevo endpoint
Write-TestSection "9. Tu Nueva API"

$tuApiResult = Test-BackendEndpoint -Url "$BackendUrl/api/tu-endpoint"
$tuApiSuccess = $tuApiResult.Success
Write-TestResult "Tu API endpoint" $tuApiSuccess $tuApiResult.Error
Add-TestResult $tuApiSuccess
```

### Configurar Diferentes Entornos

**Desarrollo:**
```powershell
.\test_complete_integration.ps1 -BackendUrl "http://localhost:8080"
```

**Testing:**
```powershell
.\test_complete_integration.ps1 -BackendUrl "http://test-server:8080"
```

**Staging:**
```powershell
.\test_complete_integration.ps1 -BackendUrl "https://staging.fintrack.com"
```

## 🏃‍♂️ Ejecución Rápida

### Validación Completa (Recomendado)
```powershell
# Ejecutar todos los tests con salida detallada
.\test_complete_integration.ps1 -Verbose
```

### Solo Verificar Frontend
```powershell
# Rápido, sin build
.\test_complete_integration.ps1 -FrontendOnly -SkipBuild
```

### Solo Verificar Backend
```powershell
# Verificar APIs
.\test_complete_integration.ps1 -BackendOnly
```

## 📈 Métricas y Reporting

Los scripts generan métricas automáticas:

- **Tasa de éxito**: Porcentaje de tests exitosos
- **Tiempo de ejecución**: Timestamps de inicio y fin
- **Detalles de errores**: Mensajes específicos para debugging
- **Recomendaciones**: Próximos pasos basados en resultados

## 🔄 Integración Continua

Para usar en pipelines de CI/CD:

```yaml
# GitHub Actions / Azure DevOps ejemplo
- name: Run Integration Tests
  run: |
    .\test_complete_integration.ps1 -BackendUrl ${{ secrets.BACKEND_URL }}
  shell: powershell
```

## 📞 Soporte

Si encuentras problemas:

1. **Verificar prerequisitos**: Backend ejecutándose, frontend compilando
2. **Ejecutar con -Verbose**: Para ver detalles de errores
3. **Revisar logs**: Cada script proporciona información detallada
4. **Probar componentes individuales**: Usar scripts específicos

---

**Nota**: Estos scripts están diseñados para ser ejecutados desde la raíz del proyecto FinTrack.