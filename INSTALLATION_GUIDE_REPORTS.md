# 🚀 Instrucciones de Instalación y Configuración del Report Service

## ✅ Backend - Report Service (Go)

### 1. Inicializar dependencias de Go

```bash
cd backend/services/report-service
go mod tidy
go mod download
```

### 2. Configurar variables de entorno

Crea un archivo `.env` basado en `.env.example`:

```bash
cp .env.example .env
```

Ajusta las variables si es necesario (las por defecto funcionan con Docker).

### 3. Probar el servicio localmente (Opcional)

```bash
go run cmd/main.go
```

El servicio estará disponible en `http://localhost:8085`

### 4. Build con Docker

```bash
# Desde la raíz del proyecto
docker-compose build report-service
```

### 5. Ejecutar con Docker Compose

```bash
# Ejecutar solo el report-service con sus dependencias
docker-compose up mysql report-service

# O ejecutar todos los servicios
docker-compose up
```

### 6. Verificar que funciona

```bash
# Health check
curl http://localhost:8085/health

# Probar endpoint de transacciones (reemplaza USER_ID con un ID real)
curl "http://localhost:8085/api/v1/reports/transactions?user_id=YOUR_USER_ID&start_date=2024-01-01&end_date=2024-01-31"
```

---

## ✅ Frontend - Components Angular

### 1. Instalar Chart.js (Recomendado para visualizaciones)

```bash
cd frontend
npm install chart.js ng2-charts
npm install --save-dev @types/chart.js
```

### 2. Actualizar environment.ts

Verifica que la URL del report-service esté configurada en `src/environments/environment.ts`:

```typescript
export const environment = {
  production: false,
  apiUrl: 'http://localhost',  // Nginx hace proxy a los servicios
  // ... otras configuraciones
};
```

Y en `proxy.conf.json`:

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

### 3. Agregar link a los reportes en el menú de navegación

Edita tu componente de navegación (sidebar/navbar) y agrega:

```html
<a routerLink="/reports" routerLinkActive="active">
  <i class="icon">📊</i>
  <span>Reportes</span>
</a>
```

### 4. Ejecutar el frontend

```bash
npm start
```

O con Docker:

```bash
docker-compose up frontend
```

---

## 📊 Componentes Frontend Creados

### 1. **ReportService** (`services/report.service.ts`)
- Servicio Angular para consumir la API de reportes
- Incluye interfaces TypeScript para todos los tipos de datos
- Métodos para cada tipo de reporte

### 2. **ReportsComponent** (`pages/reports/reports.component.ts`)
- Página principal de navegación de reportes
- Cards clickeables para cada tipo de reporte
- Filtro de reportes por rol (algunos solo admin)

### 3. **TransactionReportComponent** (`pages/reports/transaction-report/transaction-report.component.ts`)
- Reporte completo de transacciones con:
  - Filtros por fecha y tipo
  - Cards de resumen (ingresos, gastos, balance)
  - Gráfico de barras por tipo
  - Timeline de flujo de transacciones
  - Top 10 gastos
  - Exportación a CSV

### 4. Componentes Pendientes (TODO)
Puedes crear componentes similares para:
- `installment-report.component.ts` - Reporte de cuotas
- `account-report.component.ts` - Reporte de cuentas
- `expense-income-report.component.ts` - Análisis de gastos vs ingresos
- `notification-report.component.ts` - Reporte de notificaciones (admin)

---

## 🎨 Mejora con Chart.js (Opcional pero Recomendado)

### Instalación

```bash
npm install chart.js ng2-charts
```

### Ejemplo de uso en TransactionReportComponent

Agrega al componente:

```typescript
import { Chart, ChartConfiguration } from 'chart.js/auto';

export class TransactionReportComponent implements OnInit {
  private chart?: Chart;

  ngAfterViewInit(): void {
    this.createPieChart();
  }

  private createPieChart(): void {
    const ctx = document.getElementById('myChart') as HTMLCanvasElement;
    if (!ctx || !this.report) return;

    this.chart = new Chart(ctx, {
      type: 'pie',
      data: {
        labels: this.report.by_type.map(t => this.getTypeLabel(t.type)),
        datasets: [{
          data: this.report.by_type.map(t => t.amount),
          backgroundColor: [
            '#FF6384', '#36A2EB', '#FFCE56', 
            '#4BC0C0', '#9966FF', '#FF9F40'
          ]
        }]
      },
      options: {
        responsive: true,
        plugins: {
          legend: {
            position: 'bottom'
          },
          title: {
            display: true,
            text: 'Transacciones por Tipo'
          }
        }
      }
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
  <canvas id="myChart"></canvas>
</div>
```

---

## 🧪 Testing

### Backend

```bash
cd backend/services/report-service
go test ./...
```

### Frontend

```bash
cd frontend
ng test
```

---

## 📝 Endpoints Disponibles

### 1. **GET /api/v1/reports/transactions**
- Query params: `user_id`, `start_date`, `end_date`, `type`, `group_by`
- Retorna: Análisis completo de transacciones

### 2. **GET /api/v1/reports/installments**
- Query params: `user_id`, `status`
- Retorna: Reporte de cuotas y pagos

### 3. **GET /api/v1/reports/accounts**
- Query params: `user_id`
- Retorna: Resumen de cuentas y tarjetas

### 4. **GET /api/v1/reports/expenses-income**
- Query params: `user_id`, `start_date`, `end_date`, `group_by`
- Retorna: Análisis de gastos vs ingresos con tendencias

### 5. **GET /api/v1/reports/notifications**
- Query params: `start_date`, `end_date`
- Retorna: Estadísticas de notificaciones (admin)

Ver documentación completa en: `backend/services/report-service/API_DOCUMENTATION.md`

---

## 🐛 Troubleshooting

### Error: "could not import github.com/gin-gonic/gin"

```bash
cd backend/services/report-service
go mod tidy
go mod download
```

### Error: CORS en frontend

Verifica que `proxy.conf.json` esté configurado y que estés usando `npm start` (no `ng serve`).

### Error: Cannot find module 'chart.js'

```bash
cd frontend
npm install chart.js ng2-charts
```

### Error de base de datos

Asegúrate de que MySQL esté corriendo y las migraciones se hayan ejecutado:

```bash
docker-compose up mysql
# Espera a que esté healthy
docker-compose logs mysql | grep "ready for connections"
```

---

## 🔐 Seguridad

### TODO: Implementar Autenticación JWT

Los endpoints actualmente no tienen autenticación. Para producción, debes:

1. Agregar middleware de JWT en el backend
2. Validar que el `user_id` pertenezca al usuario autenticado
3. Implementar rate limiting
4. Agregar HTTPS

### Ejemplo de middleware JWT (Go):

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        // Validar token JWT
        // ...
        c.Next()
    }
}
```

---

## 📈 Próximos Pasos

1. ✅ Completar los componentes frontend faltantes
2. ✅ Agregar Chart.js para visualizaciones mejoradas
3. ✅ Implementar autenticación JWT
4. ✅ Agregar tests unitarios y e2e
5. ✅ Implementar caché para reportes pesados
6. ✅ Agregar exportación a PDF
7. ✅ Implementar reportes programados (scheduled reports)
8. ✅ Agregar soporte para múltiples monedas

---

## 📚 Recursos

- **Go Gin**: https://gin-gonic.com/
- **Chart.js**: https://www.chartjs.org/
- **Angular**: https://angular.io/
- **Docker**: https://docs.docker.com/

---

## 👥 Soporte

Si tienes problemas, revisa:
1. Los logs del servicio: `docker-compose logs report-service`
2. El health check: `curl http://localhost:8085/health`
3. La documentación de API en `API_DOCUMENTATION.md`

---

**¡Listo! Tu microservicio de reportes está configurado y funcionando.** 🎉
