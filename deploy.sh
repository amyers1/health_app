#!/bin/bash
# Production deployment script for VitalStream Health Dashboard

set -e  # Exit on error

echo "🚀 Starting VitalStream deployment..."

# Configuration
COMPOSE_FILE="docker-compose.prod.yml"
PROJECT_NAME="vitalstream"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if .env file exists
if [ ! -f .env ]; then
    log_error ".env file not found!"
    log_info "Please create a .env file with your InfluxDB configuration"
    exit 1
fi

# Validate required environment variables
log_info "Validating environment variables..."
required_vars=("INFLUX_HOST" "INFLUX_TOKEN" "INFLUX_ORG" "INFLUX_DATABASE")
for var in "${required_vars[@]}"; do
    if ! grep -q "^${var}=" .env; then
        log_error "Required variable ${var} not found in .env file"
        exit 1
    fi
done
log_info "Environment validation passed ✓"

# Create necessary networks if they don't exist
log_info "Checking Docker networks..."
docker network inspect backend >/dev/null 2>&1 || docker network create backend
docker network inspect proxy >/dev/null 2>&1 || docker network create proxy
log_info "Networks ready ✓"

# Pull latest changes (if using git)
if [ -d .git ]; then
    log_info "Pulling latest changes from git..."
    git pull
fi

# Stop existing containers
log_info "Stopping existing containers..."
docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} down

# Build images with no cache for production
log_info "Building production images (this may take a few minutes)..."
docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} build --no-cache

# Start services
log_info "Starting services..."
docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} up -d

# Wait for services to be healthy
log_info "Waiting for services to become healthy..."
sleep 5

# Check service health
log_info "Checking service health..."
if docker inspect health-app-api | grep -q '"Status": "healthy"' || docker inspect health-app-api | grep -q '"Status": "starting"'; then
    log_info "API service is starting/healthy ✓"
else
    log_warn "API service may not be healthy yet. Check logs with: docker logs health-app-api"
fi

if docker inspect health-app-frontend | grep -q '"Status": "healthy"' || docker inspect health-app-frontend | grep -q '"Status": "starting"'; then
    log_info "Frontend service is starting/healthy ✓"
else
    log_warn "Frontend service may not be healthy yet. Check logs with: docker logs health-app-frontend"
fi

# Display running containers
log_info "Running containers:"
docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} ps

# Show logs hint
echo ""
log_info "Deployment complete! 🎉"
log_info "Access your application at: http://localhost:13000"
log_info "API available at: http://localhost:13001/api/v1"
echo ""
log_info "Useful commands:"
echo "  - View logs: docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} logs -f"
echo "  - Stop services: docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} down"
echo "  - Restart services: docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} restart"
echo "  - Check health: docker-compose -f ${COMPOSE_FILE} -p ${PROJECT_NAME} ps"
