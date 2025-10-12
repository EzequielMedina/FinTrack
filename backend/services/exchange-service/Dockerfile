# Multi-stage build for Exchange Service
FROM golang:1.21-alpine AS builder

# Install ca-certificates for HTTPS requests to external APIs
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o exchange-service ./cmd/main.go

# Production stage - using alpine for debugging capabilities
FROM alpine:latest

# Install ca-certificates for HTTPS requests to DolarAPI
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/exchange-service .

# Create app user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Change ownership of the binary
RUN chown appuser:appgroup exchange-service

# Use non-root user
USER appuser:appgroup

# Expose port (8087 para exchange-service)
EXPOSE 8087

# Health check for exchange service
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8087/health || exit 1

# Run the binary
CMD ["./exchange-service"]