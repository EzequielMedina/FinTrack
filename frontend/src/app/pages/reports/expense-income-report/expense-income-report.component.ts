import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ReportService, ExpenseIncomeReport } from '../../../services/report.service';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-expense-income-report',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    RouterModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatButtonModule
  ],
  templateUrl: './expense-income-report.component.html',
  styleUrls: ['./expense-income-report.component.css']
})
export class ExpenseIncomeReportComponent implements OnInit {
  loading = false;
  error: string | null = null;
  reportData: ExpenseIncomeReport | null = null;
  downloadingPDF = false;
  startDate: string = '';
  endDate: string = '';

  constructor(
    private reportService: ReportService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    const now = new Date();
    this.startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0];
    this.endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().split('T')[0];
    this.loadReport();
  }

  applyFilters(): void {
    this.loadReport();
  }

  resetFilters(): void {
    const now = new Date();
    this.startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0];
    this.endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().split('T')[0];
    this.loadReport();
  }

  loadReport(): void {
    const user = this.authService.getCurrentUser();
    if (!user) {
      this.error = 'Usuario no autenticado';
      return;
    }

    this.loading = true;
    this.error = null;

    this.reportService.getExpenseIncomeReport(user.id, this.startDate, this.endDate).subscribe({
      next: (data) => {
        this.reportData = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error al cargar el reporte de gastos vs ingresos:', err);
        this.error = 'Error al cargar el reporte de gastos vs ingresos';
        this.loading = false;
      }
    });
  }

  downloadPDF(): void {
    const user = this.authService.getCurrentUser();
    if (!user) {
      this.error = 'Usuario no autenticado';
      return;
    }

    if (!this.reportData) {
      this.error = 'No hay datos para generar el PDF. Por favor, carga el reporte primero.';
      return;
    }

    if (!this.startDate || !this.endDate) {
      this.error = 'Por favor, selecciona un rango de fechas antes de descargar el PDF.';
      return;
    }

    this.downloadingPDF = true;
    this.error = null;

    this.reportService.downloadExpenseIncomeReportPDF(user.id, this.startDate, this.endDate).subscribe({
      next: (blob) => {
        try {
          this.reportService.downloadPDF(blob, 'gastos-ingresos', this.startDate, this.endDate);
          this.downloadingPDF = false;
        } catch (error) {
          console.error('Error al procesar el PDF:', error);
          this.error = 'Error al procesar el archivo PDF';
          this.downloadingPDF = false;
        }
      },
      error: (err) => {
        console.error('Error al descargar el PDF:', err);
        this.error = err.error?.message || 'Error al descargar el PDF. Por favor, intenta nuevamente.';
        this.downloadingPDF = false;
      }
    });
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 2
    }).format(amount);
  }

  formatPercentage(value: number): string {
    return `${value.toFixed(1)}%`;
  }

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('es-AR', {
      year: 'numeric',
      month: 'short'
    }).format(date);
  }

  getTrendIcon(trend: string): string {
    const trendMap: { [key: string]: string } = {
      'increasing': 'trending_up',
      'decreasing': 'trending_down',
      'stable': 'trending_flat'
    };
    return trendMap[trend?.toLowerCase()] || 'trending_flat';
  }

  getTrendClass(trend: string, isExpense: boolean = false): string {
    if (trend?.toLowerCase() === 'stable') return 'trend-neutral';
    
    if (isExpense) {
      // For expenses: increasing is bad (red), decreasing is good (green)
      return trend?.toLowerCase() === 'increasing' ? 'trend-negative' : 'trend-positive';
    } else {
      // For income: increasing is good (green), decreasing is bad (red)
      return trend?.toLowerCase() === 'increasing' ? 'trend-positive' : 'trend-negative';
    }
  }

  getTrendLabel(trend: string): string {
    const labelMap: { [key: string]: string } = {
      'increasing': 'En aumento',
      'decreasing': 'En descenso',
      'stable': 'Estable'
    };
    return labelMap[trend?.toLowerCase()] || trend;
  }

  getCategoryIcon(category: string): string {
    const iconMap: { [key: string]: string } = {
      // Categorías de ingresos
      'salary': 'work',
      'salario': 'work',
      'investment': 'trending_up',
      'inversión': 'trending_up',
      'inversiones': 'trending_up',
      'business': 'business_center',
      'negocio': 'business_center',
      'other_income': 'monetization_on',
      'otros_ingresos': 'monetization_on',
      
      // Categorías de gastos
      'food': 'restaurant',
      'comida': 'restaurant',
      'alimentos': 'restaurant',
      'transport': 'directions_car',
      'transporte': 'directions_car',
      'utilities': 'home',
      'servicios': 'home',
      'entertainment': 'local_activity',
      'entretenimiento': 'local_activity',
      'healthcare': 'local_hospital',
      'salud': 'local_hospital',
      'education': 'school',
      'educación': 'school',
      'educacion': 'school',
      'shopping': 'shopping_cart',
      'compras': 'shopping_cart',
      'other': 'category',
      'otros': 'category',
      
      // Tipos de transacciones
      'installment_payment': 'payments',
      'installment_charge': 'receipt',
      'credit_payment': 'credit_card',
      'credit_charge': 'credit_score',
      'debit_charge': 'payment',
      'transfer': 'swap_horiz',
      'deposit': 'arrow_downward',
      'withdrawal': 'arrow_upward',
      'refund': 'replay',
      'adjustment': 'tune',
      'fee': 'money_off',
      'interest': 'percent',
      'purchase': 'shopping_bag',
      'payment': 'paid'
    };
    return iconMap[category?.toLowerCase()] || 'category';
  }

  getCategoryLabel(category: string): string {
    const labelMap: { [key: string]: string } = {
      // Categorías de ingresos
      'salary': 'Salario',
      'investment': 'Inversiones',
      'business': 'Negocio',
      'other_income': 'Otros Ingresos',
      
      // Categorías de gastos
      'food': 'Comida',
      'transport': 'Transporte',
      'utilities': 'Servicios',
      'entertainment': 'Entretenimiento',
      'healthcare': 'Salud',
      'education': 'Educación',
      'shopping': 'Compras',
      'other': 'Otros',
      
      // Tipos de transacciones que pueden aparecer como categorías
      'account_deposit': 'Depósito en Cuenta',
      'account_withdraw': 'Retiro de Cuenta',
      'wallet_deposit': 'Depósito en Billetera',
      'wallet_withdrawal': 'Retiro de Billetera',
      'debit_purchase': 'Compra con Débito',
      'installment_payment': 'Pago de Cuota',
      'installment_charge': 'Cargo de Cuota',
      'credit_payment': 'Pago de Crédito',
      'credit_charge': 'Cargo de Crédito',
      'debit_charge': 'Cargo de Débito',
      'transfer': 'Transferencia',
      'deposit': 'Depósito',
      'withdrawal': 'Retiro',
      'refund': 'Reembolso',
      'adjustment': 'Ajuste',
      'fee': 'Comisión',
      'interest': 'Interés',
      'purchase': 'Compra',
      'payment': 'Pago'
    };
    
    const key = category?.toLowerCase();
    if (labelMap[key]) {
      return labelMap[key];
    }
    
    // Si no está en el mapa, intentar convertir de snake_case a título con capitalización
    return category
      ?.split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ') || category;
  }

  getTypeClass(type: string): string {
    return type?.toLowerCase() === 'income' ? 'type-income' : 'type-expense';
  }

  getBarWidth(amount: number, total: number): number {
    if (total === 0) return 0;
    return Math.min((amount / total) * 100, 100);
  }

  getCategoryColor(index: number): string {
    const colors = [
      '#4CAF50', '#2196F3', '#FF9800', '#F44336', '#9C27B0',
      '#00BCD4', '#FFEB3B', '#795548', '#607D8B', '#E91E63',
      '#3F51B5', '#8BC34A', '#FFC107', '#FF5722', '#673AB7'
    ];
    return colors[index % colors.length];
  }

  getPieChartData(): Array<{category: string, amount: number, percentage: number, color: string}> {
    if (!this.reportData?.by_category) return [];
    
    return this.reportData.by_category.map((cat, index) => ({
      category: this.getCategoryLabel(cat.category),
      amount: cat.amount,
      percentage: cat.percentage,
      color: this.getCategoryColor(index)
    }));
  }

  getPieOffset(index: number): number {
    if (!this.reportData?.by_category) return 70.675;
    
    let offset = 70.675;
    const data = this.getPieChartData();
    for (let i = 0; i < index; i++) {
      offset -= data[i].percentage * 2.827;
    }
    return offset;
  }
}
