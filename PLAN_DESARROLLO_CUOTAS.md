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

### 2.1 Repositories (Account-Service)
- [ ] Crear `installment_plan_repository.go`
  - [ ] Interface InstallmentPlanRepositoryInterface
  - [ ] Implementación MySQL con GORM
  - [ ] Métodos CRUD básicos
  - [ ] Consultas especializadas (GetByCard, GetActive, etc.)
- [ ] Crear `installment_repository.go`
  - [ ] Interface InstallmentRepositoryInterface
  - [ ] Implementación MySQL con GORM
  - [ ] Métodos CRUD y consultas (GetByPlan, GetOverdue, etc.)
- [ ] Actualizar `card_repository.go`
  - [ ] Agregar método GetWithInstallmentPlans
  - [ ] Preload de relaciones de cuotas

### 2.2 Services (Account-Service)
- [ ] Crear `installment_service.go`
  - [ ] Interface InstallmentServiceInterface
  - [ ] Calculadora de cuotas (CalculateInstallmentPlan)
  - [ ] Creación de planes (CreateInstallmentPlan)
  - [ ] Gestión de pagos (PayInstallment)
  - [ ] Consultas (GetInstallmentPlans, GetInstallmentDetails)
- [ ] Actualizar `card_service.go`
  - [ ] Modificar ChargeCard para soporte de cuotas
  - [ ] Agregar ChargeCardWithInstallments
  - [ ] Integrar con InstallmentService

### 2.3 DTOs y Request/Response Models
- [ ] Crear `installment_dto.go`
  - [ ] CreateInstallmentPlanRequest
  - [ ] InstallmentPreviewRequest/Response
  - [ ] PayInstallmentRequest
  - [ ] InstallmentPlanResponse, InstallmentResponse
- [ ] Actualizar `card_dto.go`
  - [ ] Agregar CreditCardChargeWithInstallmentsRequest
  - [ ] Modificar responses para incluir datos de cuotas

---

## 🌐 FASE 3: BACKEND - APIs Y ENDPOINTS

### 3.1 Handlers (Account-Service)
- [ ] Crear `installment_handler.go`
  - [ ] POST `/api/cards/{cardId}/installments/preview` - Preview de cuotas
  - [ ] POST `/api/cards/{cardId}/charge-installments` - Compra con cuotas
  - [ ] GET `/api/cards/{cardId}/installment-plans` - Planes activos por tarjeta
  - [ ] GET `/api/installment-plans/{planId}` - Detalle de plan específico
  - [ ] POST `/api/installments/{installmentId}/pay` - Pagar cuota
  - [ ] GET `/api/installments/{installmentId}` - Detalle de cuota
- [ ] Actualizar `card_handler.go`
  - [ ] Modificar respuestas para incluir información de cuotas
  - [ ] Agregar validaciones para operaciones con cuotas

### 3.2 Routing y Middleware
- [ ] Actualizar `router.go`
  - [ ] Registrar nuevas rutas de installments
  - [ ] Aplicar middleware de autenticación
  - [ ] Configurar rate limiting si es necesario

### 3.3 Integración con Transaction-Service
- [ ] Actualizar `transaction_client.go`
  - [ ] Agregar método CreateInstallmentTransaction
  - [ ] Enviar metadata de cuotas en transacciones
  - [ ] Manejar registro de pagos de cuotas
- [ ] Validar comunicación entre servicios

---

## 💾 FASE 4: FRONTEND - SERVICIOS Y MODELOS

### 4.1 Models y Interfaces
- [ ] Crear `installment.model.ts`
  - [ ] Interface InstallmentPlan
  - [ ] Interface Installment
  - [ ] Interface InstallmentPreview
  - [ ] Enums para estados (InstallmentStatus, PlanStatus)
- [ ] Actualizar `card.model.ts`
  - [ ] Agregar campos de cuotas a Card interface
  - [ ] Tipos para respuestas con cuotas

### 4.2 Services
- [ ] Crear `installment.service.ts`
  - [ ] previewInstallments() - Calcular preview
  - [ ] createInstallmentPlan() - Crear compra con cuotas
  - [ ] getInstallmentPlans() - Obtener planes por tarjeta
  - [ ] getInstallmentPlanDetails() - Detalle de plan
  - [ ] payInstallment() - Pagar cuota individual
  - [ ] Manejo de errores específicos
- [ ] Actualizar `credit-card.service.ts`
  - [ ] Integrar métodos de cuotas
  - [ ] Modificar charge() para soportar cuotas opcionales
  - [ ] Agregar métodos de consulta con cuotas

---

## 🎨 FASE 5: FRONTEND - COMPONENTES

### 5.1 Componente InstallmentCalculator
- [ ] Crear `installment-calculator.component.ts`
- [ ] Crear `installment-calculator.component.html`
  - [ ] Selector de cantidad de cuotas (dropdown)
  - [ ] Selector de fecha inicio (date picker)
  - [ ] Input de monto total
  - [ ] Tabla de preview en tiempo real
- [ ] Crear `installment-calculator.component.scss`
  - [ ] Estilos para calculadora
  - [ ] Responsive design
  - [ ] Estados de loading/error
- [ ] Implementar lógica
  - [ ] Cálculo dinámico de cuotas
  - [ ] Validaciones en tiempo real
  - [ ] Emisión de eventos para componente padre

### 5.2 Componente InstallmentPlansList
- [ ] Crear `installment-plans-list.component.ts`
- [ ] Crear `installment-plans-list.component.html`
  - [ ] Lista de planes activos (cards/acordeón)
  - [ ] Progress bar por plan
  - [ ] Próximos vencimientos destacados
  - [ ] Botones de acción (ver detalle, pagar)
- [ ] Crear `installment-plans-list.component.scss`
  - [ ] Estilos para lista de planes
  - [ ] Indicadores visuales de estado
  - [ ] Animaciones para expansión
- [ ] Implementar lógica
  - [ ] Carga de datos desde API
  - [ ] Filtros y ordenamiento
  - [ ] Paginación si es necesario

### 5.3 Componente InstallmentPlanDetail
- [ ] Crear `installment-plan-detail.component.ts`
- [ ] Crear `installment-plan-detail.component.html`
  - [ ] Información del plan (resumen)
  - [ ] Tabla detallada de cuotas
  - [ ] Estados visuales (pagada, pendiente, vencida)
  - [ ] Botones de pago individual
- [ ] Crear `installment-plan-detail.component.scss`
  - [ ] Estilos para detalle de plan
  - [ ] Estados de cuotas (colores, iconos)
- [ ] Implementar lógica
  - [ ] Carga de detalle del plan
  - [ ] Pago de cuotas individuales
  - [ ] Actualización de estados

### 5.4 Integración en Componentes Existentes
- [ ] Actualizar `credit-card-detail.component.ts`
  - [ ] Agregar pestaña "Compras en Cuotas"
  - [ ] Integrar InstallmentCalculator en modal de compra
  - [ ] Mostrar InstallmentPlansList
- [ ] Actualizar `credit-card-detail.component.html`
  - [ ] Nueva pestaña en tabs
  - [ ] Modal modificado para incluir opción de cuotas
  - [ ] Dashboard de compromisos futuros
- [ ] Actualizar `credit-card-charge.component.ts` (si existe)
  - [ ] Checkbox para "Pagar en cuotas"
  - [ ] Mostrar/ocultar InstallmentCalculator
  - [ ] Validaciones combinadas

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

## 📈 MÉTRICAS DE ÉXITO

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
**Última actualización**: 3 de Octubre 2025