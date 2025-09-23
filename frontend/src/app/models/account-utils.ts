import { AccountType, AccountTypeInfo, Account } from './account.model';
import { Card, CardStatus } from './card.model';

// Account type utilities
export class AccountUtils {
  
  static getAccountTypeInfo(accountType: AccountType): AccountTypeInfo {
    switch (accountType) {
      case AccountType.WALLET:
        return {
          canHaveCards: false,
          requiresDNI: true,
          isWallet: true,
          isBankAccount: false,
          displayName: 'Billetera Virtual',
          description: 'Para pagos digitales sin tarjetas físicas'
        };
      
      case AccountType.BANK_ACCOUNT:
        return {
          canHaveCards: true,
          requiresDNI: false,
          isWallet: false,
          isBankAccount: true,
          displayName: 'Cuenta Bancaria',
          description: 'Cuenta tradicional que puede tener múltiples tarjetas'
        };
      
      case AccountType.CHECKING:
        return {
          canHaveCards: true,
          requiresDNI: false,
          isWallet: false,
          isBankAccount: false,
          displayName: 'Cuenta Corriente',
          description: 'Cuenta para uso diario con chequera'
        };
      
      case AccountType.SAVINGS:
        return {
          canHaveCards: true,
          requiresDNI: false,
          isWallet: false,
          isBankAccount: false,
          displayName: 'Cuenta de Ahorros',
          description: 'Cuenta para ahorrar dinero'
        };
      
      case AccountType.CREDIT:
        return {
          canHaveCards: true,
          requiresDNI: false,
          isWallet: false,
          isBankAccount: false,
          displayName: 'Tarjeta de Crédito',
          description: 'Línea de crédito con límite establecido'
        };
      
      case AccountType.DEBIT:
        return {
          canHaveCards: true,
          requiresDNI: false,
          isWallet: false,
          isBankAccount: false,
          displayName: 'Tarjeta de Débito',
          description: 'Tarjeta vinculada directamente al balance'
        };
      
      default:
        return {
          canHaveCards: false,
          requiresDNI: false,
          isWallet: false,
          isBankAccount: false,
          displayName: 'Cuenta',
          description: 'Tipo de cuenta no especificado'
        };
    }
  }
  
  static canAccountHaveCards(accountType: AccountType): boolean {
    return this.getAccountTypeInfo(accountType).canHaveCards;
  }
  
  static requiresDNI(accountType: AccountType): boolean {
    return this.getAccountTypeInfo(accountType).requiresDNI;
  }
  
  static isWalletAccount(accountType: AccountType): boolean {
    return accountType === AccountType.WALLET;
  }
  
  static isBankAccount(accountType: AccountType): boolean {
    return accountType === AccountType.BANK_ACCOUNT;
  }
  
  static getActiveCards(account: Account): Card[] {
    if (!account.cards) return [];
    return account.cards.filter(card => card.status === CardStatus.ACTIVE);
  }
  
  static getDefaultCard(account: Account): Card | null {
    const activeCards = this.getActiveCards(account);
    const defaultCard = activeCards.find(card => card.isDefault);
    return defaultCard || activeCards[0] || null;
  }
  
  static hasActiveCards(account: Account): boolean {
    return this.getActiveCards(account).length > 0;
  }
  
  static getCardsByType(account: Account, cardType: 'credit' | 'debit'): Card[] {
    const activeCards = this.getActiveCards(account);
    return activeCards.filter(card => card.cardType === cardType);
  }
  
  static getTotalCreditLimit(account: Account): number {
    const creditCards = this.getCardsByType(account, 'credit');
    return creditCards.reduce((total, card) => total + (card.creditLimit || 0), 0);
  }
  
  static validateAccountForOperation(account: Account, operation: 'deposit' | 'withdraw' | 'transfer'): {
    canOperate: boolean;
    reason?: string;
  } {
    if (!account.isActive) {
      return { canOperate: false, reason: 'La cuenta está inactiva' };
    }
    
    if (operation === 'withdraw' && account.balance <= 0) {
      return { canOperate: false, reason: 'Saldo insuficiente' };
    }
    
    return { canOperate: true };
  }
  
  static getAccountIcon(accountType: AccountType): string {
    switch (accountType) {
      case AccountType.WALLET:
        return 'account_balance_wallet';
      case AccountType.BANK_ACCOUNT:
        return 'account_balance';
      case AccountType.CHECKING:
        return 'receipt_long';
      case AccountType.SAVINGS:
        return 'savings';
      case AccountType.CREDIT:
        return 'credit_card';
      case AccountType.DEBIT:
        return 'payment';
      default:
        return 'account_circle';
    }
  }
  
  static getAccountColor(accountType: AccountType): string {
    switch (accountType) {
      case AccountType.WALLET:
        return 'accent';
      case AccountType.BANK_ACCOUNT:
        return 'primary';
      case AccountType.CHECKING:
        return 'primary';
      case AccountType.SAVINGS:
        return 'primary';
      case AccountType.CREDIT:
        return 'warn';
      case AccountType.DEBIT:
        return 'primary';
      default:
        return 'primary';
    }
  }
}

// Card utilities
export class CardUtils {
  
  static getCardIcon(cardBrand: string): string {
    switch (cardBrand.toLowerCase()) {
      case 'visa':
        return 'credit_card';
      case 'mastercard':
        return 'credit_card';
      case 'amex':
      case 'american_express':
        return 'credit_card';
      default:
        return 'payment';
    }
  }
  
  static getCardColor(cardType: 'credit' | 'debit'): string {
    return cardType === 'credit' ? 'warn' : 'primary';
  }
  
  static formatCardNumber(cardNumber: string): string {
    // Convert to masked format: **** **** **** 1234
    if (cardNumber.length < 4) return cardNumber;
    const lastFour = cardNumber.slice(-4);
    return `**** **** **** ${lastFour}`;
  }
  
  static getLastFourDigits(cardNumber: string): string {
    return cardNumber.slice(-4);
  }
  
  static isCardExpired(expirationMonth: number, expirationYear: number): boolean {
    const now = new Date();
    const expiry = new Date(expirationYear, expirationMonth - 1);
    return now > expiry;
  }
  
  static detectCardBrand(cardNumber: string): string {
    const number = cardNumber.replace(/\s/g, '');
    
    if (/^4/.test(number)) return 'visa';
    if (/^5[1-5]/.test(number)) return 'mastercard';
    if (/^3[47]/.test(number)) return 'amex';
    if (/^6/.test(number)) return 'discover';
    if (/^30[0-5]/.test(number)) return 'diners';
    
    return 'other';
  }
}

// Form validation utilities
export class ValidationUtils {
  
  static validateAccountForm(formData: any, accountType: AccountType): string[] {
    const errors: string[] = [];
    

    
    if (!formData.currency) {
      errors.push('Debe seleccionar una moneda');
    }
    
    if (formData.initialBalance < 0) {
      errors.push('El balance inicial no puede ser negativo');
    }
    
    if (AccountUtils.requiresDNI(accountType) && !formData.dni) {
      errors.push('El DNI es requerido para billeteras virtuales');
    }
    
    return errors;
  }
  
  static validateCardForm(formData: any): string[] {
    const errors: string[] = [];
    
    if (!formData.cardNumber || formData.cardNumber.length < 13) {
      errors.push('Número de tarjeta inválido');
    }
    
    if (!formData.holderName || formData.holderName.trim().length < 2) {
      errors.push('Nombre del titular es requerido');
    }
    
    if (!formData.expirationMonth || formData.expirationMonth < 1 || formData.expirationMonth > 12) {
      errors.push('Mes de expiración inválido');
    }
    
    if (!formData.expirationYear || formData.expirationYear < new Date().getFullYear()) {
      errors.push('Año de expiración inválido');
    }
    
    if (!formData.cvv || formData.cvv.length < 3) {
      errors.push('CVV inválido');
    }
    
    return errors;
  }
}