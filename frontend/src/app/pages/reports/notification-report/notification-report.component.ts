import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ReportService, NotificationReport } from '../../../services/report.service';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-notification-report',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatButtonModule
  ],
  templateUrl: './notification-report.component.html',
  styleUrls: ['./notification-report.component.css']
})
export class NotificationReportComponent implements OnInit {
  loading = false;
  error: string | null = null;
  reportData: NotificationReport | null = null;
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

    // Check if user is admin
    if (user.role !== 'admin') {
      this.error = 'No tiene permisos para ver este reporte. Solo administradores.';
      return;
    }

    this.loading = true;
    this.error = null;

    // Get last 30 days by default
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(startDate.getDate() - 30);

    const startDateStr = startDate.toISOString().split('T')[0];
    const endDateStr = endDate.toISOString().split('T')[0];

    this.reportService.getNotificationReport(startDateStr, endDateStr).subscribe({
      next: (data) => {
        this.reportData = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error al cargar el reporte de notificaciones:', err);
        this.error = 'Error al cargar el reporte de notificaciones';
        this.loading = false;
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
      month: 'short',
      day: 'numeric'
    }).format(date);
  }

  formatDateTime(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('es-AR', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date);
  }

  getStatusClass(status: string): string {
    const statusMap: { [key: string]: string } = {
      'completed': 'status-completed',
      'success': 'status-success',
      'failed': 'status-failed',
      'error': 'status-error',
      'running': 'status-running'
    };
    return statusMap[status?.toLowerCase()] || 'status-default';
  }

  getStatusLabel(status: string): string {
    const labelMap: { [key: string]: string } = {
      'completed': 'Completado',
      'success': 'Éxito',
      'failed': 'Fallido',
      'error': 'Error',
      'running': 'En ejecución'
    };
    return labelMap[status?.toLowerCase()] || status;
  }

  getStatusIcon(status: string): string {
    const iconMap: { [key: string]: string } = {
      'completed': 'check_circle',
      'success': 'check_circle',
      'failed': 'error',
      'error': 'error',
      'running': 'pending'
    };
    return iconMap[status?.toLowerCase()] || 'help';
  }

  downloadPDF(): void {
    const user = this.authService.getCurrentUser();
    if (!user || user.role !== 'admin') {
      this.error = 'No tiene permisos para descargar este reporte. Solo administradores.';
      return;
    }

    if (!this.reportData) {
      this.error = 'No hay datos para generar el PDF. Por favor, carga el reporte primero.';
      return;
    }

    this.downloadingPDF = true;
    this.error = null;

    // Obtener fechas del período del reporte
    const startDate = this.reportData.period?.start_date;
    const endDate = this.reportData.period?.end_date;

    this.reportService.downloadNotificationReportPDF(startDate, endDate).subscribe({
      next: (blob) => {
        try {
          this.reportService.downloadPDF(blob, 'notificaciones', startDate, endDate);
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
}
