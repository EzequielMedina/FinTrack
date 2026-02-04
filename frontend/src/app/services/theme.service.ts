import { Injectable, signal, effect, inject, PLATFORM_ID } from '@angular/core';
import { DOCUMENT, isPlatformBrowser } from '@angular/common';

/**
 * Tipo de tema disponible en la aplicación
 */
export type Theme = 'light' | 'dark';

/**
 * Servicio para gestionar el tema de la aplicación (claro/oscuro)
 * 
 * Características:
 * - Persistencia en localStorage
 * - Signals reactivos de Angular
 * - Detección de preferencia del sistema (opcional)
 * - SSR-safe (compatible con Server-Side Rendering)
 */
@Injectable({
  providedIn: 'root'
})
export class ThemeService {
  private readonly document = inject(DOCUMENT);
  private readonly platformId = inject(PLATFORM_ID);
  private readonly THEME_KEY = 'fintrack-theme';
  private readonly defaultTheme: Theme = 'light';

  // Signal privado para el tema actual
  private readonly _currentTheme = signal<Theme>(this.loadTheme());
  
  // Signal público de solo lectura
  public readonly currentTheme = this._currentTheme.asReadonly();

  // Media query listener para preferencia del sistema (opcional)
  private prefersDarkQuery?: MediaQueryList;
  private prefersDarkListener?: (e: MediaQueryListEvent) => void;

  constructor() {
    // Aplicar tema inicial
    if (isPlatformBrowser(this.platformId)) {
      this.applyTheme(this._currentTheme());
      
      // Efecto reactivo: aplicar tema cuando cambie el signal
      effect(() => {
        const theme = this._currentTheme();
        this.applyTheme(theme);
        this.saveTheme(theme);
      });

      // Opcional: Detectar preferencia del sistema si no hay tema guardado
      this.initializeSystemPreference();
    }
  }

  /**
   * Carga el tema guardado en localStorage o devuelve el tema por defecto
   * @returns El tema cargado o el tema por defecto
   */
  private loadTheme(): Theme {
    if (!isPlatformBrowser(this.platformId)) {
      return this.defaultTheme;
    }

    try {
      const savedTheme = localStorage.getItem(this.THEME_KEY) as Theme;
      if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
        return savedTheme;
      }
    } catch (error) {
      console.warn('Error al cargar el tema desde localStorage:', error);
    }

    return this.defaultTheme;
  }

  /**
   * Guarda el tema en localStorage
   * @param theme - Tema a guardar
   */
  private saveTheme(theme: Theme): void {
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }

    try {
      localStorage.setItem(this.THEME_KEY, theme);
    } catch (error) {
      console.warn('Error al guardar el tema en localStorage:', error);
    }
  }

  /**
   * Aplica el tema al documento HTML
   * @param theme - Tema a aplicar
   */
  private applyTheme(theme: Theme): void {
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }

    const htmlElement = this.document.documentElement;
    
    if (theme === 'dark') {
      htmlElement.setAttribute('data-theme', 'dark');
    } else {
      htmlElement.removeAttribute('data-theme');
    }
  }

  /**
   * Inicializa la detección de preferencia del sistema
   * Solo se activa si no hay un tema guardado en localStorage
   */
  private initializeSystemPreference(): void {
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }

    try {
      // Solo usar preferencia del sistema si no hay tema guardado
      const savedTheme = localStorage.getItem(this.THEME_KEY);
      if (savedTheme) {
        return; // Ya hay un tema guardado, no usar preferencia del sistema
      }

      // Detectar preferencia actual del sistema
      this.prefersDarkQuery = window.matchMedia('(prefers-color-scheme: dark)');
      const prefersDark = this.prefersDarkQuery.matches;
      
      // Aplicar tema según preferencia del sistema
      if (prefersDark && this._currentTheme() === 'light') {
        this._currentTheme.set('dark');
      }

      // Escuchar cambios en la preferencia del sistema
      this.prefersDarkListener = (e: MediaQueryListEvent) => {
        // Solo cambiar si no hay tema guardado
        if (!localStorage.getItem(this.THEME_KEY)) {
          this._currentTheme.set(e.matches ? 'dark' : 'light');
        }
      };

      // Agregar listener (compatible con navegadores modernos)
      if (this.prefersDarkQuery.addEventListener) {
        this.prefersDarkQuery.addEventListener('change', this.prefersDarkListener);
      } else {
        // Fallback para navegadores antiguos
        this.prefersDarkQuery.addListener(this.prefersDarkListener);
      }
    } catch (error) {
      console.warn('Error al inicializar preferencia del sistema:', error);
    }
  }

  /**
   * Alterna entre tema claro y oscuro
   */
  toggleTheme(): void {
    const newTheme: Theme = this._currentTheme() === 'light' ? 'dark' : 'light';
    this.setTheme(newTheme);
  }

  /**
   * Establece un tema específico
   * @param theme - Tema a establecer ('light' o 'dark')
   */
  setTheme(theme: Theme): void {
    if (theme !== 'light' && theme !== 'dark') {
      console.warn(`Tema inválido: ${theme}. Usando tema por defecto.`);
      theme = this.defaultTheme;
    }
    this._currentTheme.set(theme);
  }

  /**
   * Verifica si el tema actual es oscuro
   * @returns true si el tema es oscuro, false en caso contrario
   */
  isDarkMode(): boolean {
    return this._currentTheme() === 'dark';
  }

  /**
   * Limpia el tema guardado y restaura el tema por defecto
   */
  resetTheme(): void {
    if (isPlatformBrowser(this.platformId)) {
      try {
        localStorage.removeItem(this.THEME_KEY);
      } catch (error) {
        console.warn('Error al limpiar el tema:', error);
      }
    }
    this._currentTheme.set(this.defaultTheme);
  }

  /**
   * Limpia recursos al destruir el servicio
   */
  ngOnDestroy(): void {
    if (this.prefersDarkQuery && this.prefersDarkListener) {
      if (this.prefersDarkQuery.removeEventListener) {
        this.prefersDarkQuery.removeEventListener('change', this.prefersDarkListener);
      } else {
        // Fallback para navegadores antiguos
        this.prefersDarkQuery.removeListener(this.prefersDarkListener);
      }
    }
  }
}
