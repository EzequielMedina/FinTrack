# Nginx API Gateway - FinTrack

## 📋 Descripción

API Gateway basado en Nginx que actúa como punto de entrada único para todos los microservicios de la plataforma FinTrack. Proporciona enrutamiento, balanceeo de carga, rate limiting, CORS, seguridad y monitoreo centralizado.

## 🛠️ Tecnologías

- **Servidor Web**: Nginx Alpine
- **Load Balancing**: Least Connection
- **Rate Limiting**: Token Bucket Algorithm
- **SSL/TLS**: Soporte completo
- **WebSockets**: Soporte para tiempo real
- **Monitoring**: Health checks y métricas
- **Contenedor**: Docker multi-stage

## 🏗️ Arquitectura

### Estructura del Proyecto

```
nginx/
├── Dockerfile               # Configuración Docker
├── nginx.conf               # Configuración principal
├── api-gateway.conf         # Configuración del gateway
├── ssl/                     # Certificados SSL (opcional)
│   ├── cert.pem
│   └── key.pem
└── README.md                # Este archivo
```

### Flujo de Requests

```
Cliente → Nginx Gateway → Load Balancer → Microservicio
   ↓           ↓              ↓              ↓
 CORS      Rate Limit    Health Check   Response
   ↓           ↓              ↓              ↓
Headers    Security      Monitoring     Logging
```

## 🚀 Configuración

### Variables de Entorno

```env
# Configuración básica
NGINX_PORT=80
NGINX_WORKER_PROCESSES=auto
NGINX_WORKER_CONNECTIONS=1024

# Rate Limiting
RATE_LIMIT_API=10r/s
RATE_LIMIT_LOGIN=5r/m
RATE_LIMIT_BURST_API=20
RATE_LIMIT_BURST_LOGIN=5

# Timeouts
PROXY_CONNECT_TIMEOUT=30s
PROXY_SEND_TIMEOUT=30s
PROXY_READ_TIMEOUT=30s
KEEPALIVE_TIMEOUT=65s

# SSL/TLS (opcional)
SSL_CERTIFICATE=/etc/nginx/ssl/cert.pem
SSL_CERTIFICATE_KEY=/etc/nginx/ssl/key.pem
SSL_PROTOCOLS="TLSv1.2 TLSv1.3"

# Logging
LOG_LEVEL=warn
ACCESS_LOG_FORMAT=main
```

## 🐳 Docker

### Construcción y Ejecución

```bash
# Build de la imagen
docker build -t fintrack-nginx-gateway .

# Ejecutar contenedor
docker run -d \
  --name fintrack-gateway \
  -p 80:80 \
  -p 443:443 \
  --network fintrack-network \
  fintrack-nginx-gateway

# Docker Compose
docker-compose up nginx-gateway

# Con SSL
docker run -d \
  --name fintrack-gateway \
  -p 80:80 \
  -p 443:443 \
  -v ./ssl:/etc/nginx/ssl:ro \
  fintrack-nginx-gateway
```

### Health Check

```bash
# Verificar estado del gateway
curl http://localhost/health

# Verificar desde Docker
docker exec fintrack-gateway curl -f http://localhost/health
```

## 🌐 Enrutamiento de APIs

### Microservicios Configurados

| Servicio | Ruta | Puerto | Rate Limit | Timeout |
|----------|------|--------|------------|----------|
| User Service | `/api/users` | 8080 | 10r/s | 30s |
| Transaction Service | `/api/transactions` | 8080 | 10r/s | 30s |
| Wallet Service | `/api/wallets` | 8080 | 10r/s | 30s |
| Account Service | `/api/accounts` | 8080 | 10r/s | 30s |
| Notification Service | `/api/notifications` | 8080 | 10r/s | 30s |
| Chatbot Service | `/api/chatbot` | 8080 | 10r/s | 60s |
| Report Service | `/api/reports` | 8088 | 5r/s | 60s |
| Exchange Service | `/api/exchange` | 8087 | 10r/s | 30s |

### Rutas Especiales

```nginx
# Autenticación (rate limiting estricto)
location ~ ^/api/(auth|login|register) {
    limit_req zone=login burst=5 nodelay;
    proxy_pass http://user_service;
}

# WebSockets para tiempo real
location /ws {
    proxy_pass http://notification_service;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}

# Health check
location /health {
    return 200 "API Gateway is healthy\n";
}
```

## ⚖️ Load Balancing

### Configuración de Upstream

```nginx
upstream user_service {
    least_conn;                              # Algoritmo de balanceeo
    server user-service:8080 max_fails=3 fail_timeout=30s;
    server user-service-2:8080 max_fails=3 fail_timeout=30s;  # Opcional
    keepalive 32;                           # Pool de conexiones
}

upstream transaction_service {
    least_conn;
    server transaction-service:8080 max_fails=3 fail_timeout=30s;
    keepalive 32;
}
```

### Algoritmos Disponibles

- **Round Robin** (default): Distribución secuencial
- **Least Connections**: Menor número de conexiones activas
- **IP Hash**: Basado en IP del cliente
- **Weighted**: Con pesos específicos

```nginx
# Ejemplo con pesos
upstream weighted_service {
    server service-1:8080 weight=3;
    server service-2:8080 weight=1;
}
```

## 🛡️ Seguridad

### Headers de Seguridad

```nginx
# Headers aplicados automáticamente
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header X-Request-ID $request_id always;
```

### CORS Configuration

```nginx
# CORS headers
add_header Access-Control-Allow-Origin "*" always;
add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS" always;
add_header Access-Control-Allow-Headers "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range,Authorization" always;

# Preflight requests
if ($request_method = 'OPTIONS') {
    add_header Access-Control-Max-Age 1728000;
    add_header Content-Length 0;
    return 204;
}
```

### Rate Limiting

```nginx
# Zonas de rate limiting
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;

# Aplicación en rutas
location /api/users {
    limit_req zone=api burst=20 nodelay;
    # ... resto de configuración
}

location ~ ^/api/(auth|login|register) {
    limit_req zone=login burst=5 nodelay;
    # ... resto de configuración
}
```

## 📊 Monitoreo y Logging

### Formato de Logs

```nginx
log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                '$status $body_bytes_sent "$http_referer" '
                '"$http_user_agent" "$http_x_forwarded_for" '
                'rt=$request_time uct="$upstream_connect_time" '
                'uht="$upstream_header_time" urt="$upstream_response_time"';
```

### Métricas Importantes

- **Request Time**: Tiempo total de request
- **Upstream Connect Time**: Tiempo de conexión al backend
- **Upstream Response Time**: Tiempo de respuesta del backend
- **Status Codes**: Códigos de respuesta HTTP
- **Bytes Sent**: Bytes enviados al cliente

### Comandos de Monitoreo

```bash
# Ver logs en tiempo real
docker logs -f fintrack-gateway

# Logs de acceso
docker exec fintrack-gateway tail -f /var/log/nginx/access.log

# Logs de errores
docker exec fintrack-gateway tail -f /var/log/nginx/error.log

# Estadísticas de conexiones
docker exec fintrack-gateway nginx -T | grep upstream

# Verificar configuración
docker exec fintrack-gateway nginx -t
```

## 🔧 Configuración Avanzada

### SSL/TLS Setup

```nginx
server {
    listen 443 ssl http2;
    server_name api.fintrack.com;
    
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000" always;
    
    # ... resto de configuración
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name api.fintrack.com;
    return 301 https://$server_name$request_uri;
}
```

### Caching Configuration

```nginx
# Cache para contenido estático
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}

# Cache para APIs (opcional)
location /api/exchange/rates {
    proxy_cache api_cache;
    proxy_cache_valid 200 5m;
    proxy_cache_key $scheme$proxy_host$request_uri;
    # ... resto de configuración
}
```

### Compression

```nginx
# Gzip compression
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_comp_level 6;
gzip_types
    application/json
    application/javascript
    application/xml
    text/css
    text/javascript
    text/plain
    text/xml;
```

## 🚀 Performance Tuning

### Worker Configuration

```nginx
# nginx.conf
worker_processes auto;                    # CPU cores
worker_connections 1024;                  # Connections per worker
worker_rlimit_nofile 2048;               # File descriptors

events {
    use epoll;                           # Linux: epoll, BSD: kqueue
    multi_accept on;                     # Accept multiple connections
}
```

### Buffer Optimization

```nginx
# Client buffers
client_body_buffer_size 128k;
client_header_buffer_size 1k;
large_client_header_buffers 4 4k;
client_max_body_size 10M;

# Proxy buffers
proxy_buffering on;
proxy_buffer_size 4k;
proxy_buffers 8 4k;
proxy_busy_buffers_size 8k;
```

### Keepalive Optimization

```nginx
# Client keepalive
keepalive_timeout 65;
keepalive_requests 100;

# Upstream keepalive
upstream backend {
    server backend1:8080;
    keepalive 32;                        # Pool size
    keepalive_requests 100;              # Requests per connection
    keepalive_timeout 60s;               # Connection timeout
}
```

## 🔍 Troubleshooting

### Problemas Comunes

#### 502 Bad Gateway
```bash
# Verificar que los servicios backend estén corriendo
docker ps | grep -E "(user-service|transaction-service)"

# Verificar conectividad de red
docker exec fintrack-gateway nslookup user-service

# Verificar logs
docker logs fintrack-gateway | grep "connect() failed"
```

#### Rate Limiting Issues
```bash
# Verificar configuración de rate limiting
docker exec fintrack-gateway nginx -T | grep limit_req

# Monitorear requests bloqueados
docker logs fintrack-gateway | grep "limiting requests"
```

#### SSL Certificate Issues
```bash
# Verificar certificados
docker exec fintrack-gateway openssl x509 -in /etc/nginx/ssl/cert.pem -text -noout

# Verificar configuración SSL
docker exec fintrack-gateway nginx -T | grep ssl
```

### Comandos de Diagnóstico

```bash
# Test de configuración
docker exec fintrack-gateway nginx -t

# Reload de configuración
docker exec fintrack-gateway nginx -s reload

# Verificar procesos
docker exec fintrack-gateway ps aux

# Verificar puertos
docker exec fintrack-gateway netstat -tlnp

# Test de conectividad a backends
docker exec fintrack-gateway curl -I http://user-service:8080/health
```

## 📈 Métricas y Alertas

### Métricas Clave

- **Request Rate**: Requests por segundo
- **Response Time**: Tiempo de respuesta promedio
- **Error Rate**: Porcentaje de errores 4xx/5xx
- **Upstream Health**: Estado de servicios backend
- **Connection Pool**: Uso de conexiones keepalive
- **Rate Limit Hits**: Requests bloqueados por rate limiting

### Configuración de Alertas

```yaml
# Ejemplo para Prometheus/Grafana
alerts:
  - name: nginx_high_error_rate
    condition: rate(nginx_http_requests_total{status=~"5.."}[5m]) > 0.1
    message: "High error rate detected in Nginx Gateway"
    
  - name: nginx_high_response_time
    condition: nginx_http_request_duration_seconds{quantile="0.95"} > 1
    message: "High response time detected in Nginx Gateway"
```

## 🔄 Deployment

### Rolling Updates

```bash
# Update con zero downtime
docker build -t fintrack-nginx-gateway:v2 .
docker stop fintrack-gateway-old
docker run -d --name fintrack-gateway-new fintrack-nginx-gateway:v2
docker rm fintrack-gateway-old
```

### Blue-Green Deployment

```bash
# Preparar nueva versión
docker run -d --name fintrack-gateway-green fintrack-nginx-gateway:v2

# Verificar health
curl http://gateway-green/health

# Switch traffic
docker stop fintrack-gateway-blue
docker start fintrack-gateway-green
```

## 📚 Referencias

- [Nginx Documentation](https://nginx.org/en/docs/)
- [Nginx Load Balancing](https://docs.nginx.com/nginx/admin-guide/load-balancer/)
- [Nginx Rate Limiting](https://www.nginx.com/blog/rate-limiting-nginx/)
- [Nginx Security](https://nginx.org/en/docs/http/ngx_http_ssl_module.html)
- [Docker Nginx](https://hub.docker.com/_/nginx)

---

**Nginx API Gateway** - El guardián de FinTrack 🚪🛡️