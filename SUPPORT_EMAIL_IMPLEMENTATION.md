# Implementación de Soporte por Email - FAQ

## 📋 Resumen

Se implementó un sistema completo de contacto con soporte desde la sección FAQ de FinTrack, utilizando el **microservicio de notificaciones** existente en lugar de EmailJS directo desde el frontend.

## ✅ Arquitectura

### Backend (notification-service)

**Nuevo endpoint agregado:**
```
POST /api/notifications/support
```

**Request Body:**
```json
{
  "name": "string",
  "email": "string", 
  "subject": "string",
  "message": "string"
}
```

**Response (200 OK):**
```json
{
  "message": "Support email sent successfully",
  "timestamp": "2025-10-27T..."
}
```

### Archivos modificados en Backend:

1. **`internal/core/ports/repository.go`**
   - Agregado método `SendSupportEmail` a interfaces `EmailService` y `NotificationService`

2. **`internal/core/service/notification_service.go`**
   - Implementado método `SendSupportEmail()` que delega al EmailService

3. **`internal/infrastructure/adapters/email/emailjs_client.go`**
   - Implementado método `SendSupportEmail()` usando template `template_yst8bd2`
   - Envía email a `soporte@fintrack.com`
   - Reply-to configurado con el email del usuario

4. **`internal/infrastructure/entrypoints/handlers/notification/notification_handler.go`**
   - Agregado handler `SendSupportEmail()` con validación de request
   - Struct `SupportEmailRequest` con validaciones

5. **`internal/infrastructure/entrypoints/router/router.go`**
   - Ruta `POST /api/notifications/support` agregada
   - Endpoint documentado en lista de endpoints

### Frontend

**Archivos creados:**

1. **`services/email.service.ts`** (reemplaza `support-email.service.ts`)
   ```typescript
   - sendSupportEmail(data: SupportEmailData): Observable<any>
   - Llama a POST /api/notifications/support
   ```

2. **`pages/faq/support-dialog/support-dialog.component.ts`**
   - Modal de contacto con formulario reactivo
   - Validaciones: email, longitud mínima, campos requeridos
   - Muestra mensaje de éxito/error
   - Integrado con AuthService (pre-rellena datos del usuario)

3. **`pages/faq/support-dialog/support-dialog.component.html`**
   - Formulario con: nombre, email, asunto, mensaje
   - Estados: normal, enviando, éxito, error
   - Material Design

4. **`pages/faq/support-dialog/support-dialog.component.css`**
   - Estilos consistentes con el resto de la aplicación

**Archivos modificados:**

5. **`pages/faq/faq.component.ts`**
   - Método `openSupportDialog()` abre el nuevo modal
   - Importa `SupportDialogComponent`

**Archivos eliminados (obsoletos):**
- ❌ `services/support-email.service.ts` (usaba EmailJS directo)
- ❌ `shared/components/support-contact-dialog/` (componente viejo)

## 🔧 Configuración

### EmailJS Template (template_yst8bd2)

El template debe tener las siguientes variables:

```
{{from_name}}         - Nombre del usuario
{{from_email}}        - Email del usuario  
{{subject}}           - Asunto del mensaje
{{message}}           - Mensaje del usuario
{{to_email}}          - Email de soporte (soporte@fintrack.com)
{{reply_to}}          - Email del usuario (para responder)
```

### Variables de Entorno (Backend)

Ya configuradas en `config.json`:
```json
{
  "emailjs": {
    "service_id": "service_ceg7xlp",
    "template_id": "template_yst8bd2",
    "public_key": "MSBb87-PQcXWr1gWK",
    "from_name": "FinTrack Notifications",
    "reply_to": "noreply@fintrack.com"
  }
}
```

## 🎯 Flujo de Funcionamiento

1. Usuario hace clic en **"Contacta con Soporte"** en la sección FAQ
2. Se abre el modal `SupportDialogComponent`
3. Usuario completa el formulario (nombre, email, asunto, mensaje)
4. Al enviar:
   - Frontend valida el formulario
   - Hace POST a `/api/notifications/support`
   - Backend recibe la request y valida los datos
   - `NotificationService.SendSupportEmail()` procesa la solicitud
   - `EmailJSClient.SendSupportEmail()` envía el email usando template `template_yst8bd2`
   - Email llega a `soporte@fintrack.com` (configurado en código)
   - Response 200 OK regresa al frontend
   - Modal muestra mensaje de éxito y se cierra automáticamente

## 📧 Email Enviado

**Para:** soporte@fintrack.com  
**De:** Usuario (name, email)  
**Reply-To:** email del usuario  
**Asunto:** subject del formulario  
**Contenido:** message del formulario  

## ✅ Ventajas de esta Arquitectura

1. **Seguridad:** Credenciales de EmailJS solo en backend
2. **Centralización:** Todo el email pasa por el microservicio de notificaciones
3. **Reutilización:** Usa la misma infraestructura que las notificaciones de tarjetas
4. **Mantenibilidad:** Un solo lugar para cambiar configuración de email
5. **Escalabilidad:** Fácil agregar rate limiting, logging, etc. en el backend

## 🧪 Testing

### Probar localmente:

1. Asegurar que el `notification-service` esté corriendo:
   ```bash
   docker-compose up notification-service
   ```

2. Acceder a FAQ:
   ```
   http://localhost:4200/faq
   ```

3. Hacer clic en "Contacta con Soporte"

4. Completar y enviar el formulario

5. Verificar:
   - Logs del backend: `📧 Sending support email from...`
   - Email recibido en la bandeja de entrada configurada

### Test con cURL:

```bash
curl -X POST http://localhost:8084/api/notifications/support \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Juan Pérez",
    "email": "juan@example.com",
    "subject": "Consulta sobre reportes",
    "message": "Hola, tengo una pregunta sobre cómo generar reportes mensuales..."
  }'
```

**Response esperado:**
```json
{
  "message": "Support email sent successfully",
  "timestamp": "2025-10-27T..."
}
```

## 📝 Notas Importantes

- **Template ID fijo:** Siempre usa `template_yst8bd2` (hardcoded en backend)
- **Email de soporte:** Actualmente `soporte@fintrack.com` (cambiar en código si es necesario)
- **Sin autenticación:** El endpoint `/support` no requiere token (acceso público)
- **Validaciones:** Todos los campos son requeridos, email debe ser válido

## 🔜 Mejoras Futuras

- [ ] Agregar rate limiting (máximo X emails por usuario/IP por hora)
- [ ] Guardar logs de emails de soporte en base de datos
- [ ] Agregar categorías de soporte (técnico, facturación, general)
- [ ] Sistema de tickets con ID único
- [ ] Auto-respuesta al usuario confirmando recepción
- [ ] Panel de admin para ver emails de soporte

---

**Implementado:** 27 de Octubre, 2025  
**Estado:** ✅ Completado y funcional
