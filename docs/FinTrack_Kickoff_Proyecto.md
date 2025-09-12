# FinTrack - Kickoff del Proyecto para Equipos de Trabajo

## 📋 Información General del Proyecto

### Nombre del Proyecto
**FinTrack** - Sistema de Gestión Financiera Personal

### Descripción
FinTrack es una aplicación web moderna para la gestión integral de finanzas personales que permite a los usuarios administrar sus tarjetas de crédito/débito, billeteras virtuales centralizadas, realizar seguimiento de gastos y generar reportes analíticos detallados.

### Objetivos del Proyecto
- Centralizar la gestión de múltiples billeteras virtuales en una sola plataforma
- Proporcionar carga manual de tarjetas y billeteras para mayor control
- Ofrecer análisis detallado de patrones de gasto
- Generar reportes financieros personalizados
- Implementar un chatbot inteligente para consultas financieras

## 🛠️ Stack Tecnológico

### Frontend
- **Framework**: Angular 20
- **Lenguaje**: TypeScript
- **Gestión de Estado**: NgRx
- **UI Library**: Angular Material
- **Testing**: Jasmine + Karma

### Backend
- **Lenguaje**: Go (Golang)
- **Arquitectura**: Microservicios
- **Framework**: Gin/Echo
- **Testing**: Go testing package

### Base de Datos
- **SGBD**: MySQL 8.0+
- **ORM**: GORM (Go)
- **Migraciones**: Automáticas con GORM

### Infraestructura
- **Contenedores**: Docker
- **Orquestación**: Docker Compose
- **Proxy Reverso**: Nginx
- **Monitoreo**: Prometheus + Grafana

## 🏗️ Arquitectura del Sistema

### Microservicios
1. **User Service** (Go)
   - Gestión de usuarios y autenticación
   - Puerto: 3001

2. **Account Service** (Go)
   - Gestión de cuentas y billeteras virtuales
   - Puerto: 3002

3. **Transaction Service** (Go)
   - Procesamiento de transacciones
   - Puerto: 3003

4. **Analytics Service** (Go)
   - Generación de reportes y análisis
   - Puerto: 3004

5. **Notification Service** (Go)
   - Gestión de notificaciones
   - Puerto: 3005

### Base de Datos
- **Esquema Principal**: fintrack_db
- **Tablas Principales**:
  - users
  - cards
  - virtual_wallets
  - wallet_providers
  - transactions
  - categories
  - budgets

## 🎯 Funcionalidades Principales

### RF001 - Gestión de Usuarios
- Registro y autenticación de usuarios
- Gestión de perfiles
- Configuración de preferencias

### RF002 - Dashboard Principal
- Vista consolidada de saldos
- Resumen de transacciones recientes
- Gráficos de gastos por categoría

### RF003 - Gestión de Tarjetas
- **Carga manual** de tarjetas de crédito/débito
- Visualización de saldos y límites
- Historial de transacciones por tarjeta

### RF004 - Billeteras Virtuales Centralizadas
- **Vista unificada** de todas las billeteras virtuales
- **Carga manual** de billeteras (MercadoPago, Ualá, Brubank, etc.)
- Saldo total consolidado
- Sincronización manual de saldos

### RF005 - Seguimiento de Gastos
- Categorización automática y manual
- Filtros por fecha, categoría, monto
- Búsqueda avanzada de transacciones

### RF006 - Reportes y Analytics
- Reportes mensuales/anuales
- Análisis de patrones de gasto
- Comparativas temporales
- Exportación a PDF/Excel

### RF007 - Chatbot Inteligente
- Consultas sobre saldos y gastos
- Recomendaciones financieras
- Alertas personalizadas

## 👥 Estructura del Equipo

### Roles y Responsabilidades

#### Frontend Team
- **Desarrollador Angular Senior**: Arquitectura y componentes principales
- **Desarrollador Angular Junior**: Implementación de vistas y formularios
- **UI/UX Designer**: Diseño de interfaces y experiencia de usuario

#### Backend Team
- **Desarrollador Go Senior**: Arquitectura de microservicios
- **Desarrollador Go Junior**: Implementación de APIs y servicios
- **DevOps Engineer**: Infraestructura y despliegue

#### QA Team
- **QA Lead**: Estrategia de testing y automatización
- **QA Tester**: Testing manual y casos de prueba

#### Product Team
- **Product Owner**: Definición de requisitos y prioridades
- **Scrum Master**: Facilitación y metodología ágil

## 📅 Cronograma del Proyecto

### Sprint 1 (2 semanas) - Fundación
- Configuración del entorno de desarrollo
- Implementación de autenticación
- Dashboard básico
- Estructura de base de datos

### Sprint 2 (2 semanas) - Gestión de Tarjetas
- Carga manual de tarjetas
- CRUD de tarjetas
- Visualización de saldos

### Sprint 3 (2 semanas) - Billeteras Virtuales
- Implementación de billeteras centralizadas
- Carga manual de billeteras
- Vista consolidada de saldos

### Sprint 4 (2 semanas) - Transacciones
- Registro de transacciones
- Categorización
- Filtros y búsqueda

### Sprint 5 (2 semanas) - Analytics
- Reportes básicos
- Gráficos y visualizaciones
- Exportación de datos

### Sprint 6 (2 semanas) - Chatbot
- Implementación del chatbot
- Integración con servicios
- Testing y refinamiento

## 🔧 Configuración del Entorno

### Requisitos Previos
- Node.js 22+
- Go 1.24+
- MySQL 8.0+
- Docker y Docker Compose
- Git

### Estructura de Repositorios
```
fintrack/
├── frontend/          # Aplicación Angular
├── backend/
│   ├── user-service/   # Microservicio de usuarios
│   ├── account-service/ # Microservicio de cuentas
│   ├── transaction-service/ # Microservicio de transacciones
│   ├── analytics-service/ # Microservicio de analytics
│   └── notification-service/ # Microservicio de notificaciones
├── database/          # Scripts de base de datos
├── docker/           # Configuración Docker
└── docs/             # Documentación del proyecto
```

### Variables de Entorno
```env
# Base de datos
DB_HOST=localhost
DB_PORT=3306
DB_NAME=fintrack_db
DB_USER=fintrack_user
DB_PASSWORD=secure_password

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=24h

# Servicios
USER_SERVICE_PORT=3001
ACCOUNT_SERVICE_PORT=3002
TRANSACTION_SERVICE_PORT=3003
ANALYTICS_SERVICE_PORT=3004
NOTIFICATION_SERVICE_PORT=3005
```

## 🧪 Estrategia de Testing

### Frontend Testing
- **Unit Tests**: Jasmine + Karma
- **Integration Tests**: Angular Testing Utilities
- **E2E Tests**: Cypress
- **Coverage**: Mínimo 80%

### Backend Testing
- **Unit Tests**: Go testing package
- **Integration Tests**: Testify
- **API Tests**: Postman/Newman
- **Coverage**: Mínimo 85%

### Base de Datos
- **Migration Tests**: Verificación de esquemas
- **Data Integrity Tests**: Validación de constraints
- **Performance Tests**: Optimización de queries

## 📊 Métricas y Monitoreo

### KPIs del Proyecto
- Tiempo de respuesta de APIs < 200ms
- Disponibilidad del sistema > 99.5%
- Cobertura de tests > 80%
- Tiempo de carga de frontend < 3s

### Herramientas de Monitoreo
- **Prometheus**: Métricas de aplicación
- **Grafana**: Dashboards y visualizaciones
- **ELK Stack**: Logs centralizados
- **Sentry**: Tracking de errores

## 🚀 Proceso de Despliegue

### Ambientes
1. **Development**: Desarrollo local
2. **Staging**: Testing y validación
3. **Production**: Ambiente productivo

### Pipeline CI/CD
1. **Commit**: Push a repositorio
2. **Build**: Compilación y tests
3. **Test**: Ejecución de test suite
4. **Deploy**: Despliegue automático
5. **Monitor**: Verificación post-deploy

## 📋 Definición de Terminado (DoD)

### Para cada Feature
- [ ] Código implementado y revisado
- [ ] Tests unitarios escritos y pasando
- [ ] Tests de integración pasando
- [ ] Documentación actualizada
- [ ] Code review aprobado
- [ ] Deploy en staging exitoso
- [ ] Validación del Product Owner

## 🔐 Consideraciones de Seguridad

### Autenticación y Autorización
- JWT tokens con expiración
- Refresh tokens para sesiones largas
- Rate limiting en APIs
- Validación de entrada en todos los endpoints

### Protección de Datos
- Encriptación de datos sensibles
- HTTPS obligatorio
- Sanitización de inputs
- Logs sin información sensible

## 📞 Contactos del Proyecto

### Stakeholders
- **Product Owner**: [Nombre] - [email]
- **Scrum Master**: [Nombre] - [email]
- **Tech Lead**: [Nombre] - [email]
- **DevOps Lead**: [Nombre] - [email]

### Canales de Comunicación
- **Slack**: #fintrack-project
- **Email**: fintrack-team@company.com
- **Jira**: [URL del proyecto]
- **Confluence**: [URL de documentación]

---

**Fecha de Kickoff**: [Fecha]
**Versión del Documento**: 1.0
**Última Actualización**: [Fecha actual]

*Este documento será actualizado conforme evolucione el proyecto.*