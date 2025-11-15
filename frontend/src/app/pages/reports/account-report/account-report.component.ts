import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ReportService, AccountReport } from '../../../services/report.service';
import { AuthService } from '../../../services/auth.service';

@Component({
  selector: 'app-account-report',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatProgressSpinnerModule,
    MatIconModule,
    MatButtonModule
  ],
  templateUrl: './account-report.component.html',
  styleUrls: ['./account-report.component.css']
})
export class AccountReportComponent implements OnInit {
  loading = false;
  error: string | null = null;
  reportData: AccountReport | null = null;
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

    this.reportService.getAccountReport(user.id).subscribe({
      next: (data) => {
        this.reportData = data;
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading account report:', err);
        this.error = 'Error al cargar el reporte de cuentas';
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

    this.reportService.downloadAccountReportPDF(user.id).subscribe({
      next: (blob) => {
        this.reportService.downloadPDF(blob, 'cuentas');
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
      currency: 'ARS'
    }).format(amount);
  }

  formatPercentage(value: number): string {
    return `${value.toFixed(1)}%`;
  }

  getAccountTypeLabel(type: string): string {
    const labels: { [key: string]: string } = {
      'wallet': 'Billetera Virtual',
      'savings': 'Caja de Ahorro',
      'checking': 'Cuenta Corriente',
      'credit': 'Tarjeta de Crédito',
      'debit': 'Tarjeta de Débito'
    };
    return labels[type] || type;
  }

  getTypeIcon(type: string): string {
    const icons: { [key: string]: string } = {
      'wallet': 'account_balance_wallet',
      'savings': 'savings',
      'checking': 'account_balance',
      'credit': 'credit_card',
      'debit': 'payment'
    };
    return icons[type] || 'account_balance';
  }
}
