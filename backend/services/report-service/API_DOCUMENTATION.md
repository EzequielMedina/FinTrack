# FinTrack Report Service - Documentación de Endpoints

## Base URL
```
http://localhost:8085/api/v1
```

## Endpoints

### 1. Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "report-service"
}
```

---

### 2. Reporte de Transacciones
Obtiene un reporte detallado de transacciones por período con análisis y métricas.

```http
GET /api/v1/reports/transactions
```

**Query Parameters:**
- `user_id` (string, required): ID del usuario
- `start_date` (string, optional): Fecha inicio formato YYYY-MM-DD (default: inicio del mes actual)
- `end_date` (string, optional): Fecha fin formato YYYY-MM-DD (default: fecha actual)
- `type` (string, optional): Tipo de transacción específico
- `group_by` (string, optional): Agrupar por day/week/month

**Ejemplo:**
```bash
curl "http://localhost:8085/api/v1/reports/transactions?user_id=123&start_date=2024-01-01&end_date=2024-01-31"
```

**Response:**
```json
{
  "user_id": "123",
  "period": {
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "days": 31
  },
  "summary": {
    "total_transactions": 45,
    "total_income": 150000.00,
    "total_expenses": 95000.00,
    "net_balance": 55000.00,
    "avg_transaction": 3333.33
  },
  "by_type": [
    {
      "type": "credit_charge",
      "count": 15,
      "amount": 45000.00,
      "percentage": 47.37
    }
  ],
  "by_period": [
    {
      "period": "2024-01-15",
      "date": "2024-01-15T00:00:00Z",
      "income": 5000.00,
      "expenses": 3000.00,
      "net": 2000.00,
      "count": 3
    }
  ],
  "top_expenses": [
    {
      "id": "tx-001",
      "description": "Compra en supermercado",
      "amount": 8500.00,
      "type": "debit_purchase",
      "date": "2024-01-20T10:30:00Z",
      "merchant_name": "Carrefour"
    }
  ]
}
```

---

### 3. Reporte de Cuotas
Obtiene un reporte detallado de planes de cuotas, pagos próximos y vencidos.

```http
GET /api/v1/reports/installments
```

**Query Parameters:**
- `user_id` (string, required): ID del usuario
- `status` (string, optional): Filtrar por estado (active, completed, overdue)

**Ejemplo:**
```bash
curl "http://localhost:8085/api/v1/reports/installments?user_id=123&status=active"
```

**Response:**
```json
{
  "user_id": "123",
  "summary": {
    "total_plans": 5,
    "active_plans": 3,
    "total_amount": 120000.00,
    "paid_amount": 50000.00,
    "remaining_amount": 70000.00,
    "overdue_amount": 5000.00,
    "next_payment_amount": 8000.00,
    "next_payment_date": "2024-02-15T00:00:00Z",
    "completion_percentage": 41.67
  },
  "plans": [
    {
      "id": "plan-001",
      "card_id": "card-001",
      "card_last_four": "1234",
      "total_amount": 40000.00,
      "installments_count": 12,
      "installment_amount": 3333.33,
      "paid_installments": 5,
      "remaining_amount": 23333.35,
      "status": "active",
      "description": "Smart TV Samsung",
      "merchant_name": "Garbarino",
      "start_date": "2023-09-15T00:00:00Z",
      "next_due_date": "2024-02-15T00:00:00Z",
      "completion_percentage": 41.67
    }
  ],
  "upcoming_payments": [
    {
      "installment_id": "inst-015",
      "plan_id": "plan-001",
      "card_last_four": "1234",
      "amount": 3333.33,
      "due_date": "2024-02-15T00:00:00Z",
      "days_until_due": 5,
      "description": "Smart TV Samsung",
      "merchant_name": "Garbarino"
    }
  ],
  "overdue_payments": [
    {
      "installment_id": "inst-012",
      "plan_id": "plan-002",
      "card_last_four": "5678",
      "amount": 5000.00,
      "due_date": "2024-01-05T00:00:00Z",
      "days_overdue": 35,
      "late_fee": 250.00,
      "description": "Notebook HP",
      "merchant_name": "Fravega"
    }
  ]
}
```

---

### 4. Reporte de Cuentas
Obtiene un resumen completo de todas las cuentas y tarjetas del usuario.

```http
GET /api/v1/reports/accounts
```

**Query Parameters:**
- `user_id` (string, required): ID del usuario

**Ejemplo:**
```bash
curl "http://localhost:8085/api/v1/reports/accounts?user_id=123"
```

**Response:**
```json
{
  "user_id": "123",
  "summary": {
    "total_balance": 250000.00,
    "total_accounts": 4,
    "total_cards": 3,
    "total_credit_limit": 150000.00,
    "total_credit_used": 45000.00,
    "available_credit": 105000.00,
    "credit_utilization": 30.00,
    "net_worth": 205000.00
  },
  "accounts": [
    {
      "id": "acc-001",
      "account_type": "wallet",
      "name": "Billetera Virtual",
      "currency": "ARS",
      "balance": 50000.00,
      "credit_limit": 0,
      "is_active": true
    },
    {
      "id": "acc-002",
      "account_type": "credit",
      "name": "Tarjeta Visa",
      "currency": "ARS",
      "balance": 0,
      "credit_limit": 100000.00,
      "is_active": true
    }
  ],
  "cards": [
    {
      "id": "card-001",
      "account_id": "acc-002",
      "card_type": "credit",
      "card_brand": "visa",
      "last_four_digits": "1234",
      "holder_name": "Juan Perez",
      "status": "active",
      "credit_limit": 100000.00,
      "current_balance": 45000.00,
      "available_credit": 55000.00,
      "nickname": "Visa Gold"
    }
  ],
  "distribution": [
    {
      "account_type": "wallet",
      "count": 1,
      "total_balance": 50000.00,
      "percentage": 20.00
    },
    {
      "account_type": "bank_account",
      "count": 2,
      "total_balance": 150000.00,
      "percentage": 60.00
    },
    {
      "account_type": "credit",
      "count": 1,
      "total_balance": 50000.00,
      "percentage": 20.00
    }
  ]
}
```

---

### 5. Reporte de Gastos vs Ingresos
Obtiene un análisis detallado de gastos e ingresos con tendencias y proyecciones.

```http
GET /api/v1/reports/expenses-income
```

**Query Parameters:**
- `user_id` (string, required): ID del usuario
- `start_date` (string, required): Fecha inicio formato YYYY-MM-DD
- `end_date` (string, required): Fecha fin formato YYYY-MM-DD
- `group_by` (string, optional): Agrupar por day/week/month

**Ejemplo:**
```bash
curl "http://localhost:8085/api/v1/reports/expenses-income?user_id=123&start_date=2024-01-01&end_date=2024-01-31"
```

**Response:**
```json
{
  "user_id": "123",
  "period": {
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "days": 31
  },
  "summary": {
    "total_income": 150000.00,
    "total_expenses": 95000.00,
    "net_balance": 55000.00,
    "savings_rate": 36.67,
    "expense_ratio": 63.33,
    "avg_daily_income": 4838.71,
    "avg_daily_expense": 3064.52
  },
  "by_period": [
    {
      "period": "2024-01-15",
      "date": "2024-01-15T00:00:00Z",
      "income": 5000.00,
      "expenses": 3000.00,
      "net": 2000.00,
      "savings_rate": 40.00
    }
  ],
  "by_category": [
    {
      "category": "wallet_deposit",
      "type": "income",
      "amount": 80000.00,
      "count": 10,
      "percentage": 53.33
    },
    {
      "category": "credit_charge",
      "type": "expense",
      "amount": 45000.00,
      "count": 15,
      "percentage": 47.37
    }
  ],
  "trend": {
    "incomes_trend": "increasing",
    "expenses_trend": "stable",
    "net_trend": "improving",
    "income_change": 15.5,
    "expense_change": 2.3,
    "forecast": {
      "next_month_income": 172500.00,
      "next_month_expenses": 97185.00,
      "next_month_net": 75315.00
    }
  }
}
```

---

### 6. Reporte de Notificaciones
Obtiene estadísticas de envío de notificaciones y efectividad del sistema.

```http
GET /api/v1/reports/notifications
```

**Query Parameters:**
- `start_date` (string, optional): Fecha inicio formato YYYY-MM-DD (default: últimos 30 días)
- `end_date` (string, optional): Fecha fin formato YYYY-MM-DD (default: fecha actual)

**Ejemplo:**
```bash
curl "http://localhost:8085/api/v1/reports/notifications?start_date=2024-01-01&end_date=2024-01-31"
```

**Response:**
```json
{
  "period": {
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z",
    "days": 31
  },
  "summary": {
    "total_notifications": 250,
    "total_job_runs": 31,
    "successful_sent": 245,
    "failed": 5,
    "success_rate": 98.00,
    "failure_rate": 2.00,
    "avg_emails_per_run": 8.06
  },
  "by_day": [
    {
      "date": "2024-01-15T00:00:00Z",
      "day": "2024-01-15",
      "sent": 8,
      "failed": 0,
      "success_rate": 100.00
    }
  ],
  "by_status": [
    {
      "status": "sent",
      "count": 245,
      "percentage": 98.00
    },
    {
      "status": "failed",
      "count": 5,
      "percentage": 2.00
    }
  ],
  "job_runs": [
    {
      "id": "job-run-031",
      "started_at": "2024-01-31T08:00:00Z",
      "completed_at": "2024-01-31T08:02:15Z",
      "status": "completed",
      "cards_found": 8,
      "emails_sent": 8,
      "errors": 0,
      "duration": "2m15s"
    }
  ]
}
```

---

## Códigos de Estado HTTP

- `200 OK`: Petición exitosa
- `400 Bad Request`: Parámetros inválidos o faltantes
- `500 Internal Server Error`: Error del servidor

## Notas de Uso

1. **Fechas**: Todas las fechas deben estar en formato `YYYY-MM-DD`
2. **Períodos**: Si no se especifican fechas, se usan valores por defecto razonables
3. **Paginación**: Actualmente no implementada, se devuelven todos los resultados
4. **CORS**: Configurado para permitir `http://localhost:4200`
5. **Autenticación**: Actualmente no implementada (TODO: agregar JWT)

## Para Chart.js

Los datos están estructurados para ser fácilmente consumidos por Chart.js:

### Ejemplo para Gráfico de Torta (Transacciones por Tipo)
```javascript
const chartData = {
  labels: response.by_type.map(item => item.type),
  datasets: [{
    data: response.by_type.map(item => item.amount),
    backgroundColor: ['#FF6384', '#36A2EB', '#FFCE56', '#4BC0C0', '#9966FF']
  }]
};
```

### Ejemplo para Gráfico de Líneas (Gastos vs Ingresos)
```javascript
const chartData = {
  labels: response.by_period.map(item => item.period),
  datasets: [
    {
      label: 'Ingresos',
      data: response.by_period.map(item => item.income),
      borderColor: '#36A2EB',
      fill: false
    },
    {
      label: 'Gastos',
      data: response.by_period.map(item => item.expenses),
      borderColor: '#FF6384',
      fill: false
    }
  ]
};
```
