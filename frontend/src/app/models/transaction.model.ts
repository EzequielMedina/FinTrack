// Transaction Types
export enum TransactionType {
  // Wallet transactions
  WALLET_DEPOSIT = 'wallet_deposit',
  WALLET_WITHDRAWAL = 'wallet_withdrawal',
  WALLET_TRANSFER = 'wallet_transfer',
  
  // Account transactions (for migration compatibility)
  ACCOUNT_DEPOSIT = 'account_deposit',
  ACCOUNT_WITHDRAW = 'account_withdraw',
  ACCOUNT_TRANSFER = 'account_transfer',
  
  // Credit card transactions
  CREDIT_CHARGE = 'credit_charge',
  CREDIT_PAYMENT = 'credit_payment',
  CREDIT_REFUND = 'credit_refund',
  
  // Debit card transactions
  DEBIT_PURCHASE = 'debit_purchase',
  DEBIT_WITHDRAWAL = 'debit_withdrawal',
  DEBIT_REFUND = 'debit_refund',

  // Installment transactions
  INSTALLMENT_PAYMENT = 'installment_payment',
  
  // Salary
  SALARY = 'salary',

  // Legacy types (for backward compatibility)
  DEPOSIT = 'wallet_deposit',
  WITHDRAWAL = 'wallet_withdrawal',
  TRANSFER = 'wallet_transfer',
  PAYMENT = 'credit_payment',
  REFUND = 'credit_refund',
  PURCHASE = 'debit_purchase',
  CASH_ADVANCE = 'debit_withdrawal',
  FEE = 'debit_purchase',
  INTEREST = 'credit_refund',
  DIVIDEND = 'wallet_deposit',
  INVESTMENT = 'wallet_deposit'
}

// Transaction Status
export enum TransactionStatus {
  PENDING = 'pending',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'canceled',
  REVERSED = 'reversed'
}

// Main Transaction interface
export interface Transaction {
  id: string;
  type: TransactionType;
  amount: number;
  currency: string;
  fromAccountId?: string;
  toAccountId?: string;
  description?: string;
  status: TransactionStatus;
  createdAt: string;
  updatedAt: string;
  metadata?: TransactionMetadata;
  balanceAfter?: number;
}

// Transaction metadata for additional info
export interface TransactionMetadata {
  category?: string;
  merchantName?: string;
  location?: string;
  cardLastFour?: string;
  externalTransactionId?: string;
  failureReason?: string;
  accountType?: string;
  accountName?: string;
}

// DTOs for API communication
export interface CreateTransactionDTO {
  type: TransactionType;
  amount: number;
  currency?: string;
  fromAccountId?: string;
  toAccountId?: string;
  userId?: string;
  initiatedBy?: string;
  description?: string;
  paymentMethod?: string;
  metadata?: TransactionMetadata;
}

export interface TransactionFilterDTO {
  userId?: string;
  accountId?: string;
  type?: TransactionType;
  status?: TransactionStatus;
  startDate?: string;
  endDate?: string;
  minAmount?: number;
  maxAmount?: number;
  limit?: number;
  offset?: number;
}

export interface TransactionListResponse {
  transactions: Transaction[];
  total: number;
  page: number;
  pageSize: number;
  hasNext: boolean;
}

// Response interfaces
export interface TransactionResponse {
  success: boolean;
  transaction?: Transaction;
  message?: string;
  errors?: string[];
}

// Transfer specific interface
export interface TransferRequest {
  fromAccountId: string;
  toAccountId: string;
  amount: number;
  description?: string;
}

// Payment specific interface  
export interface PaymentRequest {
  accountId: string;
  amount: number;
  merchantName?: string;
  description?: string;
  cardId?: string;
}

// Deposit/Withdrawal interface
export interface DepositWithdrawalRequest {
  accountId: string;
  amount: number;
  description?: string;
  method?: string; // 'ATM', 'BANK_TRANSFER', 'CASH', etc.
}