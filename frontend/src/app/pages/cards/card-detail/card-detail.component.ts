import { ChangeDetectionStrategy, Component, Input, OnInit, signal, inject, ViewChild } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatTabsModule } from '@angular/material/tabs';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Card, CardType, InstallmentPlan } from '../../../models';
import { CreditCardService, CreditCardBalanceResponse } from '../../../services/credit-card.service';
import { DebitCardService, DebitCardBalanceResponse } from '../../../services/debit-card.service';
import { 
  InstallmentCalculatorComponent, 
  InstallmentPlansListComponent
} from '../../../shared/components';
import type { 
  InstallmentCalculatorResult,
  InstallmentPlanAction 
} from '../../../shared/components';

@Component({
  selector: 'app-card-detail',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatChipsModule,
    MatDividerModule,
    MatTabsModule,
    InstallmentCalculatorComponent,
    InstallmentPlansListComponent
  ],
  templateUrl: './card-detail.component.html',
  styleUrls: ['./card-detail.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CardDetailComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly snackBar = inject(MatSnackBar);
  private readonly creditCardService = inject(CreditCardService);
  private readonly debitCardService = inject(DebitCardService);

  @Input() card!: Card;

  // Signals for reactive data
  creditBalance = signal<CreditCardBalanceResponse | null>(null);
  debitBalance = signal<DebitCardBalanceResponse | null>(null);
  loading = signal(false);

  // Installment properties
  installmentPlansCount = signal(0);
  currentInstallmentCalculation = signal<InstallmentCalculatorResult | null>(null);

  // ViewChild for installment calculator and plans list
  @ViewChild('installmentCalculator') installmentCalculator!: InstallmentCalculatorComponent;
  @ViewChild('installmentPlansList') installmentPlansList!: InstallmentPlansListComponent;

  // Form groups
  chargeForm!: FormGroup;
  paymentForm!: FormGroup;
  transactionForm!: FormGroup;

  // Exponer enum para el template
  readonly CardType = CardType;

  ngOnInit(): void {
    this.initializeForms();
    this.loadCardBalance();
  }

  private initializeForms(): void {
    // Formulario para cargos de tarjeta de crédito
    this.chargeForm = this.fb.group({
      amount: [null, [Validators.required, Validators.min(0.01), Validators.max(999999)]],
      description: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(255)]],
      reference: ['', [Validators.maxLength(50)]]
    });

    // Formulario para pagos de tarjeta de crédito
    this.paymentForm = this.fb.group({
      amount: [null, [Validators.required, Validators.min(0.01), Validators.max(999999)]],
      paymentMethod: ['bank_transfer', [Validators.required]],
      reference: ['', [Validators.maxLength(50)]]
    });

    // Formulario para transacciones de tarjeta de débito
    this.transactionForm = this.fb.group({
      amount: [null, [Validators.required, Validators.min(0.01), Validators.max(999999)]],
      description: ['', [Validators.required, Validators.minLength(3), Validators.maxLength(255)]],
      merchantName: ['', [Validators.maxLength(100)]],
      reference: ['', [Validators.maxLength(50)]]
    });
  }

  private loadCardBalance(): void {
    this.loading.set(true);
    
    if (this.card.cardType === CardType.CREDIT) {
      this.creditCardService.getBalance(this.card.id).subscribe({
        next: (balance) => {
          this.creditBalance.set(balance);
          this.loading.set(false);
        },
        error: (error) => {
          console.error('Error loading credit card balance:', error);
          this.snackBar.open(error.message, 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        }
      });
    } else {
      this.debitCardService.getAvailableBalance(this.card.id).subscribe({
        next: (balance) => {
          this.debitBalance.set(balance);
          this.loading.set(false);
        },
        error: (error) => {
          console.error('Error loading debit card balance:', error);
          this.snackBar.open(error.message, 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        }
      });
    }
  }

  // OPERACIONES DE TARJETA DE CRÉDITO

  onCharge(): void {
    if (this.chargeForm.valid && this.card.cardType === CardType.CREDIT) {
      this.loading.set(true);
      
      const chargeData = this.chargeForm.value;
      this.creditCardService.charge(this.card.id, chargeData).subscribe({
        next: (response) => {
          // Check if it's a regular charge (CreditCardBalanceResponse) or installment charge
          if ('cardId' in response && 'balance' in response) {
            // Regular charge response
            this.creditBalance.set(response as CreditCardBalanceResponse);
          } else {
            // Installment charge response - reload balance
            this.loadCardBalance();
          }
          this.chargeForm.reset();
          this.snackBar.open('Cargo realizado exitosamente', 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        },
        error: (error) => {
          console.error('Error processing charge:', error);
          this.snackBar.open(error.message, 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        }
      });
    }
  }

  onPayment(): void {
    if (this.paymentForm.valid && this.card.cardType === CardType.CREDIT) {
      this.loading.set(true);
      
      const paymentData = {
        ...this.paymentForm.value,
        paymentMethod: this.paymentForm.value.paymentMethod
      };

      this.creditCardService.payment(this.card.id, paymentData).subscribe({
        next: (balance) => {
          this.creditBalance.set(balance);
          this.paymentForm.reset();
          this.paymentForm.patchValue({ paymentMethod: 'bank_transfer' });
          this.snackBar.open('Pago realizado exitosamente', 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        },
        error: (error) => {
          console.error('Error processing payment:', error);
          this.snackBar.open(error.message, 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        }
      });
    }
  }

  // OPERACIONES DE TARJETA DE DÉBITO

  onTransaction(): void {
    if (this.transactionForm.valid && this.card.cardType === CardType.DEBIT) {
      this.loading.set(true);
      
      const transactionData = this.transactionForm.value;
      this.debitCardService.processTransaction(this.card.id, transactionData).subscribe({
        next: (balance) => {
          this.debitBalance.set(balance);
          this.transactionForm.reset();
          this.snackBar.open('Transacción realizada exitosamente', 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        },
        error: (error) => {
          console.error('Error processing transaction:', error);
          this.snackBar.open(error.message, 'Cerrar', { duration: 3000 });
          this.loading.set(false);
        }
      });
    }
  }

  // MÉTODOS DE UTILIDAD

  onRefreshBalance(): void {
    this.loadCardBalance();
  }

  formatCurrency(amount: number): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS',
      minimumFractionDigits: 2
    }).format(amount);
  }

  getAvailableCreditFormatted(): string {
    const balance = this.creditBalance();
    return balance ? this.formatCurrency(balance.availableCredit) : '$0.00';
  }

  getCurrentDebtFormatted(): string {
    const balance = this.creditBalance();
    return balance ? this.formatCurrency(balance.balance) : '$0.00';
  }

  getMinimumPaymentFormatted(): string {
    const balance = this.creditBalance();
    return balance ? this.formatCurrency(balance.minimumPayment) : '$0.00';
  }

  getAccountBalanceFormatted(): string {
    const balance = this.debitBalance();
    return balance ? this.formatCurrency(balance.accountBalance) : '$0.00';
  }

  isFormDisabled(): boolean {
    return this.loading() || !this.card;
  }

  payMinimumAmount(): void {
    const balance = this.creditBalance();
    if (balance && balance.minimumPayment > 0) {
      this.paymentForm.patchValue({
        amount: balance.minimumPayment
      });
    }
  }

  payFullBalance(): void {
    const balance = this.creditBalance();
    if (balance && balance.balance > 0) {
      this.paymentForm.patchValue({
        amount: balance.balance
      });
    }
  }

  // Installment Methods
  onInstallmentCalculationChanged(calculation: InstallmentCalculatorResult | null): void {
    this.currentInstallmentCalculation.set(calculation);
  }

  onInstallmentsSelected(calculation: InstallmentCalculatorResult): void {
    if (calculation && this.card.cardType === CardType.CREDIT) {
      // Crear la compra con cuotas
      const chargeRequest = {
        amount: calculation.preview.totalAmount, // Base amount without interests
        description: calculation.description || 'Compra en cuotas', // Usar descripción personalizada o por defecto
        installments: {
          count: calculation.installmentsCount,
          interestRate: calculation.preview.interestRate || 0,
          adminFee: calculation.preview.adminFee || 0,
          startDate: calculation.startDate
        }
      };

      this.loading.set(true);
      this.creditCardService.charge(this.card.id, chargeRequest).subscribe({
        next: (response) => {
          this.loading.set(false);
          this.snackBar.open(
            `Compra en ${calculation.installmentsCount} cuotas creada exitosamente`,
            'Cerrar',
            { duration: 5000 }
          );
          // Recargar balance y planes
          this.loadCardBalance();
          
          // Reset the installment calculator after successful creation
          if (this.installmentCalculator) {
            console.log('🔄 Calling resetCalculator() on installment calculator');
            this.installmentCalculator.resetCalculator();
            console.log('✅ resetCalculator() called successfully');
          } else {
            console.warn('⚠️ installmentCalculator ViewChild not found');
          }
          
          // Refresh the installment plans list to show the new plan
          if (this.installmentPlansList) {
            this.installmentPlansList.refresh();
          }
        },
        error: (error) => {
          this.loading.set(false);
          this.snackBar.open(
            `Error al crear la compra en cuotas: ${error.message}`,
            'Cerrar',
            { duration: 5000 }
          );
        }
      });
    }
  }

  onInstallmentPlanAction(action: InstallmentPlanAction): void {
    switch (action.type) {
      case 'view':
        this.viewInstallmentPlanDetail(action.plan);
        break;
      case 'pay':
        this.showPayInstallmentDialog(action.plan);
        break;
      case 'cancel':
        this.cancelInstallmentPlan(action.plan);
        break;
    }
  }

  onInstallmentPlansLoaded(plans: InstallmentPlan[]): void {
    this.installmentPlansCount.set(plans.length);
  }

  onRefreshInstallmentPlans(): void {
    if (this.installmentPlansList) {
      this.installmentPlansList.refresh();
    }
  }

  viewAllInstallmentPlans(): void {
    // Aquí podrías navegar a una página dedicada o abrir un modal
    console.log('Navigate to all installment plans view');
  }

  viewInstallmentPlanDetail(plan: InstallmentPlan): void {
    // Aquí podrías abrir un modal con el detalle del plan
    console.log('View plan detail:', plan);
  }

  showPayInstallmentDialog(plan: InstallmentPlan): void {
    // Aquí podrías abrir un dialog para pagar la cuota
    console.log('Pay installment for plan:', plan);
  }

  cancelInstallmentPlan(plan: InstallmentPlan): void {
    // Aquí podrías mostrar confirmación y cancelar el plan
    console.log('Cancel plan:', plan);
  }
}