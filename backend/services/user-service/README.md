# FinTrack - User Service

Microservicio de autenticación y gestión de usuarios para FinTrack. Provee endpoints de registro y login, emite JWTs, y valida tokens con middleware.

## Endpoints

- GET `/health` — estado del servicio
- POST `/api/auth/register` — registro de usuario
  - body: { email, password, firstName, lastName }
- POST `/api/auth/login` — login de usuario
  - body: { email, password }
- GET `/api/me` — protegido con JWT (Authorization: Bearer <token>)

Ejemplos curl:

- Registro
  curl -X POST http://localhost:8081/api/auth/register -H "Content-Type: application/json" -d '{"email":"user@example.com","password":"StrongP@ssw0rd","firstName":"Ada","lastName":"Lovelace"}'

- Login
  curl -X POST http://localhost:8081/api/auth/login -H "Content-Type: application/json" -d '{"email":"user@example.com","password":"StrongP@ssw0rd"}'

- Protected
  curl -H "Authorization: Bearer <accessToken>" http://localhost:8081/api/me

## Variables de entorno

- DB_HOST (default: localhost)
- DB_PORT (default: 3306)
- DB_NAME (default: fintrack)
- DB_USER (default: fintrack_user)
- DB_PASSWORD (default: fintrack_password)
- JWT_SECRET (requerido en producción)
- JWT_EXPIRY (default: 24h)
- JWT_REFRESH_EXPIRY (default: 168h)
- PORT (default: 8080)
- LOG_LEVEL (default: info)

## Ejecutar con Docker Compose

- Desde la raíz del repo:
  docker-compose up --build mysql user-service

- El servicio expone 8080 dentro del contenedor y se publica en 8081 en el host. Health: http://localhost:8081/health

## Desarrollo local (sin Docker)

- Iniciar MySQL local o dockerizado y exportar las variables.
- Desde `backend/services/user-service`:
  go mod tidy
  go run ./cmd/api/main.go

## Postman

Importar la colección en `docs/postman/FinTrack_UserService.postman_collection.json` y ajustar variable `baseUrl` si es necesario.