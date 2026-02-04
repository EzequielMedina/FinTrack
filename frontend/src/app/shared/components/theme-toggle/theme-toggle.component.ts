import { Component, ChangeDetectionStrategy, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ThemeService } from '../../../services/theme.service';

/**
 * Componente para alternar entre modo claro y oscuro
 * 
 * Características:
 * - Botón con icono que cambia según el tema
 * - Tooltip descriptivo
 * - Accesible (aria-label)
 * - Estilos consistentes con el tema actual
 */
@Component({
  selector: 'app-theme-toggle',
  standalone: true,
  imports: [
    CommonModule,
    MatButtonModule,
    MatIconModule,
    MatTooltipModule
  ],
  template: `
    <button
      mat-icon-button
      type="button"
      (click)="onToggleTheme()"
      [matTooltip]="tooltipText"
      [attr.aria-label]="ariaLabel"
      [attr.aria-pressed]="isDarkMode()"
      class="theme-toggle-button">
      <mat-icon>{{ iconName }}</mat-icon>
    </button>
  `,
  styles: [`
    :host {
      display: inline-block;
    }
    
    .theme-toggle-button {
      color: var(--text-primary);
      transition: color var(--transition-base), transform var(--transition-base);
    }
    
    .theme-toggle-button:hover {
      transform: scale(1.1);
    }
    
    .theme-toggle-button:focus {
      outline: 2px solid var(--accent-500);
      outline-offset: 2px;
    }
    
    .theme-toggle-button mat-icon {
      font-size: 24px;
      width: 24px;
      height: 24px;
    }
  `],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ThemeToggleComponent {
  protected readonly themeService = inject(ThemeService);

  /**
   * Obtiene el nombre del icono según el tema actual
   */
  get iconName(): string {
    return this.themeService.isDarkMode() ? 'light_mode' : 'dark_mode';
  }

  /**
   * Obtiene el texto del tooltip según el tema actual
   */
  get tooltipText(): string {
    return this.themeService.isDarkMode() 
      ? 'Cambiar a modo claro' 
      : 'Cambiar a modo oscuro';
  }

  /**
   * Obtiene el aria-label según el tema actual
   */
  get ariaLabel(): string {
    return this.themeService.isDarkMode() 
      ? 'Cambiar a modo claro' 
      : 'Cambiar a modo oscuro';
  }

  /**
   * Verifica si el tema actual es oscuro
   */
  isDarkMode(): boolean {
    return this.themeService.isDarkMode();
  }

  /**
   * Maneja el click en el botón para alternar el tema
   */
  onToggleTheme(): void {
    this.themeService.toggleTheme();
  }
}
