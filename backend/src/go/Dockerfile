# Production Dockerfile for Go API
# Multi-stage build for minimal image size and security

# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy dependency files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build with optimizations and security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -installsuffix cgo \
    -trimpath \
    -o main .

# Stage 2: Runtime stage (minimal)
FROM alpine:3.19

# Add non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set timezone data location
ENV ZONEINFO=/usr/share/zoneinfo

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/main .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:13001/api/v1/summary || exit 1

# Expose port
EXPOSE 13001

# Run the application
CMD ["./main"]
