# Configuración de EmailJS para Soporte en FinTrack

Este documento explica cómo configurar el servicio de email para el formulario de contacto de soporte.

## 🚀 Pasos para configurar EmailJS

### 1. Crear cuenta en EmailJS

1. Ve a [https://www.emailjs.com/](https://www.emailjs.com/)
2. Haz clic en "Sign Up" y crea una cuenta gratuita
3. Verifica tu email

### 2. Configurar el servicio de email

1. En el dashboard de EmailJS, ve a **"Email Services"**
2. Haz clic en **"Add New Service"**
3. Selecciona tu proveedor de email (Gmail, Outlook, etc.)
4. Sigue las instrucciones para conectar tu cuenta
5. **Guarda el Service ID** (ejemplo: `service_abc1234`)

### 3. Crear un template de email

1. Ve a **"Email Templates"** en el dashboard
2. Haz clic en **"Create New Template"**
3. Configura el template con los siguientes campos:

**Subject (Asunto):**
```
Nuevo mensaje de soporte - {{subject}}
```

**Content (Contenido):**
```
Has recibido un nuevo mensaje de soporte:

Nombre: {{from_name}}
Email: {{from_email}}
Asunto: {{subject}}
Fecha: {{timestamp}}

Mensaje:
{{message}}

---
Para responder, envía un email a: {{reply_to}}
```

4. **Guarda el Template ID** (ejemplo: `template_xyz5678`)

### 4. Obtener tu Public Key

1. Ve a **"Account"** → **"General"**
2. Encuentra tu **Public Key** (ejemplo: `AbCdEfGhIjKlMnOp`)
3. Copia esta clave

### 5. Configurar las credenciales en FinTrack

Abre el archivo `frontend/src/app/services/support-email.service.ts` y reemplaza los siguientes valores:

```typescript
private readonly SERVICE_ID = 'service_abc1234'; // Tu Service ID del paso 2
private readonly TEMPLATE_ID = 'template_xyz5678'; // Tu Template ID del paso 3
private readonly PUBLIC_KEY = 'AbCdEfGhIjKlMnOp'; // Tu Public Key del paso 4
```

### 6. Configurar el email de destino

En el mismo archivo, actualiza el email de destino de FinTrack:

```typescript
to_email: 'soporte@fintrack.com', // Cambia esto por tu email real
```

## 📧 Configuración de Gmail (Recomendado)

Si usas Gmail como servicio de email:

1. Ve a tu cuenta de Google
2. Activa la **verificación en dos pasos**
3. Genera una **contraseña de aplicación** para EmailJS
4. Usa esa contraseña al configurar el servicio en EmailJS

## 🧪 Probar el servicio

1. Instala las dependencias:
```bash
cd frontend
npm install
```

2. Inicia el servidor de desarrollo:
```bash
npm start
```

3. Ve a http://localhost:4200/faq
4. Haz clic en "Contactar Soporte"
5. Completa el formulario y envía un mensaje de prueba

## ✅ Verificación

Después de enviar un mensaje de prueba:
- Deberías ver una notificación de éxito en la aplicación
- Deberías recibir el email en la dirección configurada
- Verifica la bandeja de spam si no lo recibes

## 📊 Límites del plan gratuito

EmailJS ofrece:
- ✅ **200 emails/mes** gratis
- ✅ Sin tarjeta de crédito requerida
- ✅ Todos los templates necesarios

Para más emails, considera actualizar al plan pagado.

## 🔧 Troubleshooting

### Error: "Public Key not configured"
- Verifica que hayas reemplazado `YOUR_PUBLIC_KEY` con tu clave real

### Error: "Service not found"
- Verifica que el Service ID sea correcto
- Asegúrate de que el servicio esté activo en EmailJS

### Error: "Template not found"
- Verifica que el Template ID sea correcto
- Asegúrate de que el template esté publicado

### No recibo emails
- Revisa la carpeta de spam
- Verifica que el email de destino sea correcto
- Revisa los logs de EmailJS en el dashboard

## 🔒 Seguridad

Las credenciales de EmailJS son **públicas por diseño** (se usan en el frontend). EmailJS usa:
- CORS para limitar los dominios permitidos
- Rate limiting para prevenir abuso
- Captcha opcional para protección adicional

Para producción, considera:
- Configurar los dominios permitidos en EmailJS
- Activar reCAPTCHA si es necesario
- Monitorear el uso en el dashboard

## 📚 Documentación oficial

- [EmailJS Documentation](https://www.emailjs.com/docs/)
- [Angular Integration](https://www.emailjs.com/docs/sdk/installation/)
- [Template Variables](https://www.emailjs.com/docs/user-guide/dynamic-variables/)

---

¿Necesitas ayuda? Contacta al equipo de desarrollo de FinTrack.
