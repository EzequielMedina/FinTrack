# FinTrack - Plataforma de Gestión Financiera Personal

![FinTrack Logo](docs/dashboard_principal.svg)

## 📋 Descripción

FinTrack es una plataforma completa de gestión financiera personal desarrollada con arquitectura de microservicios. Permite a los usuarios gestionar sus finanzas personales, realizar transacciones, obtener reportes detallados y recibir asistencia a través de un chatbot inteligente.

## 🏗️ Arquitectura

### Tecnologías Principales
- **Backend**: Go 1.24+ (Microservicios)
- **Frontend**: Node.js 22+ (Angular/React)
- **Base de Datos**: MySQL 8.0
- **Contenedores**: Docker & Docker Compose
- **CI/CD**: GitHub Actions
- **Proxy**: Nginx

### Microservicios

| Servicio | Puerto | Descripción |
|----------|--------|-------------|
| **user-service** | 8081 | Gestión de usuarios y autenticación |
| **transaction-service** | 8082 | Procesamiento de transacciones |
| **wallet-service** | 8083 | Gestión de billeteras digitales |
| **account-service** | 8084 | Administración de cuentas bancarias |
| **notification-service** | 8085 | Sistema de notificaciones |
| **chatbot-service** | 8086 | Asistente virtual con IA |
| **exchange-service** | 8087 | Tipos de cambio y conversiones |
| **report-service** | 8088 | Generación de reportes y analytics |
| **frontend** | 4200 | Interfaz de usuario web |
| **mysql** | 3306 | Base de datos principal |

## 🚀 Inicio Rápido

### Prerrequisitos

- Docker 20.10+
- Docker Compose 2.0+
- Git
- Go 1.24+ (para desarrollo)
- Node.js 22+ (para desarrollo)

### Instalación

1. **Clonar el repositorio**
   ```bash
   git clone <repository-url>
   cd PS
   ```

2. **Configurar variables de entorno**
   ```bash
   cp config_example.json config.json
   # Editar config.json con tus configuraciones
   ```

3. **Levantar todos los servicios**
   ```bash
   docker-compose up -d
   ```

4. **Verificar el estado de los servicios**
   ```bash
   docker-compose ps
   ```

### Acceso a la Aplicación

- **Frontend**: http://localhost:4200
- **API Gateway**: http://localhost:80
- **Base de Datos**: localhost:3306

## 🔧 Desarrollo

### Estructura del Proyecto

```
PS/
├── .github/workflows/     # Pipelines CI/CD
├── backend/
│   ├── services/          # Microservicios Go
│   │   ├── user-service/
│   │   ├── transaction-service/
│   │   ├── wallet-service/
│   │   ├── account-service/
│   │   ├── notification-service/
│   │   ├── chatbot-service/
│   │   ├── exchange-service/
│   │   └── report-service/
│   └── shared/            # Código compartido
├── frontend/              # Aplicación web
├── database/              # Esquemas y migraciones
├── docker/                # Configuraciones Docker
├── docs/                  # Documentación
└── scripts/               # Scripts de utilidad
```

### Comandos de Desarrollo

```bash
# Desarrollo individual de servicios
docker-compose up mysql user-service  # Solo servicios específicos

# Logs de servicios
docker-compose logs -f user-service

# Rebuild de un servicio
docker-compose up --build user-service

# Ejecutar tests
docker-compose exec user-service go test ./...

# Acceder a la base de datos
docker-compose exec mysql mysql -u fintrack_user -p fintrack
```

## 📊 Monitoreo y Salud

Todos los servicios incluyen endpoints de health check:

```bash
# Verificar salud de servicios
curl http://localhost:8081/health  # user-service
curl http://localhost:8082/health  # transaction-service
# ... otros servicios
```

## 🔐 Configuración de Seguridad

### Variables de Entorno Requeridas

```env
# Base de datos
DB_HOST=mysql
DB_PORT=3306
DB_NAME=fintrack
DB_USER=fintrack_user
DB_PASSWORD=fintrack_password

# JWT
JWT_SECRET=your-jwt-secret-key

# APIs Externas
OPENAI_API_KEY=your-openai-api-key
EXCHANGE_API_KEY=your-exchange-api-key

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
```

## 🧪 Testing

### Ejecutar Tests

```bash
# Tests unitarios por servicio
docker-compose exec user-service go test ./...
docker-compose exec transaction-service go test ./...

# Tests de integración
docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# Tests del frontend
docker-compose exec frontend npm test
```

## 📚 Documentación

- [Análisis Completo](docs/FinTrack_Analisis_Completo.md)
- [Arquitectura Técnica](docs/FinTrack_Arquitectura_Tecnica.md)
- [Casos de Uso](docs/FinTrack_Casos_de_Uso.md)
- [Diseño de Base de Datos](docs/FinTrack_Diseno_Base_Datos.md)
- [Metodología de Testing](docs/FinTrack_Testing_Metodologia.md)
- [APIs Externas](docs/FinTrack_APIs_Externas_Gratuitas.md)

## 🚀 Despliegue

### Producción

```bash
# Build para producción
docker-compose -f docker-compose.prod.yml up -d

# Escalado de servicios
docker-compose up -d --scale user-service=3
```

### CI/CD

El proyecto incluye pipelines automatizados:

- **CI/CD Pipeline**: Build, test y deploy automático
- **Security & Quality Gates**: Análisis de seguridad y calidad de código

## 🤝 Contribución

1. Fork el proyecto
2. Crear una rama feature (`git checkout -b feature/AmazingFeature`)
3. Commit los cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## 📞 Soporte

Para soporte técnico o preguntas:

- 📧 Email: support@fintrack.com
- 📖 Wiki: [Documentación completa](docs/)
- 🐛 Issues: [GitHub Issues](../../issues)

## 🔄 Estado del Proyecto

![CI/CD Status](https://github.com/username/fintrack/workflows/CI-CD/badge.svg)
![Security Status](https://github.com/username/fintrack/workflows/Security%20&%20Quality%20Gates/badge.svg)

---

**FinTrack** - Gestiona tus finanzas con inteligencia 💰✨