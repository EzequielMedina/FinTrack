import { Component, EventEmitter, Input, Output, inject, ChangeDetectorRef, OnChanges, SimpleChanges } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Account, AccountType, Currency } from '../../../models';
import { AccountValidationService } from '../../../services';

@Component({
  selector: 'app-account-list',
  standalone: true,
  imports: [
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatMenuModule,
    MatTooltipModule
  ],
  template: `
    <div class="account-list-container">
      @if (accounts.length === 0) {
        <div class="no-accounts">
          <mat-icon class="no-accounts-icon">account_balance_wallet</mat-icon>
          <p>No hay cuentas en esta categoría</p>
        </div>
      } @else {
        <div class="accounts-grid">
          @for (account of accounts; track account.id) {
            <mat-card class="account-item" [class.inactive-account]="!account.isActive">
              <!-- Account Header -->
              <mat-card-header>
                <div class="account-type-icon" mat-card-avatar [style.background]="getAccountGradient(account.accountType)">
                  <mat-icon>{{ getAccountIcon(account.accountType) }}</mat-icon>
                </div>
                <mat-card-title class="account-title">
                  <span class="account-name">{{ account.name }}</span>
                  @if (!account.isActive) {
                    <mat-chip class="inactive-chip" color="warn" [selected]="true">
                      <mat-icon>cancel</mat-icon>
                      Inactiva
                    </mat-chip>
                  }
                </mat-card-title>
                <mat-card-subtitle>{{ getAccountTypeDisplayName(account.accountType) }} - {{ account.currency }}</mat-card-subtitle>

                <!-- Actions Menu -->
                <button mat-icon-button [matMenuTriggerFor]="accountMenu" class="account-menu-trigger">
                  <mat-icon>more_vert</mat-icon>
                </button>
                <mat-menu #accountMenu="matMenu">
                  <button mat-menu-item (click)="onEdit(account)">
                    <mat-icon>edit</mat-icon>
                    <span>Editar</span>
                  </button>
                  @if (account.isActive) {
                    <button mat-menu-item (click)="onActivate(account)">
                      <mat-icon>toggle_off</mat-icon>
                      <span>Desactivar</span>
                    </button>
                  } @else {
                    <button mat-menu-item (click)="onActivate(account)">
                      <mat-icon>toggle_on</mat-icon>
                      <span>Activar</span>
                    </button>
                  }
                  <mat-divider></mat-divider>
                  <button mat-menu-item (click)="onDelete(account)" class="delete-action">
                    <mat-icon>delete</mat-icon>
                    <span>Eliminar</span>
                  </button>
                </mat-menu>
              </mat-card-header>

              <!-- Account Content -->
              <mat-card-content class="account-content">
                <!-- Account ID -->
                <div class="account-id">
                  ID: {{ account.id.slice(0, 8) }}...
                </div>

                <!-- Account Details -->
                <div class="account-details">
                  <div class="account-detail">
                    <span class="label">Saldo:</span>
                    <span class="value balance" [class.negative]="account.balance < 0">
                      {{ formatBalance(account.balance, account.currency) }}
                    </span>
                  </div>

                  @if (account.accountType === AccountType.CREDIT && account.creditLimit) {
                    <div class="account-detail">
                      <span class="label">Límite de Crédito:</span>
                      <span class="value">{{ formatBalance(account.creditLimit, account.currency) }}</span>
                    </div>
                  }

                  <div class="account-detail">
                    <span class="label">Descripción:</span>
                    <span class="value description">{{ getAccountTypeDescription(account) }}</span>
                  </div>

                  @if (account.dni) {
                    <div class="account-detail">
                      <span class="label">DNI:</span>
                      <span class="value">{{ account.dni }}</span>
                    </div>
                  }

                  @if (account.accountType === AccountType.CREDIT) {
                    @if (account.closingDate) {
                      <div class="account-detail">
                        <span class="label">Fecha de Cierre:</span>
                        <span class="value">{{ account.closingDate }}</span>
                      </div>
                    }
                    @if (account.dueDate) {
                      <div class="account-detail">
                        <span class="label">Fecha de Vencimiento:</span>
                        <span class="value">{{ account.dueDate }}</span>
                      </div>
                    }
                  }

                  <div class="account-detail">
                    <span class="label">Creada:</span>
                    <span class="value">{{ account.createdAt | date:'dd/MM/yyyy' }}</span>
                  </div>
                </div>

                <!-- Account Status -->
                <div class="account-status">
                  <mat-chip [color]="getStatusColor(account)" [selected]="true">
                    <mat-icon>{{ account.isActive ? 'check_circle' : 'cancel' }}</mat-icon>
                    {{ getStatusLabel(account) }}
                  </mat-chip>
                </div>
              </mat-card-content>

              <!-- Account Actions -->
              <mat-card-actions class="account-actions">
                @if (account.isActive) {
                  @if (canAddFunds(account)) {
                    <button mat-button color="primary" (click)="onManageWallet(account)">
                      <mat-icon>account_balance_wallet</mat-icon>
                      Gestionar Fondos
                    </button>
                  }
                  @if (canManageCredit(account)) {
                    <button mat-button color="accent" (click)="onManageCredit(account)">
                      <mat-icon>credit_card</mat-icon>
                      Gestionar Crédito
                    </button>
                  }
                } @else {
                  <button mat-button color="primary" (click)="onActivate(account)">
                    <mat-icon>toggle_on</mat-icon>
                    Activar Cuenta
                  </button>
                }
                
                <button mat-button (click)="onEdit(account)">
                  <mat-icon>edit</mat-icon>
                  Editar
                </button>
              </mat-card-actions>
            </mat-card>
          }
        </div>
      }
    </div>
  `,
  styles: [`
    /* Account List Component Styles */
    .account-list-container {
      width: 100%;
    }

    /* No Accounts State */
    .no-accounts {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 60px 20px;
      text-align: center;
      color: #6b7280;
    }

    .no-accounts-icon {
      font-size: 48px !important;
      width: 48px !important;
      height: 48px !important;
      color: #d1d5db !important;
      margin-bottom: 20px;
    }

    .no-accounts p {
      color: #6b7280;
      font-size: 1rem;
      margin: 0;
    }

    /* Accounts Grid */
    .accounts-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
      gap: 20px;
      padding: 10px 0;
    }

    .account-item {
      transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
      position: relative;
      border-radius: 12px;
      overflow: hidden;
      background: white;
      border: 1px solid #e5e7eb;
    }

    .account-item:hover {
      transform: translateY(-2px);
      box-shadow: 0 4px 12px rgba(0,0,0,0.15);
    }

    .account-item.inactive-account {
      opacity: 0.8;
      border-color: #fca5a5;
      background: #fef2f2;
    }

    /* Account Header */
    .account-type-icon {
      border-radius: 50%;
      color: white !important;
      font-size: 20px !important;
      width: 40px !important;
      height: 40px !important;
      display: flex !important;
      align-items: center !important;
      justify-content: center !important;
    }

    .account-title {
      display: flex !important;
      align-items: center !important;
      gap: 10px !important;
      flex-wrap: wrap !important;
    }

    .account-name {
      font-weight: 600;
      color: #1f2937;
      font-size: 1.1rem;
    }

    .inactive-chip {
      font-size: 0.75rem;
      height: 24px;
    }

    .inactive-chip mat-icon {
      font-size: 14px !important;
      width: 14px !important;
      height: 14px !important;
    }

    .account-menu-trigger {
      margin-left: auto;
    }

    /* Account Content */
    .account-content {
      padding: 16px !important;
    }

    .account-id {
      font-family: 'Courier New', monospace;
      font-size: 0.85rem;
      color: #6b7280;
      margin-bottom: 15px;
      padding: 5px 10px;
      background: #f9fafb;
      border-radius: 4px;
      border: 1px solid #e5e7eb;
    }

    /* Account Details */
    .account-details {
      display: flex;
      flex-direction: column;
      gap: 8px;
      margin-bottom: 15px;
    }

    .account-detail {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      font-size: 0.9rem;
      line-height: 1.4;
    }

    .account-detail .label {
      font-weight: 500;
      color: #6b7280;
      min-width: 120px;
    }

    .account-detail .value {
      color: #1f2937;
      font-weight: 500;
      text-align: right;
      flex: 1;
      word-break: break-word;
    }

    .account-detail .value.balance {
      font-weight: 700;
      font-size: 1rem;
      color: #059669;
    }

    .account-detail .value.balance.negative {
      color: #dc2626;
    }

    .account-detail .value.description {
      font-weight: 400;
      color: #6b7280;
      font-style: italic;
      font-size: 0.85rem;
    }

    /* Account Status */
    .account-status {
      display: flex;
      justify-content: center;
      margin-bottom: 10px;
    }

    .account-status mat-chip {
      font-weight: 500;
    }

    .account-status mat-chip mat-icon {
      font-size: 16px !important;
      width: 16px !important;
      height: 16px !important;
    }

    /* Account Actions */
    .account-actions {
      display: flex !important;
      flex-wrap: wrap !important;
      gap: 8px !important;
      padding: 12px 16px 16px 16px !important;
      border-top: 1px solid #e5e7eb;
      background: #fafafa;
    }

    .account-actions button {
      flex: 1;
      min-width: 120px;
      font-size: 0.85rem;
      font-weight: 500;
    }

    .account-actions button mat-icon {
      font-size: 18px !important;
      width: 18px !important;
      height: 18px !important;
      margin-right: 5px;
    }

    /* Menu Styles */
    .delete-action {
      color: #dc2626 !important;
    }

    .delete-action mat-icon {
      color: #dc2626 !important;
    }

    /* Responsive Design */
    @media (max-width: 1024px) {
      .accounts-grid {
        grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
        gap: 15px;
      }
    }

    @media (max-width: 768px) {
      .accounts-grid {
        grid-template-columns: 1fr;
        gap: 15px;
      }

      .account-detail {
        flex-direction: column;
        gap: 2px;
      }

      .account-detail .label {
        min-width: unset;
        font-size: 0.85rem;
      }

      .account-detail .value {
        text-align: left;
        font-size: 0.9rem;
      }

      .account-actions {
        flex-direction: column;
      }

      .account-actions button {
        min-width: unset;
        width: 100%;
      }
    }

    @media (max-width: 480px) {
      .accounts-grid {
        gap: 10px;
      }

      .account-content {
        padding: 12px !important;
      }

      .account-details {
        gap: 6px;
      }

      .account-detail .label,
      .account-detail .value {
        font-size: 0.8rem;
      }

      .account-actions {
        padding: 10px 12px 12px 12px !important;
      }
    }

    /* Animation for inactive accounts */
    .inactive-account {
      animation: pulse-warning 2s infinite;
    }

    @keyframes pulse-warning {
      0% {
        box-shadow: 0 0 0 0 rgba(252, 165, 165, 0.4);
      }
      70% {
        box-shadow: 0 0 0 10px rgba(252, 165, 165, 0);
      }
      100% {
        box-shadow: 0 0 0 0 rgba(252, 165, 165, 0);
      }
    }
  `]
})
export class AccountListComponent implements OnChanges {
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);
  private readonly validationService = inject(AccountValidationService);
  private readonly cdr = inject(ChangeDetectorRef);

  @Input() accounts: Account[] = [];
  @Output() editAccount = new EventEmitter<Account>();
  @Output() deleteAccount = new EventEmitter<Account>();
  @Output() activateAccount = new EventEmitter<Account>();
  @Output() manageWallet = new EventEmitter<Account>();
  @Output() manageCredit = new EventEmitter<Account>();
  @Output() accountStatusChanged = new EventEmitter<Account>();

  // Exponer enums para usar en el// Constants for template
  readonly AccountType = AccountType;
  readonly Currency = Currency;

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['accounts']) {
      console.log('AccountListComponent - accounts changed:', changes['accounts'].currentValue);
      this.cdr.detectChanges();
    }
  }

  // TrackBy function para mejorar el rendimiento del @for
  trackByAccountId(index: number, account: Account): string {
    return account.id;
  }

  onEdit(account: Account): void {
    this.editAccount.emit(account);
  }

  onDelete(account: Account): void {
    this.deleteAccount.emit(account);
  }

  onActivate(account: Account): void {
    this.activateAccount.emit(account);
  }

  onManageWallet(account: Account): void {
    this.manageWallet.emit(account);
  }

  onManageCredit(account: Account): void {
    this.manageCredit.emit(account);
  }

  getAccountIcon(accountType: AccountType): string {
    const icons = {
      [AccountType.SAVINGS]: 'savings',
      [AccountType.CHECKING]: 'account_balance',
      [AccountType.CREDIT]: 'credit_card',
      [AccountType.WALLET]: 'account_balance_wallet'
    };
    return icons[accountType as keyof typeof icons] || 'account_balance_wallet';
  }

  getAccountTypeColor(accountType: AccountType): string {
    const colors = {
      [AccountType.SAVINGS]: 'primary',
      [AccountType.CHECKING]: 'accent',
      [AccountType.CREDIT]: 'warn',
      [AccountType.WALLET]: 'primary'
    };
    return colors[accountType as keyof typeof colors] || 'primary';
  }

  getStatusColor(account: Account): string {
    return account.isActive ? 'primary' : 'warn';
  }

  getStatusLabel(account: Account): string {
    return account.isActive ? 'Activa' : 'Inactiva';
  }

  getCurrencySymbol(currency: Currency): string {
    const symbols = {
      [Currency.ARS]: '$',
      [Currency.USD]: 'US$',
      [Currency.EUR]: '€'
    };
    return symbols[currency as keyof typeof symbols] || currency;
  }

  formatBalance(balance: number, currency: Currency): string {
    return this.validationService.formatCurrency(balance, currency);
  }

  getAccountGradient(accountType: AccountType): string {
    const gradients = {
      [AccountType.SAVINGS]: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      [AccountType.CHECKING]: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
      [AccountType.CREDIT]: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
      [AccountType.WALLET]: 'linear-gradient(135deg, #a8edea 0%, #fed6e3 100%)'
    };
    return gradients[accountType as keyof typeof gradients] || gradients[AccountType.CHECKING];
  }

  getAccountTypeDisplayName(accountType: AccountType): string {
    return this.validationService.getAccountTypeDisplayName(accountType);
  }

  getCurrencyDisplayName(currency: Currency): string {
    return this.validationService.getCurrencyDisplayName(currency);
  }

  getAccountTypeDescription(account: Account): string {
    const typeDescriptions = {
      [AccountType.SAVINGS]: 'Cuenta para guardar dinero con intereses',
      [AccountType.CHECKING]: 'Cuenta para movimientos diarios',
      [AccountType.CREDIT]: 'Tarjeta para compras a crédito',
      [AccountType.WALLET]: 'Billetera digital para operaciones rápidas'
    };
    return typeDescriptions[account.accountType as keyof typeof typeDescriptions] || 'Cuenta bancaria';
  }

  hasAvailableCredit(account: Account): boolean {
    return account.accountType === AccountType.CREDIT && !!account.creditLimit;
  }

  getAvailableCredit(account: Account): number {
    if (!this.hasAvailableCredit(account)) return 0;
    return (account.creditLimit || 0) - Math.abs(account.balance);
  }

  isLowBalance(account: Account): boolean {
    if (account.accountType === AccountType.CREDIT) {
      const availableCredit = this.getAvailableCredit(account);
      const creditLimit = account.creditLimit || 0;
      return creditLimit > 0 && (availableCredit / creditLimit) < 0.2; // Less than 20% available
    }
    return account.balance < 1000; // Less than $1000 for regular accounts
  }

  canAddFunds(account: Account): boolean {
    return account.isActive && account.accountType !== AccountType.CREDIT;
  }

  canWithdrawFunds(account: Account): boolean {
    return account.isActive && account.balance > 0 && account.accountType !== AccountType.CREDIT;
  }

  canManageCredit(account: Account): boolean {
    return account.isActive && account.accountType === AccountType.CREDIT;
  }
}