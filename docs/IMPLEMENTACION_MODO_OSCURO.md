# Documentación: Implementación de Modo Oscuro en FinTrack

## 📋 Tabla de Contenidos

1. [Introducción](#introducción)
2. [Arquitectura Actual](#arquitectura-actual)
3. [Estrategia de Implementación](#estrategia-de-implementación)
4. [Pasos de Implementación](#pasos-de-implementación)
5. [Variables CSS para Modo Oscuro](#variables-css-para-modo-oscuro)
6. [Servicio de Tema](#servicio-de-tema)
7. [Componente Toggle](#componente-toggle)
8. [Integración con Angular Material](#integración-con-angular-material)
9. [Consideraciones Especiales](#consideraciones-especiales)
10. [Testing y Verificación](#testing-y-verificación)

---

## 🎯 Introducción

### Objetivo
Implementar un sistema de modo oscuro completo que permita a los usuarios alternar entre tema claro y oscuro, mejorando la experiencia visual especialmente en presentaciones donde los colores claros pueden no notarse bien.

### Beneficios
- **Mejor visibilidad en presentaciones**: Los colores oscuros destacan mejor en proyecciones
- **Reducción de fatiga visual**: El modo oscuro es más cómodo en ambientes con poca luz
- **Modernidad**: Característica esperada en aplicaciones modernas
- **Persistencia**: La preferencia del usuario se guarda para futuras sesiones

---

## 🏗️ Arquitectura Actual

### Sistema de Diseño
FinTrack utiliza un sistema de diseño basado en **variables CSS** definidas en:
- `frontend/src/design-system.css` - Variables principales
- `frontend/src/styles.css` - Estilos globales
- `frontend/src/components.css` - Componentes reutilizables

### Variables Clave Actuales

```css
/* Fondos */
--bg-primary: #ffffff;        /* Fondo principal (blanco) */
--bg-secondary: #f9fafb;     /* Fondo secundario (gris claro) */
--bg-tertiary: #f3f4f6;      /* Fondo terciario */

/* Textos */
--text-primary: #111827;      /* Texto principal (negro) */
--text-secondary: #4b5563;   /* Texto secundario */
--text-inverse: #ffffff;     /* Texto inverso (blanco) */

/* Bordes */
--border-light: #e5e7eb;      /* Borde claro */
```

### Uso de Variables
Todos los componentes utilizan estas variables CSS mediante `var(--variable-name)`, lo que facilita la implementación del modo oscuro.

---

## 🎨 Estrategia de Implementación

### Enfoque: CSS Variables + Data Attribute

Utilizaremos un enfoque híbrido:
1. **Data Attribute** en el `<html>` o `<body>` para indicar el tema activo
2. **CSS Variables** con valores diferentes según el tema
3. **Servicio Angular** para gestionar el estado y persistencia
4. **Toggle Component** para cambiar entre temas

### Ventajas de este Enfoque
- ✅ No requiere recargar la página
- ✅ Transiciones suaves entre temas
- ✅ Fácil mantenimiento
- ✅ Compatible con Angular Material
- ✅ Persistencia en localStorage

---

## 📝 Pasos de Implementación

### Paso 1: Actualizar `design-system.css` con Variables de Modo Oscuro

**Archivo**: `frontend/src/design-system.css`

Agregar después de las variables `:root` existentes:

```css
/* ==========================================
   MODO OSCURO - VARIABLES DE TEMA
   ========================================== */

[data-theme="dark"] {
  /* Fondos - Modo Oscuro */
  --bg-primary: #1f2937;           /* Gris oscuro principal */
  --bg-secondary: #111827;         /* Gris muy oscuro */
  --bg-tertiary: #0f172a;          /* Gris extremadamente oscuro */
  --bg-dark: #0a0f1a;              /* Fondo más oscuro */

  /* Textos - Modo Oscuro */
  --text-primary: #f9fafb;          /* Texto principal (casi blanco) */
  --text-secondary: #d1d5db;       /* Texto secundario (gris claro) */
  --text-tertiary: #9ca3af;        /* Texto terciario */
  --text-inverse: #111827;         /* Texto inverso (negro) */

  /* Bordes - Modo Oscuro */
  --border-light: #374151;          /* Borde claro (gris oscuro) */
  --border-medium: #4b5563;         /* Borde medio */
  --border-dark: #6b7280;          /* Borde oscuro */
  --border-color: #374151;

  /* Colores de Acento - Ajustados para contraste */
  --accent-600: #60a5fa;            /* Azul más claro para mejor visibilidad */
  --accent-500: #3b82f6;
  --accent-400: #93c5fd;
  --accent-300: #bfdbfe;
  --accent-50: #1e3a8a;             /* Fondo de acento oscuro */

  /* Sombras - Ajustadas para modo oscuro */
  --shadow-xs: 0 1px 2px 0 rgba(0, 0, 0, 0.3);
  --shadow-sm: 0 1px 3px 0 rgba(0, 0, 0, 0.4), 0 1px 2px -1px rgba(0, 0, 0, 0.4);
  --shadow-base: 0 4px 6px -1px rgba(0, 0, 0, 0.4), 0 2px 4px -2px rgba(0, 0, 0, 0.4);
  --shadow-md: 0 10px 15px -3px rgba(0, 0, 0, 0.4), 0 4px 6px -4px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 8px 10px -6px rgba(0, 0, 0, 0.5);
  --shadow-xl: 0 25px 50px -12px rgba(0, 0, 0, 0.6);
}

/* Transición suave para cambios de tema */
* {
  transition: background-color 0.3s ease, color 0.3s ease, border-color 0.3s ease;
}
```

### Paso 2: Crear Servicio de Tema

**Archivo**: `frontend/src/app/services/theme.service.ts`

```typescript
import { Injectable, signal, effect } from '@angular/core';

export type Theme = 'light' | 'dark';

@Injectable({
  providedIn: 'root'
})
export class ThemeService {
  private readonly THEME_KEY = 'fintrack-theme';
  private readonly defaultTheme: Theme = 'light';

  // Signal para el tema actual
  private _currentTheme = signal<Theme>(this.loadTheme());
  public currentTheme = this._currentTheme.asReadonly();

  constructor() {
    // Aplicar tema al inicializar
    this.applyTheme(this._currentTheme());

    // Efecto para aplicar tema cuando cambie
    effect(() => {
      const theme = this._currentTheme();
      this.applyTheme(theme);
      this.saveTheme(theme);
    });
  }

  /**
   * Carga el tema guardado en localStorage o devuelve el tema por defecto
   */
  private loadTheme(): Theme {
    if (typeof window === 'undefined') {
      return this.defaultTheme;
    }

    const savedTheme = localStorage.getItem(this.THEME_KEY) as Theme;
    return savedTheme && (savedTheme === 'light' || savedTheme === 'dark')
      ? savedTheme
      : this.defaultTheme;
  }

  /**
   * Guarda el tema en localStorage
   */
  private saveTheme(theme: Theme): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem(this.THEME_KEY, theme);
    }
  }

  /**
   * Aplica el tema al documento
   */
  private applyTheme(theme: Theme): void {
    if (typeof document === 'undefined') {
      return;
    }

    const htmlElement = document.documentElement;
    
    if (theme === 'dark') {
      htmlElement.setAttribute('data-theme', 'dark');
    } else {
      htmlElement.removeAttribute('data-theme');
    }
  }

  /**
   * Cambia el tema
   */
  toggleTheme(): void {
    const newTheme: Theme = this._currentTheme() === 'light' ? 'dark' : 'light';
    this._currentTheme.set(newTheme);
  }

  /**
   * Establece un tema específico
   */
  setTheme(theme: Theme): void {
    this._currentTheme.set(theme);
  }

  /**
   * Verifica si el tema actual es oscuro
   */
  isDarkMode(): boolean {
    return this._currentTheme() === 'dark';
  }
}
```

### Paso 3: Crear Componente Toggle de Tema

**Archivo**: `frontend/src/app/shared/components/theme-toggle/theme-toggle.component.ts`

```typescript
import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import { ThemeService } from '../../../services/theme.service';

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
      (click)="toggleTheme()"
      [matTooltip]="themeService.isDarkMode() ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'"
      [attr.aria-label]="themeService.isDarkMode() ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'">
      <mat-icon>{{ themeService.isDarkMode() ? 'light_mode' : 'dark_mode' }}</mat-icon>
    </button>
  `,
  styles: [`
    :host {
      display: inline-block;
    }
    
    button {
      color: var(--text-primary);
    }
  `]
})
export class ThemeToggleComponent {
  protected readonly themeService = inject(ThemeService);

  toggleTheme(): void {
    this.themeService.toggleTheme();
  }
}
```

### Paso 4: Integrar Toggle en el Header

**Archivo**: `frontend/src/app/app.component.html`

Agregar el componente en el toolbar:

```html
<mat-toolbar color="primary">
  <!-- ... contenido existente ... -->
  
  <span class="spacer"></span>
  
  <!-- Toggle de tema -->
  <app-theme-toggle></app-theme-toggle>
  
  <!-- ... resto del contenido ... -->
</mat-toolbar>
```

**Archivo**: `frontend/src/app/app.component.ts`

Agregar el import:

```typescript
import { ThemeToggleComponent } from './shared/components/theme-toggle/theme-toggle.component';

@Component({
  // ...
  imports: [
    // ... imports existentes
    ThemeToggleComponent
  ]
})
```

### Paso 5: Actualizar Estilos de Angular Material

**Archivo**: `frontend/src/styles.css`

Agregar estilos para Material en modo oscuro:

```css
/* ==========================================
   MATERIAL DESIGN - MODO OSCURO
   ========================================== */

[data-theme="dark"] {
  /* Toolbar */
  .mat-toolbar.mat-primary {
    --mat-toolbar-container-background-color: var(--primary-900);
    --mat-toolbar-container-text-color: var(--text-primary);
  }

  /* Cards Material */
  .mat-mdc-card {
    background-color: var(--bg-primary) !important;
    color: var(--text-primary) !important;
    border-color: var(--border-light) !important;
  }

  /* Form Fields */
  .mat-mdc-form-field {
    --mdc-filled-text-field-container-color: var(--bg-primary);
    --mdc-filled-text-field-input-text-color: var(--text-primary);
    --mdc-filled-text-field-label-text-color: var(--text-secondary);
    --mdc-filled-text-field-focus-label-text-color: var(--accent-500);
  }

  /* Buttons */
  .mat-mdc-raised-button.mat-primary {
    --mdc-protected-button-container-color: var(--accent-600);
    --mdc-protected-button-label-text-color: var(--text-inverse);
  }

  .mat-mdc-outlined-button.mat-primary {
    --mdc-outlined-button-outline-color: var(--accent-600);
    --mdc-outlined-button-label-text-color: var(--accent-600);
  }

  /* Dialogs */
  .mat-mdc-dialog-container {
    background-color: var(--bg-primary) !important;
    color: var(--text-primary) !important;
  }

  /* Snackbar */
  .mat-mdc-snack-bar-container {
    --mdc-snackbar-container-color: var(--gray-800);
    --mdc-snackbar-supporting-text-color: var(--text-primary);
  }

  /* Tabs */
  .mat-mdc-tab-group {
    --mdc-tab-indicator-active-indicator-color: var(--accent-600);
    --mat-tab-header-active-label-text-color: var(--accent-600);
    --mat-tab-header-inactive-label-text-color: var(--text-secondary);
  }

  /* Chips */
  .mat-mdc-chip {
    background-color: var(--bg-secondary);
    color: var(--text-primary);
  }

  /* Progress Spinner */
  .mat-mdc-progress-spinner {
    --mdc-circular-progress-active-indicator-color: var(--accent-600);
  }
}
```

### Paso 6: Ajustar Componentes Específicos

#### 6.1 Balance Items (Card Detail)

**Archivo**: `frontend/src/app/pages/cards/card-detail/card-detail.component.css`

Asegurar que los balance items mantengan fondo claro en modo oscuro (si es necesario):

```css
/* Mantener fondo claro en balance items incluso en modo oscuro */
[data-theme="dark"] .balance-item {
  background: #ffffff !important;
  color: var(--text-primary);
}

[data-theme="dark"] .due-date {
  background: #ffffff !important;
  border-color: var(--border-light);
  color: var(--text-secondary);
}
```

**Nota**: Si prefieres que los balance items también cambien a modo oscuro, elimina estos estilos específicos.

#### 6.2 Ajustar Colores Hardcodeados

Buscar y reemplazar colores hardcodeados en los archivos CSS:

```bash
# Buscar colores hardcodeados
grep -r "#[0-9a-fA-F]\{6\}" frontend/src/app --include="*.css"
```

Reemplazar con variables CSS cuando sea posible.

---

## 🎨 Variables CSS para Modo Oscuro

### Paleta de Colores Recomendada

```css
[data-theme="dark"] {
  /* Fondos */
  --bg-primary: #1f2937;        /* Gris oscuro principal */
  --bg-secondary: #111827;       /* Gris muy oscuro */
  --bg-tertiary: #0f172a;       /* Gris extremadamente oscuro */

  /* Textos */
  --text-primary: #f9fafb;      /* Casi blanco */
  --text-secondary: #d1d5db;    /* Gris claro */
  --text-tertiary: #9ca3af;     /* Gris medio */

  /* Bordes */
  --border-light: #374151;      /* Gris oscuro */
  --border-medium: #4b5563;    /* Gris medio-oscuro */
  --border-dark: #6b7280;       /* Gris claro-oscuro */

  /* Acentos - Más brillantes para contraste */
  --accent-600: #60a5fa;        /* Azul claro */
  --accent-500: #3b82f6;        /* Azul medio */
}
```

### Contraste Mínimo Recomendado
- **Texto sobre fondo**: Ratio de contraste mínimo 4.5:1 (WCAG AA)
- **Texto grande**: Ratio de contraste mínimo 3:1 (WCAG AA)

---

## 🔧 Servicio de Tema

### Métodos Principales

```typescript
// Obtener tema actual (readonly signal)
themeService.currentTheme()

// Verificar si es modo oscuro
themeService.isDarkMode() // boolean

// Cambiar tema
themeService.toggleTheme()

// Establecer tema específico
themeService.setTheme('dark')
themeService.setTheme('light')
```

### Uso en Componentes

```typescript
import { Component, inject } from '@angular/core';
import { ThemeService } from './services/theme.service';

@Component({...})
export class MyComponent {
  private themeService = inject(ThemeService);
  
  // Acceder al tema actual
  currentTheme = this.themeService.currentTheme;
  
  // Verificar si es oscuro
  isDark = this.themeService.isDarkMode();
  
  // Cambiar tema
  toggleTheme() {
    this.themeService.toggleTheme();
  }
}
```

---

## 🎛️ Componente Toggle

### Características
- ✅ Icono que cambia según el tema (🌙/☀️)
- ✅ Tooltip descriptivo
- ✅ Accesible (aria-label)
- ✅ Estilos consistentes con el tema

### Personalización

Puedes personalizar el componente agregando:
- Animaciones de transición
- Texto en lugar de iconos
- Diferentes posiciones en el header

---

## 🎨 Integración con Angular Material

### Temas de Material

Angular Material tiene su propio sistema de temas. Para una integración completa, considera:

1. **Crear temas Material separados**:
   - `light-theme.scss`
   - `dark-theme.scss`

2. **Aplicar tema Material dinámicamente**:
   ```typescript
   import { inject } from '@angular/core';
   import { DOCUMENT } from '@angular/common';
   
   // En el servicio de tema
   private document = inject(DOCUMENT);
   
   private applyMaterialTheme(theme: Theme) {
     const link = this.document.getElementById('material-theme') as HTMLLinkElement;
     if (link) {
       link.href = theme === 'dark' 
         ? 'dark-theme.css' 
         : 'light-theme.css';
     }
   }
   ```

### Alternativa Simple
El enfoque con CSS variables funciona bien para la mayoría de casos. Angular Material respetará las variables CSS si están bien definidas.

---

## ⚠️ Consideraciones Especiales

### 1. Imágenes y Logos
- Considera tener versiones claras y oscuras de logos
- Usa filtros CSS para invertir imágenes si es necesario:
  ```css
  [data-theme="dark"] .logo {
    filter: invert(1) brightness(1.2);
  }
  ```

### 2. Gráficos y Charts
Si usas bibliotecas de gráficos (Chart.js, D3, etc.):
- Configura colores dinámicos según el tema
- Actualiza los gráficos cuando cambie el tema

### 3. Transiciones
Las transiciones están habilitadas por defecto. Si quieres deshabilitarlas temporalmente:
```css
* {
  transition: none !important;
}
```

### 4. Persistencia
El tema se guarda automáticamente en `localStorage`. Para limpiar:
```typescript
localStorage.removeItem('fintrack-theme');
```

### 5. Detección de Preferencia del Sistema
Puedes agregar detección de preferencia del sistema:

```typescript
// En ThemeService constructor
constructor() {
  // Detectar preferencia del sistema si no hay tema guardado
  if (!localStorage.getItem(this.THEME_KEY)) {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    this._currentTheme.set(prefersDark ? 'dark' : 'light');
  }
  
  // Escuchar cambios en la preferencia del sistema
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    if (!localStorage.getItem(this.THEME_KEY)) {
      this._currentTheme.set(e.matches ? 'dark' : 'light');
    }
  });
}
```

---

## ✅ Testing y Verificación

### Checklist de Verificación

- [ ] Toggle funciona correctamente
- [ ] Tema se aplica inmediatamente sin recargar
- [ ] Tema persiste después de recargar la página
- [ ] Todos los componentes respetan el tema oscuro
- [ ] Contraste de texto es adecuado
- [ ] Botones y enlaces son visibles
- [ ] Formularios son legibles
- [ ] Modales y diálogos tienen fondo correcto
- [ ] Tablas son legibles
- [ ] Cards tienen fondo y texto correctos
- [ ] Iconos son visibles
- [ ] Transiciones son suaves

### Páginas a Verificar

1. **Login/Register**
   - Formularios legibles
   - Botones visibles
   - Links visibles

2. **Dashboard**
   - Cards de resumen
   - Gráficos (si hay)
   - Transacciones recientes

3. **Cuentas**
   - Lista de cuentas
   - Modales de gestión
   - Formularios

4. **Tarjetas**
   - Lista de tarjetas
   - Detalle de tarjeta
   - Balance items

5. **Transacciones**
   - Tabla de transacciones
   - Filtros
   - Formularios

6. **Reportes**
   - Gráficos
   - Tablas
   - Filtros

### Comandos de Testing

```bash
# Verificar que no hay errores de compilación
cd frontend && npm run build

# Verificar linting
npm run lint

# Ejecutar tests (si existen)
npm test
```

---

## 📚 Recursos Adicionales

### Documentación de Referencia
- [CSS Variables](https://developer.mozilla.org/en-US/docs/Web/CSS/Using_CSS_custom_properties)
- [Angular Signals](https://angular.dev/guide/signals)
- [Angular Material Theming](https://material.angular.io/guide/theming)

### Herramientas de Contraste
- [WebAIM Contrast Checker](https://webaim.org/resources/contrastchecker/)
- [Coolors Contrast Checker](https://coolors.co/contrast-checker)

---

## 🚀 Resumen de Archivos a Modificar

1. ✅ `frontend/src/design-system.css` - Agregar variables de modo oscuro
2. ✅ `frontend/src/app/services/theme.service.ts` - Crear servicio
3. ✅ `frontend/src/app/shared/components/theme-toggle/theme-toggle.component.ts` - Crear componente
4. ✅ `frontend/src/app/app.component.ts` - Importar toggle
5. ✅ `frontend/src/app/app.component.html` - Agregar toggle al header
6. ✅ `frontend/src/styles.css` - Agregar estilos Material para modo oscuro
7. ⚠️ Revisar componentes específicos que puedan necesitar ajustes

---

## 💡 Tips Finales

1. **Empieza simple**: Implementa el modo oscuro básico primero, luego ajusta componentes específicos
2. **Prueba en diferentes pantallas**: Verifica que se vea bien en diferentes dispositivos
3. **Mantén consistencia**: Usa siempre las variables CSS, evita colores hardcodeados
4. **Documenta excepciones**: Si algún componente necesita estilos especiales, documenta por qué
5. **Feedback del usuario**: Considera agregar una opción para reportar problemas de contraste

---

**Última actualización**: 2024
**Versión del documento**: 1.0.0
