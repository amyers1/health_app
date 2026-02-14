# VitalStream Makefile
# Common commands for development and deployment

.PHONY: help build deploy start stop restart logs clean health test prod-build prod-deploy

# Default target
.DEFAULT_GOAL := help

# Variables
COMPOSE_FILE := docker-compose.yml
COMPOSE_PROD := docker-compose.prod.yml
PROJECT_NAME := vitalstream

# Colors for output
CYAN := \033[0;36m
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m

help: ## Show this help message
	@echo '$(CYAN)VitalStream Health Dashboard - Available Commands:$(NC)'
	@echo ''
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ''

# Development commands
build: ## Build development containers
	@echo '$(CYAN)Building development containers...$(NC)'
	docker-compose -f $(COMPOSE_FILE) build

start: ## Start development services
	@echo '$(CYAN)Starting development services...$(NC)'
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo '$(GREEN)Services started!$(NC)'
	@echo 'Frontend: http://localhost:13000'
	@echo 'API: http://localhost:13001/api/v1'

stop: ## Stop all services
	@echo '$(CYAN)Stopping services...$(NC)'
	docker-compose -f $(COMPOSE_FILE) down

restart: stop start ## Restart all services

logs: ## View logs (all services)
	docker-compose -f $(COMPOSE_FILE) logs -f

logs-api: ## View API logs only
	docker logs -f health-app-api

logs-frontend: ## View frontend logs only
	docker logs -f health-app-frontend

# Production commands
prod-build: ## Build production containers
	@echo '$(CYAN)Building production containers...$(NC)'
	docker-compose -f $(COMPOSE_PROD) -p $(PROJECT_NAME) build --no-cache

prod-deploy: ## Deploy to production
	@echo '$(CYAN)Deploying to production...$(NC)'
	@chmod +x deploy.sh
	./deploy.sh

prod-start: ## Start production services
	@echo '$(CYAN)Starting production services...$(NC)'
	docker-compose -f $(COMPOSE_PROD) -p $(PROJECT_NAME) up -d
	@$(MAKE) health

prod-stop: ## Stop production services
	@echo '$(CYAN)Stopping production services...$(NC)'
	docker-compose -f $(COMPOSE_PROD) -p $(PROJECT_NAME) down

prod-logs: ## View production logs
	docker-compose -f $(COMPOSE_PROD) -p $(PROJECT_NAME) logs -f

# Utility commands
health: ## Check service health
	@echo '$(CYAN)Checking service health...$(NC)'
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep health-app || echo "$(YELLOW)No services running$(NC)"
	@echo ''
	@echo '$(CYAN)Health Status:$(NC)'
	@docker inspect health-app-api 2>/dev/null | grep -A5 '"Health"' || echo "$(YELLOW)API health check not available yet$(NC)"
	@docker inspect health-app-frontend 2>/dev/null | grep -A5 '"Health"' || echo "$(YELLOW)Frontend health check not available yet$(NC)"

ps: ## List running containers
	docker-compose -f $(COMPOSE_FILE) ps

stats: ## Show container resource usage
	@echo '$(CYAN)Container Resource Usage:$(NC)'
	docker stats --no-stream health-app-api health-app-frontend

clean: ## Remove all containers, images, and volumes
	@echo '$(YELLOW)Warning: This will remove all VitalStream containers, images, and volumes!$(NC)'
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo '$(RED)Cleaning up...$(NC)'; \
		docker-compose -f $(COMPOSE_FILE) down -v --rmi all; \
		docker-compose -f $(COMPOSE_PROD) -p $(PROJECT_NAME) down -v --rmi all; \
		echo '$(GREEN)Cleanup complete!$(NC)'; \
	else \
		echo '$(GREEN)Cancelled.$(NC)'; \
	fi

test-api: ## Test API endpoint
	@echo '$(CYAN)Testing API endpoint...$(NC)'
	@curl -s http://localhost:13001/api/v1/summary | python -m json.tool || echo "$(RED)API not responding$(NC)"

test-frontend: ## Test frontend accessibility
	@echo '$(CYAN)Testing frontend...$(NC)'
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:13000/ || echo "$(RED)Frontend not responding$(NC)"

shell-api: ## Open shell in API container
	docker exec -it health-app-api sh

shell-frontend: ## Open shell in frontend container
	docker exec -it health-app-frontend sh

# Network management
network-create: ## Create required Docker networks
	@echo '$(CYAN)Creating Docker networks...$(NC)'
	docker network inspect backend >/dev/null 2>&1 || docker network create backend
	docker network inspect proxy >/dev/null 2>&1 || docker network create proxy
	@echo '$(GREEN)Networks ready!$(NC)'

network-clean: ## Remove Docker networks
	@echo '$(CYAN)Removing Docker networks...$(NC)'
	docker network rm backend proxy 2>/dev/null || true

# Database management
db-backup: ## Backup InfluxDB (requires influx CLI)
	@echo '$(CYAN)Creating InfluxDB backup...$(NC)'
	@mkdir -p backups
	@echo '$(YELLOW)Note: This requires manual InfluxDB backup configuration$(NC)'
	@echo 'See: https://docs.influxdata.com/influxdb/v3/admin/backup/'

# Validation
validate-env: ## Validate .env file
	@echo '$(CYAN)Validating environment variables...$(NC)'
	@if [ ! -f .env ]; then \
		echo '$(RED)Error: .env file not found!$(NC)'; \
		echo 'Copy .env.example to .env and configure it.'; \
		exit 1; \
	fi
	@grep -q "^INFLUX_HOST=" .env || (echo "$(RED)Missing INFLUX_HOST$(NC)" && exit 1)
	@grep -q "^INFLUX_TOKEN=" .env || (echo "$(RED)Missing INFLUX_TOKEN$(NC)" && exit 1)
	@grep -q "^INFLUX_ORG=" .env || (echo "$(RED)Missing INFLUX_ORG$(NC)" && exit 1)
	@grep -q "^INFLUX_DATABASE=" .env || (echo "$(RED)Missing INFLUX_DATABASE$(NC)" && exit 1)
	@echo '$(GREEN)Environment validation passed!$(NC)'

# Full deployment workflow
deploy: validate-env network-create prod-build prod-start health ## Full production deployment
	@echo ''
	@echo '$(GREEN)====================================$(NC)'
	@echo '$(GREEN)Deployment Complete! 🎉$(NC)'
	@echo '$(GREEN)====================================$(NC)'
	@echo ''
	@echo 'Access your application:'
	@echo '  Frontend: http://localhost:13000'
	@echo '  API: http://localhost:13001/api/v1'
	@echo ''
	@echo 'Useful commands:'
	@echo '  make logs      - View all logs'
	@echo '  make health    - Check service health'
	@echo '  make stats     - View resource usage'
	@echo '  make help      - Show all commands'
