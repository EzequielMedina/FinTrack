import { Injectable } from '@angular/core';
import {
  AccountType,
  Currency,
  Account,
  CreateAccountRequest,
  UpdateAccountRequest,
  AccountValidationResult,
  AccountValidationError
} from '../models';

@Injectable({ providedIn: 'root' })
export class AccountValidationService {
  
  // Business rules constants
  private readonly MIN_BALANCE = 0;
  private readonly MAX_BALANCE = 999999999;
  private readonly MAX_NAME_LENGTH = 50;
  private readonly MIN_DESCRIPTION_LENGTH = 0;
  private readonly MAX_DESCRIPTION_LENGTH = 200;
  private readonly DNI_REGEX = /^\d{7,8}$/; // Argentine DNI format
  private readonly ALLOWED_CURRENCIES = Object.values(Currency);
  private readonly ALLOWED_ACCOUNT_TYPES = Object.values(AccountType);

  validateCreateAccountRequest(request: CreateAccountRequest): AccountValidationResult {
    const errors: AccountValidationError[] = [];

    // Validate required fields
    const userIdError = this.validateUserId(request.userId);
    if (userIdError) errors.push(userIdError);

    const accountTypeError = this.validateAccountType(request.accountType);
    if (accountTypeError) errors.push(accountTypeError);

    const nameError = this.validateAccountName(request.name);
    if (nameError) errors.push(nameError);

    const currencyError = this.validateCurrency(request.currency);
    if (currencyError) errors.push(currencyError);

    const dniError = request.dni ? this.validateDni(request.dni) : null;
    if (dniError) errors.push(dniError);

    // Validate optional fields
    if (request.description !== undefined) {
      const descriptionError = this.validateDescription(request.description);
      if (descriptionError) errors.push(descriptionError);
    }

    if (request.initialBalance !== undefined) {
      const balanceError = this.validateBalance(request.initialBalance, 'initialBalance');
      if (balanceError) errors.push(balanceError);
    }

    // Business rule validations
    const businessRuleErrors = this.validateCreateAccountBusinessRules(request);
    errors.push(...businessRuleErrors);

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  validateUpdateAccountRequest(request: UpdateAccountRequest, currentAccount?: Account): AccountValidationResult {
    const errors: AccountValidationError[] = [];

    // Validate fields if they are being updated
    if (request.name !== undefined) {
      const nameError = this.validateAccountName(request.name);
      if (nameError) errors.push(nameError);
    }

    if (request.description !== undefined) {
      const descriptionError = this.validateDescription(request.description);
      if (descriptionError) errors.push(descriptionError);
    }

    // Business rule validations for updates
    if (currentAccount) {
      const businessRuleErrors = this.validateUpdateAccountBusinessRules(request, currentAccount);
      errors.push(...businessRuleErrors);
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  validateAccount(account: Partial<Account>): AccountValidationResult {
    const errors: AccountValidationError[] = [];

    // Validate all provided fields
    if (account.id !== undefined) {
      const idError = this.validateAccountId(account.id);
      if (idError) errors.push(idError);
    }

    if (account.userId !== undefined) {
      const userIdError = this.validateUserId(account.userId);
      if (userIdError) errors.push(userIdError);
    }

    if (account.accountType !== undefined) {
      const accountTypeError = this.validateAccountType(account.accountType);
      if (accountTypeError) errors.push(accountTypeError);
    }

    if (account.name !== undefined) {
      const nameError = this.validateAccountName(account.name);
      if (nameError) errors.push(nameError);
    }

    if (account.currency !== undefined) {
      const currencyError = this.validateCurrency(account.currency);
      if (currencyError) errors.push(currencyError);
    }

    if (account.balance !== undefined) {
      const balanceError = this.validateBalance(account.balance, 'balance');
      if (balanceError) errors.push(balanceError);
    }

    if (account.dni !== undefined) {
      const dniError = this.validateDni(account.dni);
      if (dniError) errors.push(dniError);
    }

    if (account.description !== undefined) {
      const descriptionError = this.validateDescription(account.description);
      if (descriptionError) errors.push(descriptionError);
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  // Specific field validation methods
  validateAccountId(accountId: string): AccountValidationError | null {
    if (!accountId || accountId.trim().length === 0) {
      return {
        field: 'id',
        message: 'El ID de la cuenta es requerido'
      };
    }

    // UUID format validation (basic)
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    if (!uuidRegex.test(accountId)) {
      return {
        field: 'id',
        message: 'El ID de la cuenta debe ser un UUID válido'
      };
    }

    return null;
  }

  validateUserId(userId: string): AccountValidationError | null {
    if (!userId || userId.trim().length === 0) {
      return {
        field: 'userId',
        message: 'El ID del usuario es requerido'
      };
    }

    // UUID format validation (basic)
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    if (!uuidRegex.test(userId)) {
      return {
        field: 'userId',
        message: 'El ID del usuario debe ser un UUID válido'
      };
    }

    return null;
  }

  validateAccountType(accountType: AccountType): AccountValidationError | null {
    if (!accountType) {
      return {
        field: 'accountType',
        message: 'El tipo de cuenta es requerido'
      };
    }

    if (!this.ALLOWED_ACCOUNT_TYPES.includes(accountType)) {
      return {
        field: 'accountType',
        message: `Tipo de cuenta inválido. Tipos permitidos: ${this.ALLOWED_ACCOUNT_TYPES.join(', ')}`
      };
    }

    return null;
  }

  validateAccountName(name: string): AccountValidationError | null {
    if (!name || name.trim().length === 0) {
      return {
        field: 'name',
        message: 'El nombre de la cuenta es requerido'
      };
    }

    const trimmedName = name.trim();

    if (trimmedName.length > this.MAX_NAME_LENGTH) {
      return {
        field: 'name',
        message: `El nombre de la cuenta no puede exceder ${this.MAX_NAME_LENGTH} caracteres`
      };
    }

    // Validate that name contains only valid characters (letters, numbers, spaces, hyphens, underscores)
    const nameRegex = /^[a-zA-ZáéíóúÁÉÍÓÚñÑ0-9\s\-_]+$/;
    if (!nameRegex.test(trimmedName)) {
      return {
        field: 'name',
        message: 'El nombre de la cuenta solo puede contener letras, números, espacios, guiones y guiones bajos'
      };
    }

    return null;
  }

  validateDescription(description: string): AccountValidationError | null {
    if (description && description.length > this.MAX_DESCRIPTION_LENGTH) {
      return {
        field: 'description',
        message: `La descripción no puede exceder ${this.MAX_DESCRIPTION_LENGTH} caracteres`
      };
    }

    return null;
  }

  validateCurrency(currency: Currency): AccountValidationError | null {
    if (!currency) {
      return {
        field: 'currency',
        message: 'La moneda es requerida'
      };
    }

    if (!this.ALLOWED_CURRENCIES.includes(currency)) {
      return {
        field: 'currency',
        message: `Moneda inválida. Monedas permitidas: ${this.ALLOWED_CURRENCIES.join(', ')}`
      };
    }

    return null;
  }

  validateBalance(balance: number, fieldName: string = 'balance'): AccountValidationError | null {
    if (balance === null || balance === undefined) {
      return {
        field: fieldName,
        message: 'El saldo es requerido'
      };
    }

    if (typeof balance !== 'number' || isNaN(balance)) {
      return {
        field: fieldName,
        message: 'El saldo debe ser un número válido'
      };
    }

    if (balance < this.MIN_BALANCE) {
      return {
        field: fieldName,
        message: `El saldo mínimo es $${this.MIN_BALANCE.toLocaleString()}`
      };
    }

    if (balance > this.MAX_BALANCE) {
      return {
        field: fieldName,
        message: `El saldo máximo es $${this.MAX_BALANCE.toLocaleString()}`
      };
    }

    // Validate decimal places (max 2 for currency)
    const decimalPlaces = (balance.toString().split('.')[1] || '').length;
    if (decimalPlaces > 2) {
      return {
        field: fieldName,
        message: 'El saldo no puede tener más de 2 decimales'
      };
    }

    return null;
  }

  validateDni(dni: string): AccountValidationError | null {
    if (!dni || dni.trim().length === 0) {
      return {
        field: 'dni',
        message: 'El DNI es requerido'
      };
    }

    const trimmedDni = dni.trim();

    if (!this.DNI_REGEX.test(trimmedDni)) {
      return {
        field: 'dni',
        message: 'El DNI debe tener entre 7 y 8 dígitos'
      };
    }

    return null;
  }

  // Business rule validation methods
  private validateCreateAccountBusinessRules(request: CreateAccountRequest): AccountValidationError[] {
    const errors: AccountValidationError[] = [];

    // Business Rule: Credit cards require specific validations
    if (request.accountType === AccountType.CREDIT) {
      // Credit cards shouldn't have an initial balance (they start at 0)
      if (request.initialBalance && request.initialBalance !== 0) {
        errors.push({
          field: 'initialBalance',
          message: 'Las tarjetas de crédito deben comenzar con saldo cero'
        });
      }
    }

    // Business Rule: Savings accounts have minimum balance requirements
    if (request.accountType === AccountType.SAVINGS) {
      const minimumSavingsBalance = 1000; // Example minimum for savings accounts
      if (request.initialBalance && request.initialBalance < minimumSavingsBalance) {
        errors.push({
          field: 'initialBalance',
          message: `Las cuentas de ahorro requieren un saldo mínimo de $${minimumSavingsBalance.toLocaleString()}`
        });
      }
    }

    // Business Rule: Validate USD currency consistency
    // Note: USD accounts are differentiated by currency rather than account type

    return errors;
  }

  private validateUpdateAccountBusinessRules(request: UpdateAccountRequest, currentAccount: Account): AccountValidationError[] {
    const errors: AccountValidationError[] = [];

    // Business Rule: Cannot change certain fields for active accounts
    if (currentAccount.isActive) {
      // Add specific business rules for active accounts
      // For example, certain fields might be immutable once the account is active
    }

    // Business Rule: Name uniqueness within user's accounts (would need additional context)
    // This would typically require checking against existing accounts

    return errors;
  }

  // Utility methods for specific validations
  isValidAccountForOperation(account: Account, operation: string): boolean {
    // Check if account is active
    if (!account.isActive) {
      return false;
    }

    // Check operation-specific rules
    switch (operation) {
      case 'withdraw':
        return account.accountType !== AccountType.CREDIT;
      case 'credit_operation':
        return account.accountType === AccountType.CREDIT;
      case 'transfer':
        return true; // Most accounts can receive transfers
      default:
        return true;
    }
  }

  getAccountTypeDisplayName(accountType: AccountType): string {
    const displayNames = {
      [AccountType.SAVINGS]: 'Cuenta de Ahorros',
      [AccountType.CHECKING]: 'Cuenta Corriente',
      [AccountType.CREDIT]: 'Tarjeta de Crédito',
      [AccountType.WALLET]: 'Billetera Virtual'
    };

    return displayNames[accountType as keyof typeof displayNames] || accountType;
  }

  getCurrencyDisplayName(currency: Currency): string {
    const displayNames = {
      [Currency.ARS]: 'Pesos Argentinos (ARS)',
      [Currency.USD]: 'Dólares Estadounidenses (USD)',
      [Currency.EUR]: 'Euros (EUR)'
    };

    return displayNames[currency as keyof typeof displayNames] || currency;
  }

  formatCurrency(amount: number, currency: Currency): string {
    const formatter = new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    });

    return formatter.format(amount);
  }

  // Additional validation methods for the form component
  validateCreateAccount(formData: any): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];



    if (formData.accountName && formData.accountName.length > 100) {
      errors.push('El nombre de la cuenta no puede exceder 100 caracteres');
    }

    // Validate account type
    if (!formData.accountType) {
      errors.push('El tipo de cuenta es requerido');
    }

    // Validate currency
    if (!formData.currency) {
      errors.push('La moneda es requerida');
    }

    // Validate initial balance
    if (formData.initialBalance < 0) {
      errors.push('El saldo inicial no puede ser negativo');
    }

    if (formData.initialBalance > 999999999) {
      errors.push('El saldo inicial no puede exceder $999,999,999');
    }

    // Credit card specific validations
    if (formData.accountType === AccountType.CREDIT) {
      if (!formData.creditLimit || formData.creditLimit < 100) {
        errors.push('El límite de crédito debe ser al menos $100');
      }

      if (formData.creditLimit > 999999999) {
        errors.push('El límite de crédito no puede exceder $999,999,999');
      }

      if (!formData.closingDate) {
        errors.push('La fecha de cierre es requerida para tarjetas de crédito');
      }

      if (!formData.dueDate) {
        errors.push('La fecha de vencimiento es requerida para tarjetas de crédito');
      }

      // Validate that due date is after closing date
      if (formData.closingDate && formData.dueDate) {
        const closingDate = new Date(formData.closingDate);
        const dueDate = new Date(formData.dueDate);
        
        if (dueDate <= closingDate) {
          errors.push('La fecha de vencimiento debe ser posterior a la fecha de cierre');
        }
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  validateUpdateAccount(formData: any): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];



    if (formData.accountName && formData.accountName.length > 100) {
      errors.push('El nombre de la cuenta no puede exceder 100 caracteres');
    }

    // Credit card specific validations (only if it's a credit card account)
    if (formData.accountType === AccountType.CREDIT) {
      if (formData.creditLimit !== undefined && formData.creditLimit < 0) {
        errors.push('El límite de crédito no puede ser negativo');
      }

      if (formData.creditLimit !== undefined && formData.creditLimit > 999999999) {
        errors.push('El límite de crédito no puede exceder $999,999,999');
      }

      // Validate dates if provided
      if (formData.closingDate && formData.dueDate) {
        const closingDate = new Date(formData.closingDate);
        const dueDate = new Date(formData.dueDate);
        
        if (dueDate <= closingDate) {
          errors.push('La fecha de vencimiento debe ser posterior a la fecha de cierre');
        }
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }
}