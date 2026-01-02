# Healthcare Module Deployment Makefile
# This Makefile handles deployment operations for the healthcare module

.PHONY: help setup build deploy clean logs migrate migrate-down health docker-build docker-up docker-down

MODULE := healthcare
BINARY := bin/healthcare
CONFIG_DIR := config/healthcare
DEPLOYMENT_DIR := deployment/healthcare

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Complete setup for healthcare module
	@echo "🏥 Setting up Healthcare module..."
	@./setup.sh
	@echo "✅ Healthcare module setup complete"

build: ## Build the healthcare module
	@echo "🔨 Building Healthcare module..."
	@cd ../.. && go build -o $(BINARY) ./app/healthcare
	@echo "✅ Build complete: $(BINARY)"

run: ## Run the healthcare module locally
	@echo "▶️  Running Healthcare module..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/local.json

run-dev: ## Run in development mode
	@echo "▶️  Running Healthcare module (development)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/development.json

run-stage: ## Run in staging mode
	@echo "▶️  Running Healthcare module (staging)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/stage.json

run-prod: ## Run in production mode
	@echo "▶️  Running Healthcare module (production)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/production.json

migrate: ## Run database migrations
	@echo "📦 Running Healthcare module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=healthcare --direction=up
	@echo "✅ Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "⏪ Rolling back Healthcare module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=healthcare --direction=down
	@echo "✅ Rollback complete"

test: ## Run tests
	@echo "🧪 Running Healthcare module tests..."
	@cd ../.. && go test -v ./app/healthcare/...
	@echo "✅ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "🧪 Running Healthcare module tests with coverage..."
	@cd ../.. && go test -v -coverprofile=coverage-healthcare.out ./app/healthcare/...
	@cd ../.. && go tool cover -html=coverage-healthcare.out
	@echo "✅ Coverage report generated"

docker-build: ## Build Docker image
	@echo "🐳 Building Healthcare module Docker image..."
	@docker build -f Dockerfile -t evero-healthcare:latest ../..
	@echo "✅ Docker image built"

docker-up: ## Start Docker containers
	@echo "🐳 Starting Healthcare module containers..."
	@docker-compose up -d
	@echo "✅ Containers started"

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Healthcare module containers..."
	@docker-compose down
	@echo "✅ Containers stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f healthcare

deploy: build migrate ## Deploy the healthcare module (build + migrate)
	@echo "🚀 Deploying Healthcare module..."
	@echo "✅ Deployment complete"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning Healthcare module artifacts..."
	@cd ../.. && rm -rf $(BINARY)
	@cd ../.. && rm -f coverage-healthcare.out
	@echo "✅ Clean complete"

health: ## Check module health
	@echo "🏥 Checking Healthcare module health..."
	@curl -f http://localhost:3001/health || echo "❌ Health check failed"

logs: ## View application logs
	@echo "📋 Viewing Healthcare module logs..."
	@tail -f ../../logs/healthcare.log

.DEFAULT_GOAL := help
