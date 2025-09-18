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
  accountId?: string;          // Opcional por compatibilidad
  userId: string;
  cardType: CardType;
  cardBrand: CardBrand;
  lastFourDigits: string;      // Solo mostramos los últimos 4 dígitos
  maskedNumber: string;        // Número enmascarado: **** **** **** 1234
  holderName: string;
  expirationMonth: number;     // 1-12
  expirationYear: number;      // YYYY
  status: CardStatus;
  isDefault: boolean;          // Tarjeta predeterminada
  nickname?: string;           // Nombre personalizado para la tarjeta
  currency: string;            // Moneda (ARS, USD)
  balance: number;             // Balance disponible
  createdAt: string;
  updatedAt: string;
}

export interface CreateCardRequest {
  userId: string;              // ID del usuario propietario
  cardType: CardType;
  cardNumber: string;          // Número completo encriptado en el frontend
  holderName: string;
  expirationMonth: number;
  expirationYear: number;
  cvv: string;                 // CVV para validación inicial
  nickname?: string;
  currency?: string;           // Moneda de la cuenta (ARS, USD)
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
  userId: string;              // Cambiado de accountId a userId
  cardType: CardType;
  cardNumber: string;
  holderName: string;
  expirationMonth: number;
  expirationYear: number;
  cvv: string;
  nickname: string;
  currency: string;            // Añadido campo currency
}

export interface CardFormErrors {
  cardNumber?: string;
  holderName?: string;
  expirationMonth?: string;
  expirationYear?: string;
  cvv?: string;
  userId?: string;             // Cambiado de accountId a userId
  currency?: string;           // Añadido currency
}