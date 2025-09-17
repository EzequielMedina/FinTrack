# Multi-stage build for Go microservice
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/api/main.go

# Production stage
FROM alpine:latest

# Install ca-certificates and create app user
RUN apk --no-cache add ca-certificates tzdata wget && \
    addgroup -g 1001 appgroup && \
    adduser -D -u 1001 -G appgroup appuser

# Set working directory
WORKDIR /app

# Copy the binary
COPY --from=builder /app/main .

# Change ownership
RUN chown appuser:appgroup /app/main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV ENVIRONMENT=production

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]