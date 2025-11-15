# Plan de Implementación - Transacciones de Tarjetas

## Resumen Ejecutivo

Este documento presenta el plan detallado para implementar la funcionalidad completa de transacciones de tarjetas en el sistema FinTrack, integrando los servicios existentes y asegurando la consistencia de datos entre microservicios.

## Estado Actual del Sistema

### ✅ Implementado
- **Frontend**: Componentes completos para operaciones de tarjetas (`card-detail.component`)
- **Backend**: Endpoints básicos en `account-service` para:
  - Cargos de tarjeta de crédito (`ChargeCard`)
  - Pagos de tarjeta de crédito (`PaymentCard`) 
  - Transacciones de tarjeta de débito (`ProcessDebitTransaction`)
- **Servicios Frontend**: `CreditCardService` y `DebitCardService` funcionales
- **Base de Datos**: Esquemas de tarjetas y transacciones implementados

### ❌ Faltante
- **Integración con Transaction Service**: Las operaciones de tarjetas no se registran como transacciones
- **Conexión con Account Service**: Las tarjetas de débito no afectan el balance de la cuenta
- **Lógica de negocio completa**: Falta integración entre microservicios

## Arquitectura Objetivo

```
Frontend (Angular)
    ↓
Transaction Service (Maneja toda la lógica)
    ↓
Account Service (Actualización de Balances y Tarjetas)
```

**Flujo Simplificado:**
1. Frontend llama a Transaction Service
2. Transaction Service valida y procesa la transacción
3. Transaction Service llama a Account Service para actualizar balances/tarjetas
4. Transaction Service registra la transacción en su BD

## Plan de Desarrollo

### Fase 1: Crear AccountClient en Transaction Service

#### 1.1 Crear AccountClient en Transaction Service

**Archivo**: `backend/services/transaction-service/internal/clients/account_client.go`

```go
type AccountClient struct {
    baseURL string
    client  *http.Client
}

type CardChargeRequest struct {
    CardID      string  `json:"card_id"`
    Amount      float64 `json:"amount"`
    Description string  `json:"description"`
    Reference   string  `json:"reference,omitempty"`
}

type CardPaymentRequest struct {
    CardID        string  `json:"card_id"`
    Amount        float64 `json:"amount"`
    PaymentMethod string  `json:"payment_method"`
    Reference     string  `json:"reference,omitempty"`
}

type DebitTransactionRequest struct {
    CardID       string  `json:"card_id"`
    Amount       float64 `json:"amount"`
    Description  string  `json:"description"`
    MerchantName string  `json:"merchant_name"`
    Reference    string  `json:"reference,omitempty"`
}
```

**Métodos a implementar**:
- `ChargeCard(req CardChargeRequest) error`
- `PaymentCard(req CardPaymentRequest) error`
- `ProcessDebitTransaction(req DebitTransactionRequest) error`

#### 1.2 Configuración

**Archivo**: `backend/services/transaction-service/.env`
```
ACCOUNT_SERVICE_URL=http://account-service:8081
```

### Fase 2: Implementar Endpoints de Tarjetas en Transaction Service

#### 2.1 Crear nuevos endpoints en Transaction Service

**Archivo**: `backend/services/transaction-service/internal/handlers/card_handler.go` (NUEVO)

```go
func (h *CardHandler) ChargeCard(c *gin.Context) {
    var req CardChargeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    result, err := h.service.ProcessCreditCardCharge(req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, result)
}
```

#### 2.2 Modificar TransactionService

**Archivo**: `backend/services/transaction-service/internal/services/transaction_service_impl.go`

**Nuevos métodos**:
```go
func (s *TransactionServiceImpl) ProcessCreditCardCharge(req CardChargeRequest) (*Transaction, error) {
    // 1. Crear transacción
    transaction := &models.Transaction{
        CardID:      req.CardID,
        Amount:      req.Amount,
        Type:        "credit_card_charge",
        Description: req.Description,
        Reference:   req.Reference,
        Status:      "pending",
    }
    
    // 2. Guardar en BD
    if err := s.repo.CreateTransaction(transaction); err != nil {
        return nil, err
    }
    
    // 3. Procesar en Account Service
    if err := s.accountClient.ChargeCard(req); err != nil {
        s.repo.UpdateTransactionStatus(transaction.ID, "failed")
        return nil, err
    }
    
    // 4. Marcar como completada
    s.repo.UpdateTransactionStatus(transaction.ID, "completed")
    return transaction, nil
}
```

### Fase 3: Reutilización de Componentes Existentes

#### 3.1 Endpoints Reutilizables

**Account Service**:
- ✅ `POST /api/accounts/{id}/add-funds` (para pagos de crédito)
- ✅ `POST /api/accounts/{id}/withdraw-funds` (para débito)

**Transaction Service**:
- ✅ `POST /api/transactions` (base para registrar transacciones)

#### 3.2 Servicios Frontend Reutilizables

**Archivos existentes**:
- `frontend/src/app/core/services/account.service.ts`
- `frontend/src/app/core/services/transaction.service.ts`
- `frontend/src/app/core/services/wallet.service.ts`

**Integración**: Los servicios de tarjetas ya utilizan estos servicios base.

### Fase 4: Configuración y Variables de Entorno

#### 4.1 Docker Compose

**Archivo**: `docker-compose.yml`

```yaml
services:
  account-service:
    environment:
      - TRANSACTION_SERVICE_URL=http://transaction-service:8083
      
  transaction-service:
    environment:
      - ACCOUNT_SERVICE_URL=http://account-service:8081
```

#### 4.2 Configuración de Red

Asegurar que los servicios puedan comunicarse:
```yaml
networks:
  fintrack-network:
    driver: bridge
```

## Lógica Específica por Tipo de Tarjeta

### Tarjetas de Débito

**Flujo**:
1. Usuario realiza compra → `ProcessDebitTransaction` (Account Service)
2. Account Service → `ProcessDebitCardPurchase` (Transaction Service)
3. Transaction Service → `WithdrawFunds` (Account Service)
4. Actualización de balance de cuenta en tiempo real

**Validaciones**:
- Balance suficiente en cuenta asociada
- Tarjeta activa y no bloqueada
- Límites diarios/mensuales

### Tarjetas de Crédito

**Flujo Cargos**:
1. Usuario realiza compra → `ChargeCard` (Account Service)
2. Account Service → `CreateCreditCardCharge` (Transaction Service)
3. Actualización de balance de crédito utilizado

**Flujo Pagos**:
1. Usuario realiza pago → `PaymentCard` (Account Service)
2. Account Service → `CreateCreditCardPayment` (Transaction Service)
3. Reducción de deuda de tarjeta de crédito

**Validaciones**:
- Límite de crédito disponible
- Tarjeta activa y no bloqueada
- Monto mínimo de pago

## Cronograma de Implementación

### Día 1-2: Backend Integration - Transaction Service
- [ ] Crear `TransactionClient` en Account Service
- [ ] Modificar `CardService` para integrar con Transaction Service
- [ ] Configurar variables de entorno
- [ ] Testing unitario

### Día 3-4: Backend Integration - Account Service  
- [ ] Crear `AccountClient` en Transaction Service
- [ ] Implementar `ProcessDebitCardPurchase`
- [ ] Modificar endpoints existentes
- [ ] Testing de integración

### Día 5: Testing y Verificación Final
- [ ] Testing end-to-end completo
- [ ] Verificación de consistencia de datos
- [ ] Testing de rollback en caso de errores
- [ ] Documentación de APIs

## Ventajas de este Enfoque

### ✅ Reutilización Máxima
- Aprovecha endpoints existentes de `add-funds` y `withdraw-funds`
- Utiliza la infraestructura de `Transaction Service` ya implementada
- Mantiene la consistencia con flujos existentes

### ✅ No Rompe Funcionalidad Existente
- Los cambios son aditivos, no modifican lógica existente
- Mantiene compatibilidad con APIs actuales
- Frontend no requiere cambios significativos

### ✅ Arquitectura Consistente
- Sigue el patrón de microservicios establecido
- Mantiene separación de responsabilidades
- Facilita mantenimiento y escalabilidad

### ✅ Testing Sencillo
- Cada fase se puede testear independientemente
- Reutiliza casos de prueba existentes
- Permite rollback fácil en caso de problemas

## Archivos a Modificar

### Backend - Account Service
- `internal/clients/transaction_client.go` (NUEVO)
- `internal/services/card_service.go` (MODIFICAR)
- `internal/handlers/card_handler.go` (MODIFICAR)
- `.env` (MODIFICAR)

### Backend - Transaction Service
- `internal/clients/account_client.go` (NUEVO)
- `internal/services/transaction_service_impl.go` (MODIFICAR)
- `internal/handlers/transaction_handler.go` (MODIFICAR)
- `.env` (MODIFICAR)

### Configuración
- `docker-compose.yml` (MODIFICAR)

### Frontend
- Sin cambios requeridos (ya implementado)

## Comandos de Testing

```bash
# Compilar y levantar servicios
docker compose up --build

# Verificar logs
docker compose logs -f account-service
docker compose logs -f transaction-service

# Testing de integración
./test_card_integration_final.ps1
```

## Consideraciones de Seguridad

- Validación de autenticación en todos los endpoints
- Encriptación de datos sensibles de tarjetas
- Logging de todas las transacciones para auditoría
- Rate limiting para prevenir abuso
- Validación de montos y límites

## Conclusión

Este plan proporciona una implementación completa y robusta para las transacciones de tarjetas, maximizando la reutilización de código existente y manteniendo la arquitectura de microservicios. La implementación por fases permite un desarrollo incremental y testing continuo.