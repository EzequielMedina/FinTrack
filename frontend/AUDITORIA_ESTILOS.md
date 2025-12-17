# Auditoría de Estilos - FinTrack Frontend

## Resumen Ejecutivo

Este documento mapea los estilos personalizados encontrados en los componentes y su correspondencia con el sistema de diseño existente.

## Sistema de Diseño Base

- **design-system.css**: Variables CSS (colores, tipografía, espaciado, sombras)
- **components.css**: Componentes reutilizables (cards, botones, badges, inputs, tablas, alertas, KPIs)

## Mapeo de Estilos por Componente

### Dashboard Component

#### Estilos que YA usan el sistema:
- ✅ Variables CSS (--space-*, --text-*, --font-*, --accent-*, etc.)
- ✅ Container principal con padding y max-width

#### Estilos personalizados a migrar:
- `.balance-card` → Puede usar `.kpi-card` o `.info-card-wrapper` con modificadores
- `.stat-card` → Usar `.card` o `.info-card-wrapper`
- `.stat-icon-wrapper` → Usar `.info-card-icon` del sistema
- `.empty-state` → **FALTA en sistema** - Necesita agregarse
- `.loading-container` → **FALTA en sistema** - Necesita agregarse

### Accounts Component

#### Estilos que YA usan el sistema:
- ✅ Variables CSS
- ✅ Container con padding y max-width

#### Estilos personalizados a migrar:
- `.info-card` → Usar `.info-card-wrapper` del sistema
- `.filter-tabs` → **FALTA en sistema** - Necesita agregarse
- `.tab-btn` → **FALTA en sistema** - Necesita agregarse
- `.empty-state` → **FALTA en sistema** - Necesita agregarse

### Cards Component

#### Estilos que YA usan el sistema:
- ✅ Variables CSS
- ✅ Container con padding y max-width

#### Estilos personalizados a migrar:
- `.balance-card` (versión diferente) → Unificar con dashboard
- `.filter-tabs` → **FALTA en sistema** - Necesita agregarse
- `.tab-btn` → **FALTA en sistema** - Necesita agregarse
- `.empty-state` → **FALTA en sistema** - Necesita agregarse

### Transactions Component

#### Estilos que YA usan el sistema:
- ✅ Variables CSS
- ✅ Container con padding y max-width

#### Estilos personalizados a migrar:
- `.tabs` y `.tab` → **FALTA en sistema** - Necesita agregarse
- `.filters-panel` → Usar `.card` del sistema
- `.filter-group` → Usar `.form-group` del sistema
- `.transaction-card` → Usar `.list-item` del sistema con modificadores
- `.empty-state` → **FALTA en sistema** - Necesita agregarse

### Login/Register Components

#### Estilos personalizados a migrar:
- `.login-wrapper`, `.register-wrapper` → Container genérico
- `.login-card`, `.register-card` → Usar `.card` del sistema
- `.form-group` → Ya existe en sistema, verificar uso
- `.input-wrapper` → **FALTA en sistema** - Necesita agregarse
- `.error-text` → Usar `.form-error` del sistema

## Clases Faltantes en el Sistema de Diseño

### Prioridad Alta (usadas en múltiples componentes)

1. **`.empty-state`** - Estado vacío reutilizable
   - Usado en: dashboard, accounts, cards, transactions, reports
   - Incluir: icono, título, descripción, botón de acción

2. **`.loading-container`** / **`.loading-spinner`** - Estados de carga
   - Usado en: todos los componentes
   - Incluir: spinner centrado, mensaje opcional

3. **`.tabs`** y **`.tab`** - Sistema de pestañas
   - Usado en: transactions, accounts, cards
   - Incluir: estados activo, hover, disabled

4. **`.filter-tabs`** / **`.tab-btn`** - Pestañas de filtro
   - Usado en: accounts, cards
   - Variante más simple de tabs

5. **`.input-wrapper`** - Wrapper para inputs con iconos
   - Usado en: login, register
   - Incluir: icono, input, validación visual

### Prioridad Media

6. **`.page-header`** - Header de página estandarizado
   - Ya existe parcialmente en styles.css, verificar

7. **`.content-container`** - Container de contenido genérico
   - Varios componentes tienen su propia versión

## Patrones Comunes a Consolidar

### Cards de Balance
- Dashboard tiene `.balance-card` con colores específicos (ARS/USD)
- Cards tiene `.balance-card` más simple
- **Solución**: Crear variantes del sistema o usar `.kpi-card` con modificadores

### Estados Vacíos
- Cada componente tiene su propia implementación
- **Solución**: Crear `.empty-state` reutilizable en components.css

### Tabs/Filtros
- Transactions usa `.tabs` y `.tab`
- Accounts/Cards usan `.filter-tabs` y `.tab-btn`
- **Solución**: Unificar en un sistema de tabs flexible

### Loading States
- Todos usan `mat-spinner` pero con diferentes contenedores
- **Solución**: Crear `.loading-container` estándar

## Plan de Migración

1. **Extender sistema de diseño** con clases faltantes
2. **Migrar componentes principales** (dashboard, accounts, cards, transactions)
3. **Migrar formularios** (login, register, forms)
4. **Migrar diálogos y modales**
5. **Migrar reportes**
6. **Limpieza final** - Eliminar código redundante

## Métricas Objetivo

- **Reducción de código CSS**: 30%+
- **Uso de sistema de diseño**: 90%+
- **Clases personalizadas restantes**: Solo para casos muy específicos

