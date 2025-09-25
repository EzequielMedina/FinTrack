import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import {
  Account,
  CreateAccountRequest,
  UpdateAccountRequest,
  AccountsListResponse,
  BalanceResponse,
  AddFundsRequest,
  WithdrawFundsRequest,
  UpdateCreditLimitRequest,
  UpdateCreditDatesRequest,
  AvailableCreditResponse,
  AccountValidationResult,
  AccountValidationError,
  AccountType,
  Currency
} from '../models';

// Interface for Account Service following SOLID principles
export interface IAccountService {
  // Basic CRUD operations
  createAccount(accountData: CreateAccountRequest): Observable<Account>;
  getAccountById(accountId: string): Observable<Account>;
  getAccountsByUser(userId: string, page?: number, pageSize?: number): Observable<AccountsListResponse>;
  updateAccount(accountId: string, updateData: UpdateAccountRequest): Observable<Account>;
  deleteAccount(accountId: string): Observable<void>;
  
  // Balance operations
  getAccountBalance(accountId: string): Observable<BalanceResponse>;
  updateAccountBalance(accountId: string, amount: number): Observable<BalanceResponse>;
}

// Interface for Wallet Service
export interface IWalletService {
  addFunds(accountId: string, fundsData: AddFundsRequest): Observable<BalanceResponse>;
  withdrawFunds(accountId: string, withdrawData: WithdrawFundsRequest): Observable<BalanceResponse>;
  validateWalletOperation(operationData: AddFundsRequest | WithdrawFundsRequest): AccountValidationResult;
}

// Interface for Credit Service
export interface ICreditService {
  updateCreditLimit(accountId: string, limitData: UpdateCreditLimitRequest): Observable<Account>;
  updateCreditDates(accountId: string, datesData: UpdateCreditDatesRequest): Observable<Account>;
  getAvailableCredit(accountId: string): Observable<AvailableCreditResponse>;
  validateCreditLimit(creditLimit: number): AccountValidationResult;
  validateCreditDates(closingDate?: string, dueDate?: string): AccountValidationResult;
}

// Interface for Account Validation Service
export interface IAccountValidationService {
  validateAccountData(accountData: CreateAccountRequest): AccountValidationResult;
  validateAccountName(name: string): AccountValidationError | null;
  validateInitialBalance(balance: number, accountType: AccountType): AccountValidationError | null;
  validateDni(dni: string): AccountValidationError | null;
}

@Injectable({ providedIn: 'root' })
export class AccountService implements IAccountService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.accountServiceUrl}/accounts`;

  // Basic CRUD operations
  createAccount(accountData: CreateAccountRequest): Observable<Account> {
    const payload = {
      user_id: accountData.userId,
      account_type: accountData.accountType,
      name: accountData.name,
      description: accountData.description || "",
      currency: accountData.currency,
      initial_balance: accountData.initialBalance,
      is_active: accountData.isActive !== false, // Default to true
      credit_limit: accountData.creditLimit,
      closing_date: accountData.closingDate ? new Date(accountData.closingDate).toISOString() : null,
      due_date: accountData.dueDate ? new Date(accountData.dueDate).toISOString() : null,
      dni: accountData.dni
    };

    return this.http.post<any>(this.apiUrl, payload).pipe(
      map(response => this.mapBackendResponseToAccount(response)),
      catchError(error => {
        console.error('Error creating account:', error);
        throw this.handleHttpError(error);
      })
    );
  }

  getAccountById(accountId: string): Observable<Account> {
    return this.http.get<any>(`${this.apiUrl}/${accountId}`).pipe(
      map(response => this.mapBackendResponseToAccount(response)),
      catchError(error => {
        console.error('Error getting account:', error);
        throw this.handleHttpError(error);
      })
    );
  }

  getAccountsByUser(userId: string, page: number = 1, pageSize: number = 20): Observable<AccountsListResponse> {
    const params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());

    return this.http.get<any>(`${this.apiUrl}/user/${userId}`, { params }).pipe(
      map(response => {
        console.log('Backend response:', response); // Para debugging
        return this.mapBackendResponseToAccountsList(response);
      }),
      catchError(error => {
        console.error('Error getting accounts by user:', error);
        
        // Si es un error 404 (usuario sin cuentas), devolver lista vacía
        if (error.status === 404) {
          return of({
            accounts: [],
            total: 0,
            page: 1,
            pageSize: 20,
            totalPages: 1
          });
        }
        
        throw this.handleHttpError(error);
      })
    );
  }

  updateAccount(accountId: string, updateData: UpdateAccountRequest): Observable<Account> {
    const payload = {
      name: updateData.name,
      description: updateData.description,
      credit_limit: updateData.creditLimit,
      closing_date: updateData.closingDate,
      due_date: updateData.dueDate,
      dni: updateData.dni
    };

    return this.http.put<any>(`${this.apiUrl}/${accountId}`, payload).pipe(
      map(response => this.mapBackendResponseToAccount(response)),
      catchError(error => {
        console.error('Error updating account:', error);
        throw this.handleHttpError(error);
      })
    );
  }

  deleteAccount(accountId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/${accountId}`).pipe(
      catchError(error => {
        console.error('Error deleting account:', error);
        throw this.handleHttpError(error);
      })
    );
  }

  getAccountBalance(accountId: string): Observable<BalanceResponse> {
    return this.http.get<any>(`${this.apiUrl}/${accountId}/balance`).pipe(
      map(response => ({
        accountId: response.account_id,
        balance: response.balance
      })),
      catchError(error => {
        console.error('Error getting account balance:', error);
        throw this.handleHttpError(error);
      })
    );
  }

  updateAccountBalance(accountId: string, amount: number): Observable<BalanceResponse> {
    const payload = { amount };

    return this.http.put<any>(`${this.apiUrl}/${accountId}/balance`, payload).pipe(
      map(response => ({
        accountId: response.account_id,
        balance: response.balance
      })),
      catchError(error => {
        console.error('Error updating account balance:', error);
        throw this.handleHttpError(error);
      })
    );
  }

  // Helper methods
  private mapBackendResponseToAccount(response: any): Account {
    return {
      id: response.id,
      userId: response.user_id,
      accountType: response.account_type as AccountType,
      name: response.name,
      description: response.description,
      currency: response.currency as Currency,
      balance: response.balance,
      isActive: response.is_active,
      createdAt: response.created_at,
      updatedAt: response.updated_at,
      creditLimit: response.credit_limit,
      closingDate: response.closing_date,
      dueDate: response.due_date,
      dni: response.dni
    };
  }

  private mapBackendResponseToAccountsList(response: any): AccountsListResponse {
    // Handle null or undefined response
    if (!response) {
      return {
        accounts: [],
        total: 0,
        page: 1,
        pageSize: 20,
        totalPages: 1
      };
    }

    if (Array.isArray(response)) {
      // Backend returns array directly
      return {
        accounts: response.map(item => this.mapBackendResponseToAccount(item)),
        total: response.length,
        page: 1,
        pageSize: response.length,
        totalPages: 1
      };
    }
    
    // Backend returns paginated response
    const accounts = response.data ? response.data : [];
    return {
      accounts: accounts.map((item: any) => this.mapBackendResponseToAccount(item)),
      total: response.pagination?.total || accounts.length,
      page: response.pagination?.current_page || 1,
      pageSize: response.pagination?.page_size || 20,
      totalPages: response.pagination?.total_pages || 1
    };
  }

  private handleHttpError(error: any): Error {
    let errorMessage = 'Ha ocurrido un error inesperado';
    
    if (error.error?.error) {
      errorMessage = error.error.error;
    } else if (error.message) {
      errorMessage = error.message;
    }
    
    return new Error(errorMessage);
  }
}