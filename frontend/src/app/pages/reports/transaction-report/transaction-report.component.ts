import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { ReportService, TransactionReport } from '../../../services/report.service';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-transaction-report',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule, MatButtonModule, MatIconModule],
  templateUrl: './transaction-report.component.html',
  styleUrls: ['./transaction-report.component.css']
})
export class TransactionReportComponent implements OnInit {
  userId: string = '';
  report: TransactionReport | null = null;
  isLoading: boolean = false;
  error: string = '';
  downloadingPDF: boolean = false;

  // Filtros
  startDate: string = '';
  endDate: string = '';
  selectedType: string = '';

  // Tipos de transacciones
  transactionTypes = [
    { value: '', label: 'Todas' },
    { value: 'wallet_deposit', label: 'Depósito en Billetera' },
    { value: 'wallet_withdrawal', label: 'Retiro de Billetera' },
    { value: 'credit_charge', label: 'Cargo en Crédito' },
    { value: 'credit_payment', label: 'Pago de Crédito' },
    { value: 'debit_purchase', label: 'Compra con Débito' },
    { value: 'account_deposit', label: 'Depósito en Cuenta' },
    { value: 'account_withdraw', label: 'Retiro de Cuenta' }
  ];

  constructor(
    private reportService: ReportService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    const currentUser = this.authService.getCurrentUser();
    if (currentUser) {
      this.userId = currentUser.id;
      this.initializeDates();
      this.loadReport();
    }
  }

  initializeDates(): void {
    const today = new Date();
    const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    
    this.endDate = today.toISOString().split('T')[0];
    this.startDate = firstDayOfMonth.toISOString().split('T')[0];
  }

  loadReport(): void {
    if (!this.userId) return;

    this.isLoading = true;
    this.error = '';

    this.reportService.getTransactionReport(
      this.userId,
      this.startDate,
      this.endDate,
      this.selectedType
    ).subscribe({
      next: (data) => {
        this.report = data;
        this.isLoading = false;
      },
      error: (err) => {
        this.error = 'Error al cargar el reporte: ' + (err.message || 'Error desconocido');
        this.isLoading = false;
        console.error('Error al cargar el reporte:', err);
      }
    });
  }

  applyFilters(): void {
    this.loadReport();
  }

  resetFilters(): void {
    this.initializeDates();
    this.selectedType = '';
    this.loadReport();
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS'
    }).format(amount);
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    return date.toLocaleDateString('es-AR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  getTypeLabel(type: string): string {
    const typeObj = this.transactionTypes.find(t => t.value === type);
    return typeObj ? typeObj.label : type;
  }

  downloadPDF(): void {
    if (!this.userId) {
      this.error = 'Usuario no identificado';
      return;
    }

    if (!this.report) {
      this.error = 'No hay datos para generar el PDF. Por favor, carga el reporte primero.';
      return;
    }

    this.downloadingPDF = true;
    this.error = '';

    this.reportService.downloadTransactionReportPDF(
      this.userId,
      this.startDate,
      this.endDate,
      this.selectedType
    ).subscribe({
      next: (blob) => {
        try {
          this.reportService.downloadPDF(blob, 'transacciones', this.startDate, this.endDate);
          this.downloadingPDF = false;
          // Opcional: mostrar mensaje de éxito
          // this.snackBar?.open('PDF descargado exitosamente', 'Cerrar', { duration: 3000 });
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

  exportReport(): void {
    if (!this.report) return;

    const csvContent = this.generateCSV();
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    
    // Normalizar el nombre del archivo para evitar problemas con acentos
    const date = new Date().toISOString().split('T')[0];
    const start = this.startDate ? this.startDate.split('T')[0].replace(/-/g, '') : '';
    const end = this.endDate ? this.endDate.split('T')[0].replace(/-/g, '') : '';
    const filename = start && end 
      ? `reporte-transacciones-${start}-${end}.csv`
      : `reporte-transacciones-${date}.csv`;
    
    link.setAttribute('href', url);
    link.setAttribute('download', filename);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }

  private generateCSV(): string {
    if (!this.report) return '';

    const header = 'Fecha,Descripción,Tipo,Monto,Comercio\n';
    const rows = this.report.top_expenses.map(item => {
      return `${this.formatDate(item.date)},"${item.description}",${this.getTypeLabel(item.type)},${item.amount},"${item.merchant_name || 'N/A'}"`;
    }).join('\n');

    return header + rows;
  }
}
