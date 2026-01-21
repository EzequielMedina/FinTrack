# Instrucciones para Crear PowerPoint de FinTrack

## 📋 Opción 1: Usar la Presentación HTML (Recomendado)

He creado una presentación HTML interactiva que puedes usar directamente:

1. **Abrir la presentación:**
   - Navega a `docs/FinTrack_Presentacion_Marketing.html`
   - Ábrela con cualquier navegador (Chrome, Edge, Firefox)
   - Usa las flechas del teclado o los botones para navegar

2. **Convertir a PowerPoint:**
   - Abre la presentación HTML en tu navegador
   - Presiona `F11` para modo pantalla completa
   - Usa `Print Screen` o herramientas de captura para cada slide
   - O usa herramientas online como:
     - [HTML to PPT Converter](https://www.zamzar.com/convert/html-to-pptx/)
     - [CloudConvert](https://cloudconvert.com/html-to-pptx)

## 📋 Opción 2: Crear PowerPoint Manualmente

### Pasos:

1. **Abrir PowerPoint** y crear nueva presentación en blanco

2. **Configurar tamaño de diapositiva:**
   - Diseño → Tamaño de diapositiva → 16:9 (1920x1080)

3. **Crear las siguientes diapositivas:**

#### **Slide 1: Portada**
- Título: **FinTrack**
- Subtítulo: **Tu Futuro Financiero, en Tus Manos**
- Texto: "Plataforma Integral de Gestión Financiera Personal"
- Color de fondo: Azul oscuro (#0f172a)
- Color de texto: Azul claro (#2563eb)

#### **Slide 2: El Problema**
- Título: **El Problema**
- Subtítulo: **¿En qué gasté mi dinero este mes?**
- Lista:
  - 📊 Información dispersa en múltiples apps
  - 💳 Deudas que se acumulan sin control
  - 💰 Objetivos financieros que nunca se cumplen
  - 🔄 Falta de visibilidad sobre gastos e ingresos

#### **Slide 3: La Solución**
- Título: **La Solución: FinTrack**
- Subtítulo: **Control Total de tus Finanzas**
- Imagen: `docs/screenshots/01-login.png`
- Texto: "Una plataforma integral que centraliza todas tus finanzas en un solo lugar"

#### **Slide 4: Dashboard**
- Título: **Dashboard Principal**
- Subtítulo: **Vista Consolidada**
- Imagen: `docs/screenshots/02-dashboard.png`
- Texto: "Resumen completo de tus cuentas, tarjetas y transacciones en tiempo real"

#### **Slide 5: Cuentas**
- Título: **Gestión de Cuentas**
- Subtítulo: **Todo en un Solo Lugar**
- Imagen: `docs/screenshots/03-accounts.png`
- Texto: "Administra todas tus cuentas bancarias y billeteras virtuales"

#### **Slide 6: Tarjetas**
- Título: **Gestión de Tarjetas**
- Subtítulo: **Control Total**
- Imagen: `docs/screenshots/04-cards.png`
- Texto: "Gestiona tus tarjetas de crédito y débito con seguimiento detallado"

#### **Slide 7: Transacciones**
- Título: **Transacciones**
- Subtítulo: **Seguimiento Completo**
- Imagen: `docs/screenshots/05-transactions.png`
- Texto: "Historial detallado de todas tus transacciones con filtros avanzados"

#### **Slide 8: Reportes**
- Título: **Reportes Visuales**
- Subtítulo: **Análisis Detallado**
- Imagen: `docs/screenshots/06-reports.png`
- Texto: "Gráficos y análisis de tus gastos e ingresos con exportación a PDF"

#### **Slide 9: Chatbot**
- Título: **Asistente IA**
- Subtítulo: **Chatbot Inteligente**
- Imagen: `docs/screenshots/07-chatbot.png`
- Texto: "Consulta sobre tus finanzas y recibe recomendaciones personalizadas"

#### **Slide 10: Características**
- Título: **Características Principales**
- Subtítulo: **Todo lo que Necesitas**
- Lista en 2 columnas:
  - ✓ Control Total de cuentas y tarjetas
  - ✓ Asistente IA para consultas financieras
  - ✓ Reportes visuales detallados
  - ✓ Metas y seguimiento de ahorro
  - ✓ Notificaciones inteligentes
  - ✓ Soporte multi-divisa

#### **Slide 11: Call to Action**
- Título: **Toma el Control**
- Subtítulo: **Regístrate Gratis**
- Texto destacado:
  - "FinTrack te ayuda a alcanzar tus objetivos económicos con tecnología accesible e inteligente"
  - "www.fintrack.com"
- Fondo: Gradiente azul (#2563eb a #1d4ed8)

### **Consejos de Diseño:**

1. **Colores:**
   - Primario: #2563eb (Azul FinTrack)
   - Secundario: #1e293b (Azul oscuro)
   - Texto: Blanco o #e2e8f0 (Gris claro)

2. **Tipografía:**
   - Títulos: Poppins o Arial Bold, 44-64pt
   - Subtítulos: Inter o Calibri, 28-36pt
   - Cuerpo: Inter o Calibri, 18-24pt

3. **Imágenes:**
   - Insertar desde: Insertar → Imágenes → Este dispositivo
   - Seleccionar archivos de `docs/screenshots/`
   - Ajustar tamaño manteniendo proporción
   - Agregar borde sutil azul (#2563eb)

4. **Animaciones (Opcional):**
   - Entrada: Desvanecer o Aparecer
   - Transiciones: Desvanecer o Empujar

## 📋 Opción 3: Usar el Script PowerShell

Si tienes PowerPoint instalado, puedes ejecutar:

```powershell
cd c:\Facultad\Alumno\PS
powershell -ExecutionPolicy Bypass -File ".\docs\crear_presentacion.ps1"
```

Esto creará automáticamente `docs/FinTrack_Presentacion_Marketing.pptx`

## 📁 Archivos Disponibles

- **Screenshots:** `docs/screenshots/`
  - 01-login-page.png
  - 02-dashboard.png
  - 03-accounts.png
  - 04-cards.png
  - 05-transactions.png
  - 06-reports.png
  - 07-chatbot.png

- **Presentación HTML:** `docs/FinTrack_Presentacion_Marketing.html`

- **Script PowerShell:** `docs/crear_presentacion.ps1`

## ✅ Checklist

- [ ] Screenshots capturados y guardados
- [ ] Presentación HTML creada
- [ ] PowerPoint creado (manual o automático)
- [ ] Revisar diseño y colores
- [ ] Agregar animaciones (opcional)
- [ ] Exportar en formato final (.pptx)

---

**Nota:** La presentación HTML es completamente funcional y puede usarse directamente para presentaciones en pantalla o convertirla a PowerPoint usando herramientas online.


