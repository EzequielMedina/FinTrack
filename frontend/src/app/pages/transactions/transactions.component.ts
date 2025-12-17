import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Subject, takeUntil, finalize } from 'rxjs';
import { MatIconModule } from '@angular/material/icon';

import { 
  TransactionService, 
  AccountService, 
  AuthService 
} from '../../services';

import {
  Transaction,
  TransactionType,
  TransactionStatus,
  TransactionFilterDTO,
  TransferRequest,
  DepositWithdrawalRequest,
  PaymentRequest,
  Account
} from '../../models';

@Component({
  selector: 'app-transactions',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, MatIconModule],
  templateUrl: './transactions.component.html',
  styleUrls: ['./transactions.component.css']
})
export class TransactionsComponent implements OnInit, OnDestroy {
  private destroy$ = new Subject<void>();

  // Estado del componente
  loading = false;
  transactions: Transaction[] = [];
  accounts: Account[] = [];
  currentUser: any;
  
  // Filtros y paginación
  filters: TransactionFilterDTO = {
    limit: 10,
    offset: 0
  };
  totalTransactions = 0;
  currentPage = 1;

  // Formularios
  transferForm!: FormGroup;
  depositForm!: FormGroup;
  withdrawalForm!: FormGroup;
  paymentForm!: FormGroup;

  // UI State
  activeTab = 'list'; // 'list', 'transfer', 'deposit', 'withdrawal', 'payment'
  showFilters = false;

  // Enums para template
  TransactionType = TransactionType;
  TransactionStatus = TransactionStatus;

  constructor(
    private transactionService: TransactionService,
    private accountService: AccountService,
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private router: Router
  ) {
    this.initializeForms();
  }

  ngOnInit() {
    this.loadCurrentUser();
    this.loadAccounts();
    this.loadTransactions();
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // ================================
  // INICIALIZACIÓN
  // ================================

  private initializeForms() {
    this.transferForm = this.formBuilder.group({
      fromAccountId: ['', Validators.required],
      toAccountId: ['', Validators.required],
      amount: ['', [Validators.required, Validators.min(0.01)]],
      description: ['']
    });

    this.depositForm = this.formBuilder.group({
      accountId: ['', Validators.required],
      amount: ['', [Validators.required, Validators.min(0.01)]],
      description: [''],
      method: ['BANK_TRANSFER']
    });

    this.withdrawalForm = this.formBuilder.group({
      accountId: ['', Validators.required],
      amount: ['', [Validators.required, Validators.min(0.01)]],
      description: [''],
      method: ['ATM']
    });

    this.paymentForm = this.formBuilder.group({
      accountId: ['', Validators.required],
      amount: ['', [Validators.required, Validators.min(0.01)]],
      merchantName: [''],
      description: [''],
      cardId: ['']
    });
  }

  private loadCurrentUser() {
    const user = this.authService.getCurrentUser();
    if (user) {
      this.currentUser = user;
    } else {
      this.router.navigate(['/login']);
    }
  }

  private loadAccounts() {
    if (!this.currentUser) return;

    this.accountService.getAccountsByUser(this.currentUser.id)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (response: any) => {
          this.accounts = response.accounts || response;
        },
        error: (error: any) => {
          console.error('Error loading accounts:', error);
        }
      });
  }

  private loadTransactions() {
    if (!this.currentUser) return;

    this.loading = true;
    this.transactionService.getUserTransactions(this.currentUser.id, this.filters)
      .pipe(
        takeUntil(this.destroy$),
        finalize(() => this.loading = false)
      )
      .subscribe({
        next: (response) => {
          this.transactions = response.transactions;
          this.totalTransactions = response.total;
        },
        error: (error) => {
          console.error('Error loading transactions:', error);
        }
      });
  }

  // ================================
  // OPERACIONES DE TRANSACCIÓN
  // ================================

  onTransfer() {
    if (this.transferForm.valid) {
      const transferData: TransferRequest = this.transferForm.value;
      
      this.loading = true;
      this.transactionService.transfer(transferData)
        .pipe(
          takeUntil(this.destroy$),
          finalize(() => this.loading = false)
        )
        .subscribe({
          next: (response) => {
            if (response.success) {
              this.transferForm.reset();
              this.activeTab = 'list';
              this.loadTransactions();
              this.loadAccounts(); // Actualizar balances
            }
          },
          error: (error) => {
            console.error('Transfer error:', error);
          }
        });
    }
  }

  onDeposit() {
    if (this.depositForm.valid) {
      const depositData: DepositWithdrawalRequest = this.depositForm.value;
      
      this.loading = true;
      this.transactionService.deposit(depositData)
        .pipe(
          takeUntil(this.destroy$),
          finalize(() => this.loading = false)
        )
        .subscribe({
          next: (response) => {
            if (response.success) {
              this.depositForm.reset();
              this.activeTab = 'list';
              this.loadTransactions();
              this.loadAccounts(); // Actualizar balances
            }
          },
          error: (error) => {
            console.error('Deposit error:', error);
          }
        });
    }
  }

  onWithdrawal() {
    if (this.withdrawalForm.valid) {
      const withdrawalData: DepositWithdrawalRequest = this.withdrawalForm.value;
      
      this.loading = true;
      this.transactionService.withdraw(withdrawalData)
        .pipe(
          takeUntil(this.destroy$),
          finalize(() => this.loading = false)
        )
        .subscribe({
          next: (response) => {
            if (response.success) {
              this.withdrawalForm.reset();
              this.activeTab = 'list';
              this.loadTransactions();
              this.loadAccounts(); // Actualizar balances
            }
          },
          error: (error) => {
            console.error('Withdrawal error:', error);
          }
        });
    }
  }

  onPayment() {
    if (this.paymentForm.valid) {
      const paymentData: PaymentRequest = this.paymentForm.value;
      
      this.loading = true;
      this.transactionService.makePayment(paymentData)
        .pipe(
          takeUntil(this.destroy$),
          finalize(() => this.loading = false)
        )
        .subscribe({
          next: (response) => {
            if (response.success) {
              this.paymentForm.reset();
              this.activeTab = 'list';
              this.loadTransactions();
              this.loadAccounts(); // Actualizar balances
            }
          },
          error: (error) => {
            console.error('Payment error:', error);
          }
        });
    }
  }

  // ================================
  // FILTROS Y PAGINACIÓN
  // ================================

  onFilterChange() {
    this.filters.offset = 0;
    this.currentPage = 1;
    this.loadTransactions();
  }

  onPageChange(page: number) {
    this.currentPage = page;
    this.filters.offset = (page - 1) * (this.filters.limit || 10);
    this.loadTransactions();
  }

  toggleFilters() {
    this.showFilters = !this.showFilters;
  }

  clearFilters() {
    this.filters = {
      limit: 10,
      offset: 0
    };
    this.currentPage = 1;
    this.loadTransactions();
  }

  // ================================
  // UTILIDADES PARA TEMPLATE
  // ================================

  getStatusClass(status: TransactionStatus): string {
    return `status-${this.transactionService.getStatusColor(status)}`;
  }

  getTransactionIcon(type: TransactionType): string {
    return this.transactionService.getTransactionIcon(type);
  }

  getAccountName(accountId: string): string {
    const account = this.accounts.find(acc => acc.id === accountId);
    return account ? `${account.accountType} - ${account.name}` : accountId;
  }

  canCancelTransaction(transaction: Transaction): boolean {
    return this.transactionService.canCancelTransaction(transaction);
  }

  canRefundTransaction(transaction: Transaction): boolean {
    return this.transactionService.canRefundTransaction(transaction);
  }

  // ================================
  // ACCIONES DE TRANSACCIÓN
  // ================================

  cancelTransaction(transaction: Transaction) {
    if (confirm('¿Está seguro de que desea cancelar esta transacción?')) {
      this.transactionService.cancelTransaction(transaction.id, 'Cancelled by user')
        .pipe(takeUntil(this.destroy$))
        .subscribe({
          next: (response) => {
            if (response.success) {
              this.loadTransactions();
            }
          },
          error: (error) => {
            console.error('Cancel transaction error:', error);
          }
        });
    }
  }

  viewTransactionDetails(transaction: Transaction) {
    // Implementar modal o navegación a vista de detalles
    console.log('View details for transaction:', transaction.id);
  }

  // ================================
  // NAVEGACIÓN
  // ================================

  setActiveTab(tab: string) {
    this.activeTab = tab;
  }

  getTotalPages(): number {
    return Math.ceil(this.totalTransactions / (this.filters.limit || 10));
  }

  translateTransactionDescription(description: string | undefined): string {
    if (!description) return '';
    
    // Traducir términos relacionados con installments
    return description
      .replace(/installment/gi, 'cuota')
      .replace(/Installment/gi, 'Cuota')
      .replace(/INSTALLMENT/gi, 'CUOTA')
      .replace(/installments/gi, 'cuotas')
      .replace(/Installments/gi, 'Cuotas')
      .replace(/INSTALLMENTS/gi, 'CUOTAS')
      .replace(/payment/gi, 'pago')
      .replace(/Payment/gi, 'Pago')
      .replace(/PAYMENT/gi, 'PAGO')
      .replace(/plan/gi, 'plan')
      .replace(/Plan/gi, 'Plan')
      .replace(/PLAN/gi, 'PLAN');
  }

  getTransactionIconClass(type: TransactionType): string {
    return 'icon-type-' + type.toLowerCase().replace(/_/g, '-');
  }
}