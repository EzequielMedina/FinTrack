import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { AuthService } from './auth.service';
import {
  InstallmentPlan,
  Installment,
  InstallmentPreview,
  CreateInstallmentPlanRequest,
  InstallmentPreviewRequest,
  PayInstallmentRequest,
  ChargeWithInstallmentsRequest,
  ChargeWithInstallmentsResponse,
  InstallmentPlanResponse,
  InstallmentSummary,
  MonthlyInstallmentLoad,
  InstallmentPlansListResponse,
  InstallmentsListResponse,
  InstallmentError
} from '../models';

@Injectable({ providedIn: 'root' })
export class InstallmentService {
  private readonly http = inject(HttpClient);
  private readonly authService = inject(AuthService);
  private readonly apiUrl = `${environment.accountServiceUrl}`;

  /**
   * Preview installment plan calculation
   */
  previewInstallments(cardId: string, request: InstallmentPreviewRequest): Observable<InstallmentPreview> {
    const url = `${this.apiUrl}/cards/${cardId}/installments/preview`;
    return this.http.post<InstallmentPreview>(url, request).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Create a new installment plan (charge card with installments)
   */
  createInstallmentPlan(request: ChargeWithInstallmentsRequest): Observable<ChargeWithInstallmentsResponse> {
    const url = `${this.apiUrl}/cards/${request.cardId}/charge-installments`;
    
    const currentUser = this.authService.getCurrentUser();
    const headers: Record<string, string> = {};
    
    if (currentUser) {
      headers['X-User-ID'] = currentUser.id;
    }
    
    return this.http.post<ChargeWithInstallmentsResponse>(url, request, { headers }).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Get installment plans by card
   */
  getInstallmentPlansByCard(
    cardId: string, 
    page: number = 1, 
    pageSize: number = 10
  ): Observable<InstallmentPlansListResponse> {
    const url = `${this.apiUrl}/cards/${cardId}/installment-plans`;
    const params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());

    return this.http.get<{data: any[], pagination: any}>(url, { params }).pipe(
      map(response => {
        // Transform backend response to frontend format
        return {
          plans: response.data.map(plan => ({
            id: plan.id,
            transactionId: plan.transaction_id,
            cardId: plan.card_id,
            userId: plan.user_id,
            totalAmount: plan.total_amount,
            installmentsCount: plan.installments_count,
            installmentAmount: plan.installment_amount,
            startDate: plan.start_date,
            merchantName: plan.merchant_name,
            merchantId: plan.merchant_id,
            description: plan.description,
            status: plan.status,
            paidInstallments: plan.paid_installments,
            remainingAmount: plan.remaining_amount,
            interestRate: plan.interest_rate,
            totalInterest: plan.total_interest,
            adminFee: plan.admin_fee,
            createdAt: plan.created_at,
            updatedAt: plan.updated_at,
            completedAt: plan.completed_at,
            cancelledAt: plan.cancelled_at
          })),
          total: response.pagination.total_items,
          page: response.pagination.current_page,
          pageSize: response.pagination.page_size,
          totalPages: response.pagination.total_pages
        };
      }),
      catchError(this.handleError)
    );
  }

  /**
   * Get specific installment plan details
   */
  getInstallmentPlanDetails(planId: string): Observable<InstallmentPlanResponse> {
    const url = `${this.apiUrl}/installment-plans/${planId}`;
    return this.http.get<InstallmentPlanResponse>(url).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Get installment plans by user
   */
  getUserInstallmentPlans(
    userId: string,
    status?: string,
    page: number = 1,
    pageSize: number = 10
  ): Observable<InstallmentPlansListResponse> {
    const url = `${this.apiUrl}/installment-plans`;
    let params = new HttpParams()
      .set('userId', userId)
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());

    if (status) {
      params = params.set('status', status);
    }

    return this.http.get<{data: any[], pagination: any}>(url, { params }).pipe(
      map(response => {
        // Transform backend response to frontend format
        return {
          plans: response.data.map(plan => ({
            id: plan.id,
            transactionId: plan.transaction_id,
            cardId: plan.card_id,
            userId: plan.user_id,
            totalAmount: plan.total_amount,
            installmentsCount: plan.installments_count,
            installmentAmount: plan.installment_amount,
            startDate: plan.start_date,
            merchantName: plan.merchant_name,
            merchantId: plan.merchant_id,
            description: plan.description,
            status: plan.status,
            paidInstallments: plan.paid_installments,
            remainingAmount: plan.remaining_amount,
            interestRate: plan.interest_rate,
            totalInterest: plan.total_interest,
            adminFee: plan.admin_fee,
            createdAt: plan.created_at,
            updatedAt: plan.updated_at,
            completedAt: plan.completed_at,
            cancelledAt: plan.cancelled_at
          })),
          total: response.pagination.total_items,
          page: response.pagination.current_page,
          pageSize: response.pagination.page_size,
          totalPages: response.pagination.total_pages
        };
      }),
      catchError(this.handleError)
    );
  }

  /**
   * Pay a specific installment
   */
  payInstallment(request: PayInstallmentRequest): Observable<Installment> {
    const url = `${this.apiUrl}/installments/${request.installmentId}/pay`;
    return this.http.post<Installment>(url, request).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Cancel an installment plan
   */
  cancelInstallmentPlan(planId: string, reason: string): Observable<InstallmentPlan> {
    const url = `${this.apiUrl}/installment-plans/${planId}/cancel`;
    const body = { reason };
    return this.http.post<InstallmentPlan>(url, body).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Get overdue installments
   */
  getOverdueInstallments(
    userId?: string,
    page: number = 1,
    pageSize: number = 10
  ): Observable<InstallmentsListResponse> {
    const url = `${this.apiUrl}/installments/overdue`;
    let params = new HttpParams()
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());

    if (userId) {
      params = params.set('userId', userId);
    }

    return this.http.get<InstallmentsListResponse>(url, { params }).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Get upcoming installments
   */
  getUpcomingInstallments(
    days: number = 30,
    userId?: string,
    page: number = 1,
    pageSize: number = 10
  ): Observable<InstallmentsListResponse> {
    const url = `${this.apiUrl}/installments/upcoming`;
    let params = new HttpParams()
      .set('days', days.toString())
      .set('page', page.toString())
      .set('pageSize', pageSize.toString());

    if (userId) {
      params = params.set('userId', userId);
    }

    return this.http.get<InstallmentsListResponse>(url, { params }).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Get installments summary for a user
   */
  getInstallmentsSummary(userId: string): Observable<InstallmentSummary> {
    const url = `${this.apiUrl}/installments/summary`;
    const params = new HttpParams().set('userId', userId);

    return this.http.get<InstallmentSummary>(url, { params }).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Get monthly installment load
   */
  getMonthlyInstallmentLoad(
    userId: string,
    year: number,
    month: number
  ): Observable<MonthlyInstallmentLoad> {
    const url = `${this.apiUrl}/installments/monthly-load`;
    const params = new HttpParams()
      .set('userId', userId)
      .set('year', year.toString())
      .set('month', month.toString());

    return this.http.get<MonthlyInstallmentLoad>(url, { params }).pipe(
      catchError(this.handleError)
    );
  }

  /**
   * Calculate installment options for a given amount
   */
  calculateInstallmentOptions(
    amount: number,
    availableInstallments: number[] = [3, 6, 12, 18, 24]
  ): Observable<InstallmentPreview[]> {
    // This method calculates multiple preview options
    const baseDate = new Date();
    baseDate.setDate(1); // Start from first day of next month
    baseDate.setMonth(baseDate.getMonth() + 1);

    const previews$ = availableInstallments.map(installments => {
      const request: InstallmentPreviewRequest = {
        amount,
        installmentsCount: installments,
        startDate: baseDate.toISOString().split('T')[0],
        interestRate: this.getInterestRateForInstallments(installments),
        adminFee: this.getAdminFeeForAmount(amount)
      };

      // Note: This would need a cardId in a real implementation
      // For now, we'll return a mock calculation
      return this.calculateMockPreview(request);
    });

    return new Observable(observer => {
      Promise.all(previews$).then(results => {
        observer.next(results);
        observer.complete();
      }).catch(error => {
        observer.error(error);
      });
    });
  }

  /**
   * Validate installment parameters
   */
  validateInstallmentPlan(request: CreateInstallmentPlanRequest): Observable<boolean> {
    // Client-side validation
    if (request.totalAmount <= 0) {
      return throwError(() => new Error('Amount must be greater than 0'));
    }

    if (request.installmentsCount < 1 || request.installmentsCount > 24) {
      return throwError(() => new Error('Installments count must be between 1 and 24'));
    }

    if (new Date(request.startDate) <= new Date()) {
      return throwError(() => new Error('Start date must be in the future'));
    }

    return new Observable(observer => {
      observer.next(true);
      observer.complete();
    });
  }

  /**
   * Get installment configuration
   */
  getInstallmentConfig(): Observable<any> {
    // This could be an API call to get configuration from backend
    return new Observable(observer => {
      observer.next({
        maxInstallments: 24,
        minAmount: 1000,
        maxAmount: 500000,
        defaultInterestRate: 15.5,
        adminFeePercentage: 2.5,
        allowedInstallmentCounts: [3, 6, 9, 12, 15, 18, 21, 24]
      });
      observer.complete();
    });
  }

  // Private helper methods
  private getInterestRateForInstallments(installments: number): number {
    // Simple interest rate calculation based on installment count
    if (installments <= 6) return 12.5;
    if (installments <= 12) return 15.5;
    if (installments <= 18) return 18.0;
    return 21.5;
  }

  private getAdminFeeForAmount(amount: number): number {
    // Simple admin fee calculation - 2.5% with minimum of $50
    return Math.max(amount * 0.025, 50);
  }

  private calculateMockPreview(request: InstallmentPreviewRequest): Promise<InstallmentPreview> {
    // Mock calculation for demonstration
    return new Promise(resolve => {
      const interestRate = request.interestRate || this.getInterestRateForInstallments(request.installmentsCount);
      const adminFee = request.adminFee || this.getAdminFeeForAmount(request.amount);
      const totalInterest = (request.amount * interestRate * request.installmentsCount) / 1200; // Monthly rate
      const totalToPay = request.amount + totalInterest + adminFee;
      const installmentAmount = totalToPay / request.installmentsCount;

      const installments: any[] = [];
      const startDate = new Date(request.startDate);

      for (let i = 0; i < request.installmentsCount; i++) {
        const dueDate = new Date(startDate);
        dueDate.setMonth(startDate.getMonth() + i);

        installments.push({
          number: i + 1,
          amount: installmentAmount,
          dueDate: dueDate.toISOString().split('T')[0],
          principal: request.amount / request.installmentsCount,
          interest: totalInterest / request.installmentsCount,
          fee: i === 0 ? adminFee : 0 // Admin fee in first installment
        });
      }

      resolve({
        totalAmount: request.amount,
        installmentsCount: request.installmentsCount,
        installmentAmount,
        startDate: request.startDate,
        interestRate,
        totalInterest,
        adminFee,
        totalToPay,
        installments
      });
    });
  }

  private handleError = (error: any): Observable<never> => {
    console.error('InstallmentService Error:', error);
    
    let errorMessage = 'An unexpected error occurred';
    let errorCode = 'UNKNOWN_ERROR';

    if (error.error) {
      errorMessage = error.error.message || error.error.error || errorMessage;
      errorCode = error.error.code || errorCode;
    } else if (error.message) {
      errorMessage = error.message;
    }

    const installmentError: InstallmentError = {
      code: errorCode,
      message: errorMessage
    };

    return throwError(() => installmentError);
  };
}