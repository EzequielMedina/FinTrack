# Flujo de Trabajo - Transacciones de Cuenta (Agregar/Quitar Dinero)

## 📋 Resumen

Este documento describe el flujo completo de trabajo cuando un usuario entra a su cuenta y realiza una transacción (agregar o quitar dinero), desde el frontend hasta la persistencia en las bases de datos, pasando por todos los microservicios involucrados.

## 🏗️ Arquitectura del Sistema

### Microservicios Involucrados

1. **Frontend (Angular)** - Puerto 4200
2. **User Service (Go)** - Puerto 8081
3. **Account Service (Go)** - Puerto 8082  
4. **Transaction Service (Go)** - Puerto 8083
5. **Base de Datos MySQL** - Puerto 3306

## 🔄 Flujo Completo de Transacciones

### 1. Frontend - Interfaz de Usuario

#### Componentes Principales:
- **`transactions.component.ts`**: Componente principal que maneja las transacciones
- **`transactions.component.html`**: UI con formularios para depósito, retiro, transferencia y pago
- **`transaction.service.ts`**: Servicio para comunicación con APIs
- **`account.service.ts`**: Servicio para operaciones de cuenta

#### Tipos de Transacciones Disponibles:
```typescript
// Formularios disponibles en el frontend
- Transferencia (transfer)
- Depósito (deposit) 
- Retiro (withdrawal)
- Pago (payment)
```

#### Modelos de Datos:
```typescript
// transaction.model.ts
interface AddFundsRequest {
  amount: number;
  description: string;
  reference?: string;
}

interface WithdrawFundsRequest {
  amount: number;
  description: string;
  reference?: string;
}
```

### 2. Flujo de Agregar Dinero (Depósito)

#### 2.1 Frontend → Account Service
```
POST http://localhost:8082/api/accounts/{id}/add-funds
```

**Request Body:**
```json
{
  "amount": 100.00,
  "description": "Depósito desde billetera",
  "reference": "REF001"
}
```

#### 2.2 Account Service - Procesamiento

**Handler:** `account/handler.go → AddFunds()`

**Validaciones:**
- Verificar que la cuenta existe
- Validar que la cuenta está activa
- Validar que el monto es positivo
- Aplicar lógica específica por tipo de cuenta

**Tipos de Cuenta Soportados:**
```go
switch accountType {
case "wallet":
    // Billetera: Incremento directo del balance
case "savings", "checking", "bank_account":
    // Cuentas bancarias: Depósito directo
case "credit":
    // Tarjeta de crédito: Pago (reduce crédito usado)
case "debit":
    // Tarjeta de débito: Incremento directo
}
```

#### 2.3 Account Service → Transaction Service

**Comunicación HTTP:**
```
POST http://localhost:8083/api/transactions
```

**Request Body:**
```json
{
  "type": "account_deposit",
  "user_id": "user-uuid",
  "amount": 100.00,
  "currency": "USD",
  "to_account_id": "account-uuid",
  "description": "Depósito desde billetera",
  "payment_method": "wallet"
}
```

#### 2.4 Transaction Service - Procesamiento

**Handler:** `transaction_handler.go → CreateTransaction()`

**Flujo de Procesamiento:**
1. **Validación de Usuario**: Verificar que el usuario existe
2. **Creación de Transacción**: Estado inicial `PENDING`
3. **Validaciones Pre-transacción**: Verificar fondos/límites
4. **Ejecución de Balance**: Comunicación con Account Service
5. **Actualización de Estado**: Cambiar a `COMPLETED`
6. **Auditoría**: Registrar la transacción

**Integración con Account Service:**
```go
// Cliente HTTP interno
accountService := clients.NewAccountClient("http://account-service:8082")

// Métodos disponibles:
- GetAccountBalance()
- AddFunds()
- WithdrawFunds()
- ValidateAccountExists()
```

### 3. Flujo de Quitar Dinero (Retiro)

#### 3.1 Frontend → Account Service
```
POST http://localhost:8082/api/accounts/{id}/withdraw-funds
```

**Request Body:**
```json
{
  "amount": 50.00,
  "description": "Retiro de efectivo",
  "reference": "REF002"
}
```

#### 3.2 Account Service - Procesamiento

**Handler:** `account/handler.go → WithdrawFunds()`

**Validaciones:**
- Verificar que la cuenta existe
- Validar que la cuenta está activa
- Verificar fondos suficientes
- Validar tipos de cuenta permitidos

**Tipos de Cuenta Permitidos para Retiro:**
```go
switch accountType {
case "wallet", "savings", "checking", "bank_account", "debit":
    // Permitido
case "credit":
    // NO permitido - retiros no soportados en crédito
}
```

#### 3.3 Account Service → Transaction Service

**Request Body:**
```json
{
  "type": "account_withdraw",
  "user_id": "user-uuid",
  "amount": 50.00,
  "currency": "USD",
  "from_account_id": "account-uuid",
  "description": "Retiro de efectivo",
  "payment_method": "wallet"
}
```

## 🗄️ Persistencia en Base de Datos

### Estructura de Tablas

#### Tabla `accounts`
```sql
CREATE TABLE accounts (
  id VARCHAR(36) PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL,
  account_type VARCHAR(20) NOT NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT,
  currency VARCHAR(3) NOT NULL,
  balance DECIMAL(15,2) NOT NULL DEFAULT 0,
  credit_limit DECIMAL(15,2) NULL,
  closing_date DATE NULL,
  due_date DATE NULL,
  dni VARCHAR(20) NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  deleted_at DATETIME NULL,
  INDEX idx_user_id (user_id),
  INDEX idx_account_type (account_type),
  INDEX idx_currency (currency),
  INDEX idx_is_active (is_active)
);
```

#### Tabla `transactions`
```sql
CREATE TABLE transactions (
  id VARCHAR(36) PRIMARY KEY,
  reference_id VARCHAR(100),
  external_id VARCHAR(100),
  type VARCHAR(50) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  amount DECIMAL(15,2) NOT NULL,
  currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
  from_account_id VARCHAR(36),
  to_account_id VARCHAR(36),
  from_card_id VARCHAR(36),
  to_card_id VARCHAR(36),
  user_id VARCHAR(36) NOT NULL,
  initiated_by VARCHAR(36) NOT NULL,
  description TEXT,
  payment_method VARCHAR(30),
  merchant_name VARCHAR(255),
  merchant_id VARCHAR(100),
  previous_balance DECIMAL(15,2),
  new_balance DECIMAL(15,2),
  processed_at DATETIME,
  failed_at DATETIME,
  failure_reason TEXT,
  metadata JSON,
  tags JSON,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  INDEX idx_type (type),
  INDEX idx_status (status),
  INDEX idx_user_id (user_id),
  INDEX idx_from_account (from_account_id),
  INDEX idx_to_account (to_account_id),
  INDEX idx_reference_id (reference_id),
  INDEX idx_external_id (external_id)
);
```

### Operaciones de Base de Datos

#### Para Agregar Dinero:
1. **Account Service**: 
   ```sql
   UPDATE accounts 
   SET balance = balance + ?, updated_at = NOW() 
   WHERE id = ? AND is_active = 1;
   ```

2. **Transaction Service**:
   ```sql
   INSERT INTO transactions (
     id, type, status, amount, currency, to_account_id, 
     user_id, initiated_by, description, created_at, updated_at
   ) VALUES (?, 'account_deposit', 'pending', ?, ?, ?, ?, ?, ?, NOW(), NOW());
   
   -- Después del procesamiento exitoso:
   UPDATE transactions 
   SET status = 'completed', processed_at = NOW(), 
       previous_balance = ?, new_balance = ?
   WHERE id = ?;
   ```

#### Para Quitar Dinero:
1. **Account Service**:
   ```sql
   UPDATE accounts 
   SET balance = balance - ?, updated_at = NOW() 
   WHERE id = ? AND is_active = 1 AND balance >= ?;
   ```

2. **Transaction Service**:
   ```sql
   INSERT INTO transactions (
     id, type, status, amount, currency, from_account_id,
     user_id, initiated_by, description, created_at, updated_at
   ) VALUES (?, 'account_withdraw', 'pending', ?, ?, ?, ?, ?, ?, NOW(), NOW());
   ```

## 🔄 Tipos de Transacciones Soportados

### Depósitos:
- `wallet_deposit`: Depósito en billetera virtual
- `account_deposit`: Depósito en cuenta bancaria

### Retiros:
- `wallet_withdrawal`: Retiro de billetera virtual
- `account_withdraw`: Retiro de cuenta bancaria

### Transferencias:
- `wallet_transfer`: Transferencia entre billeteras
- `account_transfer`: Transferencia entre cuentas

### Compras/Pagos:
- `credit_charge`: Cargo en tarjeta de crédito
- `debit_purchase`: Compra con tarjeta de débito
- `credit_payment`: Pago de tarjeta de crédito

## 🛡️ Validaciones y Seguridad

### Validaciones de Negocio:
- **Fondos Suficientes**: Verificación antes de retiros
- **Cuentas Activas**: Solo cuentas activas pueden operar
- **Límites de Crédito**: Validación para tarjetas de crédito
- **Tipos de Cuenta**: Operaciones permitidas por tipo

### Manejo de Errores:
- **Rollback Automático**: En caso de fallo en transferencias
- **Estados de Transacción**: `pending`, `completed`, `failed`, `canceled`, `reversed`
- **Auditoría Completa**: Registro de todas las operaciones

## 🔧 Configuración de Servicios

### Variables de Entorno:
```env
# Account Service
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password
PORT=8082

# Transaction Service  
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password
PORT=8083
ACCOUNT_SERVICE_URL=http://account-service:8082
```

### Puertos de Comunicación:
- **Frontend**: 4200 → Account Service (8082)
- **Account Service**: 8082 → Transaction Service (8083)
- **Transaction Service**: 8083 → Account Service (8082)
- **Todos los servicios**: → MySQL (3306)

## 📊 Flujo de Datos Resumido

```
1. Usuario completa formulario en Frontend (Angular)
   ↓
2. Frontend envía request a Account Service (HTTP REST)
   ↓
3. Account Service valida y actualiza balance en MySQL
   ↓
4. Account Service notifica a Transaction Service (HTTP)
   ↓
5. Transaction Service crea registro de transacción en MySQL
   ↓
6. Transaction Service actualiza estado a 'completed'
   ↓
7. Respuesta exitosa regresa al Frontend
   ↓
8. Frontend actualiza la UI con el nuevo balance
```

## 🎯 Puntos Clave

- **Arquitectura de Microservicios**: Separación clara de responsabilidades
- **Comunicación HTTP**: REST APIs entre servicios
- **Transacciones ACID**: Consistencia en operaciones de base de datos
- **Auditoría Completa**: Registro detallado de todas las transacciones
- **Manejo de Errores**: Rollback automático y estados de error
- **Escalabilidad**: Servicios independientes y desplegables por separado