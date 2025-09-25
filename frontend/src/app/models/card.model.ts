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
  
  // Credit card specific fields
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;
  
  createdAt: string;
  updatedAt: string;
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
  dueDate?: string;
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