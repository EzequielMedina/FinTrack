# FinTrack - Sistema de Diseño Empresarial

> Guía completa del sistema de diseño moderno y profesional para FinTrack

## 📋 Tabla de Contenidos

1. [Introducción](#introducción)
2. [Paleta de Colores](#paleta-de-colores)
3. [Tipografía](#tipografía)
4. [Espaciado y Layout](#espaciado-y-layout)
5. [Componentes](#componentes)
6. [Iconos SVG](#iconos-svg)
7. [Mejores Prácticas](#mejores-prácticas)

---

## 🎨 Introducción

El sistema de diseño de FinTrack está basado en principios empresariales modernos:

- **Minimalismo**: Diseño limpio sin elementos innecesarios
- **Profesionalismo**: Paleta neutral que transmite confianza
- **Consistencia**: Componentes reutilizables y uniformes
- **Accesibilidad**: Contraste adecuado y legibilidad
- **Escalabilidad**: Sistema flexible para futuras expansiones

### Archivos Principales

```
frontend/src/
├── design-system.css    # Variables y sistema base
├── components.css       # Componentes reutilizables
├── styles.css          # Estilos globales y Material customization
└── assets/
    └── icons/          # Biblioteca de iconos SVG
```

---

## 🎨 Paleta de Colores

### Colores Primarios (Grises Corporativos)

Usados para textos, fondos y elementos estructurales.

```css
--primary-900: #0f172a  /* Texto principal oscuro */
--primary-800: #1e293b  /* Fondos oscuros */
--primary-700: #334155
--primary-600: #475569
--primary-500: #64748b  /* Grises medios */
--primary-400: #94a3b8
--primary-300: #cbd5e1
--primary-200: #e2e8f0  /* Bordes sutiles */
--primary-100: #f1f5f9  /* Fondos claros */
--primary-50: #f8fafc   /* Fondos muy claros */
```

### Colores de Acento (Azul Profesional)

Para acciones principales, enlaces y elementos interactivos.

```css
--accent-900: #1e3a8a  /* Azul muy oscuro */
--accent-800: #1e40af
--accent-700: #1d4ed8
--accent-600: #2563eb  /* Azul principal ⭐ */
--accent-500: #3b82f6
--accent-400: #60a5fa
--accent-300: #93c5fd
--accent-200: #bfdbfe
--accent-100: #dbeafe  /* Fondos de acento */
--accent-50: #eff6ff
```

### Colores Semánticos

#### Success (Verde) - Éxito, confirmación, positivo
```css
--success-700: #15803d  /* Texto */
--success-500: #22c55e  /* Iconos/Botones */
--success-100: #dcfce7  /* Fondos */
```

#### Warning (Naranja) - Advertencias, pendiente
```css
--warning-700: #c2410c  /* Texto */
--warning-500: #f97316  /* Iconos/Botones */
--warning-100: #ffedd5  /* Fondos */
```

#### Error (Rojo) - Errores, eliminación, negativo
```css
--error-700: #b91c1c    /* Texto */
--error-500: #ef4444    /* Iconos/Botones */
--error-100: #fee2e2    /* Fondos */
```

#### Info (Cyan) - Información, neutral
```css
--info-700: #0369a1     /* Texto */
--info-500: #0ea5e9     /* Iconos/Botones */
--info-100: #e0f2fe     /* Fondos */
```

### Ejemplos de Uso

```css
/* Botón principal */
background: var(--accent-600);
color: white;

/* Card con borde de éxito */
border-left: 4px solid var(--success-500);
background: var(--success-50);

/* Texto secundario */
color: var(--text-secondary); /* Equivale a var(--gray-600) */
```

---

## 📝 Tipografía

### Familias de Fuentes

```css
--font-primary: 'Inter'      /* Uso general */
--font-secondary: 'Poppins'  /* Títulos destacados */
--font-mono: 'JetBrains Mono' /* Código/Números */
```

### Escala de Tamaños

| Variable | Tamaño | Uso |
|----------|--------|-----|
| `--text-xs` | 12px | Labels pequeños, metadata |
| `--text-sm` | 14px | Texto secundario |
| `--text-base` | 16px | Texto principal |
| `--text-lg` | 18px | Subtítulos |
| `--text-xl` | 20px | Títulos de sección |
| `--text-2xl` | 24px | Títulos de página |
| `--text-3xl` | 30px | Títulos principales |
| `--text-4xl` | 36px | Hero titles |
| `--text-5xl` | 48px | Números grandes |

### Pesos de Fuente

```css
--font-light: 300      /* Poco uso */
--font-regular: 400    /* Texto normal */
--font-medium: 500     /* Énfasis sutil */
--font-semibold: 600   /* Títulos, botones */
--font-bold: 700       /* Números, highlights */
```

### Ejemplo

```css
.page-title {
  font-size: var(--text-3xl);
  font-weight: var(--font-bold);
  color: var(--text-primary);
}

.card-subtitle {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--text-secondary);
}
```

---

## 📏 Espaciado y Layout

### Sistema de Espaciado (Base 4px)

```css
--space-1: 4px
--space-2: 8px
--space-3: 12px
--space-4: 16px   /* Espaciado base */
--space-5: 20px
--space-6: 24px
--space-8: 32px
--space-10: 40px
--space-12: 48px
--space-16: 64px
--space-20: 80px
```

### Border Radius

```css
--radius-sm: 4px      /* Inputs pequeños */
--radius-base: 8px    /* Botones, badges */
--radius-md: 12px     /* Iconos circulares */
--radius-lg: 16px     /* Cards ⭐ */
--radius-xl: 24px     /* Modales, secciones */
--radius-full: 9999px /* Circular completo */
```

### Sombras

```css
--shadow-xs    /* Hover sutil */
--shadow-sm    /* Cards default ⭐ */
--shadow-base  /* Cards elevadas */
--shadow-md    /* Hover cards ⭐ */
--shadow-lg    /* Modales, dropdowns */
--shadow-xl    /* Overlays importantes */
```

### Breakpoints

```css
--breakpoint-sm: 640px   /* Móviles */
--breakpoint-md: 768px   /* Tablets ⭐ */
--breakpoint-lg: 1024px  /* Desktop */
--breakpoint-xl: 1280px
--breakpoint-2xl: 1536px
```

---

## 🧩 Componentes

### Botones

```html
<!-- Botón Principal -->
<button class="btn btn-primary">
  Acción Principal
</button>

<!-- Botón Secundario -->
<button class="btn btn-secondary">
  Cancelar
</button>

<!-- Botón Outline -->
<button class="btn btn-outline">
  Ver más
</button>

<!-- Tamaños -->
<button class="btn btn-primary btn-sm">Pequeño</button>
<button class="btn btn-primary">Normal</button>
<button class="btn btn-primary btn-lg">Grande</button>
```

### Cards

```html
<!-- Card básica -->
<div class="card">
  <h3>Título de Card</h3>
  <p>Contenido de la card...</p>
</div>

<!-- Card con elevación -->
<div class="card card-elevated">
  <h3>Card Elevada</h3>
</div>

<!-- Info Card con icono -->
<div class="info-card-wrapper">
  <div class="info-card-header">
    <div class="info-card-icon primary">
      <mat-icon>account_balance</mat-icon>
    </div>
    <h3 class="info-card-title">Cuentas Activas</h3>
  </div>
  <div class="info-card-content">
    <p>Tienes 3 cuentas activas</p>
  </div>
</div>
```

### Badges

```html
<span class="badge badge-success">Activo</span>
<span class="badge badge-warning">Pendiente</span>
<span class="badge badge-error">Cancelado</span>
<span class="badge badge-neutral">Neutral</span>
```

### KPI Cards

```html
<div class="kpi-card">
  <div class="kpi-header">
    <span class="kpi-label">Balance Total</span>
    <div class="kpi-icon">
      <mat-icon>account_balance_wallet</mat-icon>
    </div>
  </div>
  <div class="kpi-value">$125,430.50</div>
  <div class="kpi-footer">
    <span class="kpi-trend up">
      <mat-icon>trending_up</mat-icon>
      +12.5%
    </span>
    <span class="text-tertiary">vs mes anterior</span>
  </div>
</div>
```

### Alertas

```html
<div class="alert success">
  <div class="alert-icon">
    <mat-icon>check_circle</mat-icon>
  </div>
  <div class="alert-content">
    <h4 class="alert-title">¡Éxito!</h4>
    <p class="alert-message">La operación se completó correctamente.</p>
  </div>
</div>
```

### Tablas Empresariales

```html
<div class="table-container">
  <div class="table-wrapper">
    <table class="enterprise-table">
      <thead>
        <tr>
          <th>Nombre</th>
          <th>Monto</th>
          <th>Estado</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>Cuenta Principal</td>
          <td>$10,000.00</td>
          <td><span class="badge badge-success">Activa</span></td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
```

### List Items

```html
<div class="list-item">
  <div class="list-item-icon">
    <mat-icon>account_balance</mat-icon>
  </div>
  <div class="list-item-content">
    <h4 class="list-item-title">Cuenta Ahorro</h4>
    <p class="list-item-subtitle">Banco Nacional - **** 1234</p>
  </div>
  <div class="list-item-action">
    <button class="btn btn-sm btn-outline">Ver</button>
  </div>
</div>
```

---

## 🎯 Iconos SVG

### Ubicación

Los iconos SVG personalizados están en: `frontend/src/assets/icons/`

### Iconos Disponibles

| Icono | Archivo | Uso |
|-------|---------|-----|
| Dashboard | `dashboard.svg` | Menú dashboard |
| Cuentas | `account.svg` | Sección de cuentas |
| Wallet | `wallet.svg` | Billetera/Balance |
| Transacciones | `transaction.svg` | Lista de movimientos |
| Tarjetas | `card.svg` | Tarjetas de crédito/débito |
| Reportes | `report.svg` | Informes y documentos |
| Gráfico Up | `chart-up.svg` | Tendencia positiva |
| Gráfico Down | `chart-down.svg` | Tendencia negativa |
| Usuarios | `users.svg` | Admin/Gestión usuarios |
| Configuración | `settings.svg` | Ajustes |
| Chatbot | `chatbot.svg` | Asistente virtual |
| Calendario | `calendar.svg` | Fechas/Períodos |
| Dólar | `dollar.svg` | Moneda/Finanzas |
| Flecha Arriba | `arrow-up.svg` | Deposito/Ingreso |
| Flecha Abajo | `arrow-down.svg` | Retiro/Egreso |
| Check Circle | `check-circle.svg` | Éxito/Completado |
| Alert Circle | `alert-circle.svg` | Alerta/Advertencia |

### Uso de Iconos

```html
<!-- Icono inline -->
<img src="assets/icons/dashboard.svg" alt="Dashboard" class="icon icon-md">

<!-- Icono circular -->
<div class="icon-circle icon-circle-md" style="background: var(--accent-100); color: var(--accent-700)">
  <img src="assets/icons/wallet.svg" alt="Wallet" class="icon">
</div>
```

### Tamaños de Icono

```css
.icon-sm  /* 16x16px */
.icon-md  /* 24x24px - Default */
.icon-lg  /* 32x32px */
.icon-xl  /* 48x48px */
```

---

## ✅ Mejores Prácticas

### 1. Consistencia de Colores

```css
/* ✅ CORRECTO */
color: var(--text-primary);
background: var(--accent-600);

/* ❌ INCORRECTO */
color: #333;
background: #0e7490;
```

### 2. Espaciado

```css
/* ✅ CORRECTO - Usar variables */
padding: var(--space-6);
gap: var(--space-4);

/* ❌ INCORRECTO - Valores arbitrarios */
padding: 22px;
gap: 13px;
```

### 3. Tipografía

```css
/* ✅ CORRECTO */
font-size: var(--text-xl);
font-weight: var(--font-semibold);

/* ❌ INCORRECTO */
font-size: 19px;
font-weight: 550;
```

### 4. Border Radius

```css
/* ✅ CORRECTO */
border-radius: var(--radius-lg);

/* ❌ INCORRECTO */
border-radius: 14px;
```

### 5. Sombras

```css
/* ✅ CORRECTO */
box-shadow: var(--shadow-md);

/* ❌ INCORRECTO */
box-shadow: 0 5px 15px rgba(0,0,0,0.12);
```

### 6. Transiciones

```css
/* ✅ CORRECTO */
transition: all var(--transition-base);

/* ❌ INCORRECTO */
transition: all 0.25s ease-in-out;
```

### 7. Componentes Reutilizables

```html
<!-- ✅ CORRECTO - Usar clases del sistema -->
<div class="card card-elevated p-6">
  <h3 class="text-xl font-semibold text-primary">Título</h3>
  <p class="text-sm text-secondary">Descripción</p>
</div>

<!-- ❌ INCORRECTO - Estilos inline -->
<div style="padding: 24px; box-shadow: 0 4px 6px rgba(0,0,0,0.1)">
  <h3 style="font-size: 20px; color: #333">Título</h3>
</div>
```

### 8. Responsive Design

```css
/* ✅ CORRECTO - Mobile first */
.container {
  padding: var(--space-4);
}

@media (min-width: 768px) {
  .container {
    padding: var(--space-8);
  }
}
```

### 9. Accesibilidad

- Mantener contraste mínimo 4.5:1 para textos
- Usar colores semánticos para feedback
- Incluir labels descriptivos
- Asegurar navegación por teclado

### 10. Performance

- Minimizar uso de sombras complejas
- Usar `will-change` solo cuando necesario
- Optimizar transiciones (evitar `all` cuando posible)

---

## 🚀 Implementación

### Importar el Sistema

En `angular.json`:

```json
"styles": [
  "src/design-system.css",
  "src/components.css",
  "src/styles.css"
]
```

### Ejemplo Completo

```html
<div class="app-container">
  <!-- Header -->
  <div class="page-header">
    <div>
      <h1 class="page-title">Dashboard</h1>
      <p class="page-subtitle">Bienvenido de vuelta</p>
    </div>
    <button class="btn btn-primary">
      <mat-icon>add</mat-icon>
      Nueva Cuenta
    </button>
  </div>

  <!-- KPIs Grid -->
  <div class="grid grid-cols-3 gap-6">
    <div class="kpi-card">
      <div class="kpi-header">
        <span class="kpi-label">Balance Total</span>
        <div class="kpi-icon">
          <img src="assets/icons/wallet.svg" class="icon">
        </div>
      </div>
      <div class="kpi-value">$125,430</div>
      <div class="kpi-footer">
        <span class="kpi-trend up">+12.5%</span>
      </div>
    </div>
  </div>

  <!-- Content -->
  <div class="card card-elevated mt-8">
    <h3 class="text-2xl font-semibold mb-4">Transacciones Recientes</h3>
    <!-- Contenido... -->
  </div>
</div>
```

---

## 📚 Recursos Adicionales

- **Figma**: [Próximamente] Sistema de diseño en Figma
- **Storybook**: [Próximamente] Biblioteca de componentes interactiva
- **Testing**: Usar clases CSS para testing automatizado

---

## 🔄 Versionado

**Versión Actual**: 1.0.0

### Changelog

- **v1.0.0** (2025-10-20): Lanzamiento inicial del sistema de diseño empresarial

---

## 👥 Contribución

Para modificar o extender el sistema de diseño:

1. Crear nueva rama desde `main`
2. Modificar solo archivos del sistema de diseño
3. Documentar cambios en este README
4. Crear PR con revisión de diseño

---

**Desarrollado para FinTrack** | Sistema de Gestión Financiera Empresarial
