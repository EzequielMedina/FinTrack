import { Component, Inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef, MatDialogModule } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatDividerModule } from '@angular/material/divider';
import { MatChipsModule } from '@angular/material/chips';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { InstallmentService } from '../../../services/installment.service';
import { AccountService } from '../../../services/account.service';
import { AuthService } from '../../../services/auth.service';
import { InstallmentPlan, Installment, PayInstallmentRequest, InstallmentStatus } from '../../../models/installment.model';

export interface InstallmentPaymentModalData {
  installmentPlan: InstallmentPlan;
}

@Component({
  selector: 'app-installment-payment-modal',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatSelectModule,
    MatInputModule,
    MatCheckboxModule,
    MatProgressBarModule,
    MatSnackBarModule,
    MatDividerModule,
    MatChipsModule,
    MatProgressSpinnerModule
  ],
  templateUrl: './installment-payment-modal.component.html',
  styleUrls: ['./installment-payment-modal.component.scss']
})
export class InstallmentPaymentModalComponent implements OnInit {
  paymentForm: FormGroup;
  installments: Installment[] = [];
  selectedInstallments: Set<string> = new Set();
  isLoading = false;
  isProcessingPayment = false;
  isLoadingAccounts = false;
  
  totalToPay = 0;
  nextDueInstallment: Installment | null = null;
  overdueInstallments: Installment[] = [];
  
  // Account management
  availableAccounts: any[] = [];
  selectedPaymentMethod: any = null;
  
  // Make enum available in template
  readonly InstallmentStatus = InstallmentStatus;
  
  readonly paymentMethods = [
    { 
      value: 'credit_card', 
      label: '💳 Tarjeta de Crédito',
      accountTypes: ['credit'],
      requiresBalance: false,
      description: 'Pago en una sola cuota'
    },
    { 
      value: 'debit_card', 
      label: '💳 Tarjeta de Débito',
      accountTypes: ['debit', 'checking', 'savings', 'bank_account'],
      requiresBalance: true,
      description: 'Pago con saldo disponible'
    },
    { 
      value: 'bank_transfer', 
      label: '🏦 Transferencia Bancaria',
      accountTypes: ['checking', 'savings', 'bank_account'],
      requiresBalance: true,
      description: 'Transferencia desde cuenta bancaria'
    },
    { 
      value: 'wallet', 
      label: '👛 Billetera Digital',
      accountTypes: ['wallet'],
      requiresBalance: true,
      description: 'Pago desde billetera virtual'
    }
  ];

  constructor(
    private fb: FormBuilder,
    private installmentService: InstallmentService,
    private accountService: AccountService,
    private authService: AuthService,
    private snackBar: MatSnackBar,
    public dialogRef: MatDialogRef<InstallmentPaymentModalComponent>,
    @Inject(MAT_DIALOG_DATA) public data: InstallmentPaymentModalData
  ) {
    this.paymentForm = this.fb.group({
      paymentMethod: ['', Validators.required],
      selectedAccount: ['', Validators.required],
      paymentReference: ['', Validators.required],
      notes: ['']
    });
    
    // Listen for payment method changes
    this.paymentForm.get('paymentMethod')?.valueChanges.subscribe(method => {
      this.onPaymentMethodChange(method);
    });
  }

  ngOnInit(): void {
    this.loadInstallments();
  }

  private loadInstallments(): void {
    this.isLoading = true;
    
    this.installmentService.getInstallmentsByPlan(this.data.installmentPlan.id).subscribe({
      next: (installments: Installment[]) => {
        this.installments = installments;
        this.analyzeInstallments();
        this.isLoading = false;
      },
      error: (error: any) => {
        console.error('❌ Error loading installments:', error);
        this.snackBar.open('Error al cargar las cuotas', 'Cerrar', { duration: 3000 });
        this.isLoading = false;
      }
    });
  }

  private analyzeInstallments(): void {
    const now = new Date();
    const pendingInstallments = this.installments.filter(inst => 
      inst.status === InstallmentStatus.PENDING || inst.status === InstallmentStatus.PARTIAL
    );

    // Find overdue installments
    this.overdueInstallments = pendingInstallments.filter(inst => 
      new Date(inst.due_date) < now
    ).sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime());

    // Find next due installment
    const upcomingInstallments = pendingInstallments.filter(inst => 
      new Date(inst.due_date) >= now
    ).sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime());

    this.nextDueInstallment = upcomingInstallments.length > 0 ? upcomingInstallments[0] : null;

    // Auto-select overdue installments and next due
    if (this.overdueInstallments.length > 0) {
      this.overdueInstallments.forEach(inst => this.selectedInstallments.add(inst.id));
    } else if (this.nextDueInstallment) {
      this.selectedInstallments.add(this.nextDueInstallment.id);
    }

    this.calculateTotal();
  }

  onPaymentMethodChange(methodValue: string): void {
    console.log('🔄 Payment method changed to:', methodValue);
    console.log('🔍 Available payment methods:', this.paymentMethods);
    
    if (!methodValue) {
      console.log('❌ No method value, clearing accounts');
      this.availableAccounts = [];
      this.selectedPaymentMethod = null;
      this.paymentForm.get('selectedAccount')?.setValue('');
      return;
    }

    this.selectedPaymentMethod = this.paymentMethods.find(m => m.value === methodValue);
    console.log('📱 Selected payment method:', this.selectedPaymentMethod);
    console.log('✅ Form control value:', this.paymentForm.get('paymentMethod')?.value);
    
    if (this.selectedPaymentMethod) {
      console.log('🔍 Loading accounts for types:', this.selectedPaymentMethod.accountTypes);
      this.loadCompatibleAccounts(this.selectedPaymentMethod.accountTypes);
    } else {
      console.log('❌ Payment method not found in array');
    }
  }

  private loadCompatibleAccounts(accountTypes: string[]): void {
    const currentUser = this.authService.getCurrentUser();
    console.log('👤 Current user:', currentUser);
    
    if (!currentUser) {
      this.snackBar.open('Usuario no autenticado', 'Cerrar', { duration: 3000 });
      return;
    }

    console.log('⏳ Loading accounts for user:', currentUser.id, 'with types:', accountTypes);
    this.isLoadingAccounts = true;
    
    this.accountService.getAccountsByUserAndTypes(currentUser.id, accountTypes).subscribe({
      next: (accounts: any[]) => {
        console.log('✅ Accounts loaded:', accounts);
        this.availableAccounts = accounts.filter(account => account.isActive);
        console.log('🔄 Filtered active accounts:', this.availableAccounts);
        this.isLoadingAccounts = false;
        
        // Auto-select first account if only one available
        if (this.availableAccounts.length === 1) {
          this.paymentForm.get('selectedAccount')?.setValue(this.availableAccounts[0].id);
        }
      },
      error: (error: any) => {
        console.error('❌ Error loading accounts:', error);
        this.availableAccounts = [];
        this.isLoadingAccounts = false;
        this.snackBar.open('Error al cargar las cuentas disponibles', 'Cerrar', { duration: 3000 });
      }
    });
  }

  getSelectedAccount(): any {
    const selectedAccountId = this.paymentForm.get('selectedAccount')?.value;
    return this.availableAccounts.find(acc => acc.id === selectedAccountId);
  }

  canProcessPayment(): boolean {
    if (!this.selectedPaymentMethod || this.selectedInstallments.size === 0) {
      return false;
    }

    const selectedAccount = this.getSelectedAccount();
    if (!selectedAccount) {
      return false;
    }

    // For credit cards, check available credit limit
    if (this.selectedPaymentMethod.value === 'credit_card') {
      const availableCredit = (selectedAccount.creditLimit || 0) - (selectedAccount.balance || 0);
      return availableCredit >= this.totalToPay;
    }

    // For other methods, check available balance
    if (this.selectedPaymentMethod.requiresBalance) {
      return (selectedAccount.balance || 0) >= this.totalToPay;
    }

    return true;
  }

  getPaymentValidationMessage(): string {
    if (!this.selectedPaymentMethod || this.selectedInstallments.size === 0) {
      return '';
    }

    const selectedAccount = this.getSelectedAccount();
    if (!selectedAccount) {
      return 'Seleccione una cuenta para continuar';
    }

    if (this.selectedPaymentMethod.value === 'credit_card') {
      const availableCredit = (selectedAccount.creditLimit || 0) - (selectedAccount.balance || 0);
      if (availableCredit < this.totalToPay) {
        return `Límite de crédito insuficiente. Disponible: $${availableCredit.toFixed(2)}`;
      }
    } else if (this.selectedPaymentMethod.requiresBalance) {
      if ((selectedAccount.balance || 0) < this.totalToPay) {
        return `Saldo insuficiente. Disponible: $${(selectedAccount.balance || 0).toFixed(2)}`;
      }
    }

    return '';
  }

  toggleInstallmentSelection(installment: Installment): void {
    if (this.selectedInstallments.has(installment.id)) {
      this.selectedInstallments.delete(installment.id);
    } else {
      this.selectedInstallments.add(installment.id);
    }
    this.calculateTotal();
  }

  private calculateTotal(): void {
    this.totalToPay = this.installments
      .filter(inst => this.selectedInstallments.has(inst.id))
      .reduce((sum, inst) => sum + inst.remaining_amount, 0);
  }

  isInstallmentSelected(installmentId: string): boolean {
    return this.selectedInstallments.has(installmentId);
  }

  canSelectInstallment(installment: Installment): boolean {
    return installment.status === InstallmentStatus.PENDING || installment.status === InstallmentStatus.PARTIAL;
  }

  getInstallmentStatusColor(status: string): string {
    switch (status) {
      case InstallmentStatus.PAID: return 'success';
      case InstallmentStatus.PENDING: return 'warn';
      case InstallmentStatus.PARTIAL: return 'accent';
      case InstallmentStatus.OVERDUE: return 'warn';
      default: return 'primary';
    }
  }

  getInstallmentStatusText(status: string): string {
    switch (status) {
      case InstallmentStatus.PAID: return 'Pagado';
      case InstallmentStatus.PENDING: return 'Pendiente';
      case InstallmentStatus.PARTIAL: return 'Parcial';
      case InstallmentStatus.OVERDUE: return 'Vencido';
      case InstallmentStatus.CANCELLED: return 'Cancelado';
      default: return status;
    }
  }

  getInstallmentStatusIcon(installment: Installment): string {
    const now = new Date();
    const dueDate = new Date(installment.due_date);
    
    switch (installment.status) {
      case InstallmentStatus.PAID: return 'check_circle';
      case InstallmentStatus.PARTIAL: return 'hourglass_half';
      case InstallmentStatus.PENDING: 
        return dueDate < now ? 'warning' : 'schedule';
      default: return 'schedule';
    }
  }

  isInstallmentOverdue(installment: Installment): boolean {
    return installment.status !== InstallmentStatus.PAID && new Date(installment.due_date) < new Date();
  }

  selectNextInstallment(): void {
    if (this.nextDueInstallment) {
      this.selectedInstallments.clear();
      this.selectedInstallments.add(this.nextDueInstallment.id);
      this.calculateTotal();
    }
  }

  selectOverdueInstallments(): void {
    this.selectedInstallments.clear();
    this.overdueInstallments.forEach(inst => this.selectedInstallments.add(inst.id));
    this.calculateTotal();
  }

  selectAllPending(): void {
    this.selectedInstallments.clear();
    this.installments
      .filter(inst => this.canSelectInstallment(inst))
      .forEach(inst => this.selectedInstallments.add(inst.id));
    this.calculateTotal();
  }

  clearSelection(): void {
    this.selectedInstallments.clear();
    this.calculateTotal();
  }

  processPayment(): void {
    if (this.paymentForm.invalid || this.selectedInstallments.size === 0) {
      this.snackBar.open('Por favor complete todos los campos y seleccione al menos una cuota', 'Cerrar', { duration: 3000 });
      return;
    }

    // Validate payment capability
    if (!this.canProcessPayment()) {
      const validationMessage = this.getPaymentValidationMessage();
      this.snackBar.open(validationMessage || 'No se puede procesar el pago', 'Cerrar', { duration: 5000 });
      return;
    }

    this.isProcessingPayment = true;
    const formValue = this.paymentForm.value;
    const selectedAccount = this.getSelectedAccount();
    
    // Validate that an account is selected
    if (!selectedAccount) {
      this.snackBar.open('⚠️ Debe seleccionar una cuenta para el pago', 'Cerrar', { 
        duration: 3000,
        panelClass: ['error-snackbar']
      });
      this.isProcessingPayment = false;
      return;
    }
    
    // Process each selected installment
    const paymentPromises = Array.from(this.selectedInstallments).map(installmentId => {
      const installment = this.installments.find(inst => inst.id === installmentId);
      if (!installment) return Promise.reject(new Error('Installment not found'));

      const paymentRequest: PayInstallmentRequest = {
        installment_id: installmentId,
        amount: installment.remaining_amount,
        payment_method: formValue.paymentMethod,
        payment_reference: `${formValue.paymentReference}-${installment.installment_number}`,
        notes: formValue.notes || '',
        // Add account information - use string values directly
        account_id: selectedAccount.id,
        account_type: selectedAccount.accountType.toString() // Ensure string conversion
      };

      console.log('🔍 Sending payment request:', paymentRequest);
      console.log('🔍 Selected account details:', selectedAccount);

      return this.installmentService.payInstallment(paymentRequest).toPromise();
    });

    Promise.all(paymentPromises)
      .then(() => {
        this.snackBar.open(
          `✅ ${this.selectedInstallments.size} cuota(s) pagada(s) exitosamente`, 
          'Cerrar', 
          { duration: 5000 }
        );
        this.dialogRef.close({ success: true, paidCount: this.selectedInstallments.size });
      })
      .catch((error) => {
        console.error('❌ Error processing payments:', error);
        this.snackBar.open('Error al procesar los pagos. Inténtelo nuevamente.', 'Cerrar', { duration: 5000 });
      })
      .finally(() => {
        this.isProcessingPayment = false;
      });
  }

  onCancel(): void {
    this.dialogRef.close({ success: false });
  }

  getProgressPercentage(): number {
    const totalInstallments = this.installments.length;
    const paidInstallments = this.installments.filter(inst => inst.status === InstallmentStatus.PAID).length;
    return totalInstallments > 0 ? (paidInstallments / totalInstallments) * 100 : 0;
  }

  getPaidAmount(): number {
    return this.installments
      .filter(inst => inst.status === InstallmentStatus.PAID)
      .reduce((sum, inst) => sum + inst.amount, 0);
  }

  getRemainingAmount(): number {
    return this.installments
      .filter(inst => inst.status !== InstallmentStatus.PAID)
      .reduce((sum, inst) => sum + inst.remaining_amount, 0);
  }

  trackByInstallmentId(index: number, installment: Installment): string {
    return installment.id;
  }
}