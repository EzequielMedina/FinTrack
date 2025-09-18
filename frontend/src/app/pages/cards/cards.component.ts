import { ChangeDetectionStrategy, Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule, MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogModule, MatDialog } from '@angular/material/dialog';
import { CardService } from '../../services/card.service';
import { AuthService } from '../../services/auth.service';
import { Card, CardType, CardStatus } from '../../models';
import { CardListComponent } from './card-list/card-list.component';
import { CardFormComponent } from './card-form/card-form.component';

@Component({
  selector: 'app-cards',
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
    CardListComponent
  ],
  templateUrl: './cards.component.html',
  styleUrls: ['./cards.component.css'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class CardsComponent implements OnInit {
  private readonly cardService = inject(CardService);
  private readonly authService = inject(AuthService);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialog = inject(MatDialog);

  // Signals para el estado del componente
  cards = signal<Card[]>([]);
  loading = signal(false);
  selectedTabIndex = signal(0);

  // Computed signals para filtrar tarjetas
  creditCards = signal<Card[]>([]);
  debitCards = signal<Card[]>([]);
  activeCards = signal<Card[]>([]);
  inactiveCards = signal<Card[]>([]);

  get currentUser() {
    return this.authService.currentUserSig;
  }

  ngOnInit(): void {
    this.loadUserCards();
  }

  private loadUserCards(): void {
    const user = this.currentUser();
    if (!user) return;

    this.loading.set(true);
    this.cardService.getCardsByUser(user.id).subscribe({
      next: (response) => {
        this.cards.set(response.cards);
        this.updateFilteredCards();
        this.loading.set(false);
      },
      error: (error) => {
        console.error('Error loading cards:', error);
        this.snackBar.open('Error al cargar las tarjetas', 'Cerrar', {
          duration: 3000
        });
        this.loading.set(false);
      }
    });
  }

  private updateFilteredCards(): void {
    const allCards = this.cards();
    
    this.creditCards.set(allCards.filter(card => card.cardType === CardType.CREDIT));
    this.debitCards.set(allCards.filter(card => card.cardType === CardType.DEBIT));
    this.activeCards.set(allCards.filter(card => card.status === CardStatus.ACTIVE));
    this.inactiveCards.set(allCards.filter(card => 
      card.status === CardStatus.INACTIVE || 
      card.status === CardStatus.BLOCKED ||
      card.status === CardStatus.EXPIRED
    ));
  }

  onAddCard(): void {
    const dialogRef = this.dialog.open(CardFormComponent, {
      width: '600px',
      disableClose: true,
      data: { mode: 'create' }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.loadUserCards(); // Recargar tarjetas después de agregar
        this.snackBar.open('Tarjeta agregada exitosamente', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  onEditCard(card: Card): void {
    const dialogRef = this.dialog.open(CardFormComponent, {
      width: '600px',
      disableClose: true,
      data: { mode: 'edit', card }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.loadUserCards(); // Recargar tarjetas después de editar
        this.snackBar.open('Tarjeta actualizada exitosamente', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  onDeleteCard(card: Card): void {
    if (confirm(`¿Estás seguro de que deseas eliminar la tarjeta terminada en ${card.lastFourDigits}?`)) {
      this.cardService.deleteCard(card.accountId || card.id, card.id).subscribe({
        next: () => {
          this.loadUserCards();
          this.snackBar.open('Tarjeta eliminada exitosamente', 'Cerrar', {
            duration: 3000
          });
        },
        error: (error) => {
          console.error('Error deleting card:', error);
          this.snackBar.open('Error al eliminar la tarjeta', 'Cerrar', {
            duration: 3000
          });
        }
      });
    }
  }

  onSetDefaultCard(card: Card): void {
    this.cardService.setDefaultCard(card.accountId || card.id, card.id).subscribe({
      next: () => {
        this.loadUserCards();
        this.snackBar.open('Tarjeta establecida como predeterminada', 'Cerrar', {
          duration: 3000
        });
      },
      error: (error) => {
        console.error('Error setting default card:', error);
        this.snackBar.open('Error al establecer tarjeta predeterminada', 'Cerrar', {
          duration: 3000
        });
      }
    });
  }

  onBlockCard(card: Card): void {
    const action = card.status === CardStatus.BLOCKED ? 'desbloquear' : 'bloquear';
    
    if (confirm(`¿Estás seguro de que deseas ${action} la tarjeta terminada en ${card.lastFourDigits}?`)) {
      const operation = card.status === CardStatus.BLOCKED ? 
        this.cardService.unblockCard(card.accountId || card.id, card.id) :
        this.cardService.blockCard(card.accountId || card.id, card.id);

      operation.subscribe({
        next: () => {
          this.loadUserCards();
          this.snackBar.open(`Tarjeta ${action}da exitosamente`, 'Cerrar', {
            duration: 3000
          });
        },
        error: (error) => {
          console.error(`Error ${action}ing card:`, error);
          this.snackBar.open(`Error al ${action} la tarjeta`, 'Cerrar', {
            duration: 3000
          });
        }
      });
    }
  }

  onCardStatusChanged(updatedCard: Card): void {
    // Recargar la lista de tarjetas para reflejar el cambio de estado
    this.loadUserCards();
  }

  onTabChanged(index: number): void {
    this.selectedTabIndex.set(index);
  }
}