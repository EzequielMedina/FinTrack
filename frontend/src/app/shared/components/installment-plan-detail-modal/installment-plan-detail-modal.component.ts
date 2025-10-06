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
        Detalle del Plan de Cuotas
      </h2>
      <button mat-icon-button (click)="close()" class="close-button">
        <mat-icon>close</mat-icon>
      </button>
    </div>

    <mat-dialog-content class="modal-content">
      <!-- Plan Information -->
      <mat-card class="info-card">
        <mat-card-header>
          <mat-card-title>
            <mat-icon>info</mat-icon>
            Información General
          </mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">Descripción:</span>
              <span class="value">{{ plan.description || 'Compra en cuotas' }}</span>
            </div>
            <div class="info-item">
              <span class="label">Comercio:</span>
              <span class="value">{{ plan.merchantName || 'N/A' }}</span>
            </div>
            <div class="info-item">
              <span class="label">ID del Plan:</span>
              <span class="value code">{{ plan.id.substring(0, 8) }}...</span>
            </div>
            <div class="info-item">
              <span class="label">Estado:</span>
              <mat-chip [class]="getStatusClass(plan.status)">
                <mat-icon>{{ getStatusIcon(plan.status) }}</mat-icon>
                {{ getStatusText(plan.status) }}
              </mat-chip>
            </div>
          </div>
        </mat-card-content>
      </mat-card>

      <mat-divider></mat-divider>

      <!-- Financial Information -->
      <mat-card class="info-card">
        <mat-card-header>
          <mat-card-title>
            <mat-icon>attach_money</mat-icon>
            Información Financiera
          </mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="financial-grid">
            <div class="financial-item primary">
              <span class="label">Monto Total:</span>
              <span class="value amount">{{ formatCurrency(plan.totalAmount) }}</span>
            </div>
            <div class="financial-item">
              <span class="label">Cantidad de Cuotas:</span>
              <span class="value">{{ plan.installmentsCount }}</span>
            </div>
            <div class="financial-item accent">
              <span class="label">Monto por Cuota:</span>
              <span class="value amount">{{ formatCurrency(plan.installmentAmount) }}</span>
            </div>
            <div class="financial-item">
              <span class="label">Tasa de Interés:</span>
              <span class="value">{{ plan.interestRate }}%</span>
            </div>
            <div class="financial-item">
              <span class="label">Interés Total:</span>
              <span class="value amount">{{ formatCurrency(plan.totalInterest) }}</span>
            </div>
            <div class="financial-item">
              <span class="label">Comisión Admin:</span>
              <span class="value amount">{{ formatCurrency(plan.adminFee) }}</span>
            </div>
          </div>
        </mat-card-content>
      </mat-card>

      <mat-divider></mat-divider>

      <!-- Progress Information -->
      <mat-card class="info-card">
        <mat-card-header>
          <mat-card-title>
            <mat-icon>trending_up</mat-icon>
            Progreso del Plan
          </mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="progress-info">
            <div class="progress-stats">
              <div class="stat-item">
                <span class="stat-label">Cuotas Pagadas:</span>
                <span class="stat-value">{{ plan.paidInstallments }} de {{ plan.installmentsCount }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Progreso:</span>
                <span class="stat-value">{{ getProgressPercentage(plan) }}%</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Monto Restante:</span>
                <span class="stat-value amount">{{ formatCurrency(plan.remainingAmount) }}</span>
              </div>
            </div>
            
            <!-- Progress Bar -->
            <div class="progress-bar-container">
              <div class="progress-bar">
                <div 
                  class="progress-fill" 
                  [style.width.%]="getProgressPercentage(plan)">
                </div>
              </div>
              <span class="progress-text">{{ getProgressPercentage(plan) }}% completado</span>
            </div>
          </div>
        </mat-card-content>
      </mat-card>

      <mat-divider></mat-divider>

      <!-- Dates Information -->
      <mat-card class="info-card">
        <mat-card-header>
          <mat-card-title>
            <mat-icon>event</mat-icon>
            Fechas Importantes
          </mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="dates-grid">
            <div class="date-item">
              <span class="label">Fecha de Inicio:</span>
              <span class="value">{{ formatDate(plan.startDate) }}</span>
            </div>
            <div class="date-item">
              <span class="label">Creado:</span>
              <span class="value">{{ formatDate(plan.createdAt) }}</span>
            </div>
            <div class="date-item">
              <span class="label">Última Actualización:</span>
              <span class="value">{{ formatDate(plan.updatedAt) }}</span>
            </div>
            <div class="date-item" *ngIf="plan.completedAt">
              <span class="label">Completado:</span>
              <span class="value">{{ formatDate(plan.completedAt) }}</span>
            </div>
            <div class="date-item" *ngIf="plan.cancelledAt">
              <span class="label">Cancelado:</span>
              <span class="value">{{ formatDate(plan.cancelledAt) }}</span>
            </div>
          </div>
        </mat-card-content>
      </mat-card>
    </mat-dialog-content>

    <mat-dialog-actions class="modal-actions">
      <button mat-button (click)="close()">
        <mat-icon>close</mat-icon>
        Cerrar
      </button>
    </mat-dialog-actions>
  `,
  styles: [`
    .modal-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px 24px 0;
    }

    .modal-header h2 {
      display: flex;
      align-items: center;
      gap: 8px;
      margin: 0;
      color: #1976d2;
    }

    .close-button {
      color: #666;
    }

    .modal-content {
      padding: 20px 24px !important;
      max-height: 70vh;
      overflow-y: auto;
    }

    .info-card {
      margin-bottom: 16px;
    }

    .info-card mat-card-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 16px;
      color: #333;
    }

    .info-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 16px;
      margin-top: 16px;
    }

    .financial-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 16px;
      margin-top: 16px;
    }

    .dates-grid {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 16px;
      margin-top: 16px;
    }

    .info-item, .financial-item, .date-item {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .financial-item.primary {
      grid-column: 1 / -1;
      background: #f5f5f5;
      padding: 12px;
      border-radius: 8px;
    }

    .financial-item.accent {
      background: #e3f2fd;
      padding: 12px;
      border-radius: 8px;
    }

    .label {
      font-size: 12px;
      color: #666;
      font-weight: 500;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .value {
      font-size: 14px;
      color: #333;
      font-weight: 400;
    }

    .value.amount {
      font-weight: 600;
      color: #1976d2;
    }

    .value.code {
      font-family: 'Courier New', monospace;
      background: #f0f0f0;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 12px;
    }

    .progress-info {
      margin-top: 16px;
    }

    .progress-stats {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;
      margin-bottom: 16px;
    }

    .stat-item {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .stat-label {
      font-size: 12px;
      color: #666;
      font-weight: 500;
    }

    .stat-value {
      font-size: 16px;
      font-weight: 600;
      color: #333;
    }

    .stat-value.amount {
      color: #1976d2;
    }

    .progress-bar-container {
      display: flex;
      flex-direction: column;
      gap: 8px;
    }

    .progress-bar {
      width: 100%;
      height: 8px;
      background: #e0e0e0;
      border-radius: 4px;
      overflow: hidden;
    }

    .progress-fill {
      height: 100%;
      background: linear-gradient(90deg, #4caf50, #2196f3);
      transition: width 0.3s ease;
    }

    .progress-text {
      font-size: 12px;
      color: #666;
      text-align: center;
    }

    .modal-actions {
      padding: 16px 24px;
      border-top: 1px solid #e0e0e0;
    }

    mat-chip {
      font-size: 12px;
    }

    mat-chip.status-active {
      background: #e8f5e8;
      color: #2e7d32;
    }

    mat-chip.status-completed {
      background: #e3f2fd;
      color: #1976d2;
    }

    mat-chip.status-cancelled {
      background: #ffebee;
      color: #d32f2f;
    }

    mat-chip.status-overdue {
      background: #fff3e0;
      color: #f57c00;
    }

    @media (max-width: 600px) {
      .info-grid,
      .financial-grid,
      .dates-grid {
        grid-template-columns: 1fr;
      }

      .progress-stats {
        grid-template-columns: 1fr;
      }

      .financial-item.primary {
        grid-column: 1;
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