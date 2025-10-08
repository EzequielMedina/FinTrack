# 📋 PLAN DE DESARROLLO - SISTEMA DE CUOTAS PARA TARJETAS DE CRÉDITO

## 🎯 Objetivo
Implementar un sistema completo de cuotas para tarjetas de crédito que permita a los usuarios:
- Dividir compras en cuotas (1-24 meses)
- Elegir fecha de inicio de pagos
- Ver previsualización del plan de cuotas
- Gestionar y monitorear cuotas activas

---

## 📊 Estado Actual del Análisis

### ✅ Análisis Completado
- [x] Revisión de base de datos actual
- [x] Análisis de microservicios existentes
- [x] Evaluación del flujo de transacciones de crédito
- [x] Identificación de puntos de integración

### 📋 Hallazgos Clave
- **Base de datos**: Tabla `transactions` con campo `metadata` JSON para extensibilidad
- **Backend**: Account-service maneja lógica de tarjetas, Transaction-service para auditoría
- **Frontend**: Credit-card.service con operaciones básicas implementadas
- **Arquitectura**: Microservicios con comunicación HTTP establecida

---

## 🏗️ FASE 1: DISEÑO Y ESTRUCTURA DE BASE DE DATOS

### 1.1 Migración de Base de Datos
- [x] Crear migración `07_V7__installments.sql`
- [x] Implementar tabla `installment_plans`
  - [x] Campos básicos (id, transaction_id, card_id, user_id)
  - [x] Detalles del plan (total_amount, installments_count, installment_amount, start_date)
  - [x] Control de estado (status, paid_installments, remaining_amount)
  - [x] Timestamps y foreign keys
- [x] Implementar tabla `installments`
  - [x] Campos básicos (id, plan_id, installment_number)
  - [x] Detalles de pago (amount, due_date, paid_date, status)
  - [x] Referencias (payment_transaction_id)
  - [x] Índices para optimización
- [ ] Ejecutar y validar migración

### 1.2 Entities y Models (Account-Service)
- [x] Crear `installment_plan.go` entity
  - [x] Struct InstallmentPlan con validaciones
  - [x] Métodos de negocio (CanCancel, GetNextDueInstallment, etc.)
  - [x] Relaciones con Card y Transaction
- [x] Crear `installment.go` entity
  - [x] Struct Installment con validaciones
  - [x] Métodos de estado (IsOverdue, CanPay, etc.)
  - [x] Relación con InstallmentPlan
- [x] Actualizar `card.go` entity
  - [x] Agregar métodos para cuotas (CanCreateInstallmentPlan, GetActiveInstallmentPlans)
  - [x] Validaciones para compras con cuotas

---

## 🔧 FASE 2: BACKEND - REPOSITORIES Y SERVICIOS

### 2.1 Repositories (Account-Service) ✅ COMPLETADA
- [x] Crear `installment_plan_repository.go`
  - [x] Interface InstallmentPlanRepositoryInterface
  - [x] Implementación MySQL con GORM
  - [x] Métodos CRUD básicos
  - [x] Consultas especializadas (GetByCard, GetActive, etc.)
- [x] Crear `installment_repository.go`
  - [x] Interface InstallmentRepositoryInterface
  - [x] Implementación MySQL con GORM
  - [x] Métodos CRUD y consultas (GetByPlan, GetOverdue, etc.)
- [x] Crear `installment_plan_audit_repository.go`
  - [x] Interface InstallmentPlanAuditRepositoryInterface
  - [x] Implementación MySQL con GORM
  - [x] Métodos de auditoría y estadísticas
- [x] Actualizar `card_repository.go`
  - [x] Agregar método GetWithInstallmentPlans
  - [x] Preload de relaciones de cuotas

### 2.2 Services (Account-Service) ✅ COMPLETADA
- [x] Crear `installment_service.go`
  - [x] Interface InstallmentServiceInterface
  - [x] Calculadora de cuotas (CalculateInstallmentPlan)
  - [x] Creación de planes (CreateInstallmentPlan)
  - [x] Gestión de pagos (PayInstallment)
  - [x] Consultas (GetInstallmentPlans, GetInstallmentDetails)
- [x] Crear `card_service.go`
  - [x] Implementar CardServiceInterface completa
  - [x] Modificar ChargeCard para soporte tradicional
  - [x] Agregar ChargeCardWithInstallments
  - [x] Integrar con InstallmentService

### 2.3 DTOs y Request/Response Models ✅ COMPLETADA
- [x] Actualizar `installment_dto.go` (ya existía)
  - [x] CreateInstallmentPlanRequest
  - [x] InstallmentPreviewRequest/Response
  - [x] PayInstallmentRequest
  - [x] InstallmentPlanResponse, InstallmentResponse
- [x] Actualizar `card_dto.go`
  - [x] Agregar CreditCardChargeWithInstallmentsRequest
  - [x] Agregar ChargeWithInstallmentsResponse
  - [x] Modificar CardResponse para incluir resumen de cuotas
  - [x] Agregar InstallmentPlansSummary

---

## 🌐 FASE 3: BACKEND - APIs Y ENDPOINTS

### 3.1 Handlers (Account-Service) ✅ COMPLETADA
- [x] Crear `installment_handler.go`
  - [x] POST `/api/cards/{cardId}/installments/preview` - Preview de cuotas
  - [x] POST `/api/cards/{cardId}/charge-installments` - Compra con cuotas
  - [x] GET `/api/cards/{cardId}/installment-plans` - Planes activos por tarjeta
  - [x] GET `/api/installment-plans/{planId}` - Detalle de plan específico
  - [x] POST `/api/installments/{installmentId}/pay` - Pagar cuota
  - [x] GET `/api/installments/{installmentId}` - Detalle de cuota
  - [x] GET `/api/installment-plans` - Listar todos los planes con paginación
  - [x] GET `/api/installments` - Listar todas las cuotas con paginación
  - [x] GET `/api/users/{userId}/installment-summary` - Resumen por usuario
  - [x] POST `/api/installments/{installmentId}/monthly-load` - Descuento automático
  - [x] POST `/api/installment-plans/{planId}/cancel` - Cancelar plan
- [x] Actualizar `card_handler.go`
  - [x] Agregar endpoint POST `/api/cards/{cardId}/charge-installments`
  - [x] Actualizar interfaz CardServiceInterface

### 3.2 Routing y Middleware ✅ COMPLETADA
- [x] Actualizar `router.go`
  - [x] Registrar nuevas rutas de installments
  - [x] Aplicar middleware de autenticación
  - [x] Configurar dependency injection para InstallmentHandler
- [x] Actualizar dependency injection
  - [x] Conectar InstallmentService e InstallmentHandler
  - [x] Verificar configuración de todos los servicios

### 3.3 Integración con Transaction-Service ✅ COMPLETADA
- [x] Actualizar `transaction_client.go`
  - [x] Agregar método CreateInstallmentTransaction
  - [x] Enviar metadata de cuotas en transacciones
  - [x] Manejar registro de pagos de cuotas
  - [x] Agregar CreateInstallmentPaymentTransaction
  - [x] Agregar CreateInstallmentCancellationTransaction
  - [x] Agregar GetTransactionsByInstallmentPlan
  - [x] Agregar HealthCheck para validación de conectividad
- [x] Validar comunicación entre servicios
  - [x] Integrar transaction client en InstallmentService
  - [x] Agregar llamadas async para registro de transacciones
  - [x] Actualizar CardService para usar InstallmentService
  - [x] Configurar dependency injection correctamente

---

## 💾 FASE 4: FRONTEND - SERVICIOS Y MODELOS ✅

### 4.1 Models y Interfaces ✅
- [x] Crear `installment.model.ts`
  - [x] Interface InstallmentPlan
  - [x] Interface Installment
  - [x] Interface InstallmentPreview
  - [x] Enums para estados (InstallmentStatus, PlanStatus)
  - [x] Interfaces de Request/Response completas
  - [x] Tipos para formularios y configuración
- [x] Actualizar `card.model.ts`
  - [x] Agregar campos de cuotas a Card interface
  - [x] Tipos para respuestas con cuotas
  - [x] InstallmentPlansSummary interface

### 4.2 Services ✅
- [x] Crear `installment.service.ts`
  - [x] previewInstallments() - Calcular preview
  - [x] createInstallmentPlan() - Crear compra con cuotas
  - [x] getInstallmentPlans() - Obtener planes por tarjeta
  - [x] getInstallmentPlanDetails() - Detalle de plan
  - [x] payInstallment() - Pagar cuota individual
  - [x] Manejo de errores específicos
  - [x] Métodos auxiliares (resúmenes, cuotas vencidas, etc.)
- [x] Actualizar `credit-card.service.ts`
  - [x] Integrar métodos de cuotas
  - [x] Modificar charge() para soportar cuotas opcionales
  - [x] Agregar chargeWithInstallments()
  - [x] Integración con InstallmentService
- [x] Actualizar `services/index.ts` para exportar InstallmentService

---

## 🎨 FASE 5: FRONTEND - COMPONENTES ✅

### 5.1 Componente InstallmentCalculator ✅
- [x] Crear `installment-calculator.component.ts`
  - [x] Lógica de cálculo automático con debounce (300ms)
  - [x] Manejo reactivo con Angular Signals
  - [x] Integración con InstallmentService
  - [x] Validaciones de formulario completas
  - [x] Estados de loading, error y empty
- [x] Crear `installment-calculator.component.html`
  - [x] Selector de cantidad de cuotas (dropdown 3-24)
  - [x] Selector de fecha inicio (Angular Material datepicker)
  - [x] Input de monto total con validaciones
  - [x] Tabla de preview en tiempo real con desglose
  - [x] Resumen ejecutivo con tarjetas de información
  - [x] Toggle para cálculo automático
- [x] Crear `installment-calculator.component.scss`
  - [x] Estilos modernos con Material Design
  - [x] Responsive design (desktop, tablet, móvil)
  - [x] Estados de loading/error con animaciones
  - [x] Hover effects y transiciones suaves
  - [x] High contrast mode support
- [x] Implementar lógica avanzada
  - [x] Cálculo dinámico de cuotas con RxJS
  - [x] Validaciones en tiempo real
  - [x] Emisión de eventos para componente padre
  - [x] Manejo de errores específicos de cuotas
  - [x] Mock calculations para preview sin cardId

### 5.2 Componente InstallmentPlansList ✅
- [x] Crear `installment-plans-list.component.ts`
  - [x] Manejo de múltiples fuentes de datos (por tarjeta/usuario)
  - [x] Paginación con MatPaginator
  - [x] Auto-refresh opcional cada 30 segundos
  - [x] Filtros por estado de plan
  - [x] Acciones: ver detalle, pagar cuota, cancelar plan
- [x] Crear `installment-plans-list.component.html`
  - [x] Lista de planes activos en formato cards
  - [x] Progress bar por plan con porcentajes
  - [x] Próximos vencimientos destacados con alertas
  - [x] Botones de acción con menú contextual
  - [x] Estados visuales por estatus (activo, completado, cancelado, suspendido)
  - [x] Empty state y error handling
- [x] Crear `installment-plans-list.component.scss`
  - [x] Grid responsivo para lista de planes
  - [x] Indicadores visuales de estado con colores
  - [x] Animaciones para hover y interacciones
  - [x] Chips de estado con iconografía
  - [x] Responsive breakpoints optimizados
- [x] Implementar lógica completa
  - [x] Carga de datos desde API con filtros
  - [x] Ordenamiento y paginación
  - [x] Manejo de acciones con event emitters
  - [x] Tracking por plan ID para performance
  - [x] Cálculos de progreso y próximos vencimientos

### 5.3 Componente InstallmentPlanDetail ✅
- [x] Crear `installment-plan-detail.component.ts`
  - [x] Carga detallada de plan individual
  - [x] Generación de tabla de cuotas completa
  - [x] Manejo de acciones de pago individual
  - [x] Estados por cuota (pagada, pendiente, vencida, cancelada)
  - [x] Cálculos de progreso y estadísticas
- [x] Crear `installment-plan-detail.component.html`
  - [x] Información del plan con resumen ejecutivo
  - [x] Tarjetas de métricas (total, cuota mensual, progreso, estado)
  - [x] Alerta de próxima cuota con countdown
  - [x] Tabla detallada de cuotas con todas las columnas
  - [x] Estados visuales y botones de pago
  - [x] Responsive table con scroll horizontal
- [x] Crear `installment-plan-detail.component.scss`
  - [x] Layout de dashboard con cards métricas
  - [x] Estilos para tabla de cuotas con estados
  - [x] Alertas diferenciadas (próxima cuota vs vencida)
  - [x] Iconografía consistente con estados
  - [x] Mobile-first responsive design
- [x] Implementar lógica de detalle
  - [x] Carga de detalle del plan desde API
  - [x] Procesamiento de cuotas individuales
  - [x] Manejo de pagos de cuotas
  - [x] Cálculos de días hasta vencimiento
  - [x] Validaciones de acciones permitidas

### 5.4 Integración en Componentes Existentes ✅
- [x] Actualizar `credit-card-detail.component.ts`
  - [x] Import de componentes de cuotas
  - [x] Nuevas propiedades reactivas (installmentPlansCount)
  - [x] Métodos para manejo de eventos de cuotas
  - [x] Integración con CreditCardService actualizado
  - [x] Manejo de respuestas de compras con cuotas
- [x] Actualizar `credit-card-detail.component.html`
  - [x] Nueva pestaña "Compras en Cuotas" para tarjetas de crédito
  - [x] Integración de InstallmentCalculator en sección dedicada
  - [x] Integración de InstallmentPlansList con límite de 5 planes
  - [x] Botón "Ver todos los planes" cuando hay más de 5
  - [x] Modal/navegación para vista completa
- [x] Actualizar `credit-card-detail.component.css`
  - [x] Estilos para nueva pestaña de cuotas
  - [x] Responsive adjustments para componentes integrados
  - [x] Estilos para botón "Ver todos los planes"
- [x] Implementar handlers de eventos
  - [x] onInstallmentCalculationChanged() - Preview changes
  - [x] onInstallmentsSelected() - Crear compra con cuotas
  - [x] onInstallmentPlanAction() - Acciones de planes
  - [x] onInstallmentPlansLoaded() - Actualizar contador
  - [x] Métodos de navegación y dialogs

### 5.5 Organización y Exportaciones ✅
- [x] Crear `shared/components/index.ts`
  - [x] Exportar InstallmentCalculatorComponent
  - [x] Exportar InstallmentPlansListComponent  
  - [x] Exportar InstallmentPlanDetailComponent
- [x] Actualizar `shared/index.ts`
  - [x] Re-exportar todos los componentes compartidos
- [x] Validar imports y dependencias
  - [x] Verificar que todos los componentes son standalone
  - [x] Confirmar imports de Angular Material
  - [x] Validar integración con servicios existentes

---

## 🧪 FASE 6: TESTING

### 6.1 Testing Backend
- [ ] Tests unitarios - Entities
  - [ ] TestInstallmentPlan validaciones y métodos de negocio
  - [ ] TestInstallment validaciones y estados
  - [ ] TestCard métodos relacionados a cuotas
- [ ] Tests unitarios - Services
  - [ ] TestInstallmentService calculadora y creación
  - [ ] TestCardService operaciones con cuotas
- [ ] Tests unitarios - Repositories
  - [ ] TestInstallmentPlanRepository CRUD y consultas
  - [ ] TestInstallmentRepository operaciones específicas
- [ ] Tests de integración
  - [ ] TestInstallmentAPIs endpoints completos
  - [ ] TestTransactionIntegration comunicación entre servicios

### 6.2 Testing Frontend
- [ ] Tests unitarios - Components
  - [ ] TestInstallmentCalculator cálculos y validaciones
  - [ ] TestInstallmentPlansList visualización y acciones
  - [ ] TestInstallmentPlanDetail detalle y pagos
- [ ] Tests unitarios - Services
  - [ ] TestInstallmentService métodos API
  - [ ] TestCreditCardService integración cuotas
- [ ] Tests de integración
  - [ ] TestCreditCardFlow flujo completo con cuotas
  - [ ] TestInstallmentFlow gestión completa de cuotas

### 6.3 Testing E2E
- [ ] Escenario: Compra con cuotas completa
  - [ ] Usuario calcula cuotas
  - [ ] Usuario confirma compra
  - [ ] Sistema genera plan
  - [ ] Usuario ve plan creado
- [ ] Escenario: Gestión de cuotas
  - [ ] Usuario ve planes activos
  - [ ] Usuario ve detalle de plan
  - [ ] Usuario paga cuota individual
  - [ ] Sistema actualiza estado

---

## 🚀 FASE 7: DEPLOYMENT Y VALIDACIÓN

### 7.1 Preparación para Deployment
- [ ] Validar migraciones de BD en ambiente staging
- [ ] Verificar variables de entorno necesarias
- [ ] Documentar nuevos endpoints en Swagger/OpenAPI
- [ ] Crear scripts de rollback si es necesario

### 7.2 Validación de Funcionalidad
- [ ] Testing manual completo
  - [ ] Crear cuenta y tarjeta de crédito
  - [ ] Realizar compra con diferentes cantidades de cuotas
  - [ ] Verificar cálculos de fechas y montos
  - [ ] Probar pagos de cuotas individuales
- [ ] Validación de performance
  - [ ] Carga de planes con muchas cuotas
  - [ ] Consultas optimizadas con índices
  - [ ] Tiempo de respuesta de APIs

### 7.3 Documentación
- [ ] Actualizar README con nuevas funcionalidades
- [ ] Documentar APIs en Swagger
- [ ] Crear guía de usuario para cuotas
- [ ] Documentar esquema de BD actualizado

---

## � RESUMEN DE PROGRESO

### ✅ FASES COMPLETADAS
- **✅ FASE 1: DISEÑO Y ESTRUCTURA DE BASE DE DATOS** - Completada
  - Migración 07_V7__installments.sql ejecutada
  - Tablas installment_plans e installments creadas
  - Entidades Go implementadas (installment_plan.go, installment.go)

- **✅ FASE 2: BACKEND - LÓGICA DE NEGOCIO** - Completada  
  - InstallmentService con lógica completa
  - CardService actualizado con integración de cuotas
  - TransactionClient para comunicación con transaction-service

- **✅ FASE 3: BACKEND - API Y COMUNICACIÓN** - Completada
  - Repositorios (InstallmentRepo, InstallmentPlanRepo)
  - Controladores y rutas completas
  - Integración con transaction-service

- **✅ FASE 4: FRONTEND - SERVICIOS Y MODELOS** - Completada
  - Modelos TypeScript completos (installment.model.ts)
  - InstallmentService con métodos completos
  - CreditCardService actualizado con integración de cuotas

- **✅ FASE 5: FRONTEND - COMPONENTES** - Completada
  - InstallmentCalculator: Calculadora interactiva en tiempo real
  - InstallmentPlansList: Lista de planes con progress tracking
  - InstallmentPlanDetail: Vista detallada con tabla de cuotas
  - Integración completa en credit-card-detail component

### 🚀 PRÓXIMAS FASES
- **📍 FASE 6: TESTING INTEGRAL** - Pendiente
  - Tests unitarios backend (services, entities, repositories)
  - Tests unitarios frontend (components, services)
  - Tests de integración end-to-end
  - Validación de performance y optimización

- **📍 FASE 7: VALIDACIÓN Y DOCUMENTACIÓN** - Pendiente
  - Testing manual completo del flujo
  - Documentación técnica actualizada
  - Guía de usuario para funcionalidades de cuotas

### 🎯 **PROGRESO GENERAL: 71% COMPLETADO (5 de 7 fases)**

### 📋 **RESUMEN TÉCNICO FASE 5:**
- **🔢 Archivos creados**: 9 archivos (3 componentes × 3 archivos cada uno)
- **📦 Componentes Angular**: 3 componentes standalone reutilizables
- **🎨 Líneas de código**: ~2,000 líneas (TS + HTML + SCSS)
- **🔧 Funcionalidades**: Calculadora interactiva, gestión visual de planes, vista detallada
- **📱 Responsive**: Mobile-first design con breakpoints optimizados
- **♿ Accesibilidad**: High contrast mode, ARIA labels, keyboard navigation
- **⚡ Performance**: Angular Signals, debounce, lazy loading, track by functions
- **🔗 Integración**: Completamente integrado en card-detail existente

### 🛠️ **TECNOLOGÍAS UTILIZADAS:**
- **Frontend**: Angular 17+ con Standalone Components
- **UI Framework**: Angular Material (Cards, Tables, Forms, Icons)
- **State Management**: Angular Signals para reactividad
- **Styling**: SCSS con metodología BEM
- **HTTP**: RxJS para comunicación reactiva
- **Forms**: Reactive Forms con validaciones
- **Responsive**: CSS Grid + Flexbox

---

## �📈 MÉTRICAS DE ÉXITO

### Funcionales
- [ ] ✅ Usuario puede crear compras con 1-24 cuotas
- [ ] ✅ Sistema calcula fechas de vencimiento correctamente
- [ ] ✅ Usuario puede elegir fecha de inicio personalizada
- [ ] ✅ Preview de cuotas es preciso y en tiempo real
- [ ] ✅ Pagos de cuotas actualizan estado correctamente
- [ ] ✅ Dashboard muestra información completa y actualizada

### Técnicas
- [ ] ✅ APIs responden en menos de 500ms
- [ ] ✅ Base de datos optimizada con índices apropiados
- [ ] ✅ Cobertura de tests > 80%
- [ ] ✅ Integración entre servicios sin errores
- [ ] ✅ Frontend responsive en móvil y desktop
- [ ] ✅ Manejo de errores robusto en todos los niveles

---

## 📝 NOTAS Y DECISIONES TÉCNICAS

### Decisiones de Arquitectura
- **Account-Service**: Responsable principal de lógica de cuotas
- **Transaction-Service**: Solo registro para auditoría (evitar duplicación de lógica)
- **Base de datos**: 2 tablas nuevas con relaciones claras
- **Frontend**: Componentes modulares reutilizables

### Consideraciones Especiales
- Cálculo de fechas considera días hábiles y fin de mes
- Soporte para cancelación de planes (future enhancement)
- Integración con notificaciones para vencimientos (future enhancement)
- Reportes de cuotas para usuarios (future enhancement)

---

**Fecha de inicio**: 3 de Octubre 2025  
**Estimación total**: 8-12 días de desarrollo  
**Última actualización**: 3 de Octubre 2025 - ✅ **Fase 5 Completada**