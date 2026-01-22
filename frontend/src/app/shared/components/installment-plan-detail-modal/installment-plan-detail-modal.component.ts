import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatDividerModule } from '@angular/material/divider';
import { MatChipsModule } from '@angular/material/chips';
import { InstallmentPlan } from '../../../models';

@Component({
  selector: 'app-installment-plan-detail-modal',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatDividerModule,
    MatChipsModule
  ],
  template: `
    <div class="modal-header">
      <h2 mat-dialog-title>
        <mat-icon>receipt_long</mat-icon>
        {{ plan.description || plan.merchantName || 'Plan de Cuotas' }}
      </h2>
      <button mat-icon-button (click)="close()" class="close-button">
        <mat-icon>close</mat-icon>
      </button>
    </div>

    <mat-dialog-content class="modal-content">
      <!-- Summary Cards -->
      <div class="summary-cards">
        <div class="summary-card primary">
          <div class="card-icon">
            <mat-icon>account_balance_wallet</mat-icon>
          </div>
          <div class="card-content">
            <div class="card-label">Total del Plan</div>
            <div class="card-value">{{ formatCurrency(plan.totalAmount) }}</div>
          </div>
        </div>

        <div class="summary-card accent">
          <div class="card-icon">
            <mat-icon>payment</mat-icon>
          </div>
          <div class="card-content">
            <div class="card-label">Cuota Mensual</div>
            <div class="card-value">{{ formatCurrency(plan.installmentAmount) }}</div>
          </div>
        </div>

        <div class="summary-card">
          <div class="card-icon">
            <mat-icon>trending_up</mat-icon>
          </div>
          <div class="card-content">
            <div class="card-label">Progreso</div>
            <div class="card-value">{{ plan.paidInstallments }}/{{ plan.installmentsCount }}</div>
            <div class="card-percentage">{{ getProgressPercentage(plan) }}%</div>
          </div>
        </div>

        <div class="summary-card" [class.warn]="plan.remainingAmount > 0">
          <div class="card-icon">
            <mat-icon>schedule</mat-icon>
          </div>
          <div class="card-content">
            <div class="card-label">Restante</div>
            <div class="card-value">{{ formatCurrency(plan.remainingAmount) }}</div>
          </div>
        </div>
      </div>

      <!-- Progress Bar -->
      <div class="progress-section">
        <div class="progress-header">
          <span>Progreso del Plan</span>
          <span class="progress-percentage">{{ getProgressPercentage(plan) }}%</span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill" [style.width.%]="getProgressPercentage(plan)"></div>
        </div>
      </div>

      <!-- Key Information -->
      <div class="info-section">
        <div class="info-row">
          <span class="info-label">Estado:</span>
          <mat-chip [class]="getStatusClass(plan.status)">
            <mat-icon>{{ getStatusIcon(plan.status) }}</mat-icon>
            {{ getStatusText(plan.status) }}
          </mat-chip>
        </div>
        <div class="info-row">
          <span class="info-label">Fecha de Inicio:</span>
          <span class="info-value">{{ formatDate(plan.startDate) }}</span>
        </div>
        <div class="info-row" *ngIf="plan.merchantName">
          <span class="info-label">Comercio:</span>
          <span class="info-value">{{ plan.merchantName }}</span>
        </div>
      </div>
    </mat-dialog-content>

    <mat-dialog-actions class="modal-actions">
      <button mat-button (click)="close()">
        Cerrar
      </button>
    </mat-dialog-actions>
  `,
  styles: [`
    .modal-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px 24px;
      border-bottom: 1px solid var(--border-light);
    }

    .modal-header h2 {
      display: flex;
      align-items: center;
      gap: 12px;
      margin: 0;
      font-size: 20px;
      font-weight: 600;
      color: var(--text-primary);
    }

    .close-button {
      color: var(--text-secondary);
    }

    .modal-content {
      padding: 24px !important;
      max-height: 70vh;
      overflow-y: auto;
    }

    .summary-cards {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 16px;
      margin-bottom: 24px;
    }

    .summary-card {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px;
      border-radius: var(--radius-lg);
      background: var(--bg-secondary);
      border: 1px solid var(--border-light);
      transition: all var(--transition-base);
    }

    .summary-card:hover {
      box-shadow: var(--shadow-md);
      transform: translateY(-2px);
    }

    .summary-card.primary {
      border-left: 4px solid var(--accent-500);
    }

    .summary-card.accent {
      border-left: 4px solid var(--success-500);
    }

    .summary-card.warn {
      border-left: 4px solid var(--warning-500);
    }

    .card-icon {
      width: 48px;
      height: 48px;
      border-radius: var(--radius-md);
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--gray-100);
      flex-shrink: 0;
    }

    .summary-card.primary .card-icon {
      background: var(--accent-100);
      color: var(--accent-700);
    }

    .summary-card.accent .card-icon {
      background: var(--success-100);
      color: var(--success-700);
    }

    .summary-card.warn .card-icon {
      background: var(--warning-100);
      color: var(--warning-700);
    }

    .card-icon mat-icon {
      font-size: 24px;
      width: 24px;
      height: 24px;
    }

    .card-content {
      flex: 1;
    }

    .card-label {
      font-size: var(--text-xs);
      color: var(--text-secondary);
      text-transform: uppercase;
      letter-spacing: 0.5px;
      margin-bottom: 4px;
    }

    .card-value {
      font-size: var(--text-lg);
      font-weight: var(--font-bold);
      color: var(--text-primary);
    }

    .card-percentage {
      font-size: var(--text-sm);
      color: var(--text-secondary);
      margin-top: 4px;
    }

    .progress-section {
      margin-bottom: 24px;
      padding: 20px;
      background: var(--bg-secondary);
      border-radius: var(--radius-lg);
      border: 1px solid var(--border-light);
    }

    .progress-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
      font-size: var(--text-sm);
      font-weight: var(--font-medium);
      color: var(--text-primary);
    }

    .progress-percentage {
      font-weight: var(--font-bold);
      color: var(--accent-600);
    }

    .progress-bar {
      width: 100%;
      height: 12px;
      background: var(--gray-200);
      border-radius: 6px;
      overflow: hidden;
    }

    .progress-fill {
      height: 100%;
      background: linear-gradient(90deg, var(--success-500), var(--accent-500));
      transition: width 0.3s ease;
    }

    .info-section {
      display: flex;
      flex-direction: column;
      gap: 16px;
      padding: 20px;
      background: var(--bg-secondary);
      border-radius: var(--radius-lg);
      border: 1px solid var(--border-light);
    }

    .info-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid var(--border-light);
    }

    .info-row:last-child {
      border-bottom: none;
    }

    .info-label {
      font-size: var(--text-sm);
      color: var(--text-secondary);
      font-weight: var(--font-medium);
    }

    .info-value {
      font-size: var(--text-base);
      color: var(--text-primary);
      font-weight: var(--font-medium);
    }

    .modal-actions {
      padding: 16px 24px;
      border-top: 1px solid var(--border-light);
      display: flex;
      justify-content: flex-end;
    }

    mat-chip {
      font-size: var(--text-xs);
      height: 28px;
    }

    mat-chip.status-active {
      background: var(--success-100);
      color: var(--success-700);
    }

    mat-chip.status-completed {
      background: var(--accent-100);
      color: var(--accent-700);
    }

    mat-chip.status-cancelled {
      background: var(--error-100);
      color: var(--error-700);
    }

    mat-chip.status-overdue {
      background: var(--warning-100);
      color: var(--warning-700);
    }

    @media (max-width: 600px) {
      .summary-cards {
        grid-template-columns: 1fr;
      }

      .info-row {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
      }
    }
  `]
})
export class InstallmentPlanDetailModalComponent {
  plan: InstallmentPlan;

  constructor(
    public dialogRef: MatDialogRef<InstallmentPlanDetailModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { plan: InstallmentPlan }
  ) {
    this.plan = data.plan;
  }

  close(): void {
    this.dialogRef.close();
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS'
    }).format(amount);
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('es-AR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  getProgressPercentage(plan: InstallmentPlan): number {
    return Math.round((plan.paidInstallments / plan.installmentsCount) * 100);
  }

  getStatusText(status: string): string {
    switch (status) {
      case 'active': return 'Activo';
      case 'completed': return 'Completado';
      case 'cancelled': return 'Cancelado';
      case 'overdue': return 'Vencido';
      default: return status;
    }
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'active': return 'status-active';
      case 'completed': return 'status-completed';
      case 'cancelled': return 'status-cancelled';
      case 'overdue': return 'status-overdue';
      default: return '';
    }
  }

  getStatusIcon(status: string): string {
    switch (status) {
      case 'active': return 'play_circle';
      case 'completed': return 'check_circle';
      case 'cancelled': return 'cancel';
      case 'overdue': return 'warning';
      default: return 'help';
    }
  }
}