import { Component, inject, OnInit, signal, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { AccountService } from '../../services/account.service';
import { TransactionService } from '../../services/transaction.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatExpansionModule } from '@angular/material/expansion';
import { HasPermissionDirective } from '../../shared/directives/has-permission.directive';
import { HasRoleDirective } from '../../shared/directives/has-role.directive';
import { Permission, UserRole, Account, AccountType, Currency, Transaction, TransactionType } from '../../models';
import { AccountUtils } from '../../models/account-utils';
import { Subject, takeUntil, forkJoin } from 'rxjs';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule, 
    RouterModule,
    MatCardModule, 
    MatButtonModule, 
    MatIconModule,
    MatGridListModule,
    MatProgressSpinnerModule,
    MatExpansionModule,
    HasPermissionDirective,
    HasRoleDirective
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit, OnDestroy {
  private readonly auth = inject(AuthService);
  private readonly accountService = inject(AccountService);
  private readonly transactionService = inject(TransactionService);
  private readonly destroy$ = new Subject<void>();

  // Signals para datos del dashboard
  accounts = signal<Account[]>([]);
  recentTransactions = signal<Transaction[]>([]);
  loading = signal(false);
  loadingTransactions = signal(false);
  totalWalletBalance = signal(0);
  totalWalletBalanceUSD = signal(0);
  totalCreditLimit = signal(0);
  activeAccountsCount = signal(0);
  transactionsCount = signal(0);
  transactionsPanelExpanded = signal(true); // Panel expandido por defecto
  accountNamesMap = signal<Map<string, string>>(new Map()); // Mapeo de accountId -> nombre

  // Exponemos los enums para usar en el template
  readonly Permission = Permission;
  readonly UserRole = UserRole;
  readonly TransactionType = TransactionType;

  get currentUser() {
    return this.auth.currentUserSig;
  }

  ngOnInit(): void {
    // Cargar el perfil del usuario si no está disponible
    if (!this.currentUser()) {
      this.auth.loadUserProfile().subscribe({
        error: (err) => console.error('Error loading user profile:', err)
      });
    }
    
    // Cargar datos del dashboard
    this.loadDashboardData();
    this.loadRecentTransactions();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  private loadDashboardData(): void {
    const user = this.currentUser();
    console.log('Dashboard: Current user:', user);
    
    if (!user?.id) {
      console.log('Dashboard: No user ID available');
      return;
    }

    console.log('Dashboard: Loading data for user ID:', user.id);
    this.loading.set(true);
    
    this.accountService.getAccountsByUser(user.id).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (response) => {
        console.log('Dashboard: Loaded accounts:', response.accounts);
        this.accounts.set(response.accounts);
        this.calculateDashboardStats(response.accounts);
        this.loading.set(false);
      },
      error: (error) => {
        console.error('Error loading dashboard data:', error);
        this.loading.set(false);
        // En caso de error, resetear los valores
        this.resetDashboardStats();
      }
    });
  }

  private resetDashboardStats(): void {
    this.totalWalletBalance.set(0);
    this.totalWalletBalanceUSD.set(0);
    this.totalCreditLimit.set(0);
    this.activeAccountsCount.set(0);
    this.transactionsCount.set(0);
  }

  private calculateDashboardStats(accounts: Account[]): void {
    console.log('Dashboard: Calculating stats for accounts:', accounts);
    
    // Tasa de cambio aproximada ARS/USD (puedes obtenerla de un API en el futuro)
    const arsToUsd = 0.0011; // Ejemplo: 1 ARS = 0.0011 USD (aprox 900 ARS = 1 USD)
    
    // Calcular saldo total en ARS (todas las cuentas excepto CREDIT)
    const walletBalanceARS = accounts
      .filter(account => {
        const isNotCredit = account.accountType !== AccountType.CREDIT;
        const isARS = account.currency === Currency.ARS;
        const isActive = account.isActive;
        console.log(`Account ${account.name}: type=${account.accountType}, currency=${account.currency}, active=${isActive}, balance=${account.balance}`);
        return isNotCredit && isARS && isActive;
      })
      .reduce((total, account) => {
        console.log(`Adding ARS balance: ${account.balance}`);
        return total + (account.balance || 0);
      }, 0);
    
    // Calcular saldo total en USD (todas las cuentas excepto CREDIT)
    const walletBalanceUSD = accounts
      .filter(account => {
        const isNotCredit = account.accountType !== AccountType.CREDIT;
        const isUSD = account.currency === Currency.USD;
        const isActive = account.isActive;
        console.log(`Account ${account.name}: type=${account.accountType}, currency=${account.currency}, active=${isActive}, balance=${account.balance}`);
        return isNotCredit && isUSD && isActive;
      })
      .reduce((total, account) => {
        console.log(`Adding USD balance: ${account.balance}`);
        return total + (account.balance || 0);
      }, 0);
    
    console.log('Total ARS balance:', walletBalanceARS);
    console.log('Total USD balance:', walletBalanceUSD);
    
    this.totalWalletBalance.set(walletBalanceARS);
    this.totalWalletBalanceUSD.set(walletBalanceUSD);

    // Calcular límite total de crédito
    let creditLimit = 0;
    
    // Sumar límite de cuentas de tipo CREDIT (legacy)
    const legacyCreditLimit = accounts
      .filter(account => {
        const isCredit = account.accountType === AccountType.CREDIT;
        const isActive = account.isActive;
        console.log(`Credit account ${account.name}: isCredit=${isCredit}, isActive=${isActive}, creditLimit=${account.creditLimit}`);
        return isCredit && isActive;
      })
      .reduce((total, account) => total + (account.creditLimit || 0), 0);
    
    creditLimit += legacyCreditLimit;
    
    // Sumar límite de tarjetas de crédito en TODAS las cuentas activas con tarjetas
    // Incluye: BANK_ACCOUNT, CHECKING, y cualquier otro tipo con tarjetas
    accounts
      .filter(account => account.isActive && account.cards && account.cards.length > 0)
      .forEach(account => {
        account.cards?.forEach(card => {
          if (card.cardType === 'credit' && card.status === 'active' && card.creditLimit) {
            console.log(`Credit card ${card.lastFourDigits} in account ${account.name}: creditLimit=${card.creditLimit}`);
            creditLimit += card.creditLimit;
          }
        });
      });
    
    console.log('Total credit limit (legacy + cards):', creditLimit);
    this.totalCreditLimit.set(creditLimit);

    // Contar cuentas activas
    const activeAccounts = accounts.filter(account => account.isActive).length;
    console.log('Active accounts count:', activeAccounts);
    this.activeAccountsCount.set(activeAccounts);

    // Por ahora, transacciones simuladas
    this.transactionsCount.set(0);
    
    // Logging final para debug
    console.log('Dashboard: Stats updated - walletARS:', this.totalWalletBalance(), 'walletUSD:', this.totalWalletBalanceUSD(), 'credit:', this.totalCreditLimit(), 'accounts:', this.activeAccountsCount());
  }

  private loadRecentTransactions(): void {
    const user = this.currentUser();
    
    if (!user?.id) {
      console.log('Dashboard: No user ID available for transactions');
      return;
    }

    console.log('Dashboard: Loading recent transactions for user ID:', user.id);
    this.loadingTransactions.set(true);
    
    this.transactionService.getRecentTransactions(user.id, 10).pipe(
      takeUntil(this.destroy$)
    ).subscribe({
      next: (transactions) => {
        console.log('Dashboard: Loaded recent transactions:', transactions);
        this.recentTransactions.set(transactions);
        this.transactionsCount.set(transactions.length);
        this.loadingTransactions.set(false);
        this.loadAccountNamesForTransactions(transactions);
      },
      error: (error) => {
        console.error('Error loading recent transactions:', error);
        this.loadingTransactions.set(false);
        this.recentTransactions.set([]);
        this.transactionsCount.set(0);
      }
    });
  }

  private loadAccountNamesForTransactions(transactions: Transaction[]): void {
    // Obtener todos los IDs únicos de cuentas de las transacciones
    const accountIds = new Set<string>();
    transactions.forEach(tx => {
      if (tx.fromAccountId) accountIds.add(tx.fromAccountId);
      if (tx.toAccountId) accountIds.add(tx.toAccountId);
    });

    // Cargar nombres de las cuentas
    const accountsArray = Array.from(accountIds);
    if (accountsArray.length === 0) {
      return;
    }

    // Buscar en las cuentas ya cargadas primero
    const accountMap = new Map<string, string>();
    const currentAccounts = this.accounts();
    
    accountsArray.forEach(accountId => {
      const account = currentAccounts.find(acc => acc.id === accountId);
      if (account) {
        accountMap.set(accountId, account.name);
      }
    });

    this.accountNamesMap.set(accountMap);
  }

  getAccountNameForTransaction(transaction: Transaction): string {
    const accountMap = this.accountNamesMap();
    
    // Priorizar fromAccountId para transacciones salientes
    if (transaction.fromAccountId) {
      const name = accountMap.get(transaction.fromAccountId);
      if (name) return name;
    }
    
    // Usar toAccountId para depósitos
    if (transaction.toAccountId) {
      const name = accountMap.get(transaction.toAccountId);
      if (name) return name;
    }

    // Si hay metadata con nombre de cuenta
    if (transaction.metadata?.accountName) {
      return transaction.metadata.accountName;
    }
    
    // Fallback: mostrar tipo de transacción formateado
    return this.formatTransactionType(transaction.type);
  }

  private formatTransactionType(type: TransactionType): string {
    // Convertir el tipo de transacción a un formato más legible
    return type.replace(/_/g, ' ')
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
  }

  // Métodos para el template
  getTransactionIcon(type: TransactionType): string {
    return this.transactionService.getTransactionIcon(type);
  }

  getTransactionAmount(transaction: Transaction): number {
    // Para mostrar en el dashboard con signo correcto
    const isPositive = [
      TransactionType.DEPOSIT, 
      TransactionType.WALLET_DEPOSIT,
      TransactionType.ACCOUNT_DEPOSIT,
      TransactionType.REFUND, 
      TransactionType.CREDIT_REFUND,
      TransactionType.DEBIT_REFUND,
      TransactionType.SALARY
    ].includes(transaction.type);
    return isPositive ? transaction.amount : -transaction.amount;
  }

  // Método para obtener el icono de la cuenta asociada a la transacción
  getAccountIconForTransaction(transaction: Transaction): string {
    const accountId = transaction.fromAccountId || transaction.toAccountId;
    if (!accountId) return 'account_balance_wallet';
    
    const account = this.accounts().find(acc => acc.id === accountId);
    if (!account) return 'account_balance_wallet';
    
    return AccountUtils.getAccountIcon(account.accountType);
  }
}
