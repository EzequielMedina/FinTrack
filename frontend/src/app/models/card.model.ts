export enum CardType {
  CREDIT = 'credit',
  DEBIT = 'debit'
}

export enum CardBrand {
  VISA = 'visa',
  MASTERCARD = 'mastercard',
  AMERICAN_EXPRESS = 'amex',
  DISCOVER = 'discover',
  DINERS = 'diners',
  OTHER = 'other'
}

export enum CardStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  BLOCKED = 'blocked',
  EXPIRED = 'expired'
}

export interface Card {
  id: string;
  accountId: string;             // REQUIRED: Always linked to an account
  cardType: CardType;
  cardBrand: CardBrand;
  lastFourDigits: string;        // Only show last 4 digits
  maskedNumber: string;          // Masked number: **** **** **** 1234
  holderName: string;
  expirationMonth: number;       // 1-12
  expirationYear: number;        // YYYY
  status: CardStatus;
  isDefault: boolean;            // Default card for the account
  nickname?: string;             // Custom name for the card
  
  // Balance field - usage depends on card type:
  // - Credit cards: debt amount (positive = owed to bank)
  // - Debit cards: always 0 (uses account balance)
  balance: number;
  
  // Credit card specific fields
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;
  
  // Installment plans summary (only for credit cards)
  installmentPlans?: InstallmentPlansSummary;
  
  createdAt: string;
  updatedAt: string;
}

// Helper functions for Card
export class CardHelpers {
  /**
   * Verifica si una tarjeta está vencida
   */
  static isExpired(card: Card): boolean {
    const now = new Date();
    const currentMonth = now.getMonth() + 1; // getMonth() returns 0-11
    const currentYear = now.getFullYear();
    
    // La tarjeta vence al final del mes indicado
    if (card.expirationYear < currentYear) {
      return true;
    }
    
    if (card.expirationYear === currentYear && card.expirationMonth < currentMonth) {
      return true;
    }
    
    return false;
  }

  /**
   * Verifica si una tarjeta está activa y NO vencida
   */
  static isFullyActive(card: Card): boolean {
    return card.status === CardStatus.ACTIVE && !this.isExpired(card);
  }

  /**
   * Obtiene un mensaje descriptivo del estado de la tarjeta
   */
  static getStatusMessage(card: Card): string {
    if (this.isExpired(card)) {
      return `Tarjeta vencida (${card.expirationMonth}/${card.expirationYear})`;
    }
    if (card.status === CardStatus.BLOCKED) {
      return 'Tarjeta bloqueada';
    }
    if (card.status === CardStatus.INACTIVE) {
      return 'Tarjeta inactiva';
    }
    return 'Tarjeta activa';
  }

  /**
   * Formatea la fecha de expiración
   */
  static formatExpiration(card: Card): string {
    const month = card.expirationMonth.toString().padStart(2, '0');
    return `${month}/${card.expirationYear}`;
  }
}

export interface InstallmentPlansSummary {
  activePlans: number;
  totalDebt: number;
  monthlyPayment: number;
  nextPaymentDue?: string;
  overdueCount: number;
}

export interface CreateCardRequest {
  userId: string;              // REQUIRED: ID of the user who owns the card
  accountId: string;           // REQUIRED: ID of the account to link the card to
  cardType: CardType;
  cardNumber: string;          // Full number (encrypted on frontend)
  holderName: string;
  expirationMonth: number;
  expirationYear: number;
  cvv: string;                 // CVV for initial validation
  nickname?: string;
  
  // Credit card specific fields
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;
}

export interface UpdateCardRequest {
  holderName?: string;
  expirationMonth?: number;
  expirationYear?: number;
  nickname?: string;
  isDefault?: boolean;
  /** Límite de crédito (tarjetas de crédito); enviar para no perderlo al editar */
  creditLimit?: number;
  /** Fecha de vencimiento del pago (YYYY-MM-DD); enviar para actualizarla */
  dueDate?: string;
}

export interface CardValidationError {
  field: string;
  message: string;
}

export interface CardValidationResult {
  isValid: boolean;
  errors: CardValidationError[];
  brand?: CardBrand;
}

export interface EncryptedCardData {
  encryptedNumber: string;
  encryptedCvv: string;
  keyFingerprint: string;      // Para identificar la clave de encriptación usada
}

export interface CardsListResponse {
  cards: Card[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Tipos para formularios reactivos
export interface CardFormData {
  accountId: string;           // ID de la cuenta seleccionada en el formulario
  cardType: CardType;
  cardNumber: string;
  holderName: string;
  expirationMonth: number;
  expirationYear: number;
  cvv: string;
  nickname: string;
  
  // Credit card specific fields
  creditLimit?: number;
  closingDate?: string;
  dueDate?: Date;              // Fecha completa seleccionada por el usuario
}

export interface CardFormErrors {
  cardNumber?: string;
  holderName?: string;
  expirationMonth?: string;
  expirationYear?: string;
  cvv?: string;
  accountId?: string;
  creditLimit?: string;
  closingDate?: string;
  dueDate?: string;
}

// Installment-related responses
export interface CardWithInstallmentsResponse {
  card: Card;
  installmentPlans: any[]; // Will be typed as InstallmentPlan[] when imported
  installmentsSummary: InstallmentPlansSummary;
}

export interface CardTransactionResponse {
  transactionId: string;
  card: Card;
  amount: number;
  type: 'charge' | 'payment' | 'installment_charge';
  installmentPlan?: any; // Will be typed as InstallmentPlan when imported
  success: boolean;
  message?: string;
}