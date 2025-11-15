# FAQ con Sistema de Soporte - Implementación Completa

## 📋 Resumen

Se ha implementado una sección completa de **Preguntas Frecuentes (FAQ)** con un sistema de contacto de soporte integrado que envía emails a través del **microservicio de notificaciones**.

---

## ✅ Características Implementadas

### 1. **Página de FAQ** (`/faq`)

#### Categorías de Preguntas:
- ℹ️ **General** (4 preguntas) - FinTrack, seguridad, aplicación móvil
- 💼 **Cuentas** (4 preguntas) - Creación, tipos, eliminación
- 💳 **Tarjetas** (4 preguntas) - Gestión, límites de crédito, transacciones
- 💸 **Transacciones** (4 preguntas) - Tipos, edición, cuotas
- 📊 **Reportes** (4 preguntas) - Generación, exportación, chatbot

#### Funcionalidades:
- ✅ **20+ preguntas frecuentes** organizadas por categoría
- ✅ **Filtros por categoría** para navegación rápida
- ✅ **Acordeón Material Design** para expandir/colapsar preguntas
- ✅ **Botón de contacto con soporte** con modal de formulario
- ✅ **Diseño responsive** siguiendo el estilo de FinTrack

---

### 2. **Sistema de Contacto de Soporte**

#### Modal de Contacto:
- 📧 **Formulario completo** con validaciones:
  - Nombre (mínimo 3 caracteres)
  - Email (validación de formato)
  - Asunto (mínimo 5 caracteres)
  - Mensaje (mínimo 20 caracteres)
- ✅ **Auto-relleno** con datos del usuario autenticado
- ✅ **Indicadores visuales** de envío y éxito
- ✅ **Manejo de errores** con mensajes informativos

#### Integración Backend:
- 🔌 **Endpoint nuevo**: `POST /api/notifications/support`
- 📨 **EmailJS**: Usa template `template_yst8bd2`
- 🎯 **Email destino**: `soporte@fintrack.com`
- 🔄 **Reply-to**: Email del usuario para responder directamente

---

## 📁 Archivos Creados/Modificados

### Frontend (`frontend/src/app/`)

#### Nuevos Archivos:
```
pages/faq/
  ├── faq.component.ts              ✅ Componente principal FAQ
  ├── faq.component.html            ✅ Template con acordeón
  ├── faq.component.css             ✅ Estilos consistentes
  └── support-dialog/
      ├── support-dialog.component.ts    ✅ Modal de soporte
      ├── support-dialog.component.html  ✅ Formulario de contacto
      └── support-dialog.component.css   ✅ Estilos del modal
```

#### Modificados:
```
app.routes.ts                      ✅ Ruta /faq agregada
app.component.html                 ✅ Link en navegación
services/email.service.ts          ✅ Integración con backend
```

---

### Backend (`backend/services/notification-service/`)

#### Modificados:
```
internal/core/ports/repository.go
  ├── EmailService.SendSupportEmail()        ✅ Interface actualizada
  └── NotificationService.SendSupportEmail() ✅ Interface actualizada

internal/core/service/notification_service.go
  └── SendSupportEmail()                     ✅ Lógica de negocio

internal/infrastructure/adapters/email/emailjs_client.go
  └── SendSupportEmail()                     ✅ Implementación EmailJS

internal/infrastructure/entrypoints/handlers/notification/notification_handler.go
  ├── SupportEmailRequest                    ✅ DTO de request
  └── SendSupportEmail()                     ✅ Handler HTTP

internal/infrastructure/entrypoints/router/router.go
  └── POST /api/notifications/support        ✅ Nueva ruta
```

---

## 🔧 Configuración Técnica

### EmailJS Template (`template_yst8bd2`)

**Parámetros enviados al template:**
```javascript
{
  from_name: "Nombre del usuario",
  subject: "Asunto del mensaje",
  to_email: "soporte@fintrack.com",
  reply_to: "email@usuario.com",
  message: "Mensaje del usuario",
  user_email: "email@usuario.com",
  user_name: "Nombre del usuario"
}
```

### Backend Endpoint

**Request:**
```json
POST /api/notifications/support
Content-Type: application/json

{
  "name": "Juan Pérez",
  "email": "juan@example.com",
  "subject": "Problema con tarjetas",
  "message": "No puedo ver mis tarjetas de crédito..."
}
```

**Response (éxito):**
```json
{
  "message": "Support email sent successfully",
  "timestamp": "2025-01-27T10:30:00Z"
}
```

**Response (error):**
```json
{
  "error": "Failed to send support email",
  "details": "error details here..."
}
```

---

## 🚀 Cómo Usar

### Para Usuarios:

1. **Acceder al FAQ:**
   - Click en **"Preguntas Frecuentes"** en la navegación principal
   - O navega a: `http://localhost:4200/faq`

2. **Buscar respuestas:**
   - Usa los **botones de categoría** para filtrar preguntas
   - Click en cualquier pregunta para expandir la respuesta
   - Busca en las **20+ preguntas disponibles**

3. **Contactar soporte:**
   - Scroll hasta el final de la página
   - Click en **"Contacta con Soporte"**
   - Completa el formulario (auto-rellena tus datos)
   - Click en **"Enviar mensaje"**
   - ✅ Confirmación de envío exitoso

---

## 🎨 Diseño y UX

### Consistencia Visual:
- ✅ Usa el mismo formato de **balance-cards** del dashboard
- ✅ Colores y espaciado consistentes con el sistema de diseño
- ✅ Iconos Material Design para mejor identificación
- ✅ Responsive design para todos los dispositivos

### Accesibilidad:
- ✅ Etiquetas semánticas correctas
- ✅ Validaciones en tiempo real
- ✅ Mensajes de error claros
- ✅ Indicadores de carga durante el envío

---

## 🔒 Seguridad

### Validaciones:
- ✅ **Backend**: Validación de campos requeridos y formato email
- ✅ **Frontend**: Validaciones reactivas con Angular
- ✅ **CORS**: Configurado en el microservicio de notificaciones

### Privacidad:
- ✅ No se almacenan emails de soporte en base de datos
- ✅ EmailJS maneja el envío de forma segura
- ✅ Reply-to configurado para respuestas directas

---

## 📊 Flujo de Datos

```
Usuario → FAQ Component → Modal Soporte → Email Service
                                              ↓
                                    HTTP POST /api/notifications/support
                                              ↓
                            Notification Service (Backend)
                                              ↓
                                    EmailJS Client
                                              ↓
                              EmailJS API (template_yst8bd2)
                                              ↓
                                   soporte@fintrack.com
```

---

## ✨ Mejoras Futuras (Opcionales)

1. **Historial de tickets** en base de datos
2. **Sistema de tickets** con estados y respuestas
3. **Adjuntar archivos** en el formulario de soporte
4. **Rating de respuestas** de FAQ
5. **Búsqueda de texto** en las preguntas
6. **Panel de admin** para gestionar tickets de soporte

---

## 🧪 Testing

### Cómo probar:

1. **Asegúrate de que el microservicio de notificaciones esté corriendo:**
   ```bash
   docker-compose up notification-service
   ```

2. **Accede al FAQ:**
   ```
   http://localhost:4200/faq
   ```

3. **Prueba el formulario de soporte:**
   - Completa todos los campos
   - Envía el mensaje
   - Verifica el email en `soporte@fintrack.com`
   - Verifica que llegue con el formato correcto del template

4. **Verifica logs del backend:**
   ```
   📧 Sending support email from Juan Pérez (juan@example.com): Problema con tarjetas
   ✅ Support email sent successfully
   ```

---

## 📞 Soporte

Si tienes problemas con la implementación:
- Verifica que el **notification-service** esté corriendo
- Revisa las **credenciales de EmailJS** en el backend
- Verifica que el **template_yst8bd2** esté configurado correctamente
- Revisa los **logs del microservicio** para errores

---

**Implementado por:** GitHub Copilot  
**Fecha:** 27 de Octubre, 2025  
**Versión:** 1.0.0
