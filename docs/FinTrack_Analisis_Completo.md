# FinTrack - Análisis Completo del Proyecto
## Plataforma de Gestión de Cuentas, Tarjetas y Billetera Virtual con Chatbot Inteligente

---

## Índice

1. [Introducción](#introducción)
2. [Objetivos del Proyecto](#objetivos-del-proyecto)
3. [Alcance y Limitaciones](#alcance-y-limitaciones)
4. [Análisis de Requerimientos](#análisis-de-requerimientos)
5. [Arquitectura del Sistema](#arquitectura-del-sistema)
6. [Diseño de Base de Datos](#diseño-de-base-de-datos)
7. [Casos de Uso](#casos-de-uso)
8. [Flujos de Trabajo](#flujos-de-trabajo)
9. [Tecnologías y Herramientas](#tecnologías-y-herramientas)
10. [Plan de Desarrollo](#plan-de-desarrollo)
11. [Metodología de Testing](#metodología-de-testing)
12. [Seguridad](#seguridad)
13. [Interfaces de Usuario](#interfaces-de-usuario)
14. [Conclusiones](#conclusiones)

---

## Introducción

### Contexto
En la era digital actual, la gestión financiera personal se ha vuelto cada vez más compleja. Los usuarios manejan múltiples cuentas, tarjetas y servicios financieros dispersos en diferentes plataformas. FinTrack surge como una solución integral que centraliza la gestión financiera digital, proporcionando una experiencia unificada y segura.

### Problemática
- **Fragmentación**: Los usuarios deben acceder a múltiples aplicaciones para gestionar sus finanzas
- **Falta de visibilidad**: Dificultad para obtener una vista consolidada de la situación financiera
- **Complejidad de uso**: Interfaces poco intuitivas y procesos complicados
- **Limitaciones de acceso**: Falta de herramientas accesibles para estudiantes y startups

### Propuesta de Solución
FinTrack es una plataforma web integral que permite:
- Gestión centralizada de cuentas virtuales y tarjetas
- Billetera digital con soporte multi-moneda
- Chatbot inteligente para consultas y transacciones
- Reportes y dashboards personalizados
- Notificaciones en tiempo real
- Arquitectura escalable basada en microservicios

---

## Objetivos del Proyecto

### Objetivo General
Desarrollar una plataforma web integral de gestión financiera digital que centralice el manejo de cuentas, tarjetas y billetera virtual, incorporando un chatbot inteligente y funcionalidades avanzadas de reportes, utilizando tecnologías modernas y buenas prácticas de ingeniería de software.

### Objetivos Específicos

#### Funcionales
1. **Gestión Financiera Centralizada**
   - Implementar módulos para administración de cuentas virtuales
   - Desarrollar sistema de vinculación de tarjetas en modo sandbox
   - Crear billetera digital con funcionalidades de carga, retiro y transferencias

2. **Experiencia de Usuario Mejorada**
   - Diseñar interfaces responsivas y accesibles
   - Implementar chatbot con procesamiento de lenguaje natural
   - Desarrollar sistema de notificaciones en tiempo real

3. **Análisis y Reportes**
   - Crear dashboards interactivos personalizados por rol
   - Implementar exportación de datos en múltiples formatos
   - Desarrollar métricas y KPIs financieros

#### Técnicos
1. **Arquitectura Escalable**
   - Implementar arquitectura de microservicios
   - Utilizar contenedores Docker para deployment
   - Configurar CI/CD pipeline

2. **Seguridad y Confiabilidad**
   - Implementar autenticación y autorización robusta
   - Aplicar cifrado de datos sensibles
   - Desarrollar sistema de auditoría y trazabilidad

3. **Calidad de Software**
   - Alcanzar 80% de cobertura de testing automatizado
   - Implementar testing unitario, integración y E2E
   - Aplicar principios SOLID y patrones de diseño

---

## Alcance y Limitaciones

### Dentro del Alcance
- ✅ Gestión de cuentas virtuales y movimientos
- ✅ Integración con APIs de tarjetas en modo sandbox
- ✅ Billetera digital con transferencias entre usuarios
- ✅ Chatbot con IA para consultas básicas
- ✅ Sistema de roles y permisos
- ✅ Dashboards y reportes exportables
- ✅ Notificaciones en tiempo real
- ✅ Soporte multi-moneda (USD/ARS)
- ✅ APIs RESTful documentadas
- ✅ Testing automatizado

### Fuera del Alcance
- ❌ Integración con bancos reales (solo sandbox)
- ❌ Procesamiento de pagos con dinero real
- ❌ Aplicación móvil nativa
- ❌ Integración con criptomonedas
- ❌ Funcionalidades de inversión
- ❌ Préstamos o créditos

### Limitaciones Técnicas
- Uso exclusivo de APIs gratuitas o en modo sandbox
- Limitaciones de rate limiting de APIs externas
- Dependencia de servicios de terceros para cotizaciones
- Restricciones de almacenamiento en servicios gratuitos

---

## Análisis de Requerimientos

### Requerimientos Funcionales

#### RF001 - Gestión de Usuarios
- **Descripción**: El sistema debe permitir registro, login y gestión de perfiles de usuario
- **Prioridad**: Alta
- **Criterios de Aceptación**:
  - Registro con email y validación
  - Login con autenticación segura
  - Recuperación de contraseña
  - Edición de perfil

#### RF002 - Gestión de Cuentas Virtuales
- **Descripción**: Los usuarios pueden crear y administrar cuentas virtuales
- **Prioridad**: Alta
- **Criterios de Aceptación**:
  - Crear cuentas en diferentes monedas
  - Visualizar saldo y movimientos
  - Historial de transacciones
  - Estados de cuenta

#### RF003 - Gestión de Tarjetas
- **Descripción**: Carga manual y administración centralizada de tarjetas de débito/crédito
- **Prioridad**: Alta
- **Criterios de Aceptación**:
  - Agregar tarjetas manualmente con datos básicos
  - Visualizar información de tarjetas de forma centralizada
  - Gestionar múltiples tarjetas de diferentes bancos
  - Desactivar/eliminar tarjetas del sistema

#### RF004 - Billeteras Virtuales Centralizadas
- **Descripción**: Gestión centralizada de múltiples billeteras virtuales con carga manual
- **Prioridad**: Alta
- **Criterios de Aceptación**:
  - Agregar manualmente billeteras existentes (MercadoPago, Ualá, etc.)
  - Vista unificada de todas las billeteras
  - Sincronización manual de saldos
  - Transferencias P2P entre billeteras
  - Historial consolidado de operaciones

#### RF005 - Chatbot Inteligente
- **Descripción**: Asistente virtual para consultas y transacciones
- **Prioridad**: Media
- **Criterios de Aceptación**:
  - Consultas de saldo en lenguaje natural
  - FAQ financieras
  - Iniciación de transacciones con confirmación
  - Integración con WhatsApp/Telegram

#### RF006 - Notificaciones
- **Descripción**: Sistema de alertas en tiempo real
- **Prioridad**: Media
- **Criterios de Aceptación**:
  - Notificaciones de transacciones
  - Alertas de saldo bajo
  - Notificaciones por email/SMS
  - Configuración de preferencias

#### RF007 - Reportes y Dashboards
- **Descripción**: Generación de reportes y visualizaciones
- **Prioridad**: Media
- **Criterios de Aceptación**:
  - Dashboards personalizados por rol
  - Exportación en Excel/PDF
  - Gráficos interactivos
  - Filtros por fecha y categoría

#### RF008 - Gestión Multi-moneda
- **Descripción**: Soporte para múltiples divisas
- **Prioridad**: Baja
- **Criterios de Aceptación**:
  - Cotizaciones en tiempo real
  - Conversión automática
  - Historial de tipos de cambio
  - Cuentas en USD y ARS

### Requerimientos No Funcionales

#### RNF001 - Performance
- Tiempo de respuesta < 2 segundos para operaciones básicas
- Soporte para 1000 usuarios concurrentes
- Disponibilidad del 99.5%

#### RNF002 - Seguridad
- Cifrado AES-256 para datos sensibles
- Autenticación JWT con refresh tokens
- Logs de auditoría para todas las transacciones
- Cumplimiento con estándares PCI DSS básicos

#### RNF003 - Usabilidad
- Interfaz responsive para dispositivos móviles
- Tiempo de aprendizaje < 30 minutos
- Accesibilidad WCAG 2.1 nivel AA

#### RNF004 - Escalabilidad
- Arquitectura de microservicios
- Balanceador de carga
- Base de datos distribuida
- Cache distribuido con Redis

#### RNF005 - Mantenibilidad
- Cobertura de testing > 80%
- Documentación técnica completa
- Código siguiendo estándares de calidad
- CI/CD automatizado

---

## Arquitectura del Sistema

### Arquitectura General

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Microservicios│
│   (React)       │◄──►│   (Express)     │◄──►│   (Node.js)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Load Balancer │    │   Base de Datos │
                       │   (Nginx)       │    │   (PostgreSQL)  │
                       └─────────────────┘    └─────────────────┘
```

### Microservicios Propuestos

1. **User Service**
   - Gestión de usuarios y autenticación
   - JWT token management
   - Perfiles y preferencias

2. **Account Service**
   - Gestión de cuentas virtuales
   - Movimientos y transacciones
   - Estados de cuenta

3. **Card Service**
   - Integración con APIs de tarjetas
   - Gestión de tarjetas vinculadas
   - Transacciones con tarjetas

4. **Wallet Service**
   - Billetera digital
   - Transferencias entre usuarios
   - Gestión de saldos

5. **Notification Service**
   - Notificaciones en tiempo real
   - Email y SMS
   - WebSocket connections

6. **Report Service**
   - Generación de reportes
   - Dashboards y métricas
   - Exportación de datos

7. **Chatbot Service**
   - Procesamiento de lenguaje natural
   - Integración con IA
   - Gestión de conversaciones

8. **Currency Service**
   - Cotizaciones en tiempo real
   - Conversión de monedas
   - Historial de tipos de cambio

### Stack Tecnológico

#### Frontend
- **Framework**: React 18 con TypeScript
- **UI Library**: Material-UI o Ant Design
- **State Management**: Redux Toolkit
- **Routing**: React Router v6
- **Charts**: Chart.js o D3.js
- **Testing**: Jest + React Testing Library

#### Backend
- **Runtime**: Node.js 22+
- **Framework**: Express.js
- **Language**: TypeScript
- **API Documentation**: Swagger/OpenAPI
- **Validation**: Joi o Zod
- **Testing**: Jest + Supertest

#### Base de Datos
- **Principal**: PostgreSQL 14+
- **Cache**: Redis 6+
- **ORM**: Prisma o TypeORM
- **Migrations**: Prisma Migrate

#### DevOps
- **Containerización**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Monitoring**: Prometheus + Grafana
- **Logs**: Winston + ELK Stack

#### APIs Externas
- **Tarjetas**: Stripe API (modo test)
- **Cotizaciones**: ExchangeRate-API
- **Chatbot**: OpenAI GPT-3.5 o Dialogflow
- **Notificaciones**: SendGrid, Twilio

---

## Diseño de Base de Datos

### Modelo Entidad-Relación

```sql
-- Usuarios
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    role user_role NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cuentas virtuales
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    account_type account_type NOT NULL,
    currency currency_code NOT NULL DEFAULT 'ARS',
    balance DECIMAL(15,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tarjetas
CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    card_number_encrypted VARCHAR(255) NOT NULL,
    card_type card_type NOT NULL,
    brand VARCHAR(20) NOT NULL,
    last_four VARCHAR(4) NOT NULL,
    expiry_month INTEGER NOT NULL,
    expiry_year INTEGER NOT NULL,
    cardholder_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    stripe_card_id VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transacciones
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_account_id UUID REFERENCES accounts(id),
    to_account_id UUID REFERENCES accounts(id),
    card_id UUID REFERENCES cards(id),
    transaction_type transaction_type NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency currency_code NOT NULL,
    description TEXT,
    status transaction_status DEFAULT 'pending',
    reference_number VARCHAR(50) UNIQUE,
    external_reference VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Billetera digital
CREATE TABLE wallet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    balance_ars DECIMAL(15,2) DEFAULT 0.00,
    balance_usd DECIMAL(15,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Notificaciones
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type notification_type NOT NULL,
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Conversaciones del chatbot
CREATE TABLE chat_conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,
    response TEXT,
    intent VARCHAR(100),
    confidence DECIMAL(3,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tipos enumerados
CREATE TYPE user_role AS ENUM ('user', 'operator', 'admin', 'treasurer');
CREATE TYPE account_type AS ENUM ('checking', 'savings', 'business');
CREATE TYPE currency_code AS ENUM ('ARS', 'USD');
CREATE TYPE card_type AS ENUM ('debit', 'credit');
CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal', 'transfer', 'payment', 'refund');
CREATE TYPE transaction_status AS ENUM ('pending', 'completed', 'failed', 'cancelled');
CREATE TYPE notification_type AS ENUM ('transaction', 'balance_low', 'security', 'system');
```

### Índices y Optimizaciones

```sql
-- Índices para mejorar performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_transactions_from_account ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account ON transactions(to_account_id);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_chat_conversations_user_id ON chat_conversations(user_id);
```

---

## Casos de Uso

### CU001 - Registro de Usuario
**Actor**: Usuario no registrado
**Precondiciones**: El usuario accede a la página de registro
**Flujo Principal**:
1. El usuario ingresa email, contraseña, nombre y apellido
2. El sistema valida que el email no esté registrado
3. El sistema envía email de verificación
4. El usuario confirma su email
5. El sistema crea la cuenta y billetera virtual
6. El usuario es redirigido al dashboard

**Flujos Alternativos**:
- 2a. Email ya registrado: mostrar error
- 4a. Email no verificado en 24hs: eliminar cuenta temporal

### CU002 - Vinculación de Tarjeta
**Actor**: Usuario registrado
**Precondiciones**: Usuario autenticado
**Flujo Principal**:
1. Usuario accede a "Mis Tarjetas"
2. Hace clic en "Agregar Tarjeta"
3. Ingresa datos de la tarjeta
4. Sistema valida con API de Stripe (modo test)
5. Sistema guarda tarjeta cifrada
6. Muestra confirmación al usuario

**Flujos Alternativos**:
- 4a. Tarjeta inválida: mostrar error específico
- 4b. Error de API: mostrar mensaje de reintento

### CU003 - Transferencia entre Usuarios
**Actor**: Usuario con billetera activa
**Precondiciones**: Usuario autenticado con saldo suficiente
**Flujo Principal**:
1. Usuario accede a "Transferir"
2. Ingresa email del destinatario y monto
3. Sistema valida destinatario y saldo
4. Usuario confirma la operación
5. Sistema debita del remitente y acredita al destinatario
6. Sistema envía notificaciones a ambos usuarios
7. Sistema registra la transacción

**Flujos Alternativos**:
- 3a. Destinatario no existe: mostrar error
- 3b. Saldo insuficiente: mostrar error
- 5a. Error en transacción: rollback y notificar

### CU004 - Consulta con Chatbot
**Actor**: Usuario registrado
**Precondiciones**: Usuario autenticado
**Flujo Principal**:
1. Usuario abre el chat
2. Escribe consulta en lenguaje natural
3. Sistema procesa con NLP
4. Sistema identifica intención
5. Sistema genera respuesta apropiada
6. Sistema muestra respuesta al usuario

**Flujos Alternativos**:
- 4a. Intención no reconocida: ofrecer opciones predefinidas
- 4b. Consulta de transacción: solicitar confirmación adicional

---

## Flujos de Trabajo

### Flujo de Onboarding
```
Registro → Verificación Email → Configuración Perfil → 
Creación Billetera → Tutorial Interactivo → Dashboard Principal
```

### Flujo de Transacción
```
Selección Tipo → Ingreso Datos → Validación → 
Confirmación → Procesamiento → Notificación → Registro
```

### Flujo de Reportes
```
Selección Filtros → Consulta Datos → Generación Gráficos → 
Exportación (Opcional) → Visualización Dashboard
```

---

## Plan de Desarrollo

### Metodología
**Scrum** con sprints de 2 semanas

### Fases del Proyecto

#### Fase 1: Fundación (4 semanas)
- Setup del proyecto y repositorios
- Configuración de CI/CD
- Diseño de base de datos
- Arquitectura de microservicios
- Autenticación y autorización

#### Fase 2: Core Features (6 semanas)
- Gestión de usuarios
- Cuentas virtuales
- Transacciones básicas
- API de tarjetas (sandbox)
- Billetera digital

#### Fase 3: Features Avanzadas (4 semanas)
- Chatbot básico
- Notificaciones
- Dashboards
- Reportes y exportación

#### Fase 4: Optimización (3 semanas)
- Performance tuning
- Testing completo
- Documentación
- Deployment

#### Fase 5: Entrega (1 semana)
- Testing final
- Documentación de usuario
- Presentación

### Cronograma Detallado

| Sprint | Semanas | Objetivos Principales |
|--------|---------|----------------------|
| 1 | 1-2 | Setup, DB, Auth |
| 2 | 3-4 | Users, Accounts |
| 3 | 5-6 | Transactions, Cards |
| 4 | 7-8 | Wallet, Transfers |
| 5 | 9-10 | Notifications, Basic UI |
| 6 | 11-12 | Chatbot, Reports |
| 7 | 13-14 | Dashboards, Export |
| 8 | 15-16 | Testing, Performance |
| 9 | 17-18 | Documentation, Deploy |

---

## Metodología de Testing

### Estrategia de Testing

#### Pirámide de Testing
```
        E2E Tests (10%)
      ─────────────────
    Integration Tests (20%)
   ─────────────────────────
     Unit Tests (70%)
  ───────────────────────────
```

#### Tipos de Testing

1. **Unit Testing**
   - Cobertura objetivo: 80%
   - Framework: Jest
   - Mocks para dependencias externas
   - Testing de funciones puras

2. **Integration Testing**
   - APIs endpoints
   - Base de datos
   - Servicios externos (mocked)
   - Framework: Supertest + Jest

3. **E2E Testing**
   - Flujos críticos de usuario
   - Framework: Cypress
   - Ambientes de testing dedicados

4. **Performance Testing**
   - Load testing con Artillery
   - Stress testing
   - Memory leak detection

5. **Security Testing**
   - Penetration testing básico
   - OWASP Top 10
   - Dependency vulnerability scanning

### Plan de Testing

#### Casos de Prueba Críticos
1. Registro y autenticación de usuarios
2. Transacciones financieras
3. Integración con APIs externas
4. Seguridad de datos sensibles
5. Performance bajo carga

#### Ambientes de Testing
- **Development**: Testing local
- **Staging**: Testing de integración
- **Production**: Monitoring y alertas

---

## Seguridad

### Medidas de Seguridad Implementadas

#### Autenticación y Autorización
- JWT tokens con refresh mechanism
- Rate limiting por IP y usuario
- Roles y permisos granulares
- 2FA opcional

#### Protección de Datos
- Cifrado AES-256 para datos sensibles
- Hashing bcrypt para contraseñas
- HTTPS obligatorio
- Sanitización de inputs

#### Auditoría y Monitoreo
- Logs de todas las transacciones
- Detección de actividad sospechosa
- Alertas de seguridad
- Backup automático de datos

#### Cumplimiento
- GDPR compliance básico
- PCI DSS guidelines
- Políticas de retención de datos

---

## Conclusiones

FinTrack representa una solución integral para la gestión financiera digital, diseñada específicamente para entornos académicos y de desarrollo. El proyecto combina tecnologías modernas con buenas prácticas de ingeniería de software, proporcionando una base sólida para el aprendizaje y la experimentación.

### Beneficios Esperados
- **Educativos**: Aplicación práctica de conceptos de ingeniería de software
- **Técnicos**: Experiencia con arquitecturas modernas y microservicios
- **Profesionales**: Portfolio project para demostrar habilidades

### Próximos Pasos
1. Validación del diseño con stakeholders
2. Setup del ambiente de desarrollo
3. Inicio del desarrollo siguiendo la metodología Scrum
4. Implementación incremental de funcionalidades

---

*Documento creado para la Tecnicatura en Programación - UNT*  
*Versión 1.0 - Enero 2025*