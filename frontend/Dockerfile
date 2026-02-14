# Production Dockerfile for React Frontend
# Multi-stage build for optimized bundle and minimal runtime

# Stage 1: Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY src/package*.json ./

# Install dependencies with clean cache
RUN npm ci --only=production && \
    npm cache clean --force

# Copy source code
COPY src/ ./

# Build the application with production optimizations
ENV NODE_ENV=production
RUN npm run build

# Stage 2: Runtime stage with nginx
FROM nginx:1.27-alpine

# Install security updates
RUN apk --no-cache upgrade

# Remove default nginx config and static files
RUN rm -rf /etc/nginx/conf.d/* /usr/share/nginx/html/*

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy optimized nginx.conf with security headers
COPY nginx-security.conf /etc/nginx/conf.d/security.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Add non-root user
RUN addgroup -g 101 -S nginx-custom && \
    adduser -S -D -H -u 101 -h /var/cache/nginx -s /sbin/nologin -G nginx-custom -g nginx-custom nginx-custom

# Set proper permissions
RUN chown -R nginx-custom:nginx-custom /usr/share/nginx/html && \
    chown -R nginx-custom:nginx-custom /var/cache/nginx && \
    chown -R nginx-custom:nginx-custom /var/log/nginx && \
    touch /var/run/nginx.pid && \
    chown -R nginx-custom:nginx-custom /var/run/nginx.pid

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# Expose port 80
EXPOSE 80

# Use non-root user
USER nginx-custom

# Start nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
