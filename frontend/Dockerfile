# Production Dockerfile for React Frontend - Non-Root Version
# Uses nginx-unprivileged for true non-root operation

# Stage 1: Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY src/package*.json ./

# Install ALL dependencies (including devDependencies needed for build)
RUN npm install && \
    npm cache clean --force

# Copy source code
COPY src/ ./

# Build the application with production optimizations
ENV NODE_ENV=production
RUN npm run build

# Stage 2: Runtime stage with nginx-unprivileged
FROM nginxinc/nginx-unprivileged:1.27-alpine

# Install security updates (as root temporarily)
USER root
RUN apk --no-cache upgrade

# Remove default nginx config and static files
RUN rm -rf /etc/nginx/conf.d/* /usr/share/nginx/html/*

# Copy custom nginx configuration (modified for port 8080)
COPY nginx-unprivileged.conf /etc/nginx/conf.d/default.conf

# Copy security headers configuration
COPY nginx-security.conf /etc/nginx/conf.d/security.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist/ /usr/share/nginx/html/

# Set proper permissions
RUN chown -R nginx:nginx /usr/share/nginx/html

# Switch back to nginx user
USER nginx

# Health check (port 8080 for unprivileged)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Expose port 8080 (unprivileged port)
EXPOSE 8080

# Start nginx in foreground
CMD ["nginx", "-g", "daemon off;"]
