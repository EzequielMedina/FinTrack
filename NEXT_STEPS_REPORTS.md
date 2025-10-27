# 🚀 Próximos Pasos - Microservicio de Reportes

## ✅ Lo que está Completo

### Backend (Go)
- ✅ Estructura completa del microservicio
- ✅ 5 repositorios con queries SQL optimizadas
- ✅ 5 servicios de negocio
- ✅ 5 handlers HTTP con validación
- ✅ Router con CORS configurado
- ✅ Dockerfile y docker-compose configurados
- ✅ Documentación completa de API

### Frontend (Angular)
- ✅ ReportService con todas las interfaces TypeScript
- ✅ ReportsComponent (página de navegación)
- ✅ TransactionReportComponent completo con visualizaciones
- ✅ Rutas configuradas
- ✅ Estilos responsive

---

## 🔧 Pasos para Poner en Marcha

### 1. Inicializar el Backend

```bash
# Terminal 1: Navegar al servicio
cd c:\Facultad\Alumno\PS\backend\services\report-service

# Descargar dependencias de Go
go mod tidy
go mod download

# Esto debería descargar:
# - github.com/gin-gonic/gin
# - github.com/go-sql-driver/mysql
# - github.com/google/uuid
# - github.com/joho/godotenv
# Y sus dependencias transitivas
```

### 2. Probar el Servicio Localmente (Opcional)

```bash
# Asegúrate de que MySQL esté corriendo
docker-compose up -d mysql

# Espera a que MySQL esté listo (revisa los logs)
docker-compose logs -f mysql

# Ejecuta el servicio
go run cmd/main.go

# Deberías ver:
# 🚀 Iniciando FinTrack Report Service...
# ✅ Conexión a base de datos establecida
# 🚀 Servidor iniciado en puerto 8085
# 📊 Report Service API disponible en http://localhost:8085/api/v1
```

### 3. Build y Ejecutar con Docker

```bash
# Desde la raíz del proyecto
cd c:\Facultad\Alumno\PS

# Build del servicio
docker-compose build report-service

# Ejecutar solo este servicio con sus dependencias
docker-compose up mysql report-service

# O ejecutar todos los servicios
docker-compose up
```

### 4. Verificar que Funciona

```powershell
# Health check
Invoke-WebRequest -Uri "http://localhost:8085/health" -Method GET

# Debería retornar:
# {
#   "status": "healthy",
#   "service": "report-service"
# }

# Probar un reporte (reemplaza USER_ID con un ID real de tu BD)
$userId = "tu-user-id-aqui"
$url = "http://localhost:8085/api/v1/reports/transactions?user_id=$userId&start_date=2024-01-01&end_date=2024-01-31"
Invoke-WebRequest -Uri $url -Method GET
```

### 5. Configurar el Frontend

```bash
# Navegar al frontend
cd c:\Facultad\Alumno\PS\frontend

# Instalar Chart.js (recomendado)
npm install chart.js ng2-charts
npm install --save-dev @types/chart.js

# Ejecutar el frontend
npm start

# Acceder a los reportes
# http://localhost:4200/reports
```

### 6. Agregar Link de Navegación

Edita tu componente de navegación (probablemente en `shared/sidebar` o similar):

```html
<!-- Agregar este link al menú -->
<a routerLink="/reports" routerLinkActive="active" class="nav-link">
  <i class="icon">📊</i>
  <span>Reportes</span>
</a>
```

---

## 🎨 Completar los Componentes Faltantes

Ya tienes el componente de **Transacciones** completo. Ahora puedes crear los otros 4:

### 1. InstallmentReportComponent

```bash
cd frontend/src/app/pages/reports
# Copiar la estructura de transaction-report
# Adaptar el HTML/CSS para mostrar:
# - Planes de cuotas con barras de progreso
# - Calendario de pagos próximos
# - Alertas de pagos vencidos
```

### 2. AccountReportComponent

```bash
# Mostrar:
# - Cards de resumen (balance, crédito, etc.)
# - Lista de cuentas
# - Lista de tarjetas
# - Gráfico de distribución por tipo
```

### 3. ExpenseIncomeReportComponent

```bash
# Mostrar:
# - Timeline de ingresos vs gastos
# - Gráfico de tendencias
# - Indicadores de ahorro
# - Proyecciones futuras
```

### 4. NotificationReportComponent (Admin)

```bash
# Mostrar:
# - Métricas de notificaciones
# - Tasa de éxito
# - Timeline de envíos
# - Logs de jobs ejecutados
```

**Puedes usar TransactionReportComponent como plantilla** para todos estos componentes. La estructura es similar:

1. Service call en `ngOnInit()`
2. Loading/Error states
3. Summary cards
4. Visualizaciones específicas
5. Exportación (opcional)

---

## 📊 Mejorar con Chart.js

### Ejemplo: Agregar un gráfico de torta

```typescript
// transaction-report.component.ts
import { Chart } from 'chart.js/auto';

export class TransactionReportComponent implements OnInit {
  private chart?: Chart;

  ngAfterViewInit(): void {
    if (this.report) {
      this.createPieChart();
    }
  }

  private createPieChart(): void {
    const canvas = document.getElementById('pieChart') as HTMLCanvasElement;
    if (!canvas || !this.report) return;

    // Destruir chart anterior si existe
    if (this.chart) {
      this.chart.destroy();
    }

    this.chart = new Chart(canvas, {
      type: 'pie',
      data: {
        labels: this.report.by_type.map(t => this.getTypeLabel(t.type)),
        datasets: [{
          data: this.report.by_type.map(t => t.amount),
          backgroundColor: [
            '#FF6384', '#36A2EB', '#FFCE56', 
            '#4BC0C0', '#9966FF', '#FF9F40',
            '#FF6384', '#C9CBCF', '#4BC0C0'
          ],
          borderWidth: 2,
          borderColor: '#fff'
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            position: 'bottom',
            labels: {
              padding: 15,
              font: {
                size: 12
              }
            }
          },
          title: {
            display: true,
            text: 'Distribución de Transacciones',
            font: {
              size: 16,
              weight: 'bold'
            }
          },
          tooltip: {
            callbacks: {
              label: (context) => {
                const label = context.label || '';
                const value = this.formatCurrency(context.parsed);
                return `${label}: ${value}`;
              }
            }
          }
        }
      }
    });
  }

  loadReport(): void {
    // ... existing code ...
    this.reportService.getTransactionReport(...).subscribe({
      next: (data) => {
        this.report = data;
        this.isLoading = false;
        
        // Recrear gráfico con nuevos datos
        setTimeout(() => this.createPieChart(), 0);
      },
      // ... error handling ...
    });
  }

  ngOnDestroy(): void {
    if (this.chart) {
      this.chart.destroy();
    }
  }
}
```

En el HTML:

```html
<div class="chart-card">
  <h3>📊 Distribución de Transacciones</h3>
  <div class="chart-wrapper">
    <canvas id="pieChart" width="400" height="400"></canvas>
  </div>
</div>
```

En el CSS:

```css
.chart-wrapper {
  position: relative;
  height: 400px;
  margin: 2rem auto;
  max-width: 600px;
}
```

### Otros tipos de gráficos útiles:

#### Gráfico de Líneas (Timeline)
```typescript
type: 'line',
data: {
  labels: this.report.by_period.map(p => this.formatDate(p.date)),
  datasets: [
    {
      label: 'Ingresos',
      data: this.report.by_period.map(p => p.income),
      borderColor: '#4CAF50',
      backgroundColor: 'rgba(76, 175, 80, 0.1)',
      fill: true,
      tension: 0.4
    },
    {
      label: 'Gastos',
      data: this.report.by_period.map(p => p.expenses),
      borderColor: '#FF9800',
      backgroundColor: 'rgba(255, 152, 0, 0.1)',
      fill: true,
      tension: 0.4
    }
  ]
}
```

#### Gráfico de Barras Horizontales
```typescript
type: 'bar',
options: {
  indexAxis: 'y',  // Esto lo hace horizontal
  // ... otras opciones
}
```

#### Gráfico de Dona (Doughnut)
```typescript
type: 'doughnut',
// Similar a pie, pero con un hueco en el centro
```

---

## 🔐 Agregar Autenticación JWT

### En el Backend (Go)

```go
// internal/infrastructure/entrypoints/router/router.go

func SetupRouter(cfg *config.Config, reportService service.ReportService) *gin.Engine {
    r := gin.Default()
    
    // Middleware de autenticación
    r.Use(authMiddleware())
    
    // ... resto del código
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "No autorizado"})
            c.Abort()
            return
        }
        
        // Quitar el prefijo "Bearer "
        if len(token) > 7 && token[:7] == "Bearer " {
            token = token[7:]
        }
        
        // TODO: Validar el token JWT
        // claims, err := validateJWT(token)
        // if err != nil {
        //     c.JSON(401, gin.H{"error": "Token inválido"})
        //     c.Abort()
        //     return
        // }
        
        // Guardar el userId del token en el context
        // c.Set("userId", claims.UserID)
        
        c.Next()
    }
}
```

### En el Frontend (Angular)

El servicio ya debería estar usando el interceptor de autenticación existente. Verifica:

```typescript
// interceptors/auth.interceptor.ts
// Debería agregar automáticamente el header Authorization
```

---

## 📝 Testing

### Backend Tests

```bash
cd backend/services/report-service

# Crear archivo de test
# internal/core/service/report_service_test.go
```

```go
package service_test

import (
    "testing"
    "context"
    "time"
)

func TestGetTransactionReport(t *testing.T) {
    // TODO: Implementar tests
    // 1. Mock del repositorio
    // 2. Crear servicio con el mock
    // 3. Llamar al método
    // 4. Verificar el resultado
}
```

### Frontend Tests

```bash
cd frontend

# Los archivos .spec.ts ya están creados
# Ejecutar tests
ng test
```

---

## 🐛 Troubleshooting Común

### Problema: "could not import github.com/gin-gonic/gin"

**Solución:**
```bash
cd backend/services/report-service
go mod tidy
go mod download
go mod verify
```

### Problema: Error de CORS en el frontend

**Solución:**
Verifica `frontend/proxy.conf.json`:

```json
{
  "/report-service": {
    "target": "http://localhost:8085",
    "secure": false,
    "changeOrigin": true,
    "pathRewrite": {
      "^/report-service": ""
    }
  }
}
```

Y asegúrate de usar `npm start` (no `ng serve`) para que use el proxy.

### Problema: "dial tcp: lookup mysql: no such host"

**Solución:**
Asegúrate de que MySQL esté corriendo:

```bash
docker-compose up -d mysql
docker-compose logs mysql | grep "ready for connections"
```

### Problema: Frontend no encuentra el servicio

**Solución:**
Verifica en `environment.ts`:

```typescript
export const environment = {
  apiUrl: 'http://localhost',  // Sin puerto, Nginx hace proxy
  // ...
};
```

Y el servicio debe llamar a:
```typescript
`${environment.apiUrl}/report-service/api/v1/reports/...`
```

---

## 📚 Recursos Útiles

### Documentación
- **Go Gin**: https://gin-gonic.com/docs/
- **Chart.js**: https://www.chartjs.org/docs/latest/
- **Angular HttpClient**: https://angular.io/guide/http
- **TypeScript**: https://www.typescriptlang.org/docs/

### Tutoriales
- Chart.js con Angular: https://www.chartjs.org/docs/latest/getting-started/integration.html
- Testing en Go: https://go.dev/doc/tutorial/add-a-test
- Angular Component Testing: https://angular.io/guide/testing-components-basics

---

## ✅ Checklist Final

Antes de considerar el proyecto completo, verifica:

### Backend
- [ ] `go mod tidy` ejecutado sin errores
- [ ] Servicio compila sin errores (`go build cmd/main.go`)
- [ ] Health check responde (`curl http://localhost:8085/health`)
- [ ] Al menos un endpoint retorna datos válidos
- [ ] Docker build exitoso
- [ ] docker-compose up funciona

### Frontend
- [ ] `npm install` completado
- [ ] Servicio de reportes compilado sin errores
- [ ] Componentes creados y rutas configuradas
- [ ] Navegación a /reports funciona
- [ ] Al menos un reporte se visualiza correctamente
- [ ] Diseño responsive en móvil

### Integración
- [ ] Frontend puede llamar al backend
- [ ] CORS configurado correctamente
- [ ] Datos fluyen de BD → Backend → Frontend
- [ ] Filtros funcionan
- [ ] Exportación funciona (si implementada)

### Documentación
- [ ] README.md del servicio completo
- [ ] API_DOCUMENTATION.md con ejemplos
- [ ] Código comentado
- [ ] Variables de entorno documentadas

---

## 🎯 Siguientes Features Recomendadas

### Prioritarias
1. **Completar componentes faltantes** (4 reportes restantes)
2. **Agregar Chart.js** para visualizaciones mejoradas
3. **Implementar JWT** para seguridad
4. **Tests básicos** backend y frontend

### Secundarias
5. **Paginación** en listados largos
6. **Caché** para reportes pesados (Redis)
7. **Exportación a PDF** con gráficos
8. **Filtros avanzados** (múltiples tipos, rangos personalizados)

### Avanzadas
9. **Reportes programados** (envío automático por email)
10. **Dashboard en tiempo real** (WebSockets)
11. **Predicciones** con Machine Learning
12. **Alertas personalizables** basadas en métricas

---

## 💡 Tips Finales

1. **Empieza simple**: Haz que funcione un reporte básico antes de agregar features complejas
2. **Usa los datos de ejemplo**: Inserta datos de prueba en la BD para ver resultados realistas
3. **Revisa los logs**: `docker-compose logs report-service` te dirá qué está fallando
4. **Itera rápidamente**: Haz cambios pequeños, prueba, ajusta
5. **Reutiliza código**: Los 5 reportes son muy similares, copia y adapta
6. **Mantén la consistencia**: Usa los mismos patrones de diseño en todos los componentes

---

## 🎉 ¡Éxito!

Tienes todo lo necesario para tener un sistema de reportes completo y funcional. 

**Siguiente paso inmediato:**
```bash
cd c:\Facultad\Alumno\PS\backend\services\report-service
go mod tidy
docker-compose up mysql report-service
```

Y luego en otra terminal:
```bash
cd c:\Facultad\Alumno\PS\frontend
npm start
```

Accede a `http://localhost:4200/reports` y deberías ver la página de reportes funcionando.

**¡A trabajar!** 🚀📊💪
