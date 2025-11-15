import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import {
  Account,
  UpdateCreditLimitRequest,
  UpdateCreditDatesRequest,
  AvailableCreditResponse,
  AccountValidationResult,
  AccountValidationError
} from '../models';
import { ICreditService } from './account.service';

@Injectable({ providedIn: 'root' })
export class CreditService implements ICreditService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.accountServiceUrl}/accounts`;

  // Credit card business logic constants
  private readonly MIN_CREDIT_LIMIT = 0;
  private readonly MAX_CREDIT_LIMIT = 999999999;
  private readonly CREDIT_LIMIT_STEP = 100; // Credit limits should be multiples of 100

  updateCreditLimit(accountId: string, limitData: UpdateCreditLimitRequest): Observable<Account> {
    // Validate credit limit before making request
    const validationResult = this.validateCreditLimit(limitData.creditLimit);
    if (!validationResult.isValid) {
      throw new Error(validationResult.errors[0]?.message || 'Límite de crédito inválido');
    }

    const payload = {
      credit_limit: limitData.creditLimit
    };

    return this.http.put<any>(`${this.apiUrl}/${accountId}/credit-limit`, payload).pipe(
      map(response => this.mapBackendResponseToAccount(response)),
      catchError(error => {
        console.error('Error updating credit limit:', error);
        throw this.handleCreditOperationError(error);
      })
    );
  }

  updateCreditDates(accountId: string, datesData: UpdateCreditDatesRequest): Observable<Account> {
    // Validate credit dates before making request
    const validationResult = this.validateCreditDates(datesData.closingDate, datesData.dueDate);
    if (!validationResult.isValid) {
      throw new Error(validationResult.errors[0]?.message || 'Fechas de crédito inválidas');
    }

    const payload = {
      closing_date: datesData.closingDate,
      due_date: datesData.dueDate
    };

    return this.http.put<any>(`${this.apiUrl}/${accountId}/credit-dates`, payload).pipe(
      map(response => this.mapBackendResponseToAccount(response)),
      catchError(error => {
        console.error('Error updating credit dates:', error);
        throw this.handleCreditOperationError(error);
      })
    );
  }

  getAvailableCredit(accountId: string): Observable<AvailableCreditResponse> {
    return this.http.get<any>(`${this.apiUrl}/${accountId}/available-credit`).pipe(
      map(response => this.mapBackendResponseToAvailableCredit(response)),
      catchError(error => {
        console.error('Error getting available credit:', error);
        throw this.handleCreditOperationError(error);
      })
    );
  }

  validateCreditLimit(creditLimit: number): AccountValidationResult {
    const errors: AccountValidationError[] = [];

    const limitError = this.validateCreditLimitAmount(creditLimit);
    if (limitError) {
      errors.push(limitError);
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  validateCreditDates(closingDate?: string, dueDate?: string): AccountValidationResult {
    const errors: AccountValidationError[] = [];

    // Validate closing date if provided
    if (closingDate) {
      const closingDateError = this.validateCreditDate(closingDate, 'closing_date', 'fecha de corte');
      if (closingDateError) {
        errors.push(closingDateError);
      }
    }

    // Validate due date if provided
    if (dueDate) {
      const dueDateError = this.validateCreditDate(dueDate, 'due_date', 'fecha de vencimiento');
      if (dueDateError) {
        errors.push(dueDateError);
      }
    }

    // Validate date relationship if both dates are provided
    if (closingDate && dueDate && errors.length === 0) {
      const relationshipError = this.validateCreditDateRelationship(closingDate, dueDate);
      if (relationshipError) {
        errors.push(relationshipError);
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  // Business logic methods
  calculateMinimumPayment(usedCredit: number): number {
    // Standard minimum payment calculation (e.g., 3% of used credit, minimum $50)
    const percentagePayment = usedCredit * 0.03;
    const minimumAmount = 50;
    return Math.max(percentagePayment, minimumAmount);
  }

  calculateInterestRate(creditLimit: number, creditScore?: number): number {
    // Example business logic for interest rate calculation
    let baseRate = 0.25; // 25% annual

    // Adjust based on credit limit (higher limits get better rates)
    if (creditLimit >= 100000) {
      baseRate -= 0.05;
    } else if (creditLimit >= 50000) {
      baseRate -= 0.02;
    }

    // Adjust based on credit score if available
    if (creditScore) {
      if (creditScore >= 800) {
        baseRate -= 0.08;
      } else if (creditScore >= 700) {
        baseRate -= 0.05;
      } else if (creditScore >= 600) {
        baseRate -= 0.02;
      }
    }

    return Math.max(baseRate, 0.12); // Minimum 12% annual rate
  }

  getNextClosingDate(currentClosingDate?: string): string {
    const today = new Date();
    let nextClosing: Date;

    if (currentClosingDate) {
      const current = new Date(currentClosingDate);
      nextClosing = new Date(current.getFullYear(), current.getMonth() + 1, current.getDate());
      
      // If the calculated date is in the past, move to the following month
      if (nextClosing <= today) {
        nextClosing = new Date(current.getFullYear(), current.getMonth() + 2, current.getDate());
      }
    } else {
      // Default to last day of current month if no current closing date
      nextClosing = new Date(today.getFullYear(), today.getMonth() + 1, 0);
    }

    return nextClosing.toISOString().split('T')[0];
  }

  getNextDueDate(closingDate: string, gracePeriodDays: number = 15): string {
    const closing = new Date(closingDate);
    const dueDate = new Date(closing.getTime() + (gracePeriodDays * 24 * 60 * 60 * 1000));
    return dueDate.toISOString().split('T')[0];
  }

  // Private validation methods
  private validateCreditLimitAmount(creditLimit: number): AccountValidationError | null {
    if (creditLimit < this.MIN_CREDIT_LIMIT) {
      return {
        field: 'creditLimit',
        message: `El límite de crédito mínimo es $${this.MIN_CREDIT_LIMIT.toLocaleString()}`
      };
    }

    if (creditLimit > this.MAX_CREDIT_LIMIT) {
      return {
        field: 'creditLimit',
        message: `El límite de crédito máximo es $${this.MAX_CREDIT_LIMIT.toLocaleString()}`
      };
    }

    if (creditLimit % this.CREDIT_LIMIT_STEP !== 0) {
      return {
        field: 'creditLimit',
        message: `El límite de crédito debe ser múltiplo de $${this.CREDIT_LIMIT_STEP}`
      };
    }

    return null;
  }

  private validateCreditDate(dateString: string, field: string, displayName: string): AccountValidationError | null {
    // Validate date format (YYYY-MM-DD)
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
    if (!dateRegex.test(dateString)) {
      return {
        field,
        message: `La ${displayName} debe tener el formato YYYY-MM-DD`
      };
    }

    const date = new Date(dateString);
    
    // Validate that it's a valid date
    if (isNaN(date.getTime())) {
      return {
        field,
        message: `La ${displayName} no es válida`
      };
    }

    // Validate that the date is not in the past (allow today)
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    
    if (date < today) {
      return {
        field,
        message: `La ${displayName} no puede ser anterior al día de hoy`
      };
    }

    // Validate that the date is not too far in the future (e.g., max 2 years)
    const maxFutureDate = new Date();
    maxFutureDate.setFullYear(maxFutureDate.getFullYear() + 2);
    
    if (date > maxFutureDate) {
      return {
        field,
        message: `La ${displayName} no puede ser más de 2 años en el futuro`
      };
    }

    return null;
  }

  private validateCreditDateRelationship(closingDate: string, dueDate: string): AccountValidationError | null {
    const closing = new Date(closingDate);
    const due = new Date(dueDate);

    // Due date should be after closing date
    if (due <= closing) {
      return {
        field: 'dueDate',
        message: 'La fecha de vencimiento debe ser posterior a la fecha de corte'
      };
    }

    // Due date should not be more than 45 days after closing date
    const maxDaysBetween = 45;
    const daysDifference = Math.ceil((due.getTime() - closing.getTime()) / (1000 * 60 * 60 * 24));
    
    if (daysDifference > maxDaysBetween) {
      return {
        field: 'dueDate',
        message: `La fecha de vencimiento no puede ser más de ${maxDaysBetween} días después de la fecha de corte`
      };
    }

    return null;
  }

  private mapBackendResponseToAccount(response: any): Account {
    return {
      id: response.id,
      userId: response.user_id,
      accountType: response.account_type,
      name: response.name,
      description: response.description,
      currency: response.currency,
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

  private mapBackendResponseToAvailableCredit(response: any): AvailableCreditResponse {
    return {
      accountId: response.account_id,
      creditLimit: response.credit_limit,
      usedCredit: response.used_credit,
      availableCredit: response.available_credit
    };
  }

  private handleCreditOperationError(error: any): Error {
    let errorMessage = 'Error en la operación de crédito';

    if (error.error?.error) {
      const backendError = error.error.error.toLowerCase();
      
      // Map specific backend errors to user-friendly messages
      if (backendError.includes('credit limit') || backendError.includes('límite de crédito')) {
        errorMessage = 'Error al actualizar el límite de crédito';
      } else if (backendError.includes('credit card') || backendError.includes('tarjeta de crédito')) {
        errorMessage = 'Esta operación solo está disponible para tarjetas de crédito';
      } else if (backendError.includes('account not found') || backendError.includes('cuenta no encontrada')) {
        errorMessage = 'La cuenta no existe o no está disponible';
      } else if (backendError.includes('validation') || backendError.includes('validación')) {
        errorMessage = 'Los datos ingresados no son válidos';
      } else {
        errorMessage = error.error.error;
      }
    } else if (error.status === 404) {
      errorMessage = 'La cuenta no existe';
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