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
    { value: 'account_withdraw', label: 'Retiro de Cuenta' },
    { value: 'installment_payment', label: 'Pago de Cuota' },
    { value: 'installment_charge', label: 'Cargo de Cuota' },
    { value: 'transfer', label: 'Transferencia' },
    { value: 'refund', label: 'Reembolso' },
    { value: 'adjustment', label: 'Ajuste' }
  ];

  constructor(
    private reportService: ReportService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    this.initializeDates();
    this.loadReport();
  }

  initializeDates(): void {
    const today = new Date();
    const firstDayOfMonth = new Date(today.getFullYear(), today.getMonth(), 1);
    
    this.endDate = today.toISOString().split('T')[0];
    this.startDate = firstDayOfMonth.toISOString().split('T')[0];
  }

  loadReport(): void {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser?.id) {
      this.error = 'Usuario no autenticado. Inicia sesión para ver el reporte.';
      return;
    }
    this.userId = currentUser.id;

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
    if (typeObj) {
      return typeObj.label;
    }
    
    // Mapeo directo de tipos comunes al español
    const directTranslations: { [key: string]: string } = {
      'account_deposit': 'Depósito en Cuenta',
      'account_withdraw': 'Retiro de Cuenta',
      'wallet_deposit': 'Depósito en Billetera',
      'wallet_withdrawal': 'Retiro de Billetera',
      'debit_purchase': 'Compra con Débito',
      'credit_charge': 'Cargo de Crédito',
      'credit_payment': 'Pago de Crédito',
      'installment_payment': 'Pago de Cuota',
      'installment_charge': 'Cargo de Cuota',
      'transfer': 'Transferencia',
      'refund': 'Reembolso',
      'adjustment': 'Ajuste',
      'fee': 'Comisión',
      'interest': 'Interés',
      'purchase': 'Compra',
      'payment': 'Pago',
      'deposit': 'Depósito',
      'withdrawal': 'Retiro'
    };
    
    // Intentar traducción directa primero
    if (directTranslations[type.toLowerCase()]) {
      return directTranslations[type.toLowerCase()];
    }
    
    // Si no encuentra, formatear de forma legible
    return type
      .split('_')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
  }

  getTypeColor(type: string): string {
    const colorMap: { [key: string]: string } = {
      'wallet_deposit': '#4CAF50',
      'account_deposit': '#4CAF50',
      'credit_payment': '#8BC34A',
      'refund': '#8BC34A',
      'wallet_withdrawal': '#FF9800',
      'account_withdraw': '#FF9800',
      'credit_charge': '#F44336',
      'debit_purchase': '#E91E63',
      'installment_payment': '#9C27B0',
      'installment_charge': '#673AB7',
      'transfer': '#2196F3',
      'adjustment': '#607D8B'
    };
    return colorMap[type] || '#667eea';
  }

  getTypeIcon(type: string): string {
    const iconMap: { [key: string]: string } = {
      'wallet_deposit': '💰',
      'account_deposit': '💵',
      'credit_payment': '✅',
      'refund': '↩️',
      'wallet_withdrawal': '💸',
      'account_withdraw': '🏧',
      'credit_charge': '💳',
      'debit_purchase': '🛒',
      'installment_payment': '📅',
      'installment_charge': '📊',
      'transfer': '🔄',
      'adjustment': '⚙️'
    };
    return iconMap[type] || '📋';
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
