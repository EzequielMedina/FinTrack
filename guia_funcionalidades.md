# 🚀 Guía de Funcionalidades FinTrack - Sistema de Roles y Permisos

## 👑 **Como Usuario ADMIN, puedes hacer lo siguiente:**

### 1. **Acceso al Sistema**
- Ve a: http://localhost:4200
- Credenciales: `admin@fintrack.com` / `admin123`

### 2. **Navegación Principal**

#### **Dashboard (Página principal)**
- **Ubicación**: Botón "Dashboard" en la barra superior
- **Funciones**:
  - Ver resumen de cuentas y transacciones
  - Acciones rápidas para operaciones comunes
  - **SECCIÓN ESPECIAL DE ADMIN**: Accesos directos a funciones administrativas

#### **Panel de Administración**
- **Ubicación**: Botón "Administración" en la barra superior
- **Funciones**:
  - Dashboard administrativo con métricas del sistema
  - Vista general de usuarios registrados
  - Acceso rápido a gestión de usuarios

#### **Gestión de Usuarios**
- **Ubicación**: Dashboard → "Gestionar Usuarios" OR Administración → "Gestión de Usuarios"
- **Funciones**:
  - ✅ Ver lista completa de usuarios
  - ✅ Crear nuevos usuarios
  - ✅ Editar información de usuarios existentes
  - ✅ Cambiar roles de usuarios (admin/user/operator/treasurer)
  - ✅ Activar/desactivar usuarios
  - ✅ Filtrar y buscar usuarios
  - ✅ Ver detalles completos de cada usuario

### 3. **Funciones Específicas de Admin**

#### **En el Dashboard verás:**
- 👑 Icono de corona junto a tu nombre (indica que eres admin)
- Sección especial "Funciones de Administrador" con:
  - **Gestionar Usuarios**: Acceso directo al CRUD de usuarios
  - **Panel Admin**: Dashboard administrativo completo
  - **Reportes Avanzados**: Analytics del sistema

#### **En la Gestión de Usuarios puedes:**
- **Crear Usuario**: Botón "Nuevo Usuario" 
- **Editar Usuario**: Click en icono de edición
- **Cambiar Rol**: Dropdown con opciones admin/user/operator/treasurer
- **Toggle Status**: Activar/desactivar usuarios
- **Ver Detalles**: Click en cualquier fila de la tabla

### 4. **Permisos por Rol**

#### **Admin (tu rol actual)**
- ✅ Todos los permisos del sistema
- ✅ Crear, leer, actualizar y eliminar usuarios
- ✅ Cambiar roles de otros usuarios
- ✅ Acceso al panel de administración
- ✅ Ver reportes y analytics

#### **Operator**
- ✅ Leer y actualizar usuarios (no crear/eliminar)
- ✅ Cambiar estado de usuarios
- ✅ Ver reportes

#### **Treasurer**
- ✅ Leer usuarios
- ✅ Ver reportes y analytics
- ✅ Gestionar su propio perfil

#### **User**
- ✅ Solo gestionar su propio perfil

### 5. **Cómo Probar el Sistema**

1. **Login como Admin**: http://localhost:4200
2. **Ve al Dashboard**: Deberías ver la sección de admin con coronita 👑
3. **Click en "Gestionar Usuarios"**: Verás la tabla de usuarios (probablemente solo tú por ahora)
4. **Crea un nuevo usuario**: 
   - Click en "Nuevo Usuario"
   - Rellena los datos
   - Asigna un rol diferente (user/operator/treasurer)
5. **Edita el usuario**: Click en el icono de edición
6. **Prueba los filtros**: Busca por email o filtra por rol

### 6. **Elementos Visuales que Confirman tu Rol**
- 👑 Icono de corona en la barra superior
- Botón "Administración" visible solo para ti
- Sección especial en Dashboard
- Acceso completo a gestión de usuarios
- Todas las opciones de CRUD disponibles

### 7. **URLs Importantes**
- **Dashboard**: http://localhost:4200/dashboard
- **Admin Panel**: http://localhost:4200/admin
- **Gestión de Usuarios**: http://localhost:4200/admin/users

## 🔧 **Si algo no funciona:**
1. Verificar que estés logueado como admin
2. Comprobar que veas la coronita 👑 en tu nombre
3. Verificar que el botón "Administración" esté visible
4. Si no ves las funciones de admin, hacer logout y login nuevamente

¡Tu rol de administrador te da control completo sobre el sistema de usuarios y permisos!