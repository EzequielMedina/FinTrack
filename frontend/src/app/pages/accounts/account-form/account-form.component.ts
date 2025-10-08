import { ChangeDetectionStrategy, Component, inject, Inject, OnInit, signal, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatCardModule } from '@angular/material/card';

import { AccountService, AccountValidationService } from '../../../services';
import { Account, AccountType, Currency, CreateAccountRequest, UpdateAccountRequest, AccountUtils } from '../../../models';

export interface AccountFormDialogData {
  account?: Account; // Si existe, es edición; si no, es creación
  mode: 'create' | 'edit';
  userId?: string; // Usuario actual para nuevas cuentas
}

@Component({
  selector: 'app-account-form',
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
    MatCheckboxModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatCardModule
  ],
  template: `
    <div class="account-form-dialog">
      <mat-dialog-content>
        <!-- Dialog Header -->
        <div class="dialog-header">
          <h2 mat-dialog-title>
            <mat-icon>{{ dialogIcon }}</mat-icon>
            {{ dialogTitle }}
          </h2>
          @if (isEditMode && data.account) {
            <p class="account-info">
              Editando: {{ data.account.name }} ({{ data.account.accountType }})
            </p>
          }
        </div>

        <!-- Account Form -->
        <form [formGroup]="accountForm" (ngSubmit)="onSubmit()" class="account-form">
          
          <!-- Account Name -->
          <mat-form-field appearance="outline" class="full-width">
            <mat-label>Nombre de la Cuenta</mat-label>
            <input matInput 
                   formControlName="name"
                   placeholder="Ej: Mi cuenta de ahorros"
                   maxlength="100">
            <mat-icon matSuffix>account_balance</mat-icon>
            @if (accountForm.get('name')?.invalid && accountForm.get('name')?.touched) {
              <mat-error>{{ getFieldErrorMessage('name') }}</mat-error>
            }
          </mat-form-field>

          <!-- Account Type -->
          <mat-form-field appearance="outline" class="full-width">
            <mat-label>Tipo de Cuenta</mat-label>
            <mat-select formControlName="accountType" [disabled]="isEditMode">
              <mat-select-trigger>
                {{ getSelectedAccountTypeLabel() }}
              </mat-select-trigger>
              @for (type of accountTypes; track type.value) {
                <mat-option [value]="type.value">
                  <div class="option-content">
                    <mat-icon>{{ type.icon }}</mat-icon>
                    <span>{{ type.label }}</span>
                  </div>
                </mat-option>
              }
            </mat-select>
            <mat-icon matSuffix>{{ getSelectedAccountTypeIcon() }}</mat-icon>
            @if (accountForm.get('accountType')?.invalid && accountForm.get('accountType')?.touched) {
              <mat-error>{{ getFieldErrorMessage('accountType') }}</mat-error>
            }
            @if (isEditMode) {
              <mat-hint>El tipo de cuenta no se puede modificar</mat-hint>
            }
          </mat-form-field>

          <!-- Currency -->
          <mat-form-field appearance="outline" class="full-width">
            <mat-label>Moneda</mat-label>
            <mat-select formControlName="currency" [disabled]="isEditMode">
              <mat-select-trigger>
                {{ getSelectedCurrencyLabel() }}
              </mat-select-trigger>
              @for (currency of currencies; track currency.value) {
                <mat-option [value]="currency.value">
                  <div class="option-content">
                    <span class="currency-symbol">{{ currency.symbol }}</span>
                    <span>{{ currency.label }}</span>
                  </div>
                </mat-option>
              }
            </mat-select>
            <mat-icon matSuffix>{{ getSelectedCurrencySymbol() }}</mat-icon>
            @if (accountForm.get('currency')?.invalid && accountForm.get('currency')?.touched) {
              <mat-error>{{ getFieldErrorMessage('currency') }}</mat-error>
            }
            @if (isEditMode) {
              <mat-hint>La moneda no se puede modificar</mat-hint>
            }
          </mat-form-field>

          <!-- Initial Balance (only for creation) -->
          @if (shouldShowInitialBalance()) {
            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Saldo Inicial</mat-label>
              <input matInput 
                     type="number"
                     formControlName="initialBalance"
                     placeholder="0.00"
                     min="0"
                     max="999999999"
                     step="0.01">
              <span matPrefix>{{ getSelectedCurrencySymbol() }}&nbsp;</span>
              <mat-icon matSuffix>attach_money</mat-icon>
              @if (accountForm.get('initialBalance')?.invalid && accountForm.get('initialBalance')?.touched) {
                <mat-error>{{ getFieldErrorMessage('initialBalance') }}</mat-error>
              }
              <mat-hint>Monto inicial que tendrá la cuenta</mat-hint>
            </mat-form-field>
          }

          <!-- Credit Card Specific Fields -->
          @if (shouldShowCreditFields()) {
            <div class="credit-section">
              <h3 class="section-title">
                <mat-icon>credit_card</mat-icon>
                Configuración de Crédito
              </h3>

              <!-- Credit Limit -->
              <mat-form-field appearance="outline" class="full-width">
                <mat-label>Límite de Crédito</mat-label>
                <input matInput 
                       type="number"
                       formControlName="creditLimit"
                       placeholder="0.00"
                       min="100"
                       max="999999999"
                       step="100">
                <span matPrefix>{{ getSelectedCurrencySymbol() }}&nbsp;</span>
                <mat-icon matSuffix>credit_score</mat-icon>
                @if (accountForm.get('creditLimit')?.invalid && accountForm.get('creditLimit')?.touched) {
                  <mat-error>{{ getFieldErrorMessage('creditLimit') }}</mat-error>
                }
                <mat-hint>Límite máximo de crédito disponible</mat-hint>
              </mat-form-field>

              <!-- Credit Dates Row -->
              <div class="dates-row">
                <!-- Closing Date -->
                <mat-form-field appearance="outline" class="half-width">
                  <mat-label>Fecha de Cierre</mat-label>
                  <input matInput 
                         [matDatepicker]="closingPicker"
                         formControlName="closingDate"
                         readonly>
                  <mat-datepicker-toggle matSuffix [for]="closingPicker"></mat-datepicker-toggle>
                  <mat-datepicker #closingPicker></mat-datepicker>
                  @if (accountForm.get('closingDate')?.invalid && accountForm.get('closingDate')?.touched) {
                    <mat-error>{{ getFieldErrorMessage('closingDate') }}</mat-error>
                  }
                  <mat-hint>Día de cierre del ciclo</mat-hint>
                </mat-form-field>

                <!-- Due Date -->
                <mat-form-field appearance="outline" class="half-width">
                  <mat-label>Fecha de Vencimiento</mat-label>
                  <input matInput 
                         [matDatepicker]="duePicker"
                         formControlName="dueDate"
                         readonly>
                  <mat-datepicker-toggle matSuffix [for]="duePicker"></mat-datepicker-toggle>
                  <mat-datepicker #duePicker></mat-datepicker>
                  @if (accountForm.get('dueDate')?.invalid && accountForm.get('dueDate')?.touched) {
                    <mat-error>{{ getFieldErrorMessage('dueDate') }}</mat-error>
                  }
                  <mat-hint>Día de vencimiento del pago</mat-hint>
                </mat-form-field>
              </div>
            </div>
          }

          <!-- Virtual Wallet Specific Fields -->
          @if (shouldShowWalletFields()) {
            <div class="wallet-section">
              <h3 class="section-title">
                <mat-icon>account_balance_wallet</mat-icon>
                Configuración de Billetera Virtual
              </h3>

              <!-- DNI -->
              <mat-form-field appearance="outline" class="full-width">
                <mat-label>DNI</mat-label>
                <input matInput 
                       formControlName="dni"
                       placeholder="12345678"
                       maxlength="8">
                <mat-icon matSuffix>person</mat-icon>
                @if (accountForm.get('dni')?.invalid && accountForm.get('dni')?.touched) {
                  <mat-error>{{ getFieldErrorMessage('dni') }}</mat-error>
                }
                <mat-hint>DNI requerido para billeteras virtuales</mat-hint>
              </mat-form-field>
            </div>
          }

          <!-- Bank Account Specific Fields -->
          @if (shouldShowBankAccountFields()) {
            <div class="bank-account-section">
              <h3 class="section-title">
                <mat-icon>account_balance</mat-icon>
                Configuración de Cuenta Bancaria
              </h3>
              
              <mat-card class="info-card">
                <mat-card-content>
                  <div class="info-content">
                    <mat-icon>info</mat-icon>
                    <div>
                      <p><strong>Cuenta Bancaria con Tarjetas</strong></p>
                      <p>Podrás agregar múltiples tarjetas de crédito y débito después de crear la cuenta.</p>
                    </div>
                  </div>
                </mat-card-content>
              </mat-card>
            </div>
          }

          <!-- Account Status -->
          <div class="status-section">
            <mat-checkbox formControlName="isActive" color="primary">
              <div class="checkbox-content">
                <span class="checkbox-label">Cuenta Activa</span>
                <small class="checkbox-hint">
                  @if (accountForm.get('isActive')?.value) {
                    La cuenta estará disponible para operaciones
                  } @else {
                    La cuenta estará deshabilitada para nuevas operaciones
                  }
                </small>
              </div>
            </mat-checkbox>
          </div>

        </form>
      </mat-dialog-content>

      <!-- Dialog Actions -->
      <mat-dialog-actions align="end" class="dialog-actions">
        <button mat-button 
                type="button" 
                (click)="onCancel()"
                [disabled]="loading()">
          Cancelar
        </button>
        
        <button mat-raised-button 
                color="primary"
                type="submit"
                (click)="onSubmit()"
                [disabled]="loading() || accountForm.invalid">
          @if (loading()) {
            <mat-spinner diameter="20"></mat-spinner>
          } @else {
            <mat-icon>{{ isEditMode ? 'save' : 'add' }}</mat-icon>
          }
          {{ submitButtonText }}
        </button>
      </mat-dialog-actions>
    </div>
  `,
  styles: [`
    .account-form-dialog {
      min-width: 500px;
      max-width: 700px;
    }

    .dialog-header {
      margin-bottom: 24px;
      padding-bottom: 16px;
      border-bottom: 1px solid #e0e0e0;
    }

    .dialog-header h2 {
      margin: 0 0 8px 0;
      display: flex;
      align-items: center;
      gap: 12px;
      color: #1976d2;
      font-weight: 500;
    }

    .dialog-header h2 mat-icon {
      font-size: 28px;
      width: 28px;
      height: 28px;
    }

    .account-info {
      margin: 0;
      color: #666;
      font-size: 14px;
      font-style: italic;
    }

    .account-form {
      display: flex;
      flex-direction: column;
      gap: 20px;
      min-height: 300px;
    }

    .full-width {
      width: 100%;
    }

    .half-width {
      flex: 1;
      margin-right: 8px;
    }

    .half-width:last-child {
      margin-right: 0;
      margin-left: 8px;
    }

    .dates-row {
      display: flex;
      gap: 16px;
      width: 100%;
    }

    .option-content {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .option-content mat-icon {
      font-size: 20px;
      width: 20px;
      height: 20px;
      color: #666;
    }

    .option-icon {
      display: inline-flex;
      align-items: center;
      margin-right: 12px;
      vertical-align: middle;
    }

    .option-icon mat-icon {
      font-size: 20px;
      width: 20px;
      height: 20px;
      color: #666;
    }

    .currency-symbol {
      font-weight: bold;
      color: #1976d2;
      min-width: 30px;
      text-align: center;
      margin-right: 12px;
      display: inline-block;
    }

    .credit-section {
      background: #f8f9fa;
      padding: 20px;
      border-radius: 8px;
      border-left: 4px solid #4caf50;
      margin: 16px 0;
    }

    .section-title {
      margin: 0 0 20px 0;
      display: flex;
      align-items: center;
      gap: 8px;
      color: #4caf50;
      font-weight: 500;
      font-size: 16px;
    }

    .section-title mat-icon {
      font-size: 20px;
      width: 20px;
      height: 20px;
    }

    .status-section {
      background: #f5f5f5;
      padding: 16px;
      border-radius: 8px;
      margin: 16px 0;
    }

    .checkbox-content {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .checkbox-label {
      font-weight: 500;
      color: #333;
    }

    .checkbox-hint {
      color: #666;
      font-size: 12px;
      line-height: 1.4;
    }

    .dialog-actions {
      padding: 16px 24px;
      border-top: 1px solid #e0e0e0;
      margin-top: 24px;
      gap: 12px;
    }

    .dialog-actions button {
      min-width: 120px;
      height: 40px;
    }

    .dialog-actions button mat-spinner {
      margin-right: 8px;
    }

    .dialog-actions button mat-icon {
      margin-right: 8px;
    }

    /* Form validation styles */
    .mat-mdc-form-field.mat-form-field-invalid .mat-mdc-text-field-wrapper {
      background-color: #ffeaea;
    }

    .mat-mdc-form-field.mat-form-field-invalid .mat-mdc-floating-label {
      color: #f44336 !important;
    }

    /* Responsive design */
    @media (max-width: 768px) {
      .account-form-dialog {
        min-width: 90vw;
        max-width: 95vw;
      }
      
      .dates-row {
        flex-direction: column;
        gap: 20px;
      }
      
      .half-width {
        width: 100%;
        margin: 0;
      }
      
      .dialog-actions {
        flex-direction: column-reverse;
      }
      
      .dialog-actions button {
        width: 100%;
      }
    }

    /* Animation enhancements */
    .credit-section {
      animation: slideIn 0.3s ease-out;
    }

    @keyframes slideIn {
      from {
        opacity: 0;
        transform: translateY(-10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    /* Loading state */
    .dialog-actions button[disabled] {
      opacity: 0.6;
      cursor: not-allowed;
    }

    /* Custom Material Design adjustments */
    .mat-mdc-dialog-content {
      padding: 24px !important;
    }

    .mat-mdc-form-field-hint-wrapper {
      padding-top: 4px;
    }

    .mat-mdc-form-field-hint {
      font-size: 11px;
      color: #666;
    }

    /* Accessibility improvements */
    .mat-mdc-form-field-error {
      font-size: 12px;
      line-height: 1.4;
    }

    .option-content:hover {
      background-color: rgba(25, 118, 210, 0.04);
    }

    /* Focus styles */
    .mat-mdc-form-field.mat-focused .mat-mdc-floating-label {
      color: #1976d2 !important;
    }

    .mat-mdc-form-field.mat-focused .mat-mdc-text-field-wrapper {
      border-color: #1976d2 !important;
    }
  `],
  changeDetection: ChangeDetectionStrategy.Default
})
export class AccountFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly accountService = inject(AccountService);
  private readonly validationService = inject(AccountValidationService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialogRef = inject(MatDialogRef<AccountFormComponent>);
  private readonly cdr = inject(ChangeDetectorRef);

  // Signals
  loading = signal(false);
  
  // Form
  accountForm!: FormGroup;

  // Enums for template
  readonly AccountType = AccountType;
  readonly Currency = Currency;

  // Options for selects
  readonly accountTypes = [
    { value: AccountType.WALLET, label: 'Billetera Virtual', icon: 'account_balance_wallet', description: 'Para pagos digitales sin tarjetas físicas' },
    { value: AccountType.BANK_ACCOUNT, label: 'Cuenta Bancaria', icon: 'account_balance', description: 'Cuenta tradicional que puede tener múltiples tarjetas' },
    { value: AccountType.SAVINGS, label: 'Cuenta de Ahorro', icon: 'savings', description: 'Cuenta para ahorrar dinero' },
    { value: AccountType.CHECKING, label: 'Cuenta Corriente', icon: 'receipt_long', description: 'Cuenta para uso diario con chequera' },
    { value: AccountType.CREDIT, label: 'Tarjeta de Crédito', icon: 'credit_card', description: 'Línea de crédito con límite establecido' },
    { value: AccountType.DEBIT, label: 'Tarjeta de Débito', icon: 'payment', description: 'Tarjeta vinculada directamente al balance' }
  ];

  readonly currencies = [
    { value: Currency.ARS, label: 'Peso Argentino (ARS)', symbol: '$' },
    { value: Currency.USD, label: 'Dólar Estadounidense (USD)', symbol: 'US$' }
  ];

  // Computed values
  get isEditMode(): boolean {
    return this.data.mode === 'edit';
  }

  get dialogTitle(): string {
    return this.isEditMode ? 'Editar Cuenta' : 'Crear Nueva Cuenta';
  }

  get submitButtonText(): string {
    return this.isEditMode ? 'Actualizar Cuenta' : 'Crear Cuenta';
  }

  get dialogIcon(): string {
    return this.isEditMode ? 'edit' : 'add';
  }

  constructor(@Inject(MAT_DIALOG_DATA) public data: AccountFormDialogData) {}

  ngOnInit(): void {
    // Verificar que se haya proporcionado un userId para crear cuentas
    if (!this.isEditMode && !this.data.userId) {
      console.error('Error: userId is required for creating accounts');
      this.dialogRef.close();
      return;
    }
    this.initializeForm();
  }

  private initializeForm(): void {
    const account = this.data.account;
    
    this.accountForm = this.fb.group({
      userId: [this.data.userId || account?.userId || '', [Validators.required]],
      name: [account?.name || '', [
        Validators.required,
        Validators.maxLength(100)
      ]],
      accountType: [account?.accountType || AccountType.SAVINGS, [
        Validators.required
      ]],
      currency: [account?.currency || Currency.ARS, [
        Validators.required
      ]],
      initialBalance: [account?.balance || 0, [
        Validators.required,
        Validators.min(0),
        Validators.max(999999999)
      ]],
      creditLimit: [account?.creditLimit || 0, [
        Validators.min(0),
        Validators.max(999999999)
      ]],
      closingDate: [account?.closingDate ? new Date(account.closingDate) : null],
      dueDate: [account?.dueDate ? new Date(account.dueDate) : null],
      dni: [account?.dni || ''],
      isActive: [account?.isActive !== false] // Default true for new accounts
    });

    // Set up conditional validators based on account type
    this.setupConditionalValidators();
  }

  private setupConditionalValidators(): void {
    const accountTypeControl = this.accountForm.get('accountType');
    const creditLimitControl = this.accountForm.get('creditLimit');
    const closingDateControl = this.accountForm.get('closingDate');
    const dueDateControl = this.accountForm.get('dueDate');
    const dniControl = this.accountForm.get('dni');

    accountTypeControl?.valueChanges.subscribe(accountType => {
      // Credit card validations
      if (accountType === AccountType.CREDIT) {
        // Credit card requires credit limit and dates
        creditLimitControl?.setValidators([
          Validators.required,
          Validators.min(100),
          Validators.max(999999999)
        ]);
        closingDateControl?.setValidators([Validators.required]);
        dueDateControl?.setValidators([Validators.required]);
      } else {
        // Regular accounts don't need credit fields
        creditLimitControl?.setValidators([
          Validators.min(0),
          Validators.max(999999999)
        ]);
        closingDateControl?.clearValidators();
        dueDateControl?.clearValidators();
        
        // Reset values for non-credit accounts
        if (accountType !== AccountType.CREDIT) {
          creditLimitControl?.setValue(0);
          closingDateControl?.setValue(null);
          dueDateControl?.setValue(null);
        }
      }

      // Virtual wallet validations
      if (accountType === AccountType.WALLET) {
        dniControl?.setValidators([
          Validators.required,
          Validators.minLength(7),
          Validators.maxLength(8),
          Validators.pattern(/^\d+$/) // Only numbers
        ]);
      } else {
        dniControl?.clearValidators();
      }

      creditLimitControl?.updateValueAndValidity();
      closingDateControl?.updateValueAndValidity();
      dueDateControl?.updateValueAndValidity();
      dniControl?.updateValueAndValidity();
    });
  }

  onSubmit(): void {
    if (this.accountForm.invalid) {
      this.markFormGroupTouched(this.accountForm);
      console.log('Form is invalid. Errors:', this.getFormErrors());
      return;
    }

    const formValue = this.accountForm.value;

    // Validate using the service
    const validation = this.isEditMode 
      ? this.validationService.validateUpdateAccount(formValue)
      : this.validationService.validateCreateAccount(formValue);

    if (!validation.isValid) {
      this.snackBar.open(validation.errors.join(', '), 'Cerrar', {
        duration: 5000,
        panelClass: ['error-snackbar']
      });
      return;
    }

    this.loading.set(true);

    if (this.isEditMode) {
      this.updateAccount(formValue);
    } else {
      this.createAccount(formValue);
    }
  }

  private createAccount(formValue: any): void {
    const request: CreateAccountRequest = {
      userId: formValue.userId,
      name: formValue.name,
      accountType: formValue.accountType,
      currency: formValue.currency,
      initialBalance: formValue.initialBalance,
      isActive: formValue.isActive,
      creditLimit: formValue.accountType === AccountType.CREDIT ? formValue.creditLimit : undefined,
      closingDate: formValue.accountType === AccountType.CREDIT && formValue.closingDate 
        ? formValue.closingDate.toISOString().split('T')[0] : undefined,
      dueDate: formValue.accountType === AccountType.CREDIT && formValue.dueDate 
        ? formValue.dueDate.toISOString().split('T')[0] : undefined
    };
    // Solo agregar dni si es billetera virtual
    if (formValue.accountType === AccountType.WALLET && formValue.dni) {
      (request as any).dni = formValue.dni;
    }

    this.accountService.createAccount(request).subscribe({
      next: (newAccount) => {
        this.loading.set(false);
        this.snackBar.open('Cuenta creada exitosamente', 'Cerrar', {
          duration: 3000,
          panelClass: ['success-snackbar']
        });
        this.dialogRef.close(newAccount);
      },
      error: (error) => {
        this.loading.set(false);
        console.error('Error creating account:', error);
        this.snackBar.open('Error al crear la cuenta. Inténtalo de nuevo.', 'Cerrar', {
          duration: 5000,
          panelClass: ['error-snackbar']
        });
      }
    });
  }

  private updateAccount(formValue: any): void {
    const request: UpdateAccountRequest = {
      name: formValue.name,
      creditLimit: formValue.accountType === AccountType.CREDIT ? formValue.creditLimit : undefined,
      closingDate: formValue.accountType === AccountType.CREDIT && formValue.closingDate 
        ? formValue.closingDate.toISOString().split('T')[0] : undefined,
      dueDate: formValue.accountType === AccountType.CREDIT && formValue.dueDate 
        ? formValue.dueDate.toISOString().split('T')[0] : undefined
    };
    // Solo agregar dni si es billetera virtual
    if (formValue.accountType === AccountType.WALLET && formValue.dni) {
      (request as any).dni = formValue.dni;
    }

    this.accountService.updateAccount(this.data.account!.id, request).subscribe({
      next: (updatedAccount) => {
        this.loading.set(false);
        this.snackBar.open('Cuenta actualizada exitosamente', 'Cerrar', {
          duration: 3000,
          panelClass: ['success-snackbar']
        });
        this.dialogRef.close(updatedAccount);
      },
      error: (error) => {
        this.loading.set(false);
        console.error('Error updating account:', error);
        this.snackBar.open('Error al actualizar la cuenta. Inténtalo de nuevo.', 'Cerrar', {
          duration: 5000,
          panelClass: ['error-snackbar']
        });
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
      
      if (control instanceof FormGroup) {
        this.markFormGroupTouched(control);
      }
    });
  }

  getFieldErrorMessage(fieldName: string): string {
    const control = this.accountForm.get(fieldName);
    if (!control || !control.errors) {
      return '';
    }

    const errors = control.errors;
    
    if (errors['required']) {
      return `${this.getFieldDisplayName(fieldName)} es requerido`;
    }
    
    if (errors['minlength']) {
      return `${this.getFieldDisplayName(fieldName)} debe tener al menos ${errors['minlength'].requiredLength} caracteres`;
    }
    
    if (errors['maxlength']) {
      return `${this.getFieldDisplayName(fieldName)} no puede exceder ${errors['maxlength'].requiredLength} caracteres`;
    }
    
    if (errors['min']) {
      return `${this.getFieldDisplayName(fieldName)} debe ser mayor o igual a ${errors['min'].min}`;
    }
    
    if (errors['max']) {
      return `${this.getFieldDisplayName(fieldName)} no puede ser mayor a ${errors['max'].max}`;
    }

    if (errors['pattern']) {
      if (fieldName === 'dni') {
        return 'El DNI debe contener solo números';
      }
      return 'Formato inválido';
    }

    return 'Campo inválido';
  }

  private getFieldDisplayName(fieldName: string): string {
    const displayNames: { [key: string]: string } = {
      name: 'Nombre de la cuenta',
      accountType: 'Tipo de cuenta',
      currency: 'Moneda',
      initialBalance: 'Saldo inicial',
      creditLimit: 'Límite de crédito',
      closingDate: 'Fecha de cierre',
      dueDate: 'Fecha de vencimiento',
      dni: 'DNI'
    };
    
    return displayNames[fieldName] || fieldName;
  }

  // Template helper methods
  shouldShowCreditFields(): boolean {
    return this.accountForm.get('accountType')?.value === AccountType.CREDIT;
  }

  shouldShowWalletFields(): boolean {
    return this.accountForm.get('accountType')?.value === AccountType.WALLET;
  }

  shouldShowBankAccountFields(): boolean {
    return this.accountForm.get('accountType')?.value === AccountType.BANK_ACCOUNT;
  }

  shouldShowInitialBalance(): boolean {
    return !this.isEditMode; // Only show for creation
  }

  getSelectedAccountTypeIcon(): string {
    const accountType = this.accountForm.get('accountType')?.value;
    const typeConfig = this.accountTypes.find(type => type.value === accountType);
    return typeConfig?.icon || 'account_balance';
  }

  getSelectedCurrencySymbol(): string {
    const currency = this.accountForm.get('currency')?.value;
    const currencyConfig = this.currencies.find(curr => curr.value === currency);
    return currencyConfig?.symbol || '$';
  }

  getSelectedAccountTypeLabel(): string {
    const accountType = this.accountForm.get('accountType')?.value;
    const typeConfig = this.accountTypes.find(type => type.value === accountType);
    return typeConfig?.label || '';
  }

  getSelectedCurrencyLabel(): string {
    const currency = this.accountForm.get('currency')?.value;
    const currencyConfig = this.currencies.find(curr => curr.value === currency);
    return currencyConfig?.label || '';
  }

  // Helper method for debugging form validation
  private getFormErrors(): any {
    const formErrors: any = {};
    
    Object.keys(this.accountForm.controls).forEach(key => {
      const controlErrors = this.accountForm.get(key)?.errors;
      if (controlErrors) {
        formErrors[key] = controlErrors;
      }
    });
    
    return formErrors;
  }
}