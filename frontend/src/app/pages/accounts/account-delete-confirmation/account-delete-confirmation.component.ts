import { Component, inject, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { Account, AccountType, Currency } from '../../../models';
import { AccountValidationService } from '../../../services';

export interface DeleteConfirmationDialogData {
  account: Account;
}

@Component({
  selector: 'app-account-delete-confirmation',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule
  ],
  template: `
    <div class="delete-confirmation-dialog">
      <mat-dialog-content>
        <div class="dialog-header">
          <mat-icon class="warning-icon">warning</mat-icon>
          <h2>Confirmar Eliminación</h2>
        </div>
        
        <div class="confirmation-content">
          <p>¿Estás seguro de que deseas eliminar la siguiente cuenta?</p>
          
          <div class="account-details">
            <div class="account-info">
              <mat-icon>{{ getAccountIcon() }}</mat-icon>
              <div class="account-text">
                <strong>{{ data.account.name }}</strong>
                <span class="account-type">{{ getAccountTypeDisplayName() }}</span>
                <span class="account-balance">{{ getFormattedBalance() }}</span>
              </div>
            </div>
          </div>
          
          <div class="warning-message">
            <mat-icon>info</mat-icon>
            <p>Esta acción no se puede deshacer. Todos los datos asociados a esta cuenta se perderán permanentemente.</p>
          </div>
        </div>
      </mat-dialog-content>

      <mat-dialog-actions align="end">
        <button mat-button (click)="onCancel()">
          Cancelar
        </button>
        <button mat-raised-button color="warn" (click)="onConfirm()">
          <mat-icon>delete</mat-icon>
          Eliminar Cuenta
        </button>
      </mat-dialog-actions>
    </div>
  `,
  styles: [`
    .delete-confirmation-dialog {
      min-width: 400px;
    }

    .dialog-header {
      display: flex;
      align-items: center;
      gap: 12px;
      margin-bottom: 20px;
    }

    .dialog-header h2 {
      margin: 0;
      color: #f44336;
    }

    .warning-icon {
      color: #ff9800;
      font-size: 32px;
      width: 32px;
      height: 32px;
    }

    .confirmation-content {
      margin-bottom: 20px;
    }

    .account-details {
      background: #f5f5f5;
      border-radius: 8px;
      padding: 16px;
      margin: 16px 0;
    }

    .account-info {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .account-info mat-icon {
      color: #1976d2;
      font-size: 24px;
      width: 24px;
      height: 24px;
    }

    .account-text {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .account-type {
      color: #666;
      font-size: 14px;
    }

    .account-balance {
      color: #4caf50;
      font-weight: 500;
    }

    .warning-message {
      display: flex;
      align-items: flex-start;
      gap: 8px;
      background: #fff3cd;
      border: 1px solid #ffeaa7;
      border-radius: 4px;
      padding: 12px;
      margin-top: 16px;
    }

    .warning-message mat-icon {
      color: #856404;
      font-size: 20px;
      width: 20px;
      height: 20px;
      flex-shrink: 0;
      margin-top: 2px;
    }

    .warning-message p {
      margin: 0;
      color: #856404;
      font-size: 14px;
      line-height: 1.4;
    }

    mat-dialog-actions {
      padding: 16px 0 0 0;
      gap: 12px;
    }

    mat-dialog-actions button {
      min-width: 100px;
    }
  `]
})
export class AccountDeleteConfirmationComponent {
  private readonly dialogRef = inject(MatDialogRef<AccountDeleteConfirmationComponent>);
  private readonly validationService = inject(AccountValidationService);

  constructor(@Inject(MAT_DIALOG_DATA) public data: DeleteConfirmationDialogData) {}

  onConfirm(): void {
    this.dialogRef.close(true);
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }

  getAccountIcon(): string {
    switch (this.data.account.accountType) {
      case AccountType.SAVINGS:
        return 'savings';
      case AccountType.CHECKING:
        return 'account_balance';
      case AccountType.CREDIT:
        return 'credit_card';
      default:
        return 'account_balance_wallet';
    }
  }

  getAccountTypeDisplayName(): string {
    return this.validationService.getAccountTypeDisplayName(this.data.account.accountType);
  }

  getFormattedBalance(): string {
    return this.validationService.formatCurrency(this.data.account.balance, this.data.account.currency);
  }
}