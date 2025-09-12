# FinTrack - Diagramas de Flujo de Trabajo

## Índice
1. [Flujo de Registro de Usuario](#flujo-de-registro-de-usuario)
2. [Flujo de Autenticación](#flujo-de-autenticación)
3. [Flujo de Vinculación de Tarjeta](#flujo-de-vinculación-de-tarjeta)
4. [Flujo de Transferencia entre Usuarios](#flujo-de-transferencia-entre-usuarios)
5. [Flujo de Carga de Billetera](#flujo-de-carga-de-billetera)
6. [Flujo de Interacción con Chatbot](#flujo-de-interacción-con-chatbot)
7. [Flujo de Generación de Reportes](#flujo-de-generación-de-reportes)
8. [Flujo de Notificaciones](#flujo-de-notificaciones)

---

## Flujo de Registro de Usuario

```mermaid
flowchart TD
    A[Usuario accede a registro] --> B[Completa formulario]
    B --> C{¿Email válido?}
    C -->|No| D[Mostrar error]
    D --> B
    C -->|Sí| E{¿Email ya existe?}
    E -->|Sí| F[Mostrar error: Email registrado]
    F --> B
    E -->|No| G[Validar contraseña]
    G --> H{¿Contraseña segura?}
    H -->|No| I[Mostrar requisitos]
    I --> B
    H -->|Sí| J[Crear usuario temporal]
    J --> K[Enviar email verificación]
    K --> L[Mostrar mensaje: Verificar email]
    L --> M[Usuario hace clic en enlace]
    M --> N{¿Token válido?}
    N -->|No| O[Mostrar error: Token inválido]
    N -->|Sí| P[Activar cuenta]
    P --> Q[Crear billetera virtual]
    Q --> R[Crear cuentas por defecto]
    R --> S[Redirigir a dashboard]
    S --> T[Mostrar tutorial onboarding]
```

**Descripción**: Este flujo maneja el proceso completo de registro, desde la captura de datos hasta la activación de la cuenta y creación de recursos iniciales.

**Puntos Críticos**:
- Validación de email único
- Verificación por email obligatoria
- Creación automática de billetera y cuentas
- Tutorial de onboarding

---

## Flujo de Autenticación

```mermaid
flowchart TD
    A[Usuario ingresa credenciales] --> B{¿Email existe?}
    B -->|No| C[Error: Usuario no encontrado]
    B -->|Sí| D{¿Cuenta activa?}
    D -->|No| E[Error: Cuenta desactivada]
    D -->|Sí| F{¿Contraseña correcta?}
    F -->|No| G[Incrementar intentos fallidos]
    G --> H{¿Máximo intentos?}
    H -->|Sí| I[Bloquear cuenta temporalmente]
    H -->|No| J[Error: Contraseña incorrecta]
    F -->|Sí| K[Resetear intentos fallidos]
    K --> L{¿2FA habilitado?}
    L -->|Sí| M[Solicitar código 2FA]
    M --> N{¿Código válido?}
    N -->|No| O[Error: Código inválido]
    N -->|Sí| P[Generar JWT token]
    L -->|No| P
    P --> Q[Generar refresh token]
    Q --> R[Actualizar último login]
    R --> S[Redirigir a dashboard]
    S --> T[Cargar datos usuario]
```

**Descripción**: Proceso de autenticación con validaciones de seguridad, soporte para 2FA y generación de tokens JWT.

**Puntos Críticos**:
- Protección contra ataques de fuerza bruta
- Soporte opcional para 2FA
- Gestión de tokens JWT y refresh
- Logging de intentos de acceso

---

## Flujo de Carga Manual de Tarjeta

```mermaid
flowchart TD
    A[Usuario accede a 'Mis Tarjetas'] --> B[Clic en 'Agregar Tarjeta']
    B --> C[Formulario de carga manual]
    C --> D[Ingresa datos básicos]
    D --> E[Nombre personalizado]
    E --> F[Banco emisor]
    F --> G[Tipo de tarjeta]
    G --> H[Últimos 4 dígitos]
    H --> I[Fecha de vencimiento]
    I --> J[Nombre del titular]
    J --> K[Límite de crédito opcional]
    K --> L[Validación de formato]
    L --> M{¿Datos válidos?}
    M -->|No| N[Mostrar errores]
    N --> D
    M -->|Sí| O[Encriptar datos sensibles]
    O --> P[Guardar en MySQL]
    P --> Q[Mostrar confirmación]
    Q --> R[Actualizar lista de tarjetas]
    R --> S[Enviar notificación]
```

**Descripción**: Proceso de carga manual de tarjetas sin integración con APIs externas, almacenando datos básicos de forma segura.

**Puntos Críticos**:
- Validación de formato de datos
- Encriptación de información sensible
- Almacenamiento seguro en MySQL
- Interfaz intuitiva para carga manual

---

## Flujo de Transferencia entre Usuarios

```mermaid
flowchart TD
    A[Usuario accede a 'Transferir'] --> B[Seleccionar cuenta origen]
    B --> C[Ingresar email destinatario]
    C --> D[Ingresar monto]
    D --> E{¿Destinatario existe?}
    E -->|No| F[Error: Usuario no encontrado]
    E -->|Sí| G{¿Saldo suficiente?}
    G -->|No| H[Error: Saldo insuficiente]
    G -->|Sí| I[Mostrar resumen operación]
    I --> J[Usuario confirma]
    J --> K[Iniciar transacción]
    K --> L[Bloquear fondos]
    L --> M{¿Débito exitoso?}
    M -->|No| N[Rollback y error]
    M -->|Sí| O[Acreditar a destinatario]
    O --> P{¿Crédito exitoso?}
    P -->|No| Q[Rollback débito]
    Q --> R[Error: Transacción fallida]
    P -->|Sí| S[Confirmar transacción]
    S --> T[Generar comprobante]
    T --> U[Notificar a ambos usuarios]
    U --> V[Registrar en auditoría]
    V --> W[Mostrar confirmación]
```

**Descripción**: Proceso de transferencia con validaciones, transacciones atómicas y notificaciones.

**Puntos Críticos**:
- Validación de destinatario y saldo
- Transacciones atómicas (rollback en caso de error)
- Notificaciones a ambas partes
- Registro de auditoría completo

---

## Flujo de Gestión de Billeteras Virtuales Centralizadas

```mermaid
flowchart TD
    A[Usuario accede a Billeteras] --> B[Vista centralizada]
    B --> C{¿Qué acción?}
    C -->|Agregar| D[Selecciona 'Agregar Billetera']
    C -->|Sincronizar| E[Selecciona billetera existente]
    C -->|Transferir| F[Selecciona transferencia P2P]
    
    D --> G[Lista de proveedores]
    G --> H[Selecciona proveedor]
    H --> I[MercadoPago/Ualá/Brubank/etc]
    I --> J[Formulario de carga manual]
    J --> K[Nombre personalizado]
    K --> L[Identificador de billetera]
    L --> M[Saldo inicial manual]
    M --> N[Validar datos]
    N --> O{¿Datos válidos?}
    O -->|No| P[Mostrar errores]
    P --> J
    O -->|Sí| Q[Guardar en MySQL]
    Q --> R[Actualizar vista consolidada]
    
    E --> S[Actualización manual de saldo]
    S --> T[Usuario ingresa nuevo saldo]
    T --> U[Actualizar timestamp]
    U --> V[Guardar cambios]
    V --> R
    
    F --> W[Selecciona billetera origen]
    W --> X[Selecciona billetera destino]
    X --> Y[Ingresa monto]
    Y --> Z[Confirmar transferencia]
    Z --> AA[Actualizar saldos]
    AA --> R
    
    R --> BB[Fin]
```

**Descripción**: Proceso de gestión centralizada de múltiples billeteras virtuales con carga manual de datos y sincronización controlada por el usuario.

**Puntos Críticos**:
- Gestión centralizada de múltiples billeteras
- Carga manual de datos
- Sincronización manual de saldos
- Vista consolidada de activos digitales

---

## Flujo de Interacción con Chatbot

```mermaid
flowchart TD
    A[Usuario abre chat] --> B[Escribir mensaje]
    B --> C[Enviar al servidor]
    C --> D[Procesar con NLP]
    D --> E{¿Intención reconocida?}
    E -->|No| F[Mostrar opciones predefinidas]
    F --> G[Usuario selecciona opción]
    G --> H[Procesar selección]
    E -->|Sí| I{¿Requiere autenticación?}
    I -->|Sí| J{¿Usuario autenticado?}
    J -->|No| K[Solicitar login]
    J -->|Sí| L[Verificar permisos]
    I -->|No| M[Procesar consulta general]
    L --> N{¿Tiene permisos?}
    N -->|No| O[Error: Sin permisos]
    N -->|Sí| P{¿Tipo de consulta?}
    P -->|Saldo| Q[Consultar saldo]
    P -->|Movimientos| R[Consultar transacciones]
    P -->|Transferencia| S[Iniciar flujo transferencia]
    P -->|FAQ| T[Buscar en base conocimiento]
    Q --> U[Formatear respuesta]
    R --> U
    S --> V[Solicitar confirmación]
    T --> U
    M --> U
    U --> W[Enviar respuesta]
    W --> X[Registrar conversación]
    X --> Y[Mostrar al usuario]
```

**Descripción**: Flujo de procesamiento de consultas del chatbot con NLP y diferentes tipos de respuestas.

**Puntos Críticos**:
- Procesamiento de lenguaje natural
- Autenticación para operaciones sensibles
- Diferentes tipos de consultas
- Registro de conversaciones

---

## Flujo de Generación de Reportes

```mermaid
flowchart TD
    A[Usuario accede a 'Reportes'] --> B[Seleccionar tipo reporte]
    B --> C[Configurar filtros]
    C --> D[Seleccionar rango fechas]
    D --> E[Elegir formato exportación]
    E --> F[Clic en 'Generar']
    F --> G{¿Permisos suficientes?}
    G -->|No| H[Error: Sin permisos]
    G -->|Sí| I[Validar parámetros]
    I --> J{¿Parámetros válidos?}
    J -->|No| K[Mostrar errores]
    J -->|Sí| L[Consultar base datos]
    L --> M{¿Datos encontrados?}
    M -->|No| N[Mostrar: Sin datos]
    M -->|Sí| O[Procesar datos]
    O --> P{¿Formato solicitado?}
    P -->|Dashboard| Q[Generar gráficos]
    P -->|Excel| R[Crear archivo XLSX]
    P -->|PDF| S[Generar documento PDF]
    Q --> T[Mostrar dashboard]
    R --> U[Descargar archivo]
    S --> U
    T --> V[Permitir exportación]
    U --> W[Registrar descarga]
    V --> X[Fin]
    W --> X
```

**Descripción**: Proceso de generación de reportes con múltiples formatos y filtros personalizables.

**Puntos Críticos**:
- Validación de permisos por rol
- Filtros flexibles de datos
- Múltiples formatos de exportación
- Optimización de consultas grandes

---

## Flujo de Notificaciones

```mermaid
flowchart TD
    A[Evento del sistema] --> B{¿Tipo de evento?}
    B -->|Transacción| C[Crear notificación transacción]
    B -->|Saldo bajo| D[Crear notificación saldo]
    B -->|Seguridad| E[Crear notificación seguridad]
    B -->|Sistema| F[Crear notificación sistema]
    C --> G[Determinar destinatarios]
    D --> G
    E --> G
    F --> G
    G --> H[Consultar preferencias usuario]
    H --> I{¿Notificación habilitada?}
    I -->|No| J[Fin - No enviar]
    I -->|Sí| K{¿Canales preferidos?}
    K -->|Email| L[Enviar email]
    K -->|SMS| M[Enviar SMS]
    K -->|Push| N[Enviar push notification]
    K -->|In-app| O[Guardar en BD]
    L --> P[Registrar envío]
    M --> P
    N --> P
    O --> P
    P --> Q{¿Envío exitoso?}
    Q -->|No| R[Registrar error]
    Q -->|Sí| S[Marcar como enviado]
    R --> T[Programar reintento]
    S --> U[Actualizar estadísticas]
    T --> V[Fin]
    U --> V
```

**Descripción**: Sistema de notificaciones multi-canal con preferencias de usuario y manejo de errores.

**Puntos Críticos**:
- Múltiples canales de notificación
- Respeto a preferencias de usuario
- Manejo de errores y reintentos
- Registro de estadísticas

---

## Arquitectura de Flujos

### Patrones Comunes Identificados

1. **Validación en Capas**
   - Validación frontend (Angular UX)
   - Validación backend (Go seguridad)
   - Validación de negocio (lógica)

2. **Manejo de Errores**
   - Errores específicos por contexto
   - Logging detallado con Go
   - Recuperación automática cuando es posible

3. **Transacciones Atómicas**
   - Operaciones financieras con rollback MySQL
   - Consistencia de datos
   - Auditoría completa

4. **Carga Manual de Datos**
   - Formularios intuitivos en Angular
   - Validación client-side y server-side
   - Almacenamiento seguro en MySQL
   - Sincronización manual controlada por usuario

5. **Notificaciones Asíncronas**
   - Procesamiento en background
   - Múltiples canales
   - Preferencias de usuario

### Consideraciones de Performance

- **Caching**: Redis para datos de sesión y billeteras consolidadas
- **Queue System**: Para procesamiento asíncrono con Go
- **Database Optimization**: Índices en user_id y provider_id, consultas optimizadas con GORM
- **API Rate Limiting**: Protección contra abuso con middleware Go
- **Microservicios**: Arquitectura independiente por funcionalidad

### Monitoreo y Alertas

- **Métricas de Negocio**: Transacciones, usuarios activos
- **Métricas Técnicas**: Latencia, errores, throughput
- **Alertas Automáticas**: Fallos críticos, performance degradation
- **Dashboards**: Tiempo real para operaciones

---

*Diagramas de flujo para FinTrack - Tecnicatura en Programación UNT*  
*Versión 1.0 - Enero 2025*