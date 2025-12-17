import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError, map, switchMap } from 'rxjs/operators';

import {
  Transaction,
  CreateTransactionDTO,
  TransactionFilterDTO,
  TransactionListResponse,
  TransactionResponse,
  TransferRequest,
  PaymentRequest,
  DepositWithdrawalRequest,
  TransactionType,
  TransactionStatus
} from '../models/transaction.model';

import { AuthService } from './auth.service';
import { AccountService } from './account.service';
import { AccountType } from '../models/account.model';

@Injectable({
  providedIn: 'root'
})
export class TransactionService {
  private apiUrl = '/api/v1/transactions';

  constructor(
    private http: HttpClient, 
    private authService: AuthService,
    private accountService: AccountService
  ) {}

  // Método helper para obtener el ID del usuario actual
  private getUserId(): string {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser?.id) {
      throw new Error('Usuario no autenticado');
    }
    return currentUser.id;
  }

  // Método helper para determinar el tipo de transacción basado en el tipo de cuenta
  private getTransactionTypeForDeposit(accountType: AccountType): TransactionType {
    switch (accountType) {
      case AccountType.WALLET:
        return TransactionType.WALLET_DEPOSIT;
      case AccountType.SAVINGS:
      case AccountType.CHECKING:
      case AccountType.BANK_ACCOUNT:
        return TransactionType.ACCOUNT_DEPOSIT;
      case AccountType.CREDIT:
        return TransactionType.CREDIT_PAYMENT; // Para crédito, un depósito es un pago
      default:
        return TransactionType.WALLET_DEPOSIT; // Fallback
    }
  }

  private getTransactionTypeForWithdrawal(accountType: AccountType): TransactionType {
    switch (accountType) {
      case AccountType.WALLET:
        return TransactionType.WALLET_WITHDRAWAL;
      case AccountType.SAVINGS:
      case AccountType.CHECKING:
      case AccountType.BANK_ACCOUNT:
        return TransactionType.ACCOUNT_WITHDRAW;
      case AccountType.DEBIT:
        return TransactionType.DEBIT_WITHDRAWAL;
      default:
        return TransactionType.WALLET_WITHDRAWAL; // Fallback
    }
  }

  // ================================
  // CREAR TRANSACCIONES
  // ================================

  /**
   * Crear una transacción genérica
   */
  createTransaction(transaction: CreateTransactionDTO): Observable<TransactionResponse> {
    return this.http.post<TransactionResponse>(`${this.apiUrl}`, transaction)
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * Realizar transferencia entre cuentas
   */
  transfer(transferData: TransferRequest): Observable<TransactionResponse> {
    const transaction: CreateTransactionDTO = {
      type: TransactionType.WALLET_TRANSFER,
      amount: transferData.amount,
      fromAccountId: transferData.fromAccountId,
      toAccountId: transferData.toAccountId,
      description: transferData.description
    };

    return this.createTransaction(transaction);
  }

  /**
   * Realizar depósito - Versión dinámica que obtiene el tipo de cuenta
   */
  deposit(depositData: DepositWithdrawalRequest): Observable<TransactionResponse> {
    const userId = this.getUserId();
    
    // Primero obtenemos la información de la cuenta para determinar el tipo correcto de transacción
    return this.accountService.getAccountById(depositData.accountId).pipe(
      switchMap(account => {
        if (!account) {
          throw new Error(`Cuenta ${depositData.accountId} no encontrada`);
        }

        if (!account.isActive) {
          throw new Error(`La cuenta ${account.name} está inactiva`);
        }

        const transactionType = this.getTransactionTypeForDeposit(account.accountType);
        
        const transaction: CreateTransactionDTO = {
          type: transactionType,
          amount: depositData.amount,
          currency: account.currency || 'ARS', // ✅ Usar la moneda de la cuenta
          toAccountId: depositData.accountId,
          userId: userId,
          initiatedBy: userId,
          description: depositData.description || `Depósito en ${account.name}`,
          paymentMethod: depositData.method === 'BANK_TRANSFER' ? 'bank_transfer' : 'wallet',
          metadata: {
            category: depositData.method || 'BANK_TRANSFER',
            accountType: account.accountType,
            accountName: account.name
          }
        };

        return this.createTransaction(transaction);
      }),
      catchError(this.handleError)
    );
  }

  /**
   * Realizar retiro - Versión dinámica que obtiene el tipo de cuenta
   */
  withdraw(withdrawalData: DepositWithdrawalRequest): Observable<TransactionResponse> {
    const userId = this.getUserId();
    
    // Primero obtenemos la información de la cuenta para determinar el tipo correcto de transacción
    return this.accountService.getAccountById(withdrawalData.accountId).pipe(
      switchMap(account => {
        if (!account) {
          throw new Error(`Cuenta ${withdrawalData.accountId} no encontrada`);
        }

        if (!account.isActive) {
          throw new Error(`La cuenta ${account.name} está inactiva`);
        }

        // Validar que hay fondos suficientes
        if (account.balance < withdrawalData.amount) {
          throw new Error(`Fondos insuficientes. Saldo disponible: $${account.balance.toFixed(2)}`);
        }

        const transactionType = this.getTransactionTypeForWithdrawal(account.accountType);
        
        const transaction: CreateTransactionDTO = {
          type: transactionType,
          amount: withdrawalData.amount,
          currency: account.currency || 'ARS', // ✅ Usar la moneda de la cuenta
          fromAccountId: withdrawalData.accountId,
          userId: userId,
          initiatedBy: userId,
          description: withdrawalData.description || `Retiro de ${account.name}`,
          paymentMethod: withdrawalData.method === 'ATM' ? 'wallet' : 'cash',
          metadata: {
            category: withdrawalData.method || 'ATM',
            accountType: account.accountType,
            accountName: account.name
          }
        };

        return this.createTransaction(transaction);
      }),
      catchError(this.handleError)
    );
  }

  /**
   * Realizar pago
   */
  makePayment(paymentData: PaymentRequest): Observable<TransactionResponse> {
    const transaction: CreateTransactionDTO = {
      type: TransactionType.CREDIT_PAYMENT,
      amount: paymentData.amount,
      fromAccountId: paymentData.accountId,
      description: paymentData.description,
      metadata: {
        merchantName: paymentData.merchantName,
        cardLastFour: paymentData.cardId ? paymentData.cardId.slice(-4) : undefined
      }
    };

    return this.createTransaction(transaction);
  }

  // ================================
  // CONSULTAR TRANSACCIONES
  // ================================

  /**
   * Obtener transacción por ID
   */
  getTransaction(id: string): Observable<Transaction> {
    return this.http.get<Transaction>(`${this.apiUrl}/${id}`)
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * Obtener transacciones de un usuario
   */
  getUserTransactions(userId: string, filters?: TransactionFilterDTO): Observable<TransactionListResponse> {
    let params = new HttpParams();
    
    if (filters) {
      Object.keys(filters).forEach(key => {
        const value = (filters as any)[key];
        if (value !== undefined && value !== null && value !== '') {
          params = params.set(key, value.toString());
        }
      });
    }

    // El backend espera el userId en el header X-User-ID, no en la URL
    const headers = { 'X-User-ID': userId };

    return this.http.get<TransactionListResponse>(`${this.apiUrl}`, { params, headers })
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * Obtener transacciones de una cuenta específica
   */
  getAccountTransactions(accountId: string, filters?: TransactionFilterDTO): Observable<TransactionListResponse> {
    const accountFilters = { ...filters, accountId };
    
    let params = new HttpParams();
    Object.keys(accountFilters).forEach(key => {
      const value = (accountFilters as any)[key];
      if (value !== undefined && value !== null && value !== '') {
        params = params.set(key, value.toString());
      }
    });

    return this.http.get<TransactionListResponse>(`${this.apiUrl}`, { params })
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * Obtener historial reciente de transacciones
   */
  getRecentTransactions(userId: string, limit: number = 10): Observable<Transaction[]> {
    const filters: TransactionFilterDTO = {
      limit,
      offset: 0
    };

    return this.getUserTransactions(userId, filters).pipe(
      map(response => response.transactions)
    );
  }

  // ================================
  // ACTUALIZAR TRANSACCIONES
  // ================================

  /**
   * Actualizar transacción (principalmente para cancelar)
   */
  updateTransaction(id: string, updates: Partial<Transaction>): Observable<TransactionResponse> {
    return this.http.put<TransactionResponse>(`${this.apiUrl}/${id}`, updates)
      .pipe(
        catchError(this.handleError)
      );
  }

  /**
   * Cancelar transacción
   */
  cancelTransaction(id: string, reason?: string): Observable<TransactionResponse> {
    const updates = {
      status: TransactionStatus.CANCELLED,
      metadata: {
        failureReason: reason || 'Cancelled by user'
      }
    };

    return this.updateTransaction(id, updates);
  }

  // ================================
  // UTILIDADES Y FILTROS
  // ================================

  /**
   * Obtener transacciones por tipo
   */
  getTransactionsByType(userId: string, type: TransactionType, limit?: number): Observable<Transaction[]> {
    const filters: TransactionFilterDTO = {
      type,
      limit,
      offset: 0
    };

    return this.getUserTransactions(userId, filters).pipe(
      map(response => response.transactions)
    );
  }

  /**
   * Obtener transacciones por rango de fechas
   */
  getTransactionsByDateRange(
    userId: string, 
    startDate: Date, 
    endDate: Date, 
    limit?: number
  ): Observable<Transaction[]> {
    const filters: TransactionFilterDTO = {
      startDate: startDate.toISOString(),
      endDate: endDate.toISOString(),
      limit,
      offset: 0
    };

    return this.getUserTransactions(userId, filters).pipe(
      map(response => response.transactions)
    );
  }

  /**
   * Obtener estadísticas de transacciones
   */
  getTransactionStats(userId: string, startDate?: Date, endDate?: Date): Observable<any> {
    // Esta funcionalidad se puede implementar en el backend o calcular en el frontend
    const filters: TransactionFilterDTO = {
      startDate: startDate?.toISOString(),
      endDate: endDate?.toISOString()
    };

    return this.getUserTransactions(userId, filters).pipe(
      map(response => {
        const transactions = response.transactions;
        return this.calculateStats(transactions);
      })
    );
  }

  // ================================
  // MÉTODOS AUXILIARES
  // ================================

  private calculateStats(transactions: Transaction[]) {
    const stats = {
      total: transactions.length,
      totalIncome: 0,
      totalExpenses: 0,
      byType: {} as any,
      byStatus: {} as any
    };

    transactions.forEach(transaction => {
      // Calcular ingresos y gastos
      if ([TransactionType.WALLET_DEPOSIT, TransactionType.SALARY, TransactionType.CREDIT_REFUND, TransactionType.DEBIT_REFUND].includes(transaction.type)) {
        stats.totalIncome += transaction.amount;
      } else {
        stats.totalExpenses += transaction.amount;
      }

      // Agrupar por tipo
      stats.byType[transaction.type] = (stats.byType[transaction.type] || 0) + 1;

      // Agrupar por estado
      stats.byStatus[transaction.status] = (stats.byStatus[transaction.status] || 0) + 1;
    });

    return stats;
  }

  private handleError(error: any) {
    console.error('Transaction Service Error:', error);
    let errorMessage = 'Error en el servicio de transacciones';

    if (error.error && error.error.message) {
      errorMessage = error.error.message;
    } else if (error.message) {
      errorMessage = error.message;
    }

    return throwError(() => new Error(errorMessage));
  }

  // ================================
  // MÉTODOS DE VALIDACIÓN
  // ================================

  /**
   * Validar si una transacción se puede cancelar
   */
  canCancelTransaction(transaction: Transaction): boolean {
    return transaction.status === TransactionStatus.PENDING;
  }

  /**
   * Validar si una transacción se puede revertir
   */
  canRefundTransaction(transaction: Transaction): boolean {
    return transaction.status === TransactionStatus.COMPLETED && 
           [TransactionType.CREDIT_PAYMENT, TransactionType.DEBIT_PURCHASE].includes(transaction.type);
  }

  /**
   * Obtener color para el estado de transacción (para UI)
   */
  getStatusColor(status: TransactionStatus): string {
    switch (status) {
      case TransactionStatus.COMPLETED: return 'success';
      case TransactionStatus.PENDING: return 'warning';
      case TransactionStatus.FAILED: return 'danger';
      case TransactionStatus.CANCELLED: return 'secondary';
      default: return 'primary';
    }
  }

  /**
   * Obtener icono para el tipo de transacción (para UI)
   */
  getTransactionIcon(type: TransactionType): string {
    switch (type) {
      case TransactionType.WALLET_DEPOSIT:
      case TransactionType.ACCOUNT_DEPOSIT:
        return 'arrow_downward';
      case TransactionType.WALLET_WITHDRAWAL:
      case TransactionType.DEBIT_WITHDRAWAL:
        return 'arrow_upward';
      case TransactionType.WALLET_TRANSFER:
      case TransactionType.ACCOUNT_TRANSFER:
        return 'swap_horiz';
      case TransactionType.CREDIT_PAYMENT:
        return 'credit_card';
      case TransactionType.DEBIT_PURCHASE:
        return 'shopping_bag';
      case TransactionType.CREDIT_REFUND:
      case TransactionType.DEBIT_REFUND:
        return 'keyboard_return';
      case TransactionType.SALARY:
        return 'work';
      default:
        return 'swap_horiz';
    }
  }
}