# Access Module Deployment Makefile
# This Makefile handles deployment operations for the access module

.PHONY: help setup build deploy clean logs migrate migrate-down health docker-build docker-up docker-down

MODULE := access
BINARY := bin/access
CONFIG_DIR := config/access
DEPLOYMENT_DIR := deployment/access

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Complete setup for access module
	@echo "🚀 Setting up Access module..."
	@./setup.sh
	@echo "✅ Access module setup complete"

build: ## Build the access module
	@echo "🔨 Building Access module..."
	@cd ../.. && go build -o $(BINARY) ./modules/access/cmd/server
	@echo "✅ Build complete: $(BINARY)"

run: ## Run the access module locally
	@echo "▶️  Running Access module..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/local.json

run-dev: ## Run in development mode
	@echo "▶️  Running Access module (development)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/development.json

run-stage: ## Run in staging mode
	@echo "▶️  Running Access module (staging)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/stage.json

run-prod: ## Run in production mode
	@echo "▶️  Running Access module (production)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/production.json

migrate: ## Run database migrations
	@echo "📦 Running Access module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=access --direction=up
	@echo "✅ Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "⏪ Rolling back Access module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=access --direction=down
	@echo "✅ Rollback complete"

test: ## Run tests
	@echo "🧪 Running Access module tests..."
	@cd ../.. && go test -v ./modules/access/...
	@echo "✅ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "🧪 Running Access module tests with coverage..."
	@cd ../.. && go test -v -coverprofile=coverage-access.out ./modules/access/...
	@cd ../.. && go tool cover -html=coverage-access.out
	@echo "✅ Coverage report generated"

docker-build: ## Build Docker image
	@echo "🐳 Building Access module Docker image..."
	@docker build -f Dockerfile -t evero-access:latest ../..
	@echo "✅ Docker image built"

docker-up: ## Start Docker containers
	@echo "🐳 Starting Access module containers..."
	@docker-compose up -d
	@echo "✅ Containers started"

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Access module containers..."
	@docker-compose down
	@echo "✅ Containers stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f access

deploy: build migrate ## Deploy the access module (build + migrate)
	@echo "🚀 Deploying Access module..."
	@echo "✅ Deployment complete"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning Access module artifacts..."
	@cd ../.. && rm -rf $(BINARY)
	@cd ../.. && rm -f coverage-access.out
	@echo "✅ Clean complete"

health: ## Check module health
	@echo "🏥 Checking Access module health..."
	@curl -f http://localhost:3000/health || echo "❌ Health check failed"

logs: ## View application logs
	@echo "📋 Viewing Access module logs..."
	@tail -f ../../logs/access.log

.DEFAULT_GOAL := help
