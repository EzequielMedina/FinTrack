import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ReportService, InstallmentReport } from '../../../services/report.service';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-installment-report',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatButtonModule
  ],
  templateUrl: './installment-report.component.html',
  styleUrls: ['./installment-report.component.css']
})
export class InstallmentReportComponent implements OnInit {
  loading = false;
  error: string | null = null;
  reportData: InstallmentReport | null = null;
  downloadingPDF = false;

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

    this.reportService.getInstallmentReport(user.id).subscribe({
      next: (data) => {
        this.reportData = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error al cargar el reporte de cuotas:', err);
        this.error = 'Error al cargar el reporte de cuotas';
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

    this.downloadingPDF = true;
    this.error = null;

    this.reportService.downloadInstallmentReportPDF(user.id).subscribe({
      next: (blob) => {
        try {
          this.reportService.downloadPDF(blob, 'cuotas');
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

  formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('es-AR', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }).format(date);
  }

  formatPercentage(value: number): string {
    return `${value.toFixed(1)}%`;
  }

  getDaysUntilDueClass(days: number): string {
    if (days <= 3) return 'urgent';
    if (days <= 7) return 'warning';
    return 'normal';
  }

  getDaysUntilDueIcon(days: number): string {
    if (days <= 3) return 'priority_high';
    if (days <= 7) return 'warning';
    return 'schedule';
  }

  getStatusLabel(status: string): string {
    const labels: { [key: string]: string } = {
      'active': 'Activo',
      'completed': 'Completado',
      'pending': 'Pendiente',
      'paid': 'Pagada',
      'overdue': 'Vencido',
      'cancelled': 'Cancelado'
    };
    return labels[status?.toLowerCase()] || status;
  }

  getStatusClass(status: string): string {
    const classes: { [key: string]: string } = {
      'active': 'status-active',
      'completed': 'status-completed',
      'pending': 'status-pending',
      'paid': 'status-paid',
      'overdue': 'status-overdue',
      'cancelled': 'status-cancelled'
    };
    return classes[status?.toLowerCase()] || '';
  }
}
