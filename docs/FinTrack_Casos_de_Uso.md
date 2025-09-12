# FinTrack - Casos de Uso e Historias de Usuario

## 📋 Información del Documento

- **Proyecto:** FinTrack - Plataforma de Gestión Financiera
- **Versión:** 1.0
- **Fecha:** Enero 2024
- **Autor:** Estudiante UNT - Tecnicatura en Programación

---

## 👥 Actores del Sistema

### Actores Principales

#### 1. Usuario Final (User)
- **Descripción:** Persona que utiliza la plataforma para gestionar sus finanzas personales
- **Características:**
  - Tiene una o más cuentas en el sistema
  - Puede realizar transacciones
  - Accede a reportes personales
  - Interactúa con el chatbot

#### 2. Operador (Operator)
- **Descripción:** Personal de soporte que asiste a los usuarios
- **Características:**
  - Puede ver información de usuarios (limitada)
  - Puede asistir en resolución de problemas
  - Acceso a herramientas de soporte
  - No puede realizar transacciones por usuarios

#### 3. Administrador (Admin)
- **Descripción:** Administrador del sistema con acceso completo
- **Características:**
  - Gestión completa de usuarios
  - Configuración del sistema
  - Acceso a todos los reportes
  - Gestión de roles y permisos

#### 4. Tesorero (Treasurer)
- **Descripción:** Responsable de la gestión financiera y reportes ejecutivos
- **Características:**
  - Acceso a reportes financieros avanzados
  - Análisis de flujos de dinero
  - Gestión de liquidez
  - Reportes regulatorios

### Actores Secundarios

#### 5. Sistema Bancario (Banking API)
- **Descripción:** APIs externas de bancos y procesadores de pago
- **Interacciones:**
  - Validación de tarjetas
  - Procesamiento de transacciones
  - Consulta de saldos

#### 6. Servicio de Cotizaciones (Exchange API)
- **Descripción:** API externa para obtener tipos de cambio
- **Interacciones:**
  - Consulta de cotizaciones en tiempo real
  - Histórico de tipos de cambio

#### 7. Servicio de Notificaciones (Notification Service)
- **Descripción:** Servicios externos para envío de notificaciones
- **Interacciones:**
  - Envío de emails
  - Envío de SMS
  - Push notifications

---

## 📊 Casos de Uso por Módulo

### 🔐 Módulo de Autenticación

#### CU-001: Registro de Usuario

**Actor Principal:** Usuario Final

**Precondiciones:**
- El usuario no tiene cuenta en el sistema
- El email no está registrado previamente

**Flujo Principal:**
1. El usuario accede a la página de registro
2. El usuario completa el formulario con:
   - Nombre completo
   - Email
   - Contraseña
   - Confirmación de contraseña
   - Número de teléfono
   - Fecha de nacimiento
3. El sistema valida los datos ingresados
4. El sistema envía un email de verificación
5. El usuario confirma su email
6. El sistema crea la cuenta y perfil del usuario
7. El sistema redirige al dashboard

**Flujos Alternativos:**
- **3a.** Email ya registrado:
  - El sistema muestra error
  - Ofrece opción de recuperar contraseña
- **3b.** Contraseña no cumple criterios:
  - El sistema muestra los requisitos
  - El usuario corrige la contraseña
- **5a.** Email no verificado en 24h:
  - El sistema permite reenvío de verificación

**Postcondiciones:**
- Usuario registrado en el sistema
- Cuenta de usuario creada
- Email verificado

**Historia de Usuario:**
```
Como usuario nuevo
Quiero registrarme en FinTrack
Para poder gestionar mis finanzas digitales

Criterios de Aceptación:
- Puedo registrarme con email y contraseña
- Recibo confirmación por email
- Mi cuenta se activa tras verificar email
- Puedo acceder al dashboard tras registro
```

#### CU-002: Inicio de Sesión

**Actor Principal:** Usuario Final

**Precondiciones:**
- El usuario tiene cuenta registrada
- La cuenta está activa

**Flujo Principal:**
1. El usuario accede a la página de login
2. El usuario ingresa email y contraseña
3. El sistema valida las credenciales
4. El sistema genera token JWT
5. El sistema redirige al dashboard

**Flujos Alternativos:**
- **3a.** Credenciales incorrectas:
  - El sistema muestra error genérico
  - Incrementa contador de intentos fallidos
- **3b.** Cuenta bloqueada:
  - El sistema muestra mensaje de cuenta bloqueada
  - Ofrece opción de contactar soporte
- **3c.** Demasiados intentos fallidos:
  - El sistema bloquea temporalmente la IP
  - Requiere CAPTCHA para siguientes intentos

**Postcondiciones:**
- Usuario autenticado en el sistema
- Sesión activa creada
- Token JWT válido generado

**Historia de Usuario:**
```
Como usuario registrado
Quiero iniciar sesión en FinTrack
Para acceder a mi información financiera

Criterios de Aceptación:
- Puedo ingresar con email y contraseña
- El sistema me recuerda por 30 días (opcional)
- Recibo error claro si las credenciales son incorrectas
- Mi sesión se mantiene activa por tiempo razonable
```

### 🏦 Módulo de Gestión de Cuentas

#### CU-003: Crear Cuenta Virtual

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Usuario no excede límite de cuentas (5 por usuario)

**Flujo Principal:**
1. El usuario accede a "Crear Nueva Cuenta"
2. El usuario selecciona tipo de cuenta:
   - Cuenta Corriente Virtual
   - Cuenta de Ahorros Virtual
   - Cuenta en USD
3. El usuario ingresa:
   - Nombre de la cuenta
   - Moneda (ARS/USD)
   - Descripción (opcional)
4. El sistema valida los datos
5. El sistema genera número de cuenta único
6. El sistema crea la cuenta con saldo inicial $0
7. El sistema muestra confirmación

**Flujos Alternativos:**
- **2a.** Límite de cuentas alcanzado:
  - El sistema muestra mensaje de límite
  - Ofrece upgrade a plan premium
- **4a.** Nombre de cuenta duplicado:
  - El sistema solicita nombre diferente

**Postcondiciones:**
- Nueva cuenta virtual creada
- Cuenta visible en dashboard
- Número de cuenta único asignado

**Historia de Usuario:**
```
Como usuario
Quiero crear cuentas virtuales
Para organizar mi dinero por categorías o monedas

Criterios de Aceptación:
- Puedo crear hasta 5 cuentas gratuitas
- Puedo elegir entre ARS y USD
- Cada cuenta tiene un número único
- Puedo personalizar el nombre de la cuenta
```

#### CU-004: Vincular Tarjeta Bancaria

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Usuario tiene al menos una cuenta virtual
- Tarjeta válida y activa

**Flujo Principal:**
1. El usuario accede a "Vincular Tarjeta"
2. El usuario ingresa datos de la tarjeta:
   - Número de tarjeta
   - Fecha de vencimiento
   - CVV
   - Nombre del titular
3. El sistema valida formato de datos
4. El sistema envía datos a API bancaria (sandbox)
5. La API bancaria valida la tarjeta
6. El sistema realiza cargo de verificación ($1)
7. El sistema solicita confirmación del cargo
8. El usuario confirma el cargo
9. El sistema vincula la tarjeta a la cuenta
10. El sistema tokeniza los datos de la tarjeta

**Flujos Alternativos:**
- **5a.** Tarjeta inválida:
  - El sistema muestra error de validación
  - Permite reintentar con datos correctos
- **6a.** Cargo de verificación falla:
  - El sistema muestra error de procesamiento
  - Sugiere contactar al banco
- **8a.** Usuario no confirma cargo:
  - El sistema cancela la vinculación
  - Revierte el cargo de verificación

**Postcondiciones:**
- Tarjeta vinculada y tokenizada
- Tarjeta disponible para transacciones
- Datos sensibles no almacenados

**Historia de Usuario:**
```
Como usuario
Quiero vincular mis tarjetas bancarias
Para poder cargar dinero a mis cuentas virtuales

Criterios de Aceptación:
- Puedo vincular tarjetas de débito y crédito
- Mis datos están seguros (tokenizados)
- Recibo confirmación de vinculación exitosa
- Puedo desvincular tarjetas cuando quiera
```

### 💰 Módulo de Billetera Digital

#### CU-005: Cargar Saldo desde Tarjeta

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Tarjeta vinculada y activa
- Cuenta virtual seleccionada

**Flujo Principal:**
1. El usuario selecciona "Cargar Saldo"
2. El usuario selecciona:
   - Cuenta destino
   - Tarjeta origen
   - Monto a cargar
3. El sistema valida:
   - Monto mínimo ($100 ARS)
   - Monto máximo diario ($500,000 ARS)
   - Límites de la tarjeta
4. El sistema muestra resumen de la operación
5. El usuario confirma la transacción
6. El sistema procesa el pago con API bancaria
7. El sistema actualiza el saldo de la cuenta
8. El sistema envía notificación de confirmación
9. El sistema registra la transacción

**Flujos Alternativos:**
- **3a.** Monto excede límites:
  - El sistema muestra límites aplicables
  - Permite ajustar el monto
- **6a.** Pago rechazado:
  - El sistema muestra motivo del rechazo
  - Sugiere acciones correctivas
- **6b.** Error de conectividad:
  - El sistema reintenta automáticamente
  - Si falla, marca transacción como pendiente

**Postcondiciones:**
- Saldo actualizado en cuenta virtual
- Transacción registrada
- Notificación enviada
- Comprobante disponible

**Historia de Usuario:**
```
Como usuario
Quiero cargar dinero desde mis tarjetas
Para tener saldo disponible en mi billetera virtual

Criterios de Aceptación:
- Puedo cargar desde cualquier tarjeta vinculada
- El dinero se refleja inmediatamente
- Recibo confirmación de la operación
- Puedo ver el comprobante de la transacción
```

#### CU-006: Transferir entre Usuarios

**Actor Principal:** Usuario Final (Emisor)

**Precondiciones:**
- Usuario autenticado
- Saldo suficiente en cuenta origen
- Usuario destinatario existe en el sistema

**Flujo Principal:**
1. El usuario selecciona "Transferir Dinero"
2. El usuario ingresa:
   - Email del destinatario
   - Monto a transferir
   - Cuenta origen
   - Descripción (opcional)
3. El sistema valida:
   - Existencia del destinatario
   - Saldo suficiente
   - Límites de transferencia
4. El sistema muestra resumen de la transferencia
5. El usuario confirma la operación
6. El sistema solicita autenticación adicional (SMS/Email)
7. El usuario ingresa código de verificación
8. El sistema procesa la transferencia:
   - Debita cuenta origen
   - Acredita cuenta destino
9. El sistema notifica a ambos usuarios
10. El sistema registra la transacción

**Flujos Alternativos:**
- **3a.** Destinatario no existe:
  - El sistema sugiere invitar al usuario
  - Permite enviar invitación por email
- **3b.** Saldo insuficiente:
  - El sistema muestra saldo disponible
  - Sugiere cargar saldo
- **7a.** Código de verificación incorrecto:
  - Permite reintentar (máximo 3 veces)
  - Bloquea operación tras 3 fallos

**Postcondiciones:**
- Transferencia completada
- Saldos actualizados
- Ambos usuarios notificados
- Transacción registrada

**Historia de Usuario:**
```
Como usuario
Quiero transferir dinero a otros usuarios
Para pagar o enviar dinero de forma rápida y segura

Criterios de Aceptación:
- Puedo transferir a cualquier usuario registrado
- La transferencia es instantánea
- Ambos recibimos notificación
- Puedo agregar una descripción a la transferencia
```

### 🤖 Módulo de Chatbot

#### CU-007: Consultar Saldo via Chatbot

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Chatbot disponible

**Flujo Principal:**
1. El usuario abre el chatbot
2. El usuario escribe consulta sobre saldo:
   - "¿Cuál es mi saldo?"
   - "Saldo de mi cuenta en USD"
   - "¿Cuánto dinero tengo?"
3. El chatbot procesa la consulta con NLP
4. El chatbot identifica la intención (consulta_saldo)
5. El chatbot consulta la información de cuentas
6. El chatbot formatea la respuesta
7. El chatbot muestra saldos de todas las cuentas
8. El chatbot ofrece acciones relacionadas

**Flujos Alternativos:**
- **3a.** Consulta ambigua:
  - El chatbot solicita clarificación
  - Ofrece opciones específicas
- **5a.** Error al consultar datos:
  - El chatbot se disculpa
  - Sugiere intentar más tarde
  - Ofrece contactar soporte

**Postcondiciones:**
- Usuario informado sobre sus saldos
- Conversación registrada
- Métricas de uso actualizadas

**Historia de Usuario:**
```
Como usuario
Quiero consultar mi saldo mediante el chatbot
Para obtener información rápida sin navegar por la app

Criterios de Aceptación:
- Puedo preguntar en lenguaje natural
- El chatbot entiende diferentes formas de preguntar
- Recibo información clara y actualizada
- El chatbot me ofrece acciones relacionadas
```

#### CU-008: Iniciar Transferencia via Chatbot

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Chatbot disponible
- Usuario tiene saldo suficiente

**Flujo Principal:**
1. El usuario solicita transferencia al chatbot:
   - "Quiero transferir $1000 a juan@email.com"
   - "Enviar dinero a María"
2. El chatbot extrae información:
   - Monto
   - Destinatario
   - Cuenta origen (si especificada)
3. El chatbot valida la información
4. El chatbot solicita confirmación de datos faltantes
5. El usuario confirma o corrige información
6. El chatbot muestra resumen de la transferencia
7. El chatbot solicita confirmación final
8. El usuario confirma
9. El chatbot redirige al flujo de transferencia estándar
10. El chatbot notifica el resultado

**Flujos Alternativos:**
- **2a.** Información incompleta:
  - El chatbot solicita datos faltantes
  - Guía al usuario paso a paso
- **3a.** Destinatario no válido:
  - El chatbot sugiere verificar el email
  - Ofrece buscar en contactos
- **8a.** Usuario cancela:
  - El chatbot confirma cancelación
  - Ofrece ayuda adicional

**Postcondiciones:**
- Transferencia iniciada o cancelada
- Usuario informado del resultado
- Conversación registrada

**Historia de Usuario:**
```
Como usuario
Quiero iniciar transferencias mediante el chatbot
Para realizar pagos de forma conversacional y rápida

Criterios de Aceptación:
- Puedo especificar monto y destinatario en lenguaje natural
- El chatbot me guía si falta información
- Puedo confirmar o cancelar antes de procesar
- Recibo confirmación del resultado
```

### 📊 Módulo de Reportes

#### CU-009: Generar Reporte de Movimientos

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Usuario tiene transacciones en el período

**Flujo Principal:**
1. El usuario accede a "Reportes"
2. El usuario selecciona "Reporte de Movimientos"
3. El usuario configura filtros:
   - Rango de fechas
   - Cuentas específicas
   - Tipos de transacción
   - Moneda
4. El usuario selecciona formato de exportación:
   - PDF
   - Excel
   - CSV
5. El sistema valida los parámetros
6. El sistema genera el reporte
7. El sistema muestra vista previa
8. El usuario confirma la exportación
9. El sistema genera el archivo
10. El sistema envía enlace de descarga

**Flujos Alternativos:**
- **3a.** Rango de fechas muy amplio:
  - El sistema sugiere reducir el rango
  - Advierte sobre tiempo de procesamiento
- **6a.** No hay datos para los filtros:
  - El sistema informa que no hay datos
  - Sugiere ajustar filtros
- **9a.** Error en generación:
  - El sistema reintenta automáticamente
  - Si falla, notifica al usuario

**Postcondiciones:**
- Reporte generado exitosamente
- Archivo disponible para descarga
- Actividad registrada

**Historia de Usuario:**
```
Como usuario
Quiero generar reportes de mis movimientos
Para llevar control de mis finanzas y presentar documentación

Criterios de Aceptación:
- Puedo filtrar por fechas, cuentas y tipos de transacción
- Puedo exportar en PDF, Excel o CSV
- El reporte incluye toda la información relevante
- Puedo descargar el archivo generado
```

### 🔔 Módulo de Notificaciones

#### CU-010: Configurar Alertas de Saldo

**Actor Principal:** Usuario Final

**Precondiciones:**
- Usuario autenticado
- Usuario tiene al menos una cuenta

**Flujo Principal:**
1. El usuario accede a "Configuración de Notificaciones"
2. El usuario selecciona "Alertas de Saldo Bajo"
3. El usuario configura para cada cuenta:
   - Umbral de saldo mínimo
   - Método de notificación (Email/SMS/Push)
   - Frecuencia de alertas
4. El usuario guarda la configuración
5. El sistema valida los parámetros
6. El sistema activa las alertas
7. El sistema confirma la configuración

**Flujos Alternativos:**
- **3a.** Umbral muy bajo:
  - El sistema sugiere un mínimo razonable
  - Permite continuar con advertencia
- **5a.** Configuración inválida:
  - El sistema muestra errores específicos
  - Permite corregir antes de guardar

**Postcondiciones:**
- Alertas configuradas y activas
- Sistema monitoreando saldos
- Usuario notificado de la configuración

**Historia de Usuario:**
```
Como usuario
Quiero configurar alertas de saldo bajo
Para ser notificado cuando necesite cargar dinero

Criterios de Aceptación:
- Puedo configurar diferentes umbrales por cuenta
- Puedo elegir cómo recibir las notificaciones
- Puedo activar/desactivar alertas individualmente
- Recibo confirmación de la configuración
```

---

## 📈 Casos de Uso Administrativos

### 👨‍💼 Casos de Uso del Administrador

#### CU-011: Gestionar Usuarios

**Actor Principal:** Administrador

**Precondiciones:**
- Administrador autenticado
- Permisos de gestión de usuarios

**Flujo Principal:**
1. El administrador accede al panel de usuarios
2. El administrador puede:
   - Ver lista de todos los usuarios
   - Buscar usuarios por criterios
   - Ver detalles de un usuario específico
   - Activar/desactivar cuentas
   - Cambiar roles de usuario
   - Ver historial de actividad
3. El administrador selecciona una acción
4. El sistema solicita confirmación para acciones críticas
5. El administrador confirma la acción
6. El sistema ejecuta la acción
7. El sistema registra la actividad administrativa
8. El sistema notifica al usuario afectado (si aplica)

**Flujos Alternativos:**
- **4a.** Acción requiere justificación:
  - El sistema solicita motivo
  - El administrador proporciona justificación
- **6a.** Error en ejecución:
  - El sistema revierte cambios parciales
  - Notifica al administrador del error

**Postcondiciones:**
- Acción administrativa completada
- Actividad registrada en audit log
- Usuario notificado si corresponde

#### CU-012: Monitorear Transacciones Sospechosas

**Actor Principal:** Administrador

**Precondiciones:**
- Administrador autenticado
- Sistema de detección de fraude activo

**Flujo Principal:**
1. El sistema detecta transacción sospechosa
2. El sistema genera alerta automática
3. El administrador recibe notificación
4. El administrador revisa la transacción:
   - Detalles de la transacción
   - Historial del usuario
   - Patrones de comportamiento
5. El administrador toma una decisión:
   - Aprobar transacción
   - Rechazar transacción
   - Solicitar información adicional
   - Bloquear cuenta temporalmente
6. El sistema ejecuta la decisión
7. El sistema notifica al usuario
8. El sistema actualiza modelos de detección

**Flujos Alternativos:**
- **5a.** Información insuficiente:
  - El administrador solicita más datos
  - El sistema recopila información adicional
- **5b.** Caso complejo:
  - El administrador escala a supervisor
  - Se inicia investigación formal

**Postcondiciones:**
- Transacción procesada según decisión
- Usuario informado del resultado
- Caso documentado para análisis

### 💼 Casos de Uso del Tesorero

#### CU-013: Generar Reporte de Liquidez

**Actor Principal:** Tesorero

**Precondiciones:**
- Tesorero autenticado
- Datos financieros disponibles

**Flujo Principal:**
1. El tesorero accede a "Reportes Ejecutivos"
2. El tesorero selecciona "Reporte de Liquidez"
3. El tesorero configura parámetros:
   - Período de análisis
   - Monedas a incluir
   - Nivel de detalle
4. El sistema recopila datos de:
   - Saldos totales por moneda
   - Flujos de entrada y salida
   - Proyecciones de liquidez
   - Reservas requeridas
5. El sistema genera análisis avanzado
6. El sistema presenta dashboard interactivo
7. El tesorero puede exportar el reporte

**Flujos Alternativos:**
- **4a.** Datos incompletos:
  - El sistema identifica gaps de información
  - Proporciona estimaciones con disclaimers
- **5a.** Cálculos complejos toman tiempo:
  - El sistema muestra progreso
  - Permite continuar en background

**Postcondiciones:**
- Reporte de liquidez generado
- Dashboard actualizado
- Datos disponibles para análisis

---

## 🔄 Casos de Uso de Integración

### 🏛️ Integración con APIs Bancarias

#### CU-014: Sincronizar Estado de Tarjeta

**Actor Principal:** Sistema (Proceso Automático)

**Precondiciones:**
- Tarjetas vinculadas en el sistema
- API bancaria disponible

**Flujo Principal:**
1. El sistema inicia sincronización programada
2. Para cada tarjeta vinculada:
   - Consulta estado en API bancaria
   - Verifica límites actuales
   - Comprueba fecha de vencimiento
3. El sistema compara con datos locales
4. Si hay cambios:
   - Actualiza información local
   - Notifica al usuario si es relevante
5. El sistema registra resultado de sincronización

**Flujos Alternativos:**
- **2a.** API bancaria no disponible:
  - El sistema reintenta con backoff exponencial
  - Registra fallo para monitoreo
- **2b.** Tarjeta desactivada en banco:
  - El sistema marca tarjeta como inactiva
  - Notifica al usuario

**Postcondiciones:**
- Estados de tarjetas actualizados
- Usuarios notificados de cambios relevantes
- Logs de sincronización registrados

### 💱 Integración con API de Cotizaciones

#### CU-015: Actualizar Tipos de Cambio

**Actor Principal:** Sistema (Proceso Automático)

**Precondiciones:**
- API de cotizaciones configurada
- Monedas soportadas definidas

**Flujo Principal:**
1. El sistema consulta API de cotizaciones cada 15 minutos
2. El sistema obtiene cotizaciones para:
   - USD/ARS
   - EUR/ARS
   - BRL/ARS
3. El sistema valida coherencia de datos
4. El sistema actualiza cache de cotizaciones
5. El sistema notifica cambios significativos (>5%)
6. El sistema actualiza conversiones en tiempo real

**Flujos Alternativos:**
- **1a.** API no disponible:
  - El sistema usa última cotización válida
  - Marca datos como "no actualizados"
- **3a.** Datos inconsistentes:
  - El sistema rechaza actualización
  - Mantiene cotizaciones anteriores
  - Alerta al administrador

**Postcondiciones:**
- Cotizaciones actualizadas
- Cache actualizado
- Usuarios informados de cambios significativos

---

## 📱 Casos de Uso Móviles

### 📲 Funcionalidades Específicas Mobile

#### CU-016: Autenticación Biométrica

**Actor Principal:** Usuario Final (Mobile)

**Precondiciones:**
- App móvil instalada
- Dispositivo con sensor biométrico
- Usuario ha configurado biometría

**Flujo Principal:**
1. El usuario abre la aplicación móvil
2. La app detecta biometría configurada
3. La app solicita autenticación biométrica
4. El usuario proporciona huella/Face ID
5. El sistema valida la biometría
6. El sistema autentica al usuario
7. La app redirige al dashboard

**Flujos Alternativos:**
- **5a.** Biometría no reconocida:
  - La app permite reintentar (3 veces)
  - Tras fallos, solicita PIN/contraseña
- **5b.** Sensor biométrico no disponible:
  - La app usa autenticación tradicional

**Postcondiciones:**
- Usuario autenticado
- Sesión segura establecida
- Experiencia de usuario optimizada

#### CU-017: Notificaciones Push

**Actor Principal:** Sistema de Notificaciones

**Precondiciones:**
- Usuario tiene app móvil instalada
- Notificaciones push habilitadas
- Evento que requiere notificación

**Flujo Principal:**
1. El sistema detecta evento notificable:
   - Transacción recibida
   - Saldo bajo
   - Login desde nuevo dispositivo
2. El sistema determina usuarios a notificar
3. El sistema compone mensaje de notificación
4. El sistema envía push notification
5. El dispositivo recibe y muestra notificación
6. El usuario puede:
   - Ver detalles en la app
   - Ignorar notificación
   - Configurar preferencias

**Flujos Alternativos:**
- **4a.** Dispositivo offline:
  - El sistema reintenta envío
  - Almacena para entrega posterior
- **5a.** Notificaciones deshabilitadas:
  - El sistema registra intento
  - No envía notificación

**Postcondiciones:**
- Usuario informado del evento
- Engagement con la aplicación
- Métricas de notificación registradas

---

## 🧪 Casos de Uso de Testing

### 🔍 Casos de Uso para QA

#### CU-018: Ejecutar Suite de Pruebas Automatizadas

**Actor Principal:** Sistema de CI/CD

**Precondiciones:**
- Código committeado en repositorio
- Pipeline de CI/CD configurado
- Entorno de testing disponible

**Flujo Principal:**
1. El sistema detecta nuevo commit
2. El sistema inicia pipeline de testing:
   - Unit tests
   - Integration tests
   - API tests
   - E2E tests
3. El sistema ejecuta pruebas en paralelo
4. El sistema recopila resultados
5. El sistema genera reporte de cobertura
6. El sistema notifica resultados al equipo

**Flujos Alternativos:**
- **3a.** Pruebas fallan:
  - El sistema detiene deployment
  - Notifica detalles del fallo
  - Permite re-ejecutar pruebas
- **3b.** Timeout en pruebas:
  - El sistema cancela ejecución
  - Marca como fallo
  - Investiga causa del timeout

**Postcondiciones:**
- Calidad del código validada
- Reporte de cobertura generado
- Equipo informado de resultados

---

## 📊 Métricas y KPIs de Casos de Uso

### Métricas de Éxito

| Caso de Uso | Métrica Principal | Target | Crítico |
|-------------|-------------------|--------|----------|
| CU-001: Registro | Tasa de conversión | >60% | >40% |
| CU-002: Login | Tiempo de autenticación | <3s | <5s |
| CU-005: Carga de saldo | Tasa de éxito | >95% | >90% |
| CU-006: Transferencias | Tiempo de procesamiento | <10s | <30s |
| CU-007: Chatbot consultas | Tasa de comprensión | >85% | >70% |
| CU-009: Reportes | Tiempo de generación | <30s | <60s |

### Métricas de Usabilidad

| Aspecto | Métrica | Target |
|---------|---------|--------|
| Facilidad de uso | SUS Score | >80 |
| Satisfacción | CSAT | >4.5/5 |
| Eficiencia | Tareas completadas | >90% |
| Errores de usuario | Error rate | <5% |

---

## 🔮 Casos de Uso Futuros

### Roadmap de Funcionalidades

#### Fase 2: Funcionalidades Avanzadas
- **CU-019:** Inversiones simples (plazo fijo virtual)
- **CU-020:** Préstamos entre usuarios
- **CU-021:** Cashback y recompensas
- **CU-022:** Análisis de gastos con IA

#### Fase 3: Expansión
- **CU-023:** Integración con más bancos
- **CU-024:** Pagos con QR
- **CU-025:** Marketplace de servicios financieros
- **CU-026:** API pública para terceros

---

## 📞 Contacto y Validación

### Proceso de Validación
1. **Revisión con stakeholders**
2. **Prototipado de flujos críticos**
3. **Testing con usuarios reales**
4. **Iteración basada en feedback**

### Información de Contacto
- **Autor:** Estudiante UNT - Tecnicatura en Programación
- **Email:** fintrack.requirements@example.com
- **Versión:** 1.0
- **Última actualización:** Enero 2024

---

*Documento de Casos de Uso para Tesis de Tecnicatura en Programación - UNT*
*Este documento será actualizado iterativamente durante el desarrollo del proyecto*