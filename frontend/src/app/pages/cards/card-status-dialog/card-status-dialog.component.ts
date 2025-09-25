import { Component, Inject, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { Card, CardStatus } from '../../../models';

interface DialogData {
  card: Card;
  action: 'activate' | 'deactivate';
}

@Component({
  selector: 'app-card-status-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule
  ],
  template: `
    <div class="status-dialog-container">
      <div class="dialog-header">
        <mat-icon [class]="getIconClass()">{{ getIcon() }}</mat-icon>
        <h2 mat-dialog-title>{{ getTitle() }}</h2>
      </div>

      <div mat-dialog-content class="dialog-content">
        <p class="card-info">
          <strong>Tarjeta:</strong> {{ data.card.nickname || 'Tarjeta ' + getCardTypeLabel() }}
        </p>
        <p class="card-number">
          <strong>Número:</strong> {{ data.card.maskedNumber }}
        </p>
        
        <div class="confirmation-message">
          <p>{{ getConfirmationMessage() }}</p>
          @if (data.action === 'deactivate') {
            <div class="warning-notice">
              <mat-icon class="warning-icon">warning</mat-icon>
              <small>Una vez desactivada, no podrás usar esta tarjeta hasta que la vuelvas a activar.</small>
            </div>
          }
        </div>
      </div>

      <div mat-dialog-actions class="dialog-actions">
        <button 
          mat-button 
          (click)="onCancel()"
          class="cancel-button">
          Cancelar
        </button>
        <button 
          mat-raised-button 
          [color]="getButtonColor()"
          (click)="onConfirm()"
          class="confirm-button">
          <mat-icon>{{ getConfirmIcon() }}</mat-icon>
          {{ getConfirmButtonText() }}
        </button>
      </div>
    </div>
  `,
  styles: [`
    .status-dialog-container {
      min-width: 400px;
      max-width: 500px;
    }

    .dialog-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 16px;
    }

    .dialog-header mat-icon {
      font-size: 32px;
      width: 32px;
      height: 32px;
    }

    .activate-icon {
      color: #4caf50;
    }

    .deactivate-icon {
      color: #f44336;
    }

    .dialog-content {
      margin: 16px 0;
    }

    .card-info, .card-number {
      margin: 8px 0;
      color: #666;
    }

    .confirmation-message {
      margin: 20px 0;
      padding: 16px;
      background-color: #f5f5f5;
      border-radius: 8px;
    }

    .confirmation-message p {
      margin: 0 0 12px 0;
      font-weight: 500;
    }

    .warning-notice {
      display: flex;
      align-items: flex-start;
      gap: 8px;
      padding: 12px;
      background-color: #fff3cd;
      border: 1px solid #ffeeba;
      border-radius: 6px;
      margin-top: 12px;
    }

    .warning-icon {
      color: #856404;
      font-size: 18px;
      width: 18px;
      height: 18px;
      margin-top: 1px;
    }

    .warning-notice small {
      color: #856404;
      line-height: 1.4;
    }

    .dialog-actions {
      display: flex;
      justify-content: flex-end;
      gap: 12px;
      margin-top: 24px;
    }

    .confirm-button {
      min-width: 120px;
    }

    .confirm-button mat-icon {
      margin-right: 8px;
      font-size: 18px;
      width: 18px;
      height: 18px;
    }
  `]
})
export class CardStatusDialogComponent {
  private readonly dialogRef = inject(MatDialogRef<CardStatusDialogComponent>);

  constructor(@Inject(MAT_DIALOG_DATA) public data: DialogData) {}

  getTitle(): string {
    return this.data.action === 'activate' ? 'Activar Tarjeta' : 'Desactivar Tarjeta';
  }

  getIcon(): string {
    return this.data.action === 'activate' ? 'check_circle' : 'cancel';
  }

  getIconClass(): string {
    return this.data.action === 'activate' ? 'activate-icon' : 'deactivate-icon';
  }

  getCardTypeLabel(): string {
    return this.data.card.cardType === 'credit' ? 'de Crédito' : 'de Débito';
  }

  getConfirmationMessage(): string {
    const action = this.data.action === 'activate' ? 'activar' : 'desactivar';
    return `¿Estás seguro de que deseas ${action} esta tarjeta?`;
  }

  getButtonColor(): string {
    return this.data.action === 'activate' ? 'primary' : 'warn';
  }

  getConfirmIcon(): string {
    return this.data.action === 'activate' ? 'check' : 'block';
  }

  getConfirmButtonText(): string {
    return this.data.action === 'activate' ? 'Activar' : 'Desactivar';
  }

  onConfirm(): void {
    this.dialogRef.close(true);
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}