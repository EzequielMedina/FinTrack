import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';

export interface DebitCardTransactionRequest {
  amount: number;
  description: string;
  merchantName?: string;
  reference?: string;
}

export interface DebitCardBalanceResponse {
  cardId: string;
  accountBalance: number;
  availableBalance: number;
}

@Injectable({ providedIn: 'root' })
export class DebitCardService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.accountServiceUrl}/cards`;

  /**
   * Process a debit card transaction (purchase/withdrawal)
   */
  processTransaction(cardId: string, transactionData: DebitCardTransactionRequest, accountId?: string): Observable<DebitCardBalanceResponse> {
    return this.http.post<any>(`${this.apiUrl}/${cardId}/transaction`, transactionData).pipe(
      map(response => this.mapToDebitCardBalanceResponse(response)),
      catchError(error => {
        console.error('Error processing debit card transaction:', error);
        throw this.handleDebitCardError(error);
      })
    );
  }

  /**
   * Get available balance for debit card (account balance)
   */
  getAvailableBalance(cardId: string): Observable<DebitCardBalanceResponse> {
    return this.http.get<any>(`${this.apiUrl}/${cardId}/balance`).pipe(
      map(response => this.mapToDebitCardBalanceResponse(response)),
      catchError(error => {
        console.error('Error getting debit card balance:', error);
        throw this.handleDebitCardError(error);
      })
    );
  }

  // Helper methods
  private mapToDebitCardBalanceResponse(response: any): DebitCardBalanceResponse {
    return {
      cardId: response.card_id,
      accountBalance: response.account_balance,
      availableBalance: response.available_balance
    };
  }

  private handleDebitCardError(error: any): Error {
    let errorMessage = 'Error en la operación de tarjeta de débito';

    if (error.error?.error) {
      const backendError = error.error.error.toLowerCase();
      
      if (backendError.includes('insufficient funds') || backendError.includes('fondos insuficientes')) {
        errorMessage = 'Fondos insuficientes en la cuenta';
      } else if (backendError.includes('card not found') || backendError.includes('tarjeta no encontrada')) {
        errorMessage = 'La tarjeta no existe';
      } else if (backendError.includes('expired') || backendError.includes('vencida')) {
        errorMessage = 'La tarjeta está vencida';
      } else if (backendError.includes('blocked') || backendError.includes('bloqueada')) {
        errorMessage = 'La tarjeta está bloqueada';
      } else {
        errorMessage = error.error.error;
      }
    } else if (error.status === 404) {
      errorMessage = 'La tarjeta no existe';
    } else if (error.status === 400) {
      errorMessage = 'Los datos ingresados no son válidos';
    } else if (error.status === 401) {
      errorMessage = 'No tiene permisos para realizar esta operación';
    } else if (error.status >= 500) {
      errorMessage = 'Error del servidor. Intente nuevamente más tarde';
    }

    return new Error(errorMessage);
  }
}