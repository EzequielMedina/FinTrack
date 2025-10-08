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
import { MatStepperModule } from '@angular/material/stepper';
import { MatTabsModule } from '@angular/material/tabs';
import { WalletService, AccountValidationService } from '../../../services';
import { Account, AddFundsRequest, WithdrawFundsRequest } from '../../../models';

export interface WalletDialogData {
  account: Account;
  operation?: 'add' | 'withdraw';
}

@Component({
  selector: 'app-wallet-dialog',
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
    MatStepperModule,
    MatTabsModule
  ],
  templateUrl: './wallet-dialog.component.html',
  styleUrls: ['./wallet-dialog.component.css']
})
export class WalletDialogComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly walletService = inject(WalletService);
  private readonly validationService = inject(AccountValidationService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialogRef = inject(MatDialogRef<WalletDialogComponent>);

  // Signals
  loading = signal(false);
  selectedOperation = signal<'add' | 'withdraw'>('add');
  currentAccount = signal<Account>(this.data.account);

  // Forms
  addFundsForm!: FormGroup;
  withdrawFundsForm!: FormGroup;

  // Computed values
  currentBalance = computed(() => this.currentAccount().balance);
  formattedBalance = computed(() => 
    this.validationService.formatCurrency(this.currentBalance(), this.currentAccount().currency)
  );

  // Constants for validation
  readonly MIN_AMOUNT = 1;
  readonly MAX_AMOUNT = 999999999;
  readonly DAILY_LIMIT = 50000;

  constructor(@Inject(MAT_DIALOG_DATA) public data: WalletDialogData) {
    if (data.operation) {
      this.selectedOperation.set(data.operation);
    }
  }

  ngOnInit(): void {
    this.initializeForms();
  }

  private initializeForms(): void {
    // Add Funds Form
    this.addFundsForm = this.fb.group({
      amount: [null, [
        Validators.required,
        Validators.min(this.MIN_AMOUNT),
        Validators.max(this.MAX_AMOUNT)
      ]],
      description: ['', [
        Validators.maxLength(200)
      ]],
      reference: ['', [
        Validators.maxLength(50)
      ]]
    });

    // Withdraw Funds Form
    this.withdrawFundsForm = this.fb.group({
      amount: [null, [
        Validators.required,
        Validators.min(this.MIN_AMOUNT),
        Validators.max(this.currentBalance()) // Can't withdraw more than balance
      ]],
      description: ['', [
        Validators.maxLength(200)
      ]],
      reference: ['', [
        Validators.maxLength(50)
      ]]
    });

    // Update withdraw form max validator when balance changes
    this.withdrawFundsForm.get('amount')?.setValidators([
      Validators.required,
      Validators.min(this.MIN_AMOUNT),
      Validators.max(this.currentBalance())
    ]);
  }

  onOperationChange(operation: 'add' | 'withdraw'): void {
    this.selectedOperation.set(operation);
  }

  onAddFunds(): void {
    if (this.addFundsForm.invalid) {
      this.markFormGroupTouched(this.addFundsForm);
      return;
    }

    const formValue = this.addFundsForm.value;
    const request: AddFundsRequest = {
      amount: formValue.amount,
      description: formValue.description || undefined,
      reference: formValue.reference || undefined
    };

    // Validate using the service
    const validation = this.walletService.validateWalletOperation(request);
    if (!validation.isValid) {
      this.snackBar.open(validation.errors[0]?.message || 'Datos inválidos', 'Cerrar', {
        duration: 3000
      });
      return;
    }

    this.loading.set(true);

    this.walletService.addFunds(this.data.account.id, request).subscribe({
      next: (balanceResponse) => {
        // Update the account object with the new balance
        const updatedAccount = { ...this.data.account, balance: balanceResponse.balance };
        this.currentAccount.set(updatedAccount); // Update the signal for immediate UI update
        this.snackBar.open(
          `Fondos agregados exitosamente. Nuevo saldo: ${this.validationService.formatCurrency(updatedAccount.balance, this.data.account.currency)}`,
          'Cerrar',
          { duration: 5000 }
        );
        this.dialogRef.close(updatedAccount);
      },
      error: (error) => {
        console.error('Error adding funds:', error);
        this.snackBar.open(error.message || 'Error al agregar fondos', 'Cerrar', {
          duration: 3000
        });
        this.loading.set(false);
      }
    });
  }

  onWithdrawFunds(): void {
    if (this.withdrawFundsForm.invalid) {
      this.markFormGroupTouched(this.withdrawFundsForm);
      return;
    }

    const formValue = this.withdrawFundsForm.value;
    const request: WithdrawFundsRequest = {
      amount: formValue.amount,
      description: formValue.description || undefined,
      reference: formValue.reference || undefined
    };

    // Validate using the service
    const validation = this.walletService.validateWalletOperation(request);
    if (!validation.isValid) {
      this.snackBar.open(validation.errors[0]?.message || 'Datos inválidos', 'Cerrar', {
        duration: 3000
      });
      return;
    }

    // Additional validation for withdrawal
    if (request.amount > this.currentBalance()) {
      this.snackBar.open('No puedes retirar más dinero del que tienes disponible', 'Cerrar', {
        duration: 3000
      });
      return;
    }

    this.loading.set(true);

    this.walletService.withdrawFunds(this.data.account.id, request).subscribe({
      next: (balanceResponse) => {
        // Update the account object with the new balance
        const updatedAccount = { ...this.data.account, balance: balanceResponse.balance };
        this.currentAccount.set(updatedAccount); // Update the signal for immediate UI update
        this.snackBar.open(
          `Fondos retirados exitosamente. Nuevo saldo: ${this.validationService.formatCurrency(updatedAccount.balance, this.data.account.currency)}`,
          'Cerrar',
          { duration: 5000 }
        );
        this.dialogRef.close(updatedAccount);
      },
      error: (error) => {
        console.error('Error withdrawing funds:', error);
        this.snackBar.open(error.message || 'Error al retirar fondos', 'Cerrar', {
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

  getAmountErrorMessage(form: FormGroup): string {
    const amountControl = form.get('amount');
    if (!amountControl) return '';

    if (amountControl.hasError('required')) {
      return 'El monto es requerido';
    }
    if (amountControl.hasError('min')) {
      return `El monto mínimo es $${this.MIN_AMOUNT.toLocaleString()}`;
    }
    if (amountControl.hasError('max')) {
      const maxValue = form === this.withdrawFundsForm ? this.currentBalance() : this.MAX_AMOUNT;
      return `El monto máximo es $${maxValue.toLocaleString()}`;
    }
    return '';
  }

  getDescriptionErrorMessage(): string {
    const addDescControl = this.addFundsForm?.get('description');
    const withdrawDescControl = this.withdrawFundsForm?.get('description');
    
    const control = this.selectedOperation() === 'add' ? addDescControl : withdrawDescControl;
    
    if (control?.hasError('maxlength')) {
      return 'La descripción no puede exceder 200 caracteres';
    }
    return '';
  }

  getReferenceErrorMessage(): string {
    const addRefControl = this.addFundsForm?.get('reference');
    const withdrawRefControl = this.withdrawFundsForm?.get('reference');
    
    const control = this.selectedOperation() === 'add' ? addRefControl : withdrawRefControl;
    
    if (control?.hasError('maxlength')) {
      return 'La referencia no puede exceder 50 caracteres';
    }
    return '';
  }

  // Quick amount buttons
  setQuickAmount(amount: number): void {
    const form = this.selectedOperation() === 'add' ? this.addFundsForm : this.withdrawFundsForm;
    form.get('amount')?.setValue(amount);
  }

  setMaxWithdrawAmount(): void {
    if (this.selectedOperation() === 'withdraw') {
      this.withdrawFundsForm.get('amount')?.setValue(this.currentBalance());
    }
  }

  // Utility methods for template
  getQuickAmounts(): number[] {
    if (this.selectedOperation() === 'add') {
      return [1000, 5000, 10000, 20000];
    } else {
      const balance = this.currentBalance();
      const amounts = [1000, 5000, 10000, 20000];
      return amounts.filter(amount => amount <= balance);
    }
  }

  canWithdraw(): boolean {
    return this.currentBalance() > 0;
  }

  getAccountDisplayName(): string {
    return this.currentAccount().name || `${this.validationService.getAccountTypeDisplayName(this.currentAccount().accountType)}`;
  }

  // Calculate estimated fees (if any)
  calculateEstimatedFee(amount: number): number {
    return this.walletService.calculateOperationFee(amount, 'withdrawal');
  }

  getEstimatedTotal(): number {
    const form = this.selectedOperation() === 'add' ? this.addFundsForm : this.withdrawFundsForm;
    const amount = form.get('amount')?.value || 0;
    const fee = this.calculateEstimatedFee(amount);
    
    if (this.selectedOperation() === 'add') {
      return amount; // No fees for adding funds typically
    } else {
      return amount + fee; // Include withdrawal fee
    }
  }
}