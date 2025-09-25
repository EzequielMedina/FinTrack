import { Card } from './card.model';

export enum AccountType {
  SAVINGS = 'savings',
  CHECKING = 'checking',
  CREDIT = 'credit',
  DEBIT = 'debit',
  WALLET = 'wallet',
  BANK_ACCOUNT = 'bank_account'  // New integrated type: can have multiple cards
}

export enum Currency {
  ARS = 'ARS',
  USD = 'USD',
  EUR = 'EUR'
}

export enum AccountStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended'
}

export interface Account {
  id: string;
  userId: string;
  accountType: AccountType;
  name: string;
  description: string;
  currency: Currency;
  balance: number;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;

  // Cards array (optional - only for bank_account type)
  cards?: Card[];

  // Credit card specific fields (legacy - for backward compatibility)
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;

  // Personal identification (for virtual wallets)
  dni?: string;
}

export interface CreateAccountRequest {
  userId: string;
  accountType: AccountType;
  name: string;
  description?: string;
  currency: Currency;
  initialBalance: number;
  isActive?: boolean;

  // Cards array (optional - for bank_account type)
  cards?: Partial<Card>[];

  // Credit card specific fields (legacy - for backward compatibility)
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;

  // Personal identification (for virtual wallets)
  dni?: string;
}

export interface UpdateAccountRequest {
  name: string;
  description?: string;

  // Credit card specific fields
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;

  // Personal identification (for virtual wallets)
  dni?: string;
}

export interface AddFundsRequest {
  amount: number;
  description: string;
  reference?: string;
}

export interface WithdrawFundsRequest {
  amount: number;
  description: string;
  reference?: string;
}

export interface UpdateCreditLimitRequest {
  creditLimit: number;
}

export interface UpdateCreditDatesRequest {
  closingDate?: string;
  dueDate?: string;
}

export interface BalanceResponse {
  accountId: string;
  balance: number;
}

export interface AvailableCreditResponse {
  accountId: string;
  creditLimit: number;
  usedCredit: number;
  availableCredit: number;
}

export interface AccountsListResponse {
  accounts: Account[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

// Validation interfaces
export interface AccountValidationError {
  field: string;
  message: string;
}

export interface AccountValidationResult {
  isValid: boolean;
  errors: AccountValidationError[];
}

// Form interfaces
export interface AccountFormData {
  userId: string;
  accountType: AccountType;
  name: string;
  description: string;
  currency: Currency;
  initialBalance: number;

  // Cards data (optional - for bank_account type)
  cards?: Partial<Card>[];

  // Credit card specific fields (legacy)
  creditLimit?: number;
  closingDate?: string;
  dueDate?: string;

  // Personal identification (for virtual wallets)
  dni?: string;
}

export interface AccountFormErrors {
  userId?: string;
  accountType?: string;
  name?: string;
  description?: string;
  currency?: string;
  initialBalance?: string;
  creditLimit?: string;
  closingDate?: string;
  dueDate?: string;
  dni?: string;
}

export interface WalletOperationFormData {
  amount: number;
  description: string;
  reference: string;
}

export interface WalletOperationFormErrors {
  amount?: string;
  description?: string;
  reference?: string;
}

export interface CreditLimitFormData {
  creditLimit: number;
}

export interface CreditLimitFormErrors {
  creditLimit?: string;
}

export interface CreditDatesFormData {
  closingDate: string;
  dueDate: string;
}

export interface CreditDatesFormErrors {
  closingDate?: string;
  dueDate?: string;
}

// Helper interfaces for business logic
export interface AccountSummary {
  id: string;
  name: string;
  accountType: AccountType;
  balance: number;
  currency: Currency;
  isActive: boolean;
  availableCredit?: number; // Only for credit accounts
  cardsCount?: number;      // Number of active cards
  hasCards?: boolean;       // Whether account supports/has cards
}

// Account type utilities
export interface AccountTypeInfo {
  canHaveCards: boolean;
  requiresDNI: boolean;
  isWallet: boolean;
  isBankAccount: boolean;
  displayName: string;
  description: string;
}

// Helper for account creation workflow
export interface AccountCreationStep {
  step: 'type' | 'details' | 'cards' | 'confirmation';
  title: string;
  description: string;
  isRequired: boolean;
  canSkip: boolean;
}

export interface WalletTransactionHistory {
  id: string;
  accountId: string;
  type: 'deposit' | 'withdrawal';
  amount: number;
  description: string;
  reference?: string;
  timestamp: string;
  balance: number;
}

export interface CreditCardStatement {
  accountId: string;
  period: string;
  creditLimit: number;
  usedCredit: number;
  availableCredit: number;
  minimumPayment: number;
  dueDate: string;
  closingDate: string;
}