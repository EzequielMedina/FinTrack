import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
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

    // Get current month date range
    const now = new Date();
    this.startDate = new Date(now.getFullYear(), now.getMonth(), 1).toISOString().split('T')[0];
    this.endDate = new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().split('T')[0];

    this.reportService.getExpenseIncomeReport(user.id, this.startDate, this.endDate).subscribe({
      next: (data) => {
        this.reportData = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading expense-income report:', err);
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

    this.downloadingPDF = true;

    this.reportService.downloadExpenseIncomeReportPDF(user.id, this.startDate, this.endDate).subscribe({
      next: (blob) => {
        this.reportService.downloadPDF(blob, 'gastos-ingresos');
        this.downloadingPDF = false;
      },
      error: (err) => {
        console.error('Error downloading PDF:', err);
        this.error = 'Error al descargar el PDF';
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
      'salary': 'work',
      'investment': 'trending_up',
      'business': 'business_center',
      'other_income': 'monetization_on',
      'food': 'restaurant',
      'transport': 'directions_car',
      'utilities': 'home',
      'entertainment': 'local_activity',
      'healthcare': 'local_hospital',
      'education': 'school',
      'shopping': 'shopping_cart',
      'other': 'category'
    };
    return iconMap[category?.toLowerCase()] || 'category';
  }

  getTypeClass(type: string): string {
    return type?.toLowerCase() === 'income' ? 'type-income' : 'type-expense';
  }

  getBarWidth(amount: number, total: number): number {
    if (total === 0) return 0;
    return Math.min((amount / total) * 100, 100);
  }
}
