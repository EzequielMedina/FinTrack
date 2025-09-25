import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Card } from '../models';

export interface CreditCardChargeRequest {
  amount: number;
  description: string;
  reference?: string;
}

export interface CreditCardPaymentRequest {
  amount: number;
  paymentMethod: 'bank_transfer' | 'debit_card' | 'cash';
  reference?: string;
}

export interface CreditCardStatement {
  cardId: string;
  period: string;
  previousBalance: number;
  purchases: number;
  payments: number;
  currentBalance: number;
  creditLimit: number;
  availableCredit: number;
  minimumPayment: number;
  dueDate: string;
  closingDate: string;
}

export interface CreditCardBalanceResponse {
  cardId: string;
  balance: number;
  creditLimit: number;
  availableCredit: number;
  minimumPayment: number;
  dueDate?: string;
}

@Injectable({ providedIn: 'root' })
export class CreditCardService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.accountServiceUrl}/cards`;

  /**
   * Make a charge to a credit card
   */
  charge(cardId: string, chargeData: CreditCardChargeRequest): Observable<CreditCardBalanceResponse> {
    return this.http.post<any>(`${this.apiUrl}/${cardId}/charge`, chargeData).pipe(
      map(response => this.mapToCreditCardBalanceResponse(response)),
      catchError(error => {
        console.error('Error charging credit card:', error);
        throw this.handleCreditCardError(error);
      })
    );
  }

  /**
   * Make a payment to a credit card
   */
  payment(cardId: string, paymentData: CreditCardPaymentRequest): Observable<CreditCardBalanceResponse> {
    return this.http.post<any>(`${this.apiUrl}/${cardId}/payment`, paymentData).pipe(
      map(response => this.mapToCreditCardBalanceResponse(response)),
      catchError(error => {
        console.error('Error processing credit card payment:', error);
        throw this.handleCreditCardError(error);
      })
    );
  }

  /**
   * Get credit card balance and available credit
   */
  getBalance(cardId: string): Observable<CreditCardBalanceResponse> {
    return this.http.get<any>(`${this.apiUrl}/${cardId}/balance`).pipe(
      map(response => this.mapToCreditCardBalanceResponse(response)),
      catchError(error => {
        console.error('Error getting credit card balance:', error);
        throw this.handleCreditCardError(error);
      })
    );
  }

  /**
   * Get credit card statement for a specific period
   */
  getStatement(cardId: string, year: number, month: number): Observable<CreditCardStatement> {
    return this.http.get<any>(`${this.apiUrl}/${cardId}/statement/${year}/${month}`).pipe(
      map(response => this.mapToCreditCardStatement(response)),
      catchError(error => {
        console.error('Error getting credit card statement:', error);
        throw this.handleCreditCardError(error);
      })
    );
  }

  /**
   * Update credit limit for a credit card
   */
  updateCreditLimit(cardId: string, newLimit: number): Observable<Card> {
    return this.http.put<any>(`${this.apiUrl}/${cardId}/credit-limit`, { credit_limit: newLimit }).pipe(
      catchError(error => {
        console.error('Error updating credit limit:', error);
        throw this.handleCreditCardError(error);
      })
    );
  }

  /**
   * Update closing and due dates for a credit card
   */
  updateBillingDates(cardId: string, closingDate: string, dueDate: string): Observable<Card> {
    return this.http.put<any>(`${this.apiUrl}/${cardId}/billing-dates`, { 
      closing_date: closingDate,
      due_date: dueDate 
    }).pipe(
      catchError(error => {
        console.error('Error updating billing dates:', error);
        throw this.handleCreditCardError(error);
      })
    );
  }

  // Helper methods
  private mapToCreditCardBalanceResponse(response: any): CreditCardBalanceResponse {
    return {
      cardId: response.card_id,
      balance: response.balance,
      creditLimit: response.credit_limit,
      availableCredit: response.available_credit,
      minimumPayment: response.minimum_payment,
      dueDate: response.due_date
    };
  }

  private mapToCreditCardStatement(response: any): CreditCardStatement {
    return {
      cardId: response.card_id,
      period: response.period,
      previousBalance: response.previous_balance,
      purchases: response.purchases,
      payments: response.payments,
      currentBalance: response.current_balance,
      creditLimit: response.credit_limit,
      availableCredit: response.available_credit,
      minimumPayment: response.minimum_payment,
      dueDate: response.due_date,
      closingDate: response.closing_date
    };
  }

  private handleCreditCardError(error: any): Error {
    let errorMessage = 'Error en la operación de tarjeta de crédito';

    if (error.error?.error) {
      const backendError = error.error.error.toLowerCase();
      
      if (backendError.includes('insufficient credit') || backendError.includes('límite excedido')) {
        errorMessage = 'Límite de crédito insuficiente';
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