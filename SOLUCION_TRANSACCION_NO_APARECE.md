# 🔧 Solución: Transacción No Aparece en Reportes

## 🔍 Problema Identificado

**Tu transacción NO aparece en el reporte porque tiene un `user_id` incorrecto.**

### Detalles de la Transacción

```sql
ID: txn_1760973884851164923
Tipo: wallet_deposit
Monto: $20,001.00
Descripción: test
Fecha: 2025-10-20 15:24:44
Status: completed ✅
user_id: default-user ❌ (PROBLEMA)
```

**Tu user_id real es:** `6a67040e-79fe-4b98-8980-1929f2b5b8bb`  
**El user_id guardado es:** `default-user`

---

## 🎯 Causa Raíz

El problema está en el **servicio que crea las transacciones** (transaction-service o account-service), NO en el servicio de reportes.

Cuando agregas saldo a una cuenta, el servicio está guardando la transacción con `user_id = "default-user"` en lugar de usar el ID del usuario autenticado.

---

## ✅ Verificación

### El servicio de reportes funciona correctamente:

```bash
# Búsqueda con tu user_id real (NO encuentra la transacción)
curl "http://localhost:8085/api/v1/reports/transactions?user_id=6a67040e-79fe-4b98-8980-1929f2b5b8bb&start_date=2025-10-20&end_date=2025-10-20"
# Resultado: 0 transacciones ✅ (correcto, porque no hay ninguna con tu user_id)

# La transacción existe en la BD pero con user_id incorrecto
SELECT * FROM transactions WHERE id = 'txn_1760973884851164923';
# user_id: default-user ❌
```

---

## 🔧 Soluciones

### Solución 1: Corregir la Transacción Existente (Temporal)

Actualiza el user_id de la transacción manualmente:

```bash
docker-compose exec mysql mysql -u fintrack_user -pfintrack_password fintrack -e "UPDATE transactions SET user_id = '6a67040e-79fe-4b98-8980-1929f2b5b8bb', initiated_by = '6a67040e-79fe-4b98-8980-1929f2b5b8bb' WHERE id = 'txn_1760973884851164923';"
```

Después de esto, recarga el reporte y debería aparecer.

---

### Solución 2: Corregir el Código del Servicio (Permanente)

Necesitas encontrar dónde se crea la transacción de tipo `wallet_deposit` y asegurarte de que use el user_id correcto.

#### Paso 1: Identificar el Servicio Responsable

Busca en el código:

```bash
# Desde la raíz del proyecto
grep -r "wallet_deposit" backend/services/
grep -r "default-user" backend/services/
```

#### Paso 2: Buscar la Función que Crea la Transacción

Probablemente esté en:
- `backend/services/account-service/` (cuando agregas saldo a una cuenta)
- `backend/services/transaction-service/` (cuando creas transacciones)

Busca algo como:

```go
// ❌ INCORRECTO
transaction := &Transaction{
    UserID: "default-user",  // ← Hardcodeado!
    Type: "wallet_deposit",
    // ...
}

// ✅ CORRECTO
transaction := &Transaction{
    UserID: authenticatedUserID,  // ← Desde el contexto/token JWT
    Type: "wallet_deposit",
    // ...
}
```

#### Paso 3: Obtener el User ID Correcto

El user_id debería venir de:

1. **Token JWT** (si están usando autenticación):
```go
// En el middleware o handler
userID := c.GetString("user_id")  // Gin
// o
userID := ctx.Value("user_id").(string)  // Context estándar
```

2. **Parámetro de la request**:
```go
userID := c.Query("user_id")
// o
userID := request.UserID
```

3. **De la cuenta destino** (to_account_id):
```sql
SELECT user_id FROM accounts WHERE id = ?
```

---

## 🔍 Dónde Buscar

### Archivos Probables:

```
backend/services/account-service/
├── internal/
│   ├── handlers/
│   │   └── account_handler.go  ← Busca la función que agrega saldo
│   ├── service/
│   │   └── account_service.go  ← Lógica de negocio
│   └── repository/
│       └── account_repository.go
```

O en:

```
backend/services/transaction-service/
├── internal/
│   ├── handlers/
│   │   └── transaction_handler.go
│   ├── service/
│   │   └── transaction_service.go
│   └── repository/
│       └── transaction_repository.go
```

### Busca funciones como:

- `CreateTransaction`
- `AddBalance`
- `Deposit`
- `CreateWalletDeposit`

---

## 📝 Ejemplo de Corrección

### Antes (Incorrecto):

```go
func (s *AccountService) AddBalance(accountID string, amount float64) error {
    // Crear transacción
    transaction := &Transaction{
        ID:          generateID(),
        Type:        "wallet_deposit",
        UserID:      "default-user",  // ❌ PROBLEMA
        InitiatedBy: "default-user",  // ❌ PROBLEMA
        ToAccountID: accountID,
        Amount:      amount,
        Status:      "completed",
        CreatedAt:   time.Now(),
    }
    
    return s.repo.CreateTransaction(transaction)
}
```

### Después (Correcto):

```go
func (s *AccountService) AddBalance(ctx context.Context, accountID string, amount float64) error {
    // Obtener el user_id del contexto (del JWT)
    userID := ctx.Value("user_id").(string)
    
    // O obtenerlo de la cuenta
    account, err := s.repo.GetAccount(accountID)
    if err != nil {
        return err
    }
    userID := account.UserID
    
    // Crear transacción con el user_id correcto
    transaction := &Transaction{
        ID:          generateID(),
        Type:        "wallet_deposit",
        UserID:      userID,        // ✅ CORRECTO
        InitiatedBy: userID,        // ✅ CORRECTO
        ToAccountID: accountID,
        Amount:      amount,
        Status:      "completed",
        CreatedAt:   time.Now(),
    }
    
    return s.repo.CreateTransaction(transaction)
}
```

---

## ✅ Verificación Final

Después de corregir el código:

### 1. Crea una nueva transacción
```
Agrega saldo a una cuenta desde la UI
```

### 2. Verifica en la BD
```sql
SELECT id, type, user_id, amount, created_at 
FROM transactions 
ORDER BY created_at DESC 
LIMIT 1;
```

Debería mostrar tu user_id real, no "default-user"

### 3. Verifica en el reporte
```
Ve a http://localhost:4200/reports/transactions
```

La transacción debería aparecer ahora.

---

## 🚀 Solución Rápida (Para Probar Ahora)

Ejecuta este comando para corregir la transacción actual:

```bash
docker-compose exec mysql mysql -u fintrack_user -pfintrack_password fintrack -e "UPDATE transactions SET user_id = '6a67040e-79fe-4b98-8980-1929f2b5b8bb', initiated_by = '6a67040e-79fe-4b98-8980-1929f2b5b8bb' WHERE id = 'txn_1760973884851164923';"
```

Luego:
1. Ve a http://localhost:4200/reports/transactions
2. Presiona Ctrl + Shift + R para limpiar cache
3. Selecciona fecha: 2025-10-20
4. Deberías ver tu transacción de $20,001.00

---

## 📊 Resumen

| Item | Estado | Nota |
|------|--------|------|
| **Servicio de Reportes** | ✅ Funciona correctamente | Busca transacciones por user_id |
| **Base de Datos** | ✅ Transacción existe | Pero con user_id incorrecto |
| **Problema Real** | ⚠️ Bug en creación de transacciones | Usa "default-user" en lugar del user_id real |
| **Solución Temporal** | 🔧 UPDATE manual | Corrige esta transacción |
| **Solución Permanente** | 🔧 Fix en el código | Corrige el servicio que crea transacciones |

---

## 🎯 Próximos Pasos

1. ✅ **Inmediato:** Ejecuta el UPDATE para ver la transacción en el reporte
2. 🔧 **Corto plazo:** Encuentra y corrige el código que crea transacciones
3. ✅ **Verificación:** Crea una nueva transacción y verifica que use el user_id correcto
4. 🧪 **Testing:** Agrega tests para verificar que el user_id se guarde correctamente

---

**Nota:** El servicio de reportes está funcionando perfectamente. El problema está en otro servicio que no está guardando el user_id correctamente cuando crea transacciones.
