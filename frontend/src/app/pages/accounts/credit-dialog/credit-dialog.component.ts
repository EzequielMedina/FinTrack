import { Component, Inject, inject, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA, MatDialogModule } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatCardModule } from '@angular/material/card';
import { CreditService, AccountValidationService } from '../../../services';
import { Account, UpdateCreditLimitRequest, UpdateCreditDatesRequest, AvailableCreditResponse } from '../../../models';

export interface CreditDialogData {
  account: Account;
  operation?: 'limit' | 'dates' | 'info';
}

@Component({
  selector: 'app-credit-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatTabsModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatCardModule
  ],
  templateUrl: './credit-dialog.component.html',
  styleUrls: ['./credit-dialog.component.css']
})
export class CreditDialogComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly creditService = inject(CreditService);
  private readonly validationService = inject(AccountValidationService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialogRef = inject(MatDialogRef<CreditDialogComponent>);

  // Signals
  loading = signal(false);
  selectedOperation = signal<'limit' | 'dates' | 'info'>('info');
  availableCredit = signal<AvailableCreditResponse | null>(null);
  loadingCredit = signal(false);

  // Forms
  creditLimitForm!: FormGroup;
  creditDatesForm!: FormGroup;

  // Computed values
  currentCreditLimit = computed(() => this.data.account.creditLimit || 0);
  currentClosingDate = computed(() => this.data.account.closingDate);
  currentDueDate = computed(() => this.data.account.dueDate);
  
  formattedCreditLimit = computed(() => 
    this.validationService.formatCurrency(this.currentCreditLimit(), this.data.account.currency)
  );

  // Constants for validation
  readonly MIN_CREDIT_LIMIT = 0;
  readonly MAX_CREDIT_LIMIT = 999999999;
  readonly CREDIT_LIMIT_STEP = 100;

  constructor(@Inject(MAT_DIALOG_DATA) public data: CreditDialogData) {
    if (data.operation) {
      this.selectedOperation.set(data.operation);
    }
  }

  ngOnInit(): void {
    this.initializeForms();
    this.loadAvailableCredit();
  }

  private initializeForms(): void {
    // Credit Limit Form
    this.creditLimitForm = this.fb.group({
      creditLimit: [this.currentCreditLimit(), [
        Validators.required,
        Validators.min(this.MIN_CREDIT_LIMIT),
        Validators.max(this.MAX_CREDIT_LIMIT)
      ]]
    });

    // Credit Dates Form
    this.creditDatesForm = this.fb.group({
      closingDate: [this.currentClosingDate() ? new Date(this.currentClosingDate()!) : null, [
        Validators.required
      ]],
      dueDate: [this.currentDueDate() ? new Date(this.currentDueDate()!) : null, [
        Validators.required
      ]]
    });
  }

  private loadAvailableCredit(): void {
    this.loadingCredit.set(true);
    
    this.creditService.getAvailableCredit(this.data.account.id).subscribe({
      next: (creditInfo) => {
        this.availableCredit.set(creditInfo);
        this.loadingCredit.set(false);
      },
      error: (error) => {
        console.error('Error loading available credit:', error);
        this.loadingCredit.set(false);
      }
    });
  }

  onOperationChange(operation: 'limit' | 'dates' | 'info'): void {
    this.selectedOperation.set(operation);
  }

  onUpdateCreditLimit(): void {
    if (this.creditLimitForm.invalid) {
      this.markFormGroupTouched(this.creditLimitForm);
      return;
    }

    const formValue = this.creditLimitForm.value;
    const request: UpdateCreditLimitRequest = {
      creditLimit: formValue.creditLimit
    };

    // Validate using the service
    const validation = this.creditService.validateCreditLimit(request.creditLimit);
    if (!validation.isValid) {
      this.snackBar.open(validation.errors[0]?.message || 'Límite de crédito inválido', 'Cerrar', {
        duration: 3000
      });
      return;
    }

    this.loading.set(true);

    this.creditService.updateCreditLimit(this.data.account.id, request).subscribe({
      next: (updatedAccount) => {
        this.snackBar.open(
          `Límite de crédito actualizado exitosamente a ${this.validationService.formatCurrency(updatedAccount.creditLimit || 0, updatedAccount.currency)}`,
          'Cerrar',
          { duration: 5000 }
        );
        this.dialogRef.close(updatedAccount);
      },
      error: (error) => {
        console.error('Error updating credit limit:', error);
        this.snackBar.open(error.message || 'Error al actualizar el límite de crédito', 'Cerrar', {
          duration: 3000
        });
        this.loading.set(false);
      }
    });
  }

  onUpdateCreditDates(): void {
    if (this.creditDatesForm.invalid) {
      this.markFormGroupTouched(this.creditDatesForm);
      return;
    }

    const formValue = this.creditDatesForm.value;
    const request: UpdateCreditDatesRequest = {
      closingDate: this.formatDateForBackend(formValue.closingDate),
      dueDate: this.formatDateForBackend(formValue.dueDate)
    };

    // Validate using the service
    const validation = this.creditService.validateCreditDates(request.closingDate, request.dueDate);
    if (!validation.isValid) {
      this.snackBar.open(validation.errors[0]?.message || 'Fechas de crédito inválidas', 'Cerrar', {
        duration: 3000
      });
      return;
    }

    this.loading.set(true);

    this.creditService.updateCreditDates(this.data.account.id, request).subscribe({
      next: (updatedAccount) => {
        this.snackBar.open(
          'Fechas de crédito actualizadas exitosamente',
          'Cerrar',
          { duration: 5000 }
        );
        this.dialogRef.close(updatedAccount);
      },
      error: (error) => {
        console.error('Error updating credit dates:', error);
        this.snackBar.open(error.message || 'Error al actualizar las fechas de crédito', 'Cerrar', {
          duration: 3000
        });
        this.loading.set(false);
      }
    });
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  // Helper methods
  private markFormGroupTouched(formGroup: FormGroup): void {
    Object.keys(formGroup.controls).forEach(key => {
      const control = formGroup.get(key);
      control?.markAsTouched();
    });
  }

  private formatDateForBackend(date: Date): string {
    if (!date) return '';
    return date.toISOString().split('T')[0];
  }

  // Quick credit limit buttons
  setQuickCreditLimit(limit: number): void {
    this.creditLimitForm.get('creditLimit')?.setValue(limit);
  }

  // Auto-calculate due date based on closing date
  onClosingDateChange(): void {
    const closingDate = this.creditDatesForm.get('closingDate')?.value;
    if (closingDate) {
      const suggestedDueDate = this.creditService.getNextDueDate(this.formatDateForBackend(closingDate));
      this.creditDatesForm.get('dueDate')?.setValue(new Date(suggestedDueDate));
    }
  }

  // Validation methods for template
  getCreditLimitErrorMessage(): string {
    const control = this.creditLimitForm.get('creditLimit');
    if (!control) return '';

    if (control.hasError('required')) {
      return 'El límite de crédito es requerido';
    }
    if (control.hasError('min')) {
      return `El límite mínimo es $${this.MIN_CREDIT_LIMIT.toLocaleString()}`;
    }
    if (control.hasError('max')) {
      return `El límite máximo es $${this.MAX_CREDIT_LIMIT.toLocaleString()}`;
    }
    return '';
  }

  getClosingDateErrorMessage(): string {
    const control = this.creditDatesForm.get('closingDate');
    if (!control) return '';

    if (control.hasError('required')) {
      return 'La fecha de corte es requerida';
    }
    return '';
  }

  getDueDateErrorMessage(): string {
    const control = this.creditDatesForm.get('dueDate');
    if (!control) return '';

    if (control.hasError('required')) {
      return 'La fecha de vencimiento es requerida';
    }
    return '';
  }

  // Utility methods for template
  getQuickCreditLimits(): number[] {
    return [50000, 100000, 200000, 500000, 1000000];
  }

  getAccountDisplayName(): string {
    return this.data.account.name || `${this.validationService.getAccountTypeDisplayName(this.data.account.accountType)}`;
  }

  formatDate(dateString?: string): string {
    if (!dateString) return 'No establecida';
    
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  formatDateShort(dateString?: string): string {
    if (!dateString) return 'No establecida';
    
    const date = new Date(dateString);
    return date.toLocaleDateString('es-AR', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    });
  }

  calculateMinimumPayment(): number {
    const creditInfo = this.availableCredit();
    if (!creditInfo) return 0;
    
    return this.creditService.calculateMinimumPayment(creditInfo.usedCredit);
  }

  calculateInterestRate(): number {
    const creditLimit = this.currentCreditLimit();
    return this.creditService.calculateInterestRate(creditLimit);
  }

  getNextClosingDate(): string {
    return this.creditService.getNextClosingDate(this.currentClosingDate());
  }

  getCreditUtilizationPercentage(): number {
    const creditInfo = this.availableCredit();
    if (!creditInfo || creditInfo.creditLimit === 0) return 0;
    
    return (creditInfo.usedCredit / creditInfo.creditLimit) * 100;
  }

  getCreditUtilizationColor(): string {
    const percentage = this.getCreditUtilizationPercentage();
    if (percentage >= 90) return 'warn';
    if (percentage >= 70) return 'accent';
    return 'primary';
  }
}