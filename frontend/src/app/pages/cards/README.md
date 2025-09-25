# Gestión de Tarjetas - FinTrack

## 📋 Descripción

Módulo de gestión de tarjetas de crédito y débito para FinTrack. Implementa funcionalidades completas de CRUD, validaciones de seguridad, encriptación de datos sensibles y detección automática de marcas de tarjetas.

## 🏗️ Arquitectura

### Componentes Principales

- **CardsComponent**: Componente principal que actúa como contenedor
- **CardListComponent**: Lista y visualización de tarjetas con filtros
- **CardFormComponent**: Formulario para crear/editar tarjetas
- **CardDetailComponent**: Vista detallada de una tarjeta específica

### Servicios

- **CardService**: Manejo de operaciones CRUD y validaciones
- **EncryptionService**: Encriptación de datos sensibles
- **UserService**: Integración con el sistema de usuarios

### Modelos

- **Card**: Entidad principal de tarjeta
- **CardType**: Enum para tipos (crédito/débito)
- **CardBrand**: Enum para marcas (Visa, Mastercard, etc.)
- **CardStatus**: Enum para estados (activa, bloqueada, etc.)

## 🔒 Seguridad

### Encriptación de Datos
- Números de tarjeta encriptados antes del envío al backend
- CVV encriptado usando algoritmos seguros
- Implementación con Web Crypto API en producción
- Fallback a encriptación mock en desarrollo

### Validaciones
- Algoritmo de Luhn para validación de números de tarjeta
- Validación de CVV según la marca de tarjeta
- Validación de fechas de expiración
- Sanitización de inputs

### Detección de Marcas
- Detección automática basada en patrones de números
- Validación de longitud según la marca
- Iconos específicos para cada marca

## 🎨 UI/UX

### Características
- Diseño responsive con Angular Material
- Números de tarjeta enmascarados para seguridad
- Indicadores visuales de estado
- Formularios reactivos con validación en tiempo real
- Animaciones suaves y feedback visual

### Estados Visuales
- Tarjetas predeterminadas destacadas
- Indicadores de expiración
- Estados de carga y guardado
- Mensajes de error contextuales

## 🔄 Funcionalidades

### CRUD Completo
- ✅ Crear nuevas tarjetas
- ✅ Listar tarjetas por usuario/cuenta
- ✅ Editar información de tarjetas
- ✅ Eliminar tarjetas
- ✅ Establecer tarjeta predeterminada
- ✅ Bloquear/desbloquear tarjetas

### Filtros y Búsqueda
- Filtro por tipo (crédito/débito)
- Filtro por estado (activa/inactiva)
- Agrupación por cuenta asociada
- Búsqueda por nombre personalizado

### Validaciones en Tiempo Real
- Detección automática de marca mientras se escribe
- Validación de formato de número
- Validación de CVV según la marca
- Validación de fecha de expiración

## 📱 Responsive Design

### Breakpoints
- **Desktop**: > 768px - Grid de tarjetas con múltiples columnas
- **Tablet**: 481px - 768px - Grid adaptativo
- **Mobile**: ≤ 480px - Lista vertical con optimizaciones táctiles

### Optimizaciones Móviles
- Formularios de pantalla completa en móviles
- Botones de acción más grandes
- Navegación por tabs optimizada
- Inputs numéricos nativos para CVV

## 🧪 Testing

### Tests Unitarios
```bash
# Ejecutar tests de servicios
ng test --include="**/*card*.service.spec.ts"

# Ejecutar tests de componentes
ng test --include="**/*card*.component.spec.ts"
```

### Tests de Integración
```bash
# Tests E2E para flujo completo
ng e2e --spec="**/card-management.e2e-spec.ts"
```

## 🚀 Configuración y Uso

### Instalación
```bash
# El módulo está integrado en la aplicación principal
npm install
ng serve
```

### Rutas
- `/cards` - Página principal de gestión de tarjetas
- `/cards/new` - Formulario de nueva tarjeta (modal)
- `/cards/:id/edit` - Edición de tarjeta (modal)

### Variables de Entorno
```typescript
// environment.ts
export const environment = {
  accountServiceUrl: '/api/accounts',  // Endpoint del microservicio
  encryptionEnabled: true,             // Habilitar encriptación
  mockEncryption: false               // Usar encriptación real vs mock
};
```

## 🔧 Configuración del Backend

### Endpoints Esperados
```
POST   /api/accounts/{accountId}/cards          - Crear tarjeta
GET    /api/accounts/{accountId}/cards          - Listar tarjetas por cuenta
GET    /api/accounts/user/{userId}/cards        - Listar tarjetas por usuario
GET    /api/accounts/{accountId}/cards/{cardId} - Obtener tarjeta específica
PUT    /api/accounts/{accountId}/cards/{cardId} - Actualizar tarjeta
DELETE /api/accounts/{accountId}/cards/{cardId} - Eliminar tarjeta
PUT    /api/accounts/{accountId}/cards/{cardId}/set-default - Establecer como predeterminada
PUT    /api/accounts/{accountId}/cards/{cardId}/block       - Bloquear tarjeta
PUT    /api/accounts/{accountId}/cards/{cardId}/unblock     - Desbloquear tarjeta
```

### Formato de Datos
```typescript
interface CreateCardRequest {
  accountId: string;
  cardType: 'credit' | 'debit';
  cardNumber: string;        // Encriptado
  holderName: string;
  expirationMonth: number;
  expirationYear: number;
  cvv: string;              // Encriptado
  nickname?: string;
}
```

## 📚 Casos de Uso

### Usuario Estándar
1. **Agregar primera tarjeta**: Registro inicial con validaciones
2. **Gestionar múltiples tarjetas**: Organización por cuentas y tipos
3. **Establecer predeterminada**: Configuración de tarjeta principal
4. **Bloquear en emergencia**: Bloqueo temporal de seguridad

### Usuario Avanzado
1. **Múltiples cuentas**: Gestión de tarjetas por cuenta bancaria
2. **Nombres personalizados**: Organización con aliases descriptivos
3. **Filtros avanzados**: Búsqueda y categorización
4. **Exportación de datos**: Descarga de información de tarjetas

## 🔍 Troubleshooting

### Problemas Comunes

**Error de encriptación**
```
Error al procesar datos de la tarjeta
```
- Verificar Web Crypto API disponible
- Revisar configuración de HTTPS en producción

**Validación fallida**
```
El número de tarjeta no es válido
```
- Verificar algoritmo de Luhn
- Comprobar patrones de detección de marca

**Problemas de conectividad**
```
Error al cargar las tarjetas
```
- Verificar endpoint del account-service
- Revisar configuración de CORS

## 📈 Métricas y Monitoreo

### KPIs Implementados
- Tiempo de carga de lista de tarjetas
- Tasa de éxito en validaciones
- Errores de encriptación/desencriptación
- Uso de diferentes marcas de tarjetas

### Logs de Auditoria
- Creación/edición/eliminación de tarjetas
- Cambios de estado (bloqueo/desbloqueo)
- Establecimiento de tarjetas predeterminadas
- Errores de validación y seguridad

## 🛡️ Consideraciones de Seguridad

### Datos Sensibles
- Nunca almacenar números completos en localStorage
- Limpiar datos de memoria después del uso
- Usar HTTPS obligatorio en producción
- Implementar timeout de sesión

### Compliance
- Cumplimiento PCI DSS para manejo de tarjetas
- Encriptación end-to-end
- Logs de auditoría para compliance
- Validación de entrada estricta

## 🔮 Roadmap

### Próximas Funcionalidades
- [ ] Integración con procesadores de pago
- [ ] Tokenización de tarjetas
- [ ] Verificación biométrica
- [ ] Integración con wallet digital
- [ ] Análisis de gastos por tarjeta
- [ ] Alertas de fraude en tiempo real

---

**Versión**: 1.0.0  
**Última Actualización**: Septiembre 2025  
**Mantenido por**: FinTrack Development Team