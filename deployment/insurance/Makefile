# Insurance Module Deployment Makefile
# This Makefile handles deployment operations for the insurance module

.PHONY: help setup build deploy clean logs migrate migrate-down health docker-build docker-up docker-down

MODULE := insurance
BINARY := bin/insurance
CONFIG_DIR := config/insurance
DEPLOYMENT_DIR := deployment/insurance

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Complete setup for insurance module
	@echo "🛡️  Setting up Insurance module..."
	@./setup.sh
	@echo "✅ Insurance module setup complete"

build: ## Build the insurance module
	@echo "🔨 Building Insurance module..."
	@cd ../.. && go build -o $(BINARY) ./app/insurance
	@echo "✅ Build complete: $(BINARY)"

run: ## Run the insurance module locally
	@echo "▶️  Running Insurance module..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/local.json

run-dev: ## Run in development mode
	@echo "▶️  Running Insurance module (development)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/development.json

run-stage: ## Run in staging mode
	@echo "▶️  Running Insurance module (staging)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/stage.json

run-prod: ## Run in production mode
	@echo "▶️  Running Insurance module (production)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/production.json

migrate: ## Run database migrations
	@echo "📦 Running Insurance module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=insurance --direction=up
	@echo "✅ Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "⏪ Rolling back Insurance module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=insurance --direction=down
	@echo "✅ Rollback complete"

test: ## Run tests
	@echo "🧪 Running Insurance module tests..."
	@cd ../.. && go test -v ./app/insurance/...
	@echo "✅ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "🧪 Running Insurance module tests with coverage..."
	@cd ../.. && go test -v -coverprofile=coverage-insurance.out ./app/insurance/...
	@cd ../.. && go tool cover -html=coverage-insurance.out
	@echo "✅ Coverage report generated"

docker-build: ## Build Docker image
	@echo "🐳 Building Insurance module Docker image..."
	@docker build -f Dockerfile -t evero-insurance:latest ../..
	@echo "✅ Docker image built"

docker-up: ## Start Docker containers
	@echo "🐳 Starting Insurance module containers..."
	@docker-compose up -d
	@echo "✅ Containers started"

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Insurance module containers..."
	@docker-compose down
	@echo "✅ Containers stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f insurance

deploy: build migrate ## Deploy the insurance module (build + migrate)
	@echo "🚀 Deploying Insurance module..."
	@echo "✅ Deployment complete"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning Insurance module artifacts..."
	@cd ../.. && rm -rf $(BINARY)
	@cd ../.. && rm -f coverage-insurance.out
	@echo "✅ Clean complete"

health: ## Check module health
	@echo "🛡️  Checking Insurance module health..."
	@curl -f http://localhost:3002/health || echo "❌ Health check failed"

logs: ## View application logs
	@echo "📋 Viewing Insurance module logs..."
	@tail -f ../../logs/insurance.log

.DEFAULT_GOAL := help
