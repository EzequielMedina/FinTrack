# FinTrack - Cronograma de Sprints y Backlog

## 📅 Cronograma General

**Fecha de Inicio:** 10 de Septiembre 2024  
**Fecha de Finalización:** 3 de Diciembre 2024  
**Duración Total:** 12 semanas y 2 días (84 días)  
**Duración por Sprint:** 15 días  
**Total de Sprints:** 5 sprints completos + 1 sprint final de 9 días

---

## 🗓️ División de Sprints

### Sprint 1: Fundación y Configuración
**📅 Duración:** 10 Sep - 24 Sep (15 días)  
**🎯 Objetivo:** Establecer la base técnica del proyecto

#### Tareas del Sprint 1

**Configuración del Entorno (5 días)**
- [ ] **TASK-001** - Configurar repositorio Git con estructura de microservicios
  - Estimación: 1 día
  - Prioridad: Alta
  - Criterios de aceptación: Repo con estructura frontend/backend/docs

- [ ] **TASK-002** - Configurar Docker y Docker Compose para desarrollo
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: Containers funcionando para MySQL, Redis, Go, Angular

- [ ] **TASK-003** - Configurar CI/CD pipeline básico con GitHub Actions
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: Pipeline ejecutando tests y build automático

**Backend - Autenticación (5 días)**
- [ ] **TASK-004** - Implementar microservicio de autenticación en Go
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: JWT, registro, login, middleware de auth

- [ ] **TASK-005** - Configurar base de datos MySQL con migraciones
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: Esquema inicial, migraciones automáticas con GORM

**Frontend - Base (5 días)**
- [ ] **TASK-006** - Configurar proyecto Angular 20 con arquitectura base
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: Proyecto con routing, guards, interceptors

- [ ] **TASK-007** - Implementar componentes de autenticación (login/registro)
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: Formularios reactivos, validaciones, integración con API

---

### Sprint 2: Gestión de Usuarios y Dashboard Básico
**📅 Duración:** 25 Sep - 9 Oct (15 días)  
**🎯 Objetivo:** Sistema de usuarios y dashboard inicial

#### Tareas del Sprint 2

**Backend - Gestión de Usuarios (6 días)**
- [ ] **TASK-008** - Implementar microservicio de gestión de usuarios
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: CRUD usuarios, perfiles, roles

- [ ] **TASK-009** - Implementar sistema de roles y permisos
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: Roles (admin, user), middleware de autorización

**Frontend - Dashboard (6 días)**
- [ ] **TASK-010** - Crear layout principal con navegación
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: Sidebar, header, routing funcional

- [ ] **TASK-011** - Implementar dashboard básico con widgets
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: Resumen de cuentas, gráficos básicos, responsive

**Testing y Documentación (3 días)**
- [ ] **TASK-012** - Implementar tests unitarios para autenticación
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: Cobertura >80% en auth service

- [ ] **TASK-013** - Documentar APIs con Swagger/OpenAPI
  - Estimación: 1 día
  - Prioridad: Media
  - Criterios de aceptación: Documentación interactiva disponible

---

### Sprint 3: Gestión de Cuentas y Tarjetas
**📅 Duración:** 10 Oct - 24 Oct (15 días)  
**🎯 Objetivo:** Funcionalidades core de gestión financiera

#### Tareas del Sprint 3

**Backend - Cuentas y Tarjetas (8 días)**
- [ ] **TASK-014** - Implementar microservicio de cuentas virtuales
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: CRUD cuentas, tipos de cuenta, saldos

- [ ] **TASK-015** - Implementar gestión de tarjetas (Front)
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: Vinculación tarjetas, validaciones, encriptación

**Frontend - Gestión Financiera (5 días)**
- [ ] **TASK-016** - Crear módulo de gestión de cuentas
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: Lista, crear, editar, eliminar cuentas

- [ ] **TASK-017** - Crear módulo de gestión de tarjetas
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: Formulario seguro, lista enmascarada

**Integración APIs (2 días)**
- [ ] **TASK-018** - Integrar API de conversión de divisas (ExchangeRates)
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: Conversión USD/ARS automática, cache

---

### Sprint 4: Billetera Digital y Transacciones
**📅 Duración:** 25 Oct - 8 Nov (15 días)  
**🎯 Objetivo:** Sistema de billetera y transacciones

#### Tareas del Sprint 4

**Backend - Billetera y Transacciones (8 días)**
- [ ] **TASK-019** - Implementar microservicio de billetera digital
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: Carga, retiro, transferencias entre usuarios

- [ ] **TASK-020** - Implementar sistema de transacciones
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: Historial, categorización, validaciones

**Frontend - Billetera (5 días)**
- [ ] **TASK-021** - Crear módulo de billetera digital
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: Vista de saldo, operaciones, historial

- [ ] **TASK-022** - Implementar formularios de transacciones
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: Transferencias, validaciones en tiempo real

**Notificaciones (2 días)**
- [ ] **TASK-023** - Integrar Firebase Cloud Messaging
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: Notificaciones de transacciones en tiempo real

---

### Sprint 5: Reportes y Analytics
**📅 Duración:** 9 Nov - 23 Nov (15 días)  
**🎯 Objetivo:** Sistema de reportes y análisis

#### Tareas del Sprint 5

**Backend - Reportes (6 días)**
- [ ] **TASK-024** - Implementar microservicio de reportes
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: Generación PDF/Excel, filtros, agregaciones

- [ ] **TASK-025** - Implementar analytics y métricas
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: KPIs financieros, patrones de gasto

**Frontend - Reportes (6 días)**
- [ ] **TASK-026** - Crear módulo de reportes con gráficos
  - Estimación: 4 días
  - Prioridad: Alta
  - Criterios de aceptación: Charts interactivos, filtros, exportación

- [ ] **TASK-027** - Implementar dashboard de analytics
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: Métricas en tiempo real, comparativas

**Integración APIs Analytics (3 días)**
- [ ] **TASK-028** - Integrar Google Analytics 4
  - Estimación: 1 día
  - Prioridad: Baja
  - Criterios de aceptación: Tracking de eventos, métricas de uso

- [ ] **TASK-029** - Integrar Alpha Vantage para datos de mercado
  - Estimación: 2 días
  - Prioridad: Media
  - Criterios de aceptación: Cotizaciones en dashboard, cache por rate limit

---

### Sprint 6: Chatbot y Finalización
**📅 Duración:** 24 Nov - 3 Dic (9 días)  
**🎯 Objetivo:** Chatbot inteligente y cierre del proyecto

#### Tareas del Sprint 6

**Chatbot (6 días)**
- [ ] **TASK-030** - Implementar integración con OpenAI GPT-3.5
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: API integration, context management

- [ ] **TASK-031** - Crear interfaz de chatbot en frontend
  - Estimación: 3 días
  - Prioridad: Alta
  - Criterios de aceptación: Chat UI, historial, respuestas en tiempo real

**Finalización (3 días)**
- [ ] **TASK-032** - Testing integral y corrección de bugs
  - Estimación: 2 días
  - Prioridad: Alta
  - Criterios de aceptación: E2E tests, cobertura >80%

- [ ] **TASK-033** - Documentación final y deployment
  - Estimación: 1 día
  - Prioridad: Alta
  - Criterios de aceptación: README, guías de instalación, demo funcional

---

## 📊 Resumen de Estimaciones

| Sprint | Duración | Tareas | Días Estimados | Complejidad |
|--------|----------|--------|----------------|-------------|
| Sprint 1 | 15 días | 7 tareas | 15 días | Alta |
| Sprint 2 | 15 días | 6 tareas | 15 días | Media |
| Sprint 3 | 15 días | 5 tareas | 15 días | Alta |
| Sprint 4 | 15 días | 5 tareas | 15 días | Alta |
| Sprint 5 | 15 días | 6 tareas | 15 días | Media |
| Sprint 6 | 9 días | 4 tareas | 9 días | Media |
| **Total** | **84 días** | **33 tareas** | **84 días** | - |

---

## 🎯 Criterios de Éxito por Sprint

### Sprint 1 ✅
- [ ] Entorno de desarrollo completamente funcional
- [ ] Autenticación JWT implementada y probada
- [ ] Frontend base con login/registro operativo

### Sprint 2 ✅
- [ ] Sistema de usuarios con roles funcionando
- [ ] Dashboard básico responsive
- [ ] Tests unitarios con cobertura >80%

### Sprint 3 ✅
- [ ] Gestión completa de cuentas y tarjetas
- [ ] Integración con API de divisas
- [ ] Validaciones de seguridad implementadas

### Sprint 4 ✅
- [ ] Billetera digital completamente funcional
- [ ] Sistema de transacciones con historial
- [ ] Notificaciones en tiempo real

### Sprint 5 ✅
- [ ] Reportes exportables en PDF/Excel
- [ ] Dashboard de analytics con gráficos
- [ ] Integración con APIs de datos financieros

### Sprint 6 ✅
- [ ] Chatbot inteligente operativo
- [ ] Testing integral completado
- [ ] Documentación y deployment finalizados

---

## 🔄 Metodología de Trabajo

### Daily Standups
- **Frecuencia:** Diaria (15 min)
- **Horario:** 9:00 AM
- **Formato:** ¿Qué hice ayer? ¿Qué haré hoy? ¿Hay impedimentos?

### Sprint Planning
- **Duración:** 2 horas al inicio de cada sprint
- **Participantes:** Todo el equipo
- **Entregables:** Sprint backlog refinado y estimado

### Sprint Review
- **Duración:** 1 hora al final de cada sprint
- **Formato:** Demo de funcionalidades completadas
- **Stakeholders:** Product Owner, usuarios finales

### Sprint Retrospective
- **Duración:** 1 hora después del review
- **Formato:** ¿Qué funcionó bien? ¿Qué mejorar? ¿Acciones?

---

## 📈 Métricas de Seguimiento

### Velocity
- **Sprint 1:** Baseline (primera medición)
- **Objetivo:** Mantener velocity consistente
- **Métrica:** Story points completados por sprint

### Burndown
- **Tracking diario:** Tareas restantes vs tiempo
- **Alertas:** Si el burndown se desvía >20%

### Quality Metrics
- **Code Coverage:** >80% en backend, >70% en frontend
- **Bug Rate:** <5 bugs por sprint
- **Technical Debt:** <10% del tiempo de desarrollo

---

## 🚨 Riesgos y Mitigaciones

### Riesgos Técnicos
1. **Rate Limits de APIs externas**
   - Mitigación: Implementar cache y fallbacks
   - Contingencia: APIs alternativas identificadas

2. **Complejidad de microservicios**
   - Mitigación: Documentación detallada, tests de integración
   - Contingencia: Simplificar a monolito si es necesario

### Riesgos de Cronograma
1. **Subestimación de tareas**
   - Mitigación: Buffer del 20% en estimaciones
   - Contingencia: Priorizar features core

2. **Dependencias externas**
   - Mitigación: Identificar dependencias críticas temprano
   - Contingencia: Mocks y simuladores

---

*Documento generado para el proyecto FinTrack - Tecnicatura en Programación UNT*  
*Versión 1.0 - Cronograma actualizado para período 10/09 - 3/12/2024*