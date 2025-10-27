# Template de Email de Soporte - Diseño Profesional

## 📧 Nuevo Diseño HTML para Emails de Soporte

Se ha implementado un template HTML profesional y visualmente atractivo para los emails de soporte de FinTrack.

---

## 🎨 Características del Diseño

### 1. **Header con Gradiente**
- Fondo degradado violeta (667eea → 764ba2)
- Título: "💬 Nuevo Mensaje de Soporte"
- Subtítulo: "FinTrack Support System"
- Colores: Texto blanco con sombra para mejor legibilidad

### 2. **Banner de Alerta**
- Fondo amarillo claro (#fef3c7)
- Borde izquierdo naranja (#f59e0b)
- Mensaje: "⚡ Acción requerida: Un usuario necesita asistencia"
- Propósito: Llamar la atención del equipo de soporte

### 3. **Tarjeta de Información del Usuario**
- Fondo gris claro (#f9fafb)
- Bordes redondeados
- Información organizada:
  - 👤 Nombre (en negrita)
  - 📧 Email (link clickeable)
  - 📅 Fecha y hora

### 4. **Sección de Asunto**
- Fondo morado claro (#ede9fe)
- Borde izquierdo morado (#7c3aed)
- Texto destacado en morado oscuro

### 5. **Sección de Mensaje**
- Fondo gris muy claro
- Bordes sutiles
- Texto formateado (respeta saltos de línea)
- Fuente clara y legible

### 6. **Botón de Acción**
- Fondo con gradiente violeta
- Botón blanco con texto violeta
- Texto: "✉️ Responder al Usuario"
- Link mailto pre-configurado con:
  - Destinatario: email del usuario
  - Asunto: "Re: [asunto original]"

### 7. **Footer**
- Información del sistema
- Copyright
- Texto en gris claro

---

## 📝 Estructura del Template

```html
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="background: #f5f5f5;">
    <table width="600px" style="max-width: 95%; background: white;">
        
        <!-- HEADER -->
        <tr>
            <td style="background: linear-gradient(135deg, #667eea, #764ba2);">
                <h1>💬 Nuevo Mensaje de Soporte</h1>
                <p>FinTrack Support System</p>
            </td>
        </tr>
        
        <!-- ALERT BANNER -->
        <tr>
            <td style="background: #fef3c7; border-left: 4px solid #f59e0b;">
                <p>⚡ Acción requerida: Usuario necesita asistencia</p>
            </td>
        </tr>
        
        <!-- CONTENT -->
        <tr>
            <td style="padding: 40px 30px;">
                
                <!-- User Info Card -->
                <table style="background: #f9fafb; border-radius: 8px;">
                    <tr><td><strong>Nombre:</strong> Juan Pérez</td></tr>
                    <tr><td><strong>Email:</strong> juan@example.com</td></tr>
                    <tr><td><strong>Fecha:</strong> 27/10/2025 12:00:00</td></tr>
                </table>
                
                <!-- Subject -->
                <div style="background: #ede9fe; border-left: 4px solid #7c3aed;">
                    <h3>📋 Asunto</h3>
                    <p>Consulta sobre reportes</p>
                </div>
                
                <!-- Message -->
                <div style="background: #f9fafb; border: 1px solid #e5e7eb;">
                    <p>Mensaje del usuario aquí...</p>
                </div>
                
                <!-- Action Button -->
                <a href="mailto:usuario@email.com?subject=Re: Asunto" 
                   style="background: white; color: #667eea; padding: 12px 30px;">
                    ✉️ Responder al Usuario
                </a>
                
            </td>
        </tr>
        
        <!-- FOOTER -->
        <tr>
            <td style="background: #f9fafb; text-align: center;">
                <p>Email automático del sistema de soporte</p>
                <p>© 2025 FinTrack</p>
            </td>
        </tr>
        
    </table>
</body>
</html>
```

---

## 🎨 Paleta de Colores

| Elemento | Color | Hex Code |
|----------|-------|----------|
| **Gradiente Principal** | Violeta → Morado | `#667eea` → `#764ba2` |
| **Fondo Página** | Gris muy claro | `#f5f5f5` |
| **Tarjetas** | Gris claro | `#f9fafb` |
| **Alerta Fondo** | Amarillo claro | `#fef3c7` |
| **Alerta Borde** | Naranja | `#f59e0b` |
| **Asunto Fondo** | Violeta muy claro | `#ede9fe` |
| **Asunto Borde** | Morado | `#7c3aed` |
| **Texto Principal** | Gris oscuro | `#1f2937` |
| **Texto Secundario** | Gris medio | `#6b7280` |
| **Links** | Violeta | `#667eea` |

---

## 📱 Responsividad

El template está optimizado para:
- ✅ Desktop (600px de ancho)
- ✅ Tablet (se adapta al contenedor)
- ✅ Mobile (max-width: 95vw)
- ✅ Clientes de email (Gmail, Outlook, Apple Mail)

### Técnicas Utilizadas:
- **Tables en lugar de divs** (mejor compatibilidad)
- **Inline CSS** (no todos los clientes soportan `<style>`)
- **Fuentes seguras** (Segoe UI, Tahoma, Geneva, Verdana)
- **Sin imágenes** (puro HTML/CSS para evitar spam filters)

---

## 🔧 Variables del Template

El template utiliza las siguientes variables de EmailJS:

```javascript
{
  from_name: "Nombre del usuario",
  subject: "💬 Nuevo mensaje de soporte: [Asunto]",
  to_email: "soporte@fintrack.com",
  reply_to: "email@usuario.com",
  message: "Contenido del mensaje",
  html_content: "<html>...</html>",  // Template completo
  user_email: "email@usuario.com",
  user_name: "Nombre",
  timestamp: "27/10/2025 12:00:00"
}
```

---

## 📧 Ejemplo de Email Recibido

```
De: FinTrack Support System <noreply@emailjs.com>
Para: soporte@fintrack.com
Asunto: 💬 Nuevo mensaje de soporte: Consulta sobre reportes
Reply-To: juan.perez@example.com

[EMAIL RENDERIZADO CON HTML]

┌─────────────────────────────────────────────────┐
│  💬 Nuevo Mensaje de Soporte                    │
│  FinTrack Support System                        │
│  [Fondo gradiente violeta]                      │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ ⚡ Acción requerida: Un usuario necesita        │
│    asistencia. Responde a la brevedad.          │
│ [Fondo amarillo, borde naranja]                 │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│ 👤 Información del Usuario                      │
│ ───────────────────────────────────────────────│
│ Nombre:    Juan Pérez                          │
│ Email:     juan.perez@example.com              │
│ Fecha:     27/10/2025 12:00:00                 │
└─────────────────────────────────────────────────┘

📋 Asunto
┌─────────────────────────────────────────────────┐
│ Consulta sobre reportes de gastos              │
│ [Fondo morado claro, borde morado]             │
└─────────────────────────────────────────────────┘

💬 Mensaje
┌─────────────────────────────────────────────────┐
│ Hola equipo de FinTrack,                       │
│                                                 │
│ Tengo una consulta sobre cómo generar         │
│ reportes mensuales de gastos por categoría.   │
│                                                 │
│ ¿Podrían ayudarme con esto?                   │
│                                                 │
│ Gracias por su tiempo.                         │
│ Saludos!                                        │
└─────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────┐
│     ┌─────────────────────────────┐            │
│     │ ✉️ Responder al Usuario     │            │
│     └─────────────────────────────┘            │
│ [Botón blanco sobre fondo gradiente]           │
└─────────────────────────────────────────────────┘

Este email fue enviado automáticamente por el 
sistema de soporte de FinTrack
© 2025 FinTrack. Todos los derechos reservados.
```

---

## ✨ Mejoras vs. Versión Anterior

### Antes:
```
A message by Juan Pérez has been received.

Nombre: Juan Pérez
Email: juan@example.com
Asunto: Consulta
Mensaje: Texto plano...

Para responder: juan@example.com
```

### Ahora:
- ✅ **Diseño profesional** con colores corporativos
- ✅ **Header visualmente atractivo** con gradiente
- ✅ **Información organizada** en tarjetas
- ✅ **Banner de alerta** para llamar la atención
- ✅ **Botón de acción** pre-configurado
- ✅ **Formato de mensaje** respetando saltos de línea
- ✅ **Footer corporativo**
- ✅ **Compatible con todos los clientes de email**

---

## 🚀 Cómo Funciona

1. Usuario completa formulario en `/faq`
2. Frontend envía POST a `/api/notifications/support`
3. Nginx hace proxy a `notification-service:8088`
4. Backend recibe datos y construye HTML
5. EmailJS envía email con template `template_yst8bd2`
6. Email llega a `soporte@fintrack.com` con diseño bonito
7. Equipo de soporte puede responder con un click

---

## 📝 Configuración en EmailJS

El template `template_yst8bd2` debe estar configurado con:

**Subject Line:**
```
{{subject}}
```

**HTML Content:**
```
{{{html_content}}}
```

**Reply-To:**
```
{{reply_to}}
```

**To Email:**
```
{{to_email}}
```

> **Nota:** Usar `{{{html_content}}}` con 3 llaves para que EmailJS no escape el HTML.

---

**Implementado:** 27 de Octubre, 2025  
**Versión:** 2.0 (con diseño HTML profesional)  
**Estado:** ✅ Funcional y probado
