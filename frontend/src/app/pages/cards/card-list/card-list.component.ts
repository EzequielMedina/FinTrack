import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Card, CardType, CardStatus, CardBrand } from '../../../models';
import { CardService } from '../../../services/card.service';
import { CardStatusDialogComponent } from '../card-status-dialog/card-status-dialog.component';

@Component({
  selector: 'app-card-list',
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
  templateUrl: './card-list.component.html',
  styleUrls: ['./card-list.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CardListComponent {
  private readonly cardService = inject(CardService);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);

  @Input() cards: Card[] = [];
  @Output() editCard = new EventEmitter<Card>();
  @Output() deleteCard = new EventEmitter<Card>();
  @Output() setDefaultCard = new EventEmitter<Card>();
  @Output() blockCard = new EventEmitter<Card>();
  @Output() cardStatusChanged = new EventEmitter<Card>();
  @Output() viewCardDetail = new EventEmitter<Card>();

  // Exponer enums para usar en el template
  readonly CardType = CardType;
  readonly CardStatus = CardStatus;
  readonly CardBrand = CardBrand;

  onEdit(card: Card): void {
    this.editCard.emit(card);
  }

  onDelete(card: Card): void {
    this.deleteCard.emit(card);
  }

  onSetDefault(card: Card): void {
    this.setDefaultCard.emit(card);
  }

  onChangeStatus(card: Card): void {
    const action = card.status === CardStatus.ACTIVE ? 'deactivate' : 'activate';
    
    const dialogRef = this.dialog.open(CardStatusDialogComponent, {
      width: '500px',
      data: { card, action }
    });

    dialogRef.afterClosed().subscribe(confirmed => {
      if (confirmed) {
        this.performStatusChange(card, action);
      }
    });
  }

  private performStatusChange(card: Card, action: 'activate' | 'deactivate'): void {
    const operation = action === 'activate' 
      ? this.cardService.activateCard(card.accountId, card.id)
      : this.cardService.deactivateCard(card.accountId, card.id);

    operation.subscribe({
      next: (updatedCard) => {
        this.cardStatusChanged.emit(updatedCard);
        const message = action === 'activate' 
          ? 'Tarjeta activada exitosamente' 
          : 'Tarjeta desactivada exitosamente';
        this.snackBar.open(message, 'Cerrar', { duration: 3000 });
      },
      error: (error) => {
        console.error('Error changing card status:', error);
        const message = action === 'activate' 
          ? 'Error al activar la tarjeta' 
          : 'Error al desactivar la tarjeta';
        this.snackBar.open(message, 'Cerrar', { duration: 3000 });
      }
    });
  }

  onBlock(card: Card): void {
    // Mantener el método existente para bloqueo/desbloqueo
    this.blockCard.emit(card);
  }

  getCardBrandIcon(brand: CardBrand): string {
    return this.cardService.getCardBrandIcon(brand);
  }

  getCardTypeIcon(type: CardType): string {
    return type === CardType.CREDIT ? 'credit_card' : 'payment';
  }

  getCardTypeLabel(type: CardType): string {
    return type === CardType.CREDIT ? 'Crédito' : 'Débito';
  }

  getStatusLabel(status: CardStatus): string {
    const labels = {
      [CardStatus.ACTIVE]: 'Activa',
      [CardStatus.INACTIVE]: 'Inactiva',
      [CardStatus.BLOCKED]: 'Bloqueada',
      [CardStatus.EXPIRED]: 'Vencida'
    };
    return labels[status];
  }

  getStatusColor(status: CardStatus): string {
    const colors = {
      [CardStatus.ACTIVE]: 'primary',
      [CardStatus.INACTIVE]: 'accent',
      [CardStatus.BLOCKED]: 'warn',
      [CardStatus.EXPIRED]: 'warn'
    };
    return colors[status];
  }

  getStatusIcon(status: CardStatus): string {
    const icons = {
      [CardStatus.ACTIVE]: 'check_circle',
      [CardStatus.INACTIVE]: 'cancel',
      [CardStatus.BLOCKED]: 'block',
      [CardStatus.EXPIRED]: 'warning'
    };
    return icons[status];
  }

  isCardExpired(card: Card): boolean {
    const now = new Date();
    const currentYear = now.getFullYear();
    const currentMonth = now.getMonth() + 1;
    
    return card.expirationYear < currentYear || 
           (card.expirationYear === currentYear && card.expirationMonth < currentMonth);
  }

  formatExpirationDate(month: number, year: number): string {
    const monthStr = month.toString().padStart(2, '0');
    const yearStr = year.toString().slice(-2);
    return `${monthStr}/${yearStr}`;
  }

  onViewDetail(card: Card): void {
    this.viewCardDetail.emit(card);
  }
}