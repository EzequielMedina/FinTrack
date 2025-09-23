# Diseño de Integración Account-Card

## Objetivo
Rediseñar la relación entre Account y Card para soportar:
- Cuentas con múltiples tarjetas (crédito/débito)
- Billeteras virtuales sin tarjetas
- Operaciones de extracción/depósito

## Arquitectura Propuesta

### 1. Backend Go - Account Entity (Actualizada)

```go
type Account struct {
    // Core identity fields
    ID            string `json:"id"`
    UserID        string `json:"userId"`
    AccountNumber string `json:"accountNumber"`

    // Account details  
    AccountName string      `json:"accountName"`
    AccountType AccountType `json:"accountType"` // wallet, bank_account
    Currency    Currency    `json:"currency"`

    // Financial information
    Balance          float64 `json:"balance"`
    AvailableBalance float64 `json:"availableBalance"`
    DailyLimit       float64 `json:"dailyLimit"`
    MonthlyLimit     float64 `json:"monthlyLimit"`

    // Cards (optional - only for bank accounts)
    Cards []Card `json:"cards,omitempty"`

    // Virtual wallet specific
    DNI *string `json:"dni,omitempty"` // Para billeteras virtuales

    // Status and metadata
    Status      AccountStatus `json:"status"`
    Description string        `json:"description,omitempty"`
    
    // Audit fields
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

type Card struct {
    ID              string    `json:"id"`
    AccountID       string    `json:"accountId"`
    CardType        CardType  `json:"cardType"` // credit, debit
    CardBrand       CardBrand `json:"cardBrand"` // visa, mastercard
    LastFourDigits  string    `json:"lastFourDigits"`
    MaskedNumber    string    `json:"maskedNumber"`
    HolderName      string    `json:"holderName"`
    ExpirationMonth int       `json:"expirationMonth"`
    ExpirationYear  int       `json:"expirationYear"`
    Status          CardStatus `json:"status"`
    IsDefault       bool      `json:"isDefault"`
    Nickname        string    `json:"nickname,omitempty"`
    
    // Credit card specific
    CreditLimit     *float64   `json:"creditLimit,omitempty"`
    ClosingDate     *time.Time `json:"closingDate,omitempty"`
    DueDate         *time.Time `json:"dueDate,omitempty"`
    
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
}

// Nuevos tipos
type AccountType string
const (
    AccountTypeWallet     AccountType = "wallet"      // Billetera virtual
    AccountTypeBankAccount AccountType = "bank_account" // Cuenta bancaria (puede tener tarjetas)
)

type CardType string
const (
    CardTypeCredit CardType = "credit"
    CardTypeDebit  CardType = "debit"
)
```

### 2. Frontend TypeScript - Models Actualizados

```typescript
// account.model.ts
export enum AccountType {
  WALLET = 'wallet',           // Billetera virtual
  BANK_ACCOUNT = 'bank_account' // Cuenta bancaria
}

export interface Account {
  id: string;
  userId: string;
  accountNumber: string;
  accountName: string;
  accountType: AccountType;
  currency: Currency;
  balance: number;
  availableBalance: number;
  dailyLimit: number;
  monthlyLimit: number;
  
  // Cards array (opcional, solo para bank_account)
  cards?: Card[];
  
  // Para billeteras virtuales
  dni?: string;
  
  status: AccountStatus;
  description?: string;
  createdAt: string;
  updatedAt: string;
}

// card.model.ts (actualizado)
export interface Card {
  id: string;
  accountId: string;           // Siempre vinculada a una cuenta
  cardType: CardType;
  cardBrand: CardBrand;
  lastFourDigits: string;
  maskedNumber: string;
  holderName: string;
  expirationMonth: number;
  expirationYear: number;
  status: CardStatus;
  isDefault: boolean;
  nickname?: string;
  
  // Credit card specific
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;
  
  createdAt: string;
  updatedAt: string;
}
```

### 3. Flujo de Operaciones

#### A. Crear Cuenta Bancaria con Tarjetas
1. Usuario crea Account tipo `bank_account`
2. Opcionalmente agrega Cards (credit/debit)
3. Las Cards comparten el balance de la Account

#### B. Crear Billetera Virtual
1. Usuario crea Account tipo `wallet`
2. Requiere DNI
3. Sin Cards asociadas
4. Permite depósitos/extracciones directas

#### C. Operaciones Financieras
- **Depósito/Extracción**: Siempre se hace a nivel Account
- **Tarjetas**: Solo reflejan el balance de la Account padre
- **Límites**: Se aplican a nivel Account, no Card individual

### 4. Cambios en UI

#### AccountFormComponent
- Radio buttons: "Cuenta Bancaria" vs "Billetera Virtual"
- Si es Cuenta Bancaria: permite agregar tarjetas
- Si es Billetera Virtual: solo requiere DNI

#### AccountListComponent  
- Muestra cuentas con sus tarjetas asociadas
- Iconos diferentes para wallet vs bank_account
- Acciones de depósito/extracción por cuenta

### 5. Ventajas del Diseño

1. **Simplicidad**: Una cuenta = un balance
2. **Flexibilidad**: Soporta múltiples tarjetas por cuenta
3. **Compatibilidad**: Mantiene la estructura actual del backend
4. **Escalabilidad**: Fácil agregar nuevos tipos de cuenta/tarjeta

### 6. Plan de Migración

1. Actualizar modelos Go (mantener compatibilidad)
2. Actualizar modelos TypeScript 
3. Modificar UI components
4. Migrar datos existentes
5. Testing y validación

## Próximos Pasos

1. ✅ Completar diseño técnico
2. 🔄 Implementar modelos Go actualizados
3. ⏳ Implementar modelos TypeScript
4. ⏳ Actualizar UI components
5. ⏳ Testing y validación