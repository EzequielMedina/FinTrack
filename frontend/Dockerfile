# Multi-stage build for Angular frontend
FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
# Install all dependencies (dev dependencies required for build)
RUN npm ci

# Copy source code
COPY . .

# Build the application
# Angular 17+ uses configuration flag
RUN npm run build:prod

# Production stage
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built application
# Copy Angular browser build output
COPY --from=builder /app/dist/fintrack/browser /usr/share/nginx/html

RUN apk add --no-cache curl

# Add non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S angular -u 1001

# Change ownership
RUN chown -R angular:nodejs /usr/share/nginx/html
RUN chown -R angular:nodejs /var/cache/nginx
RUN chown -R angular:nodejs /var/log/nginx
RUN chown -R angular:nodejs /etc/nginx/conf.d
RUN touch /var/run/nginx.pid
RUN chown -R angular:nodejs /var/run/nginx.pid

# Switch to non-root user
USER angular

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]