import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import {
  AddFundsRequest,
  WithdrawFundsRequest,
  BalanceResponse,
  AccountValidationResult,
  AccountValidationError
} from '../models';
import { IWalletService } from './account.service';

@Injectable({ providedIn: 'root' })
export class WalletService implements IWalletService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = `${environment.accountServiceUrl}/api/accounts`;

  // Minimum and maximum amounts for wallet operations
  private readonly MIN_OPERATION_AMOUNT = 0.01;
  private readonly MAX_OPERATION_AMOUNT = 999999.99;
  private readonly MIN_DESCRIPTION_LENGTH = 3;
  private readonly MAX_DESCRIPTION_LENGTH = 255;
  private readonly MAX_REFERENCE_LENGTH = 50;

  addFunds(accountId: string, fundsData: AddFundsRequest): Observable<BalanceResponse> {
    // Validate input before making request
    const validationResult = this.validateWalletOperation(fundsData);
    if (!validationResult.isValid) {
      throw new Error(validationResult.errors[0]?.message || 'Datos de operación inválidos');
    }

    const payload = {
      amount: fundsData.amount,
      description: fundsData.description.trim(),
      reference: fundsData.reference?.trim()
    };

    return this.http.post<any>(`${this.apiUrl}/${accountId}/add-funds`, payload).pipe(
      map(response => this.mapBackendResponseToBalanceResponse(response)),
      catchError(error => {
        console.error('Error adding funds:', error);
        throw this.handleWalletOperationError(error);
      })
    );
  }

  withdrawFunds(accountId: string, withdrawData: WithdrawFundsRequest): Observable<BalanceResponse> {
    // Validate input before making request
    const validationResult = this.validateWalletOperation(withdrawData);
    if (!validationResult.isValid) {
      throw new Error(validationResult.errors[0]?.message || 'Datos de operación inválidos');
    }

    const payload = {
      amount: withdrawData.amount,
      description: withdrawData.description.trim(),
      reference: withdrawData.reference?.trim()
    };

    return this.http.post<any>(`${this.apiUrl}/${accountId}/withdraw-funds`, payload).pipe(
      map(response => this.mapBackendResponseToBalanceResponse(response)),
      catchError(error => {
        console.error('Error withdrawing funds:', error);
        throw this.handleWalletOperationError(error);
      })
    );
  }

  validateWalletOperation(operationData: AddFundsRequest | WithdrawFundsRequest): AccountValidationResult {
    const errors: AccountValidationError[] = [];

    // Validate amount
    const amountError = this.validateOperationAmount(operationData.amount);
    if (amountError) {
      errors.push(amountError);
    }

    // Validate description
    const descriptionError = this.validateOperationDescription(operationData.description);
    if (descriptionError) {
      errors.push(descriptionError);
    }

    // Validate reference (optional)
    if (operationData.reference) {
      const referenceError = this.validateOperationReference(operationData.reference);
      if (referenceError) {
        errors.push(referenceError);
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    };
  }

  // Business logic for determining operation fees
  calculateOperationFee(amount: number, operationType: 'deposit' | 'withdrawal'): number {
    // Example business logic - could be configurable
    if (operationType === 'withdrawal' && amount > 10000) {
      return amount * 0.001; // 0.1% fee for large withdrawals
    }
    return 0; // No fee for deposits or small withdrawals
  }

  // Business logic for determining daily limits
  getRemainingDailyLimit(userId: string, operationType: 'deposit' | 'withdrawal'): Observable<number> {
    // This would typically call another endpoint to get user's daily limits
    // For now, return a fixed limit
    const dailyLimits = {
      deposit: 50000,
      withdrawal: 20000
    };
    
    // In a real implementation, you would:
    // 1. Get user's current daily usage from backend
    // 2. Calculate remaining limit
    // 3. Return the result
    
    return new Observable(observer => {
      observer.next(dailyLimits[operationType]);
      observer.complete();
    });
  }

  // Helper methods for validation
  private validateOperationAmount(amount: number): AccountValidationError | null {
    if (!amount || amount <= 0) {
      return {
        field: 'amount',
        message: 'El monto debe ser mayor a cero'
      };
    }

    if (amount < this.MIN_OPERATION_AMOUNT) {
      return {
        field: 'amount',
        message: `El monto mínimo es $${this.MIN_OPERATION_AMOUNT}`
      };
    }

    if (amount > this.MAX_OPERATION_AMOUNT) {
      return {
        field: 'amount',
        message: `El monto máximo es $${this.MAX_OPERATION_AMOUNT.toLocaleString()}`
      };
    }

    // Validate decimal places (max 2)
    if (!/^\d+(\.\d{1,2})?$/.test(amount.toString())) {
      return {
        field: 'amount',
        message: 'El monto puede tener máximo 2 decimales'
      };
    }

    return null;
  }

  private validateOperationDescription(description: string): AccountValidationError | null {
    if (!description || description.trim().length === 0) {
      return {
        field: 'description',
        message: 'La descripción es obligatoria'
      };
    }

    const trimmedDescription = description.trim();

    if (trimmedDescription.length < this.MIN_DESCRIPTION_LENGTH) {
      return {
        field: 'description',
        message: `La descripción debe tener al menos ${this.MIN_DESCRIPTION_LENGTH} caracteres`
      };
    }

    if (trimmedDescription.length > this.MAX_DESCRIPTION_LENGTH) {
      return {
        field: 'description',
        message: `La descripción no puede exceder ${this.MAX_DESCRIPTION_LENGTH} caracteres`
      };
    }

    // Validate description contains valid characters
    if (!/^[a-zA-Z0-9\s\-\.,áéíóúñÁÉÍÓÚÑ]+$/.test(trimmedDescription)) {
      return {
        field: 'description',
        message: 'La descripción contiene caracteres no válidos'
      };
    }

    return null;
  }

  private validateOperationReference(reference: string): AccountValidationError | null {
    const trimmedReference = reference.trim();

    if (trimmedReference.length > this.MAX_REFERENCE_LENGTH) {
      return {
        field: 'reference',
        message: `La referencia no puede exceder ${this.MAX_REFERENCE_LENGTH} caracteres`
      };
    }

    // Validate reference contains valid characters (more permissive than description)
    if (!/^[a-zA-Z0-9\s\-\._#]+$/.test(trimmedReference)) {
      return {
        field: 'reference',
        message: 'La referencia contiene caracteres no válidos'
      };
    }

    return null;
  }

  private mapBackendResponseToBalanceResponse(response: any): BalanceResponse {
    return {
      accountId: response.account_id,
      balance: response.balance
    };
  }

  private handleWalletOperationError(error: any): Error {
    let errorMessage = 'Error en la operación de billetera';

    if (error.error?.error) {
      const backendError = error.error.error.toLowerCase();
      
      // Map specific backend errors to user-friendly messages
      if (backendError.includes('insufficient funds') || backendError.includes('fondos insuficientes')) {
        errorMessage = 'Fondos insuficientes para realizar la operación';
      } else if (backendError.includes('account not found') || backendError.includes('cuenta no encontrada')) {
        errorMessage = 'La cuenta no existe o no está disponible';
      } else if (backendError.includes('validation') || backendError.includes('validación')) {
        errorMessage = 'Los datos ingresados no son válidos';
      } else if (backendError.includes('wallet') && backendError.includes('only')) {
        errorMessage = 'Esta operación solo está disponible para billeteras virtuales';
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