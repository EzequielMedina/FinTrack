# FinTrack Frontend

## 📋 Descripción

Interfaz de usuario web para la plataforma FinTrack, desarrollada con Angular y optimizada para producción con Nginx.

## 🛠️ Tecnologías

- **Framework**: Angular 17+
- **Runtime**: Node.js 22+
- **Servidor Web**: Nginx (Alpine)
- **Contenedor**: Docker multi-stage
- **Estilos**: CSS3, SCSS
- **Build Tool**: Angular CLI

## 🏗️ Arquitectura

### Build Multi-Stage

1. **Stage 1 (Builder)**: 
   - Base: `node:22-alpine`
   - Instala dependencias
   - Compila la aplicación Angular

2. **Stage 2 (Production)**:
   - Base: `nginx:alpine`
   - Sirve archivos estáticos
   - Configuración optimizada

## 🚀 Desarrollo Local

### Prerrequisitos

- Node.js 22+
- npm 10+
- Angular CLI 17+

### Instalación

```bash
# Navegar al directorio frontend
cd frontend

# Instalar dependencias
npm install

# Instalar Angular CLI globalmente (si no está instalado)
npm install -g @angular/cli
```

### Comandos de Desarrollo

```bash
# Servidor de desarrollo
npm start
# o
ng serve

# Acceder a la aplicación
# http://localhost:4200

# Build de desarrollo
npm run build

# Build de producción
npm run build --prod

# Tests unitarios
npm test

# Tests e2e
npm run e2e

# Linting
npm run lint

# Formateo de código
npm run format
```

## 🐳 Docker

### Build Local

```bash
# Build de la imagen
docker build -t fintrack-frontend .

# Ejecutar contenedor
docker run -p 4200:80 fintrack-frontend

# Acceder a la aplicación
# http://localhost:4200
```

### Docker Compose

```bash
# Desde el directorio raíz del proyecto
docker-compose up frontend

# Con rebuild
docker-compose up --build frontend

# Solo frontend y dependencias
docker-compose up frontend user-service transaction-service
```

## 📁 Estructura del Proyecto

```
frontend/
├── src/
│   ├── app/                 # Componentes Angular
│   │   ├── components/      # Componentes reutilizables
│   │   ├── pages/          # Páginas principales
│   │   ├── services/       # Servicios Angular
│   │   ├── guards/         # Guards de autenticación
│   │   ├── interceptors/   # HTTP interceptors
│   │   ├── models/         # Interfaces y modelos
│   │   └── shared/         # Módulos compartidos
│   ├── assets/             # Recursos estáticos
│   ├── environments/       # Configuraciones de entorno
│   └── styles/             # Estilos globales
├── Dockerfile              # Configuración Docker
├── nginx.conf              # Configuración Nginx
├── package.json            # Dependencias npm
├── angular.json            # Configuración Angular
├── tsconfig.json           # Configuración TypeScript
└── README.md               # Este archivo
```

## 🔧 Configuración

### Variables de Entorno

```typescript
// src/environments/environment.ts
export const environment = {
  production: false,
  apiUrl: 'http://localhost:8080/api',
  wsUrl: 'ws://localhost:8080/ws'
};

// src/environments/environment.prod.ts
export const environment = {
  production: true,
  apiUrl: '/api',
  wsUrl: '/ws'
};
```

### Nginx Configuration

El archivo `nginx.conf` incluye:

- Compresión Gzip
- Caché de archivos estáticos
- Proxy para APIs
- Configuración de seguridad
- Health checks

## 🔗 Integración con Backend

### Endpoints de API

```typescript
// Servicios disponibles
const API_ENDPOINTS = {
  users: '/api/users',
  transactions: '/api/transactions',
  wallets: '/api/wallets',
  accounts: '/api/accounts',
  notifications: '/api/notifications',
  chatbot: '/api/chatbot',
  exchange: '/api/exchange',
  reports: '/api/reports'
};
```

### Autenticación

```typescript
// JWT Token management
// Interceptor automático para headers Authorization
// Guards para rutas protegidas
// Refresh token automático
```

## 🧪 Testing

### Tests Unitarios

```bash
# Ejecutar tests
npm test

# Tests con coverage
npm run test:coverage

# Tests en modo watch
npm run test:watch
```

### Tests E2E

```bash
# Cypress tests
npm run e2e

# Cypress en modo interactivo
npm run e2e:open
```

## 📊 Performance

### Optimizaciones Incluidas

- **Lazy Loading**: Módulos cargados bajo demanda
- **Tree Shaking**: Eliminación de código no utilizado
- **Minificación**: CSS y JS minificados
- **Compresión Gzip**: Reducción del tamaño de transferencia
- **Service Workers**: Caché offline (PWA ready)
- **Bundle Splitting**: Separación de vendor y app bundles

### Métricas de Build

```bash
# Análisis del bundle
npm run build:analyze

# Lighthouse audit
npm run lighthouse
```

## 🔐 Seguridad

### Medidas Implementadas

- **CSP Headers**: Content Security Policy
- **HTTPS Redirect**: Redirección automática a HTTPS
- **XSS Protection**: Protección contra Cross-Site Scripting
- **CSRF Protection**: Tokens CSRF en formularios
- **Sanitización**: Sanitización automática de inputs

## 🚀 Despliegue

### Build de Producción

```bash
# Build optimizado
npm run build --prod

# Verificar archivos generados
ls -la dist/fintrack/
```

### Docker Production

```bash
# Build para producción
docker build -t fintrack-frontend:prod .

# Tag para registry
docker tag fintrack-frontend:prod registry.com/fintrack-frontend:latest

# Push al registry
docker push registry.com/fintrack-frontend:latest
```

## 🔍 Debugging

### Logs de Desarrollo

```bash
# Logs del servidor de desarrollo
ng serve --verbose

# Logs de Docker
docker-compose logs -f frontend

# Logs de Nginx
docker-compose exec frontend tail -f /var/log/nginx/access.log
docker-compose exec frontend tail -f /var/log/nginx/error.log
```

### Herramientas de Debug

- **Angular DevTools**: Extensión de Chrome
- **Redux DevTools**: Para manejo de estado
- **Network Tab**: Monitoreo de requests
- **Console Logs**: Logs estructurados

## 📚 Recursos

- [Angular Documentation](https://angular.io/docs)
- [Angular CLI](https://cli.angular.io/)
- [Nginx Documentation](https://nginx.org/en/docs/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

## 🤝 Contribución

1. Seguir las convenciones de Angular
2. Usar TypeScript estricto
3. Escribir tests para nuevas funcionalidades
4. Seguir el style guide de Angular
5. Documentar componentes complejos

---

**Frontend FinTrack** - Interfaz moderna y responsiva 🎨✨