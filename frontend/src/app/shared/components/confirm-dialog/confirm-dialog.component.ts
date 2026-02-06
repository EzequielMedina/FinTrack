import { Component, Inject, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

export interface ConfirmDialogData {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  type?: 'warn' | 'primary' | 'accent';
}

@Component({
  selector: 'app-confirm-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule
  ],
  template: `
    <div class="confirm-dialog">
      <div class="dialog-header">
        <span class="header-icon" [class]="'icon-' + (data.type || 'warn')">
          <mat-icon>{{ getIcon() }}</mat-icon>
        </span>
        <h2 mat-dialog-title class="dialog-title">{{ data.title }}</h2>
      </div>

      <div mat-dialog-content class="dialog-content">
        <p class="message">{{ data.message }}</p>
      </div>

      <div mat-dialog-actions class="dialog-actions">
        <button mat-button (click)="onCancel()" class="cancel-btn">
          {{ data.cancelText || 'Cancelar' }}
        </button>
        <button
          mat-raised-button
          [color]="data.type === 'primary' ? 'primary' : 'warn'"
          (click)="onConfirm()"
          class="confirm-btn">
          <mat-icon>{{ getConfirmIcon() }}</mat-icon>
          {{ data.confirmText || 'Confirmar' }}
        </button>
      </div>
    </div>
  `,
  styles: [`
    .confirm-dialog {
      min-width: 360px;
      max-width: 480px;
      overflow: hidden;
      box-sizing: border-box;
    }

    .dialog-header {
      display: flex;
      align-items: center;
      gap: var(--space-3);
      margin-bottom: var(--space-4);
      min-width: 0;
    }

    .header-icon {
      flex-shrink: 0;
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      border-radius: var(--radius-full);
      background: var(--bg-tertiary);
    }

    .header-icon mat-icon {
      font-size: 24px;
      width: 24px;
      height: 24px;
    }

    .dialog-title {
      margin: 0 !important;
      flex: 1;
      min-width: 0;
      font-size: var(--text-lg);
      font-weight: var(--font-semibold);
    }

    .header-icon.icon-warn {
      background: var(--error-100);
      color: var(--error-700);
    }

    .header-icon.icon-primary {
      background: var(--accent-100);
      color: var(--accent-600);
    }

    .header-icon.icon-accent {
      background: var(--accent-100);
      color: var(--accent-600);
    }

    .dialog-content {
      margin: var(--space-4) 0;
    }

    .message {
      margin: 0;
      color: var(--text-primary);
      line-height: var(--leading-relaxed);
      font-size: var(--text-sm);
    }

    .dialog-actions {
      display: flex;
      justify-content: flex-end;
      gap: var(--space-3);
      margin-top: var(--space-6);
      padding-bottom: 0;
    }

    .cancel-btn {
      color: var(--text-secondary);
    }

    .confirm-btn mat-icon {
      margin-right: var(--space-2);
      font-size: 18px;
      width: 18px;
      height: 18px;
    }
  `]
})
export class ConfirmDialogComponent {
  private readonly dialogRef = inject(MatDialogRef<ConfirmDialogComponent>);

  constructor(@Inject(MAT_DIALOG_DATA) public data: ConfirmDialogData) {}

  getIcon(): string {
    return this.data.type === 'primary' ? 'info' : 'warning';
  }

  getConfirmIcon(): string {
    return 'check';
  }

  onConfirm(): void {
    this.dialogRef.close(true);
  }

  onCancel(): void {
    this.dialogRef.close(false);
  }
}
