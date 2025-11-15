FROM golang:1.24-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git ca-certificates tzdata

# Configurar directorio de trabajo
WORKDIR /app

# Copiar go mod y sum
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-service cmd/main.go

# Etapa final
FROM alpine:latest

# Instalar certificados CA y timezone data
RUN apk --no-cache add ca-certificates tzdata

# Crear usuario no root
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Configurar directorio de trabajo
WORKDIR /root/

# Copiar binario desde builder
COPY --from=builder /app/notification-service .

# Copiar templates si existen
COPY --from=builder /app/templates ./templates

# Cambiar propietario de archivos
RUN chown -R appuser:appgroup /root/

# Cambiar a usuario no root
USER appuser

# Exponer puerto
EXPOSE 8088

# Configurar health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8088/health || exit 1

# Comando por defecto
CMD ["./notification-service"]