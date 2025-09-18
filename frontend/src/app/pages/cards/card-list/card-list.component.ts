import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatMenuModule } from '@angular/material/menu';
import { MatTooltipModule } from '@angular/material/tooltip';
import { Card, CardType, CardStatus, CardBrand } from '../../../models';
import { CardService } from '../../../services/card.service';

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

  @Input() cards: Card[] = [];
  @Output() editCard = new EventEmitter<Card>();
  @Output() deleteCard = new EventEmitter<Card>();
  @Output() setDefaultCard = new EventEmitter<Card>();
  @Output() blockCard = new EventEmitter<Card>();

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

  onBlock(card: Card): void {
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
}