# Implementación de Campo due_date en Tarjetas

## Resumen de Cambios

Se ha implementado el campo `due_date` (fecha de pago) en todo el flujo de creación y edición de tarjetas, desde la base de datos hasta el frontend.

## ✅ Cambios Realizados

### 1. **Base de Datos** ✅
- **Archivo**: `database/migrations/04_V4__cards.sql`
- **Estado**: ✅ **Ya existía** - Campo `due_date DATE NULL` en tabla `cards`
- **Funcionalidad**: Almacena la fecha de vencimiento mensual de tarjetas de crédito

### 2. **Backend - Account Service** ✅

#### **DTOs (`card/dto/dto.go`)** ✅
- `CreateCardRequest.DueDate *time.Time` - Ya existía
- `CardResponse.DueDate *time.Time` - Ya existía  
- `UpdateCardRequest` - Listo para futuras actualizaciones

#### **Entidades (`domain/entities/account.go`)** ✅
- `Card.DueDate *time.Time` - Ya existía con mapping GORM correcto

#### **Lógica de Negocio (`services/card_service.go`)** ✅ **NUEVO**
```go
// Set default due date if not provided and it's a credit card
dueDate := req.DueDate
if dueDate == nil && req.CardType == "credit" {
    // Set due date to the 5th of next month by default
    now := time.Now()
    nextMonth := now.AddDate(0, 1, 0)
    defaultDueDate := time.Date(nextMonth.Year(), nextMonth.Month(), 5, 0, 0, 0, 0, nextMonth.Location())
    dueDate = &defaultDueDate
}
```

**Funcionalidad**: Si no se proporciona fecha de pago para tarjetas de crédito, automáticamente se asigna el día 5 del mes siguiente.

### 3. **Frontend - Angular** ✅

#### **Modelos (`models/card.model.ts`)** ✅ **ACTUALIZADO**
```typescript
export interface CardFormData {
  // ... otros campos
  dueDate?: number;  // Día del mes (1-31)
}
```

#### **Formulario (`card-form.component.ts`)** ✅ **ACTUALIZADO**

**Inicialización del formulario**:
```typescript
this.cardForm = this.fb.group({
  // ... otros campos
  dueDate: [''] // Nuevo campo para fecha de pago
});
```

**Lógica de creación**:
```typescript
// Calcular dueDate si se proporciona un día del mes
let dueDate: string | undefined;
if (formData.dueDate && formData.cardType === CardType.CREDIT) {
  const dayOfMonth = parseInt(formData.dueDate.toString());
  if (dayOfMonth >= 1 && dayOfMonth <= 31) {
    const today = new Date();
    const nextMonth = new Date(today.getFullYear(), today.getMonth() + 1, dayOfMonth);
    dueDate = nextMonth.toISOString().split('T')[0]; // Format: YYYY-MM-DD
  }
}
```

**Edición de tarjetas**:
```typescript
dueDate: card.dueDate ? new Date(card.dueDate).getDate() : null // Extraer solo el día del mes
```

#### **Template (`card-form.component.html`)** ✅ **NUEVO**
```html
<!-- Due Date (only for credit cards) -->
@if (cardForm.get('cardType')?.value === 'credit') {
  <mat-form-field appearance="outline" class="full-width">
    <mat-label>Día de pago mensual</mat-label>
    <mat-select formControlName="dueDate">
      <mat-option value="">Sin especificar (se asignará el día 5)</mat-option>
      @for (day of getDaysOfMonth(); track day) {
        <mat-option [value]="day">
          Día {{ day }} de cada mes
        </mat-option>
      }
    </mat-select>
    <mat-icon matSuffix>calendar_today</mat-icon>
    <mat-hint>Día del mes en que vence el pago de la tarjeta de crédito</mat-hint>
    <mat-error>{{ getFieldError('dueDate') }}</mat-error>
  </mat-form-field>
}
```

**Método helper**:
```typescript
getDaysOfMonth(): number[] {
  return Array.from({ length: 31 }, (_, i) => i + 1);
}
```

## 🎯 Funcionalidad Implementada

### **Para el Usuario:**
1. **Creación de Tarjeta de Crédito**:
   - Puede seleccionar el día del mes (1-31) para el pago
   - Si no selecciona, automáticamente se asigna el día 5
   - Solo aparece para tarjetas de crédito

2. **Edición de Tarjeta**:
   - Puede modificar el día de pago
   - Se muestra el día actual configurado

3. **Validación**:
   - Campo opcional para tarjetas de crédito
   - No aparece para tarjetas de débito

### **Lógica de Negocio:**
- **Backend**: Si no se proporciona fecha → día 5 del mes siguiente
- **Frontend**: Usuario selecciona día (1-31) → se convierte a fecha completa del mes siguiente
- **Almacenamiento**: Fecha completa en formato `YYYY-MM-DD`
- **Visualización**: Solo se muestra el día del mes

## 📍 Flujo Completo

```
1. Usuario crea tarjeta de crédito
2. Selecciona "Día 15" en el formulario
3. Frontend convierte a fecha: "2025-11-15" (mes siguiente)
4. Backend recibe la fecha y la almacena
5. Si no selecciona día, backend asigna día 5 por defecto
6. Al editar, frontend extrae el día (15) de la fecha almacenada
```

## 🔧 Testing

### **Para Probar:**
1. **Crear Tarjeta de Crédito**:
   - Seleccionar día de pago → Verificar que se guarda
   - No seleccionar día → Verificar que se asigna día 5

2. **Editar Tarjeta**:
   - Abrir tarjeta existente → Verificar que muestra el día correcto
   - Cambiar día → Verificar que se actualiza

3. **Tarjeta de Débito**:
   - Verificar que NO aparece el campo de fecha de pago

## ✅ Estado Final

- **✅ Base de Datos**: Campo `due_date` disponible
- **✅ Backend**: Lógica de asignación automática implementada
- **✅ Frontend**: Formulario con selector de día (1-31)
- **✅ Validaciones**: Solo para tarjetas de crédito
- **✅ UX**: Campo opcional con valor por defecto
- **✅ Edición**: Extrae y muestra el día del mes correctamente

El campo `due_date` está completamente integrado en todo el flujo de tarjetas, con lógica de negocio apropiada y una experiencia de usuario intuitiva. 🚀