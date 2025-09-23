import { Component, inject, OnInit, OnDestroy, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { Subject, takeUntil, finalize } from 'rxjs';
import { AccountService, AuthService } from '../../services';
import { Account, AccountType, Currency, AccountsListResponse, User } from '../../models';
import { AccountListComponent } from './account-list/account-list.component';
import { WalletDialogComponent } from './wallet-dialog/wallet-dialog.component';
import { CreditDialogComponent } from './credit-dialog/credit-dialog.component';
import { AccountFormComponent, AccountFormDialogData } from './account-form/account-form.component';
import { AccountDeleteConfirmationComponent } from './account-delete-confirmation/account-delete-confirmation.component';

@Component({
  selector: 'app-accounts',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatTabsModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatDialogModule,
    AccountListComponent
  ],
  templateUrl: './accounts.component.html',
  styleUrls: ['./accounts.component.css']
})
export class AccountsComponent implements OnInit, OnDestroy {
  // Services injection
  private readonly accountService = inject(AccountService);
  private readonly authService = inject(AuthService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialog = inject(MatDialog);

  // Component lifecycle management
  private readonly destroy$ = new Subject<void>();

  // Core state signals - Focused on user accounts loading
  readonly accounts = signal<Account[]>([]);
  readonly loading = signal(false);
  readonly selectedTabIndex = signal(0);
  readonly error = signal<string | null>(null);

  // Simple getter instead of computed to avoid circular dependencies
  get currentUser() {
    return this.authService.currentUserSig();
  }

  get isUserAuthenticated(): boolean {
    return !!this.currentUser?.id;
  }

  // SOLUCION: Formatear valores en el component para evitar pipes reactivos
  
  get savingsAccounts(): Account[] {
    return this.accounts().filter(account => account.accountType === AccountType.SAVINGS);
  }
  
  get checkingAccounts(): Account[] {
    return this.accounts().filter(account => account.accountType === AccountType.CHECKING);
  }
  
  get creditCards(): Account[] {
    return this.accounts().filter(account => account.accountType === AccountType.CREDIT);
  }
  
  get usdAccounts(): Account[] {
    return this.accounts().filter(account => account.currency === Currency.USD);
  }
  
  get activeAccounts(): Account[] {
    return this.accounts().filter(account => account.isActive);
  }
  
  get inactiveAccounts(): Account[] {
    return this.accounts().filter(account => !account.isActive);
  }

  // Financial summary - valores ya formateados para evitar pipes
  get totalBalance(): number {
    return this.activeAccounts
      .filter(account => account.currency === Currency.ARS && account.accountType !== AccountType.CREDIT)
      .reduce((total, account) => total + account.balance, 0);
  }

  get totalBalanceFormatted(): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS'
    }).format(this.totalBalance);
  }

  get totalCreditLimit(): number {
    return this.creditCards.reduce((total, account) => total + (account.creditLimit || 0), 0);
  }

  get totalCreditLimitFormatted(): string {
    return new Intl.NumberFormat('es-AR', {
      style: 'currency',
      currency: 'ARS'
    }).format(this.totalCreditLimit);
  }

  get totalUsdBalance(): number {
    return this.usdAccounts.reduce((total, account) => total + account.balance, 0);
  }

  get totalUsdBalanceFormatted(): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(this.totalUsdBalance);
  }

  get accountsCount(): number {
    return this.accounts().length;
  }

  ngOnInit(): void {
    this.initializeUserAccounts();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  /**
   * Initialize user accounts loading
   * Main entry point for account data management
   */
  private initializeUserAccounts(): void {
    if (!this.isUserAuthenticated) {
      this.handleAuthenticationError();
      return;
    }

    // Load accounts immediately if user is available
    this.loadUserAccounts();
  }

  /**
   * Core method: Load user accounts from API
   * Centralized account loading logic with proper error handling
   */
  private loadUserAccounts(): void {
    const user = this.currentUser;
    
    if (!user?.id) {
      this.handleAuthenticationError();
      return;
    }

    // Reset error state and start loading
    this.error.set(null);
    this.loading.set(true);

    this.accountService.getAccountsByUser(user.id)
      .pipe(
        takeUntil(this.destroy$),
        finalize(() => this.loading.set(false))
      )
      .subscribe({
        next: (response: AccountsListResponse) => this.handleAccountsLoaded(response),
        error: (error: any) => this.handleAccountsLoadError(error)
      });
  }

  /**
   * Handle successful accounts loading
   */
  private handleAccountsLoaded(response: AccountsListResponse): void {
    const accounts = response.accounts || [];
    this.accounts.set(accounts);
    this.error.set(null);
    
    console.log(`Successfully loaded ${accounts.length} accounts for user`);
  }

  /**
   * Handle accounts loading error
   */
  private handleAccountsLoadError(error: any): void {
    const errorMessage = error.message || 'Error desconocido al cargar cuentas';
    
    this.accounts.set([]);
    this.error.set(errorMessage);
    
    this.snackBar.open(
      `Error al cargar las cuentas: ${errorMessage}`, 
      'Cerrar', 
      { duration: 5000, panelClass: ['error-snackbar'] }
    );
    
    console.error('Failed to load user accounts:', error);
  }

  /**
   * Handle authentication errors
   */
  private handleAuthenticationError(): void {
    const message = 'Usuario no autenticado. Por favor, inicia sesión.';
    this.error.set(message);
    this.snackBar.open(message, 'Cerrar', { 
      duration: 3000, 
      panelClass: ['error-snackbar'] 
    });
  }

  /**
   * Refresh accounts data
   * Public method to reload accounts when needed
   */
  refreshAccounts(): void {
    this.loadUserAccounts();
  }

  // =============================================================================
  // UI EVENT HANDLERS - Account Management Actions
  // =============================================================================

  onAddAccount(): void {
    if (!this.isUserAuthenticated) {
      this.handleAuthenticationError();
      return;
    }

    const dialogRef = this.dialog.open(AccountFormComponent, {
      width: '600px',
      maxWidth: '90vw',
      data: { 
        mode: 'create',
        userId: this.currentUser!.id 
      } as AccountFormDialogData
    });

    dialogRef.afterClosed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(newAccount => {
        if (newAccount) {
          this.handleAccountCreated();
        }
      });
  }

  onEditAccount(account: Account): void {
    const dialogRef = this.dialog.open(AccountFormComponent, {
      width: '600px',
      maxWidth: '90vw',
      data: { account, mode: 'edit' } as AccountFormDialogData
    });

    dialogRef.afterClosed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(updatedAccount => {
        if (updatedAccount) {
          this.handleAccountUpdated();
        }
      });
  }

  onDeleteAccount(account: Account): void {
    const dialogRef = this.dialog.open(AccountDeleteConfirmationComponent, {
      width: '400px',
      data: { account }
    });

    dialogRef.afterClosed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(confirmed => {
        if (confirmed) {
          this.performAccountDeletion(account);
        }
      });
  }

  onActivateAccount(account: Account): void {
    // TODO: Implement when account activation API is available
    console.log('Toggle account status:', account);
    this.refreshAccounts();
  }

  onManageWallet(account: Account): void {
    const dialogRef = this.dialog.open(WalletDialogComponent, {
      width: '600px',
      maxWidth: '90vw',
      data: { account }
    });

    dialogRef.afterClosed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(updatedAccount => {
        if (updatedAccount) {
          this.updateAccountInList(updatedAccount);
        }
      });
  }

  onManageCredit(account: Account): void {
    const dialogRef = this.dialog.open(CreditDialogComponent, {
      width: '700px',
      maxWidth: '90vw',
      data: { account }
    });

    dialogRef.afterClosed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(updatedAccount => {
        if (updatedAccount) {
          this.updateAccountInList(updatedAccount);
        }
      });
  }

  onAccountStatusChanged(updatedAccount: Account): void {
    this.updateAccountInList(updatedAccount);
  }

  onTabChanged(index: number): void {
    this.selectedTabIndex.set(index);
  }

  // =============================================================================
  // PRIVATE HELPER METHODS - Business Logic Support
  // =============================================================================

  /**
   * Handle successful account creation
   */
  private handleAccountCreated(): void {
    this.refreshAccounts();
    this.snackBar.open('Cuenta creada exitosamente', 'Cerrar', {
      duration: 3000,
      panelClass: ['success-snackbar']
    });
  }

  /**
   * Handle successful account update
   */
  private handleAccountUpdated(): void {
    this.refreshAccounts();
    this.snackBar.open('Cuenta actualizada exitosamente', 'Cerrar', {
      duration: 3000,
      panelClass: ['success-snackbar']
    });
  }

  /**
   * Perform account deletion with proper error handling
   */
  private performAccountDeletion(account: Account): void {
    this.accountService.deleteAccount(account.id)
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: () => {
          this.refreshAccounts();
          this.snackBar.open('Cuenta eliminada exitosamente', 'Cerrar', {
            duration: 3000,
            panelClass: ['success-snackbar']
          });
        },
        error: (error) => {
          console.error('Error deleting account:', error);
          this.snackBar.open('Error al eliminar la cuenta. Inténtalo de nuevo.', 'Cerrar', {
            duration: 5000,
            panelClass: ['error-snackbar']
          });
        }
      });
  }

  /**
   * Update a specific account in the local list
   * Optimized update without full refresh for better UX
   */
  private updateAccountInList(updatedAccount: Account): void {
    const currentAccounts = this.accounts();
    const accountIndex = currentAccounts.findIndex(acc => acc.id === updatedAccount.id);
    
    if (accountIndex !== -1) {
      const newAccounts = [...currentAccounts];
      newAccounts[accountIndex] = updatedAccount;
      this.accounts.set(newAccounts);
    }
  }

  // =============================================================================
  // REMOVED DEPRECATED METHODS - Now using direct getters for debugging
  // =============================================================================
}
