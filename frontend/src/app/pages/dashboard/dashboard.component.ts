import { Component, inject, OnInit, signal, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { AccountService } from '../../services/account.service';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { HasPermissionDirective } from '../../shared/directives/has-permission.directive';
import { HasRoleDirective } from '../../shared/directives/has-role.directive';
import { Permission, UserRole, Account, AccountType, Currency } from '../../models';
import { Subject, takeUntil } from 'rxjs';

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
    HasPermissionDirective,
    HasRoleDirective
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit, OnDestroy {
  private readonly auth = inject(AuthService);
  private readonly accountService = inject(AccountService);
  private readonly destroy$ = new Subject<void>();

  // Signals para datos del dashboard
  accounts = signal<Account[]>([]);
  loading = signal(false);
  totalWalletBalance = signal(0);
  totalWalletBalanceUSD = signal(0);
  totalCreditLimit = signal(0);
  activeAccountsCount = signal(0);
  transactionsCount = signal(0);

  // Exponemos los enums para usar en el template
  readonly Permission = Permission;
  readonly UserRole = UserRole;

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
    
    // Calcular saldo total de billeteras en ARS
    const walletBalanceARS = accounts
      .filter(account => {
        const isWallet = account.accountType === AccountType.WALLET;
        const isARS = account.currency === Currency.ARS;
        const isActive = account.isActive;
        console.log(`Account ${account.name}: type=${account.accountType}, currency=${account.currency}, active=${isActive}, balance=${account.balance}`);
        return isWallet && isARS && isActive;
      })
      .reduce((total, account) => {
        console.log(`Adding ARS wallet balance: ${account.balance}`);
        return total + (account.balance || 0);
      }, 0);
    
    // Calcular saldo total de billeteras en USD (sin conversiones)
    const walletBalanceUSD = accounts
      .filter(account => {
        const isWallet = account.accountType === AccountType.WALLET;
        const isUSD = account.currency === Currency.USD;
        const isActive = account.isActive;
        console.log(`Account ${account.name}: type=${account.accountType}, currency=${account.currency}, active=${isActive}, balance=${account.balance}`);
        return isWallet && isUSD && isActive;
      })
      .reduce((total, account) => {
        console.log(`Adding USD wallet balance: ${account.balance}`);
        return total + (account.balance || 0);
      }, 0);
    
    console.log('ARS wallet balance:', walletBalanceARS);
    console.log('USD wallet balance (direct only):', walletBalanceUSD);
    
    this.totalWalletBalance.set(walletBalanceARS);
    this.totalWalletBalanceUSD.set(walletBalanceUSD);

    // Calcular límite total de crédito
    const creditLimit = accounts
      .filter(account => {
        const isCredit = account.accountType === AccountType.CREDIT;
        const isActive = account.isActive;
        return isCredit && isActive;
      })
      .reduce((total, account) => total + (account.creditLimit || 0), 0);
    
    console.log('Total credit limit:', creditLimit);
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
}
