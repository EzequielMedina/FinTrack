import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

// Interfaces para los reportes
export interface Period {
  start_date: string;
  end_date: string;
  days: number;
}

export interface TransactionSummary {
  total_transactions: number;
  total_income: number;
  total_expenses: number;
  net_balance: number;
  avg_transaction: number;
}

export interface TransactionByType {
  type: string;
  count: number;
  amount: number;
  percentage: number;
}

export interface TransactionByPeriod {
  period: string;
  date: string;
  income: number;
  expenses: number;
  net: number;
  count: number;
}

export interface TransactionItem {
  id: string;
  description: string;
  amount: number;
  type: string;
  date: string;
  merchant_name?: string;
}

export interface TransactionReport {
  user_id: string;
  period: Period;
  summary: TransactionSummary;
  by_type: TransactionByType[];
  by_period: TransactionByPeriod[];
  top_expenses: TransactionItem[];
}

export interface InstallmentSummary {
  total_plans: number;
  active_plans: number;
  total_amount: number;
  paid_amount: number;
  remaining_amount: number;
  overdue_amount: number;
  next_payment_amount: number;
  next_payment_date?: string;
  completion_percentage: number;
}

export interface InstallmentPlan {
  id: string;
  card_id: string;
  card_last_four: string;
  total_amount: number;
  installments_count: number;
  installment_amount: number;
  paid_installments: number;
  remaining_amount: number;
  status: string;
  description?: string;
  merchant_name?: string;
  start_date: string;
  next_due_date?: string;
  completion_percentage: number;
}

export interface UpcomingPayment {
  installment_id: string;
  plan_id: string;
  card_last_four: string;
  amount: number;
  due_date: string;
  days_until_due: number;
  description?: string;
  merchant_name?: string;
}

export interface OverduePayment {
  installment_id: string;
  plan_id: string;
  card_last_four: string;
  amount: number;
  due_date: string;
  days_overdue: number;
  late_fee: number;
  description?: string;
  merchant_name?: string;
}

export interface InstallmentReport {
  user_id: string;
  summary: InstallmentSummary;
  plans: InstallmentPlan[];
  upcoming_payments: UpcomingPayment[];
  overdue_payments: OverduePayment[];
}

export interface AccountSummary {
  total_balance: number;
  total_accounts: number;
  total_cards: number;
  total_credit_limit: number;
  total_credit_used: number;
  available_credit: number;
  credit_utilization: number;
  net_worth: number;
}

export interface AccountDetail {
  id: string;
  account_type: string;
  name: string;
  currency: string;
  balance: number;
  credit_limit?: number;
  is_active: boolean;
}

export interface CardDetail {
  id: string;
  account_id: string;
  card_type: string;
  card_brand: string;
  last_four_digits: string;
  holder_name: string;
  status: string;
  credit_limit?: number;
  current_balance?: number;
  available_credit?: number;
  nickname?: string;
}

export interface AccountDistribution {
  account_type: string;
  count: number;
  total_balance: number;
  percentage: number;
}

export interface AccountReport {
  user_id: string;
  summary: AccountSummary;
  accounts: AccountDetail[];
  cards: CardDetail[];
  distribution: AccountDistribution[];
}

export interface ExpenseIncomeSummary {
  total_income: number;
  total_expenses: number;
  net_balance: number;
  savings_rate: number;
  expense_ratio: number;
  avg_daily_income: number;
  avg_daily_expense: number;
}

export interface ExpenseIncomeByPeriod {
  period: string;
  date: string;
  income: number;
  expenses: number;
  net: number;
  savings_rate: number;
}

export interface ExpenseIncomeByCategory {
  category: string;
  type: string;
  amount: number;
  count: number;
  percentage: number;
}

export interface TrendAnalysis {
  incomes_trend: string;
  expenses_trend: string;
  net_trend: string;
  income_change: number;
  expense_change: number;
  forecast?: {
    next_month_income: number;
    next_month_expenses: number;
    next_month_net: number;
  };
}

export interface ExpenseIncomeReport {
  user_id: string;
  period: Period;
  summary: ExpenseIncomeSummary;
  by_period: ExpenseIncomeByPeriod[];
  by_category: ExpenseIncomeByCategory[];
  trend: TrendAnalysis;
}

export interface NotificationSummary {
  total_notifications: number;
  total_job_runs: number;
  successful_sent: number;
  failed: number;
  success_rate: number;
  failure_rate: number;
  avg_emails_per_run: number;
}

export interface NotificationByDay {
  date: string;
  day: string;
  sent: number;
  failed: number;
  success_rate: number;
}

export interface NotificationByStatus {
  status: string;
  count: number;
  percentage: number;
}

export interface JobRunDetail {
  id: string;
  started_at: string;
  completed_at?: string;
  status: string;
  cards_found: number;
  emails_sent: number;
  errors: number;
  duration?: string;
  error_message?: string;
}

export interface NotificationReport {
  period: Period;
  summary: NotificationSummary;
  by_day: NotificationByDay[];
  by_status: NotificationByStatus[];
  job_runs: JobRunDetail[];
}

@Injectable({
  providedIn: 'root'
})
export class ReportService {
  private apiUrl = environment.reportServiceUrl || '/api/v1/reports';

  constructor(private http: HttpClient) {}

  /**
   * Obtiene el reporte de transacciones
   */
  getTransactionReport(
    userId: string,
    startDate?: string,
    endDate?: string,
    type?: string,
    groupBy?: string
  ): Observable<TransactionReport> {
    let params = new HttpParams().set('user_id', userId);
    
    if (startDate) params = params.set('start_date', startDate);
    if (endDate) params = params.set('end_date', endDate);
    if (type) params = params.set('type', type);
    if (groupBy) params = params.set('group_by', groupBy);

    return this.http.get<TransactionReport>(`${this.apiUrl}/transactions`, { params });
  }

  /**
   * Obtiene el reporte de cuotas
   */
  getInstallmentReport(userId: string, status?: string): Observable<InstallmentReport> {
    let params = new HttpParams().set('user_id', userId);
    
    if (status) params = params.set('status', status);

    return this.http.get<InstallmentReport>(`${this.apiUrl}/installments`, { params });
  }

  /**
   * Obtiene el reporte de cuentas
   */
  getAccountReport(userId: string): Observable<AccountReport> {
    const params = new HttpParams().set('user_id', userId);
    return this.http.get<AccountReport>(`${this.apiUrl}/accounts`, { params });
  }

  /**
   * Obtiene el reporte de gastos vs ingresos
   */
  getExpenseIncomeReport(
    userId: string,
    startDate: string,
    endDate: string,
    groupBy?: string
  ): Observable<ExpenseIncomeReport> {
    let params = new HttpParams()
      .set('user_id', userId)
      .set('start_date', startDate)
      .set('end_date', endDate);
    
    if (groupBy) params = params.set('group_by', groupBy);

    return this.http.get<ExpenseIncomeReport>(`${this.apiUrl}/expenses-income`, { params });
  }

  /**
   * Obtiene el reporte de notificaciones (solo admin)
   */
  getNotificationReport(startDate?: string, endDate?: string): Observable<NotificationReport> {
    let params = new HttpParams();
    
    if (startDate) params = params.set('start_date', startDate);
    if (endDate) params = params.set('end_date', endDate);

    return this.http.get<NotificationReport>(`${this.apiUrl}/notifications`, { params });
  }

  // ==================== MÉTODOS DE DESCARGA PDF ====================

  /**
   * Descarga el reporte de transacciones en PDF
   */
  downloadTransactionReportPDF(
    userId: string,
    startDate?: string,
    endDate?: string,
    type?: string
  ): Observable<Blob> {
    let params = new HttpParams().set('user_id', userId);
    
    if (startDate) params = params.set('start_date', startDate);
    if (endDate) params = params.set('end_date', endDate);
    if (type) params = params.set('type', type);

    return this.http.get(`${this.apiUrl}/transactions/pdf`, { 
      params, 
      responseType: 'blob' 
    });
  }

  /**
   * Descarga el reporte de cuotas en PDF
   */
  downloadInstallmentReportPDF(userId: string, status?: string): Observable<Blob> {
    let params = new HttpParams().set('user_id', userId);
    
    if (status) params = params.set('status', status);

    return this.http.get(`${this.apiUrl}/installments/pdf`, { 
      params, 
      responseType: 'blob' 
    });
  }

  /**
   * Descarga el reporte de cuentas en PDF
   */
  downloadAccountReportPDF(userId: string): Observable<Blob> {
    const params = new HttpParams().set('user_id', userId);
    
    return this.http.get(`${this.apiUrl}/accounts/pdf`, { 
      params, 
      responseType: 'blob' 
    });
  }

  /**
   * Descarga el reporte de gastos vs ingresos en PDF
   */
  downloadExpenseIncomeReportPDF(
    userId: string,
    startDate: string,
    endDate: string,
    groupBy?: string
  ): Observable<Blob> {
    let params = new HttpParams()
      .set('user_id', userId)
      .set('start_date', startDate)
      .set('end_date', endDate);
    
    if (groupBy) params = params.set('group_by', groupBy);

    return this.http.get(`${this.apiUrl}/expenses-income/pdf`, { 
      params, 
      responseType: 'blob' 
    });
  }

  /**
   * Método auxiliar para descargar un blob como archivo
   */
  private downloadBlob(blob: Blob, filename: string): void {
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  }

  /**
   * Descarga un reporte PDF con manejo automático del nombre de archivo
   */
  downloadPDF(blob: Blob, reportType: string): void {
    const date = new Date().toISOString().split('T')[0];
    const filename = `reporte-${reportType}-${date}.pdf`;
    this.downloadBlob(blob, filename);
  }
}
