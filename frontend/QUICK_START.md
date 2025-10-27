# FinTrack Design System - Quick Start

## 🚀 Inicio Rápido (5 minutos)

### 1. Archivos Ya Configurados ✅
- `design-system.css` - Sistema completo de variables
- `components.css` - Componentes reutilizables
- `styles.css` - Estilos globales actualizados
- `angular.json` - Configurado para importar archivos

### 2. Variables Más Usadas

#### Colores
```css
/* Textos */
var(--text-primary)      /* Negro principal #111827 */
var(--text-secondary)    /* Gris medio #4b5563 */
var(--text-tertiary)     /* Gris claro #9ca3af */

/* Fondos */
var(--bg-primary)        /* Blanco #ffffff */
var(--bg-secondary)      /* Gris muy claro #f9fafb */

/* Acento (Azul Profesional) */
var(--accent-600)        /* Azul principal #2563eb ⭐ */
var(--accent-700)        /* Azul hover #1d4ed8 */
var(--accent-100)        /* Fondo azul claro #dbeafe */

/* Estados */
var(--success-500)       /* Verde #22c55e */
var(--warning-500)       /* Naranja #f97316 */
var(--error-500)         /* Rojo #ef4444 */
var(--info-500)          /* Cyan #0ea5e9 */
```

#### Espaciado
```css
var(--space-2)   /* 8px */
var(--space-4)   /* 16px - Base ⭐ */
var(--space-6)   /* 24px - Padding cards ⭐ */
var(--space-8)   /* 32px - Separación secciones ⭐ */
```

#### Tipografía
```css
var(--text-sm)       /* 14px */
var(--text-base)     /* 16px ⭐ */
var(--text-xl)       /* 20px */
var(--text-2xl)      /* 24px */

var(--font-regular)   /* 400 */
var(--font-medium)    /* 500 */
var(--font-semibold)  /* 600 ⭐ */
var(--font-bold)      /* 700 */
```

#### Otros
```css
var(--radius-lg)         /* 16px - Cards ⭐ */
var(--shadow-sm)         /* Sombra cards ⭐ */
var(--shadow-md)         /* Sombra hover ⭐ */
var(--transition-base)   /* 200ms */
```

---

## 📦 Componentes Listos para Usar

### Botones
```html
<!-- Principal -->
<button class="btn btn-primary">Guardar</button>

<!-- Secundario -->
<button class="btn btn-secondary">Cancelar</button>

<!-- Outline -->
<button class="btn btn-outline">Ver más</button>

<!-- Tamaños -->
<button class="btn btn-primary btn-sm">Pequeño</button>
<button class="btn btn-primary btn-lg">Grande</button>
```

### Cards
```html
<!-- Card simple -->
<div class="card p-6">
  <h3>Título</h3>
  <p>Contenido</p>
</div>

<!-- Card con hover -->
<div class="card card-elevated p-6">
  <h3>Card Elevada</h3>
</div>
```

### Badges/Tags
```html
<span class="badge badge-success">Activo</span>
<span class="badge badge-warning">Pendiente</span>
<span class="badge badge-error">Cancelado</span>
```

### KPI Card
```html
<div class="kpi-card">
  <div class="kpi-header">
    <span class="kpi-label">Balance</span>
  </div>
  <div class="kpi-value">$125,430</div>
  <div class="kpi-footer">
    <span class="kpi-trend up">+12.5%</span>
  </div>
</div>
```

### Alertas
```html
<div class="alert success">
  <div class="alert-content">
    <h4 class="alert-title">¡Éxito!</h4>
    <p class="alert-message">Operación completada.</p>
  </div>
</div>
```

---

## 🎨 Patrones de Diseño Comunes

### Layout de Página
```html
<div class="app-container">
  <!-- Header -->
  <div class="page-header">
    <div>
      <h1 class="page-title">Mi Página</h1>
      <p class="page-subtitle">Descripción</p>
    </div>
    <button class="btn btn-primary">Acción</button>
  </div>

  <!-- Contenido -->
  <div class="grid grid-cols-3 gap-6">
    <div class="card p-6">...</div>
    <div class="card p-6">...</div>
    <div class="card p-6">...</div>
  </div>
</div>
```

### Formulario
```html
<div class="form-section">
  <h2 class="form-section-title">Datos Personales</h2>
  
  <div class="form-group">
    <label class="form-label required">Nombre</label>
    <input type="text" class="input" placeholder="Ingresa tu nombre">
    <span class="form-help">Nombre completo</span>
  </div>
</div>
```

### Lista
```html
<div class="list-item">
  <div class="list-item-icon">
    <mat-icon>account_balance</mat-icon>
  </div>
  <div class="list-item-content">
    <h4 class="list-item-title">Cuenta Principal</h4>
    <p class="list-item-subtitle">**** 1234</p>
  </div>
  <div class="list-item-action">
    <button class="btn btn-sm btn-outline">Ver</button>
  </div>
</div>
```

### Tabla
```html
<div class="table-container">
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
        <td>Cuenta 1</td>
        <td>$10,000</td>
        <td><span class="badge badge-success">Activa</span></td>
      </tr>
    </tbody>
  </table>
</div>
```

---

## 🎯 Clases Utilitarias

### Flexbox
```html
<div class="flex items-center justify-between gap-4">
  <span>Izquierda</span>
  <span>Derecha</span>
</div>
```

### Grid
```html
<div class="grid grid-cols-3 gap-6">
  <div>Col 1</div>
  <div>Col 2</div>
  <div>Col 3</div>
</div>
```

### Texto
```html
<h1 class="text-3xl font-bold text-primary">Título Grande</h1>
<p class="text-base text-secondary">Texto normal</p>
<span class="text-sm text-tertiary">Texto pequeño</span>
```

### Espaciado
```html
<div class="p-6 mb-4">Padding 24px, margin-bottom 16px</div>
```

---

## 🖼️ Iconos SVG

### Ubicación
`frontend/src/assets/icons/`

### Uso
```html
<!-- Icono simple -->
<img src="assets/icons/dashboard.svg" class="icon icon-md">

<!-- Icono circular con fondo -->
<div class="icon-circle icon-circle-md" 
     style="background: var(--accent-100); color: var(--accent-700)">
  <img src="assets/icons/wallet.svg" class="icon">
</div>
```

### Iconos Disponibles
- dashboard.svg
- account.svg
- wallet.svg
- transaction.svg
- card.svg
- report.svg
- chart-up.svg
- chart-down.svg
- users.svg
- settings.svg
- chatbot.svg
- calendar.svg
- dollar.svg
- arrow-up.svg
- arrow-down.svg
- check-circle.svg
- alert-circle.svg

---

## 🎨 Ejemplos de Estilo CSS

### Card Personalizada
```css
.mi-card {
  background: var(--bg-primary);
  padding: var(--space-6);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-base);
}

.mi-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}
```

### Botón Personalizado
```css
.mi-boton {
  padding: var(--space-3) var(--space-6);
  background: var(--accent-600);
  color: white;
  border: none;
  border-radius: var(--radius-base);
  font-weight: var(--font-semibold);
  cursor: pointer;
  transition: background var(--transition-base);
}

.mi-boton:hover {
  background: var(--accent-700);
}
```

### Badge de Estado
```css
.estado-badge {
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  text-transform: uppercase;
}

.estado-badge.activo {
  background: var(--success-100);
  color: var(--success-700);
}
```

---

## 📱 Responsive Design

```css
/* Mobile First */
.mi-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: var(--space-4);
}

/* Tablet y superior */
@media (min-width: 768px) {
  .mi-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--space-6);
  }
}

/* Desktop */
@media (min-width: 1024px) {
  .mi-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
```

---

## ✅ Checklist de Migración

Cuando actualices un componente existente:

- [ ] Reemplazar colores fijos por variables (`var(--accent-600)`)
- [ ] Usar espaciado del sistema (`var(--space-6)`)
- [ ] Aplicar tipografía del sistema (`var(--text-xl)`)
- [ ] Usar border-radius del sistema (`var(--radius-lg)`)
- [ ] Aplicar sombras del sistema (`var(--shadow-sm)`)
- [ ] Usar transiciones del sistema (`var(--transition-base)`)
- [ ] Agregar estados hover/focus
- [ ] Verificar responsive design
- [ ] Probar accesibilidad

---

## 🆘 Problemas Comunes

### Variables no funcionan
```
Error: Variable CSS no aplicada
Solución: Verificar que angular.json incluya design-system.css
```

### Iconos no cargan
```
Error: 404 en assets/icons/
Solución: Verificar que la carpeta icons exista en assets
```

### Estilos no se aplican
```
Error: Estilos del sistema ignorados
Solución: Verificar orden de importación en angular.json
```

---

## 📚 Recursos

- **Guía Completa**: `DESIGN_SYSTEM.md`
- **Preview Visual**: `design-system-preview.html`
- **Resumen**: `FINTRACK_DESIGN_SYSTEM_SUMMARY.md`

---

## 🎓 Mejores Prácticas

1. **Siempre usar variables** - No valores hardcoded
2. **Componentes reutilizables** - No duplicar CSS
3. **Mobile First** - Responsive desde diseño
4. **Consistencia** - Seguir el sistema
5. **Accesibilidad** - Contraste y navegación

---

**¡Listo para usar!** 🚀

Empieza aplicando el sistema en tus componentes siguiendo los ejemplos de arriba.
