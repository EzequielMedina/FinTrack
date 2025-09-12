# FinTrack - User Service

Microservicio completo de gestión de usuarios para FinTrack. Provee **CRUD completo** de usuarios, perfiles extendidos, sistema de roles jerárquico, autenticación JWT y autorización granular.

## 🎯 Funcionalidades Implementadas

### ✅ CRUD Completo de Usuarios
- **Create**: Crear nuevos usuarios con validaciones completas
- **Read**: Obtener usuarios por ID, email, rol, con paginación
- **Update**: Actualizar información de usuarios con autorización
- **Delete**: Eliminar usuarios con validaciones de seguridad

### ✅ Sistema de Perfiles Extendido
- **Información personal**: Teléfono, fecha de nacimiento, foto de perfil
- **Dirección completa**: Calle, ciudad, estado, código postal, país  
- **Preferencias**: Idioma, zona horaria, notificaciones email/SMS

### ✅ Sistema de Roles Robusto
- **Roles disponibles**: `user`, `operator`, `admin`, `treasurer`
- **Jerarquía de permisos**: Sistema de niveles con validaciones automáticas
- **Autorización granular**: Control de acceso basado en roles y propiedad

### ✅ Validaciones y Business Rules
- **Seguridad**: Los admins no pueden eliminarse a sí mismos
- **Autorización**: Control de acceso basado en roles y ownership
- **Validaciones**: Email único, contraseñas seguras, datos requeridos
- **Integridad**: Verificación de emails únicos y estados consistentes

## 🏗️ Arquitectura (Clean Architecture + SOLID)

```
internal/
├── core/
│   ├── domain/entities/user/     # Entidades de dominio
│   ├── service/                  # Lógica de negocio
│   ├── providers/user/           # Interfaces de repositorio
│   └── errors/                   # Errores de dominio
├── infrastructure/
│   ├── repositories/mysql/       # Persistencia MySQL
│   └── entrypoints/
│       ├── handlers/            # HTTP handlers + DTOs
│       ├── middleware/          # Autenticación JWT
│       └── router/              # Routing y configuración
├── config/                      # Configuración de aplicación
└── app/                         # Contenedor de dependencias
```

## 🔗 API Endpoints

### Autenticación (Públicos)
```bash
POST /api/auth/register    # Registro de usuario
POST /api/auth/login       # Autenticación
```

### Usuario Actual (Autenticado)
```bash
GET  /api/me              # Información del usuario actual
```

### Gestión de Usuarios (Autorización por roles)
```bash
POST   /api/users                    # Crear usuario (admin)
GET    /api/users                    # Listar usuarios (admin)
GET    /api/users/:id                # Obtener usuario (owner/admin)
PUT    /api/users/:id                # Actualizar usuario (owner/admin)
DELETE /api/users/:id                # Eliminar usuario (admin)
```

### Gestión de Perfiles
```bash
PUT    /api/users/:id/profile        # Actualizar perfil (owner/admin)
```

### Gestión de Roles y Estado (Solo Admins)
```bash
PUT    /api/users/:id/role           # Cambiar rol
PUT    /api/users/:id/status         # Activar/desactivar usuario
PUT    /api/users/:id/password       # Cambiar contraseña
```

### Consultas Especializadas
```bash
GET    /api/users/role/:role         # Filtrar usuarios por rol (admin)
```

### Paginación
Todos los endpoints que retornan listas soportan paginación:
```bash
GET /api/users?page=1&pageSize=20
```

## 📝 Ejemplos de Uso

### Registro de Usuario
```bash
curl -X POST http://localhost:8081/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fintrack.com",
    "password": "SecurePass123!",
    "firstName": "Admin",
    "lastName": "User"
  }'
```

### Login
```bash
curl -X POST http://localhost:8081/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fintrack.com", 
    "password": "SecurePass123!"
  }'
```

### Crear Usuario (Admin)
```bash
curl -X POST http://localhost:8081/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "email": "treasurer@fintrack.com",
    "password": "SecurePass123!",
    "firstName": "John",
    "lastName": "Treasurer",
    "role": "treasurer"
  }'
```

### Actualizar Perfil
```bash
curl -X PUT http://localhost:8081/api/users/user-id/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "phone": "+1234567890",
    "address": {
      "street": "123 Main St",
      "city": "Springfield",
      "state": "IL",
      "postalCode": "62701",
      "country": "USA"
    },
    "preferences": {
      "language": "es",
      "timezone": "America/Argentina/Buenos_Aires",
      "notificationEmail": true,
      "notificationSMS": false
    }
  }'
```

### Listar Usuarios por Rol
```bash
curl -H "Authorization: Bearer <admin_token>" \
  http://localhost:8081/api/users/role/admin?page=1&pageSize=10
```

## 🔐 Sistema de Autorización

### Roles y Permisos
```
admin      (Nivel 4) → Puede gestionar todos los usuarios y configuraciones
treasurer  (Nivel 3) → Puede ver reportes financieros y usuarios
operator   (Nivel 2) → Puede gestionar transacciones y operaciones
user       (Nivel 1) → Acceso básico, solo su propia información
```

### Reglas de Negocio
- **Auto-gestión**: Los usuarios pueden actualizar su propia información
- **Protección de admins**: Los admins no pueden eliminarse a sí mismos
- **Jerarquía**: Los roles superiores pueden gestionar roles inferiores
- **Emails únicos**: No se permiten emails duplicados en el sistema

## 🗃️ Base de Datos

### Esquema Principal (V1)
```sql
CREATE TABLE users (
  id VARCHAR(36) PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  role VARCHAR(20) DEFAULT 'user',
  is_active TINYINT(1) DEFAULT 1,
  email_verified TINYINT(1) DEFAULT 0,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
```

### Extensiones de Perfil (V2)
```sql
ALTER TABLE users ADD COLUMN profile_data JSON NULL;
ALTER TABLE users ADD COLUMN last_login_at DATETIME NULL;

-- Índices para performance
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_last_login ON users(last_login_at);
CREATE INDEX idx_users_active_role ON users(is_active, role);
```

## 🧪 Testing

### Ejecutar Tests
```bash
# Tests de servicios (lógica de negocio)
go test ./internal/core/service/... -v

# Tests de handlers (HTTP endpoints)  
go test ./internal/infrastructure/entrypoints/handlers/user/... -v

# Tests completos
go test ./... -v
```

### Cobertura de Tests
- **UserService**: 13 tests - CRUD completo, validaciones, autorización
- **UserHandler**: 5 tests - HTTP handling, DTOs, códigos de estado
- **AuthService**: 5 tests - Registro, login, tokens JWT

### Resultados Esperados
```
✅ AuthService: 5/5 tests PASS
✅ UserService: 13/13 tests PASS  
✅ UserHandler: 5/5 tests PASS
✅ Compilación: SUCCESS
```

## ⚙️ Variables de Entorno

```bash
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password

# JWT Configuration
JWT_SECRET=your-super-secure-secret-key
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h

# Servidor
PORT=8080
LOG_LEVEL=info
```

## 🚀 Ejecución

### Con Docker Compose (Recomendado)
```bash
# Desde la raíz del proyecto
docker-compose up --build mysql user-service

# El servicio estará disponible en http://localhost:8081
```

### Desarrollo Local
```bash
# 1. Configurar variables de entorno
export DB_HOST=localhost
export JWT_SECRET=your-secret-key

# 2. Instalar dependencias
cd backend/services/user-service
go mod tidy

# 3. Aplicar migraciones a la BD
# Ejecutar scripts en database/migrations/

# 4. Ejecutar servicio
go run ./cmd/api/main.go
```

### Health Check
```bash
curl http://localhost:8081/health
# Respuesta: {"status": "ok"}
```

## 📋 Estructura de Respuestas

### Usuario Completo
```json
{
  "id": "uuid-v4",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "fullName": "John Doe",
  "role": "user",
  "isActive": true,
  "emailVerified": false,
  "profile": {
    "phone": "+1234567890",
    "dateOfBirth": "1990-01-01T00:00:00Z",
    "address": {
      "street": "123 Main St",
      "city": "Springfield",
      "state": "IL",
      "postalCode": "62701",
      "country": "USA"
    },
    "profilePicture": "https://...",
    "preferences": {
      "language": "en",
      "timezone": "UTC",
      "notificationEmail": true,
      "notificationSMS": false
    }
  },
  "createdAt": "2025-09-12T10:00:00Z",
  "updatedAt": "2025-09-12T12:00:00Z",
  "lastLoginAt": "2025-09-12T11:30:00Z"
}
```

### Lista Paginada
```json
{
  "users": [...],
  "total": 100,
  "page": 1,
  "pageSize": 20,
  "totalPages": 5
}
```

## 🔧 Postman Collection

Usar la colección completa ubicada en:
```
docs/postman/FinTrack_UserService.postman_collection.json
```

La colección incluye:
- Todos los endpoints implementados
- Variables de entorno configurables
- Tests automáticos de respuesta
- Ejemplos de uso para cada endpoint

## 🛡️ Seguridad

### Autenticación
- **JWT Tokens**: HS256 con secret configurable
- **Access Tokens**: Duración corta (24h por defecto)
- **Refresh Tokens**: Duración extendida (7 días por defecto)

### Autorización  
- **Middleware robusto**: Validación completa de tokens
- **Control granular**: Permisos basados en roles y ownership
- **Validación de estado**: Verificación de usuarios activos

### Validaciones
- **Emails únicos**: Verificación en tiempo real
- **Contraseñas seguras**: Validación de complejidad
- **Sanitización**: Limpieza de inputs maliciosos
- **Business rules**: Prevención de operaciones peligrosas

## 📚 Documentación Adicional

- **Casos de Uso**: `docs/FinTrack_Casos_de_Uso.md`
- **Arquitectura**: `docs/FinTrack_Arquitectura_Tecnica.md`
- **Base de Datos**: `docs/FinTrack_Diseno_Base_Datos.md`
- **Testing**: `docs/FinTrack_Testing_Metodologia.md`

---

**Versión**: 2.0.0 - Microservicio completo de gestión de usuarios  
**Estado**: ✅ Producción Ready  
**Cobertura**: CRUD completo, perfiles, roles, autorización