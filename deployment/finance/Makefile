# Finance Module Deployment Makefile
# This Makefile handles deployment operations for the finance module

.PHONY: help setup build deploy clean logs migrate migrate-down health docker-build docker-up docker-down

MODULE := finance
BINARY := bin/finance
CONFIG_DIR := config/finance
DEPLOYMENT_DIR := deployment/finance

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

setup: ## Complete setup for finance module
	@echo "💰 Setting up Finance module..."
	@./setup.sh
	@echo "✅ Finance module setup complete"

build: ## Build the finance module
	@echo "🔨 Building Finance module..."
	@cd ../.. && go build -o $(BINARY) ./app/finance
	@echo "✅ Build complete: $(BINARY)"

run: ## Run the finance module locally
	@echo "▶️  Running Finance module..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/local.json

run-dev: ## Run in development mode
	@echo "▶️  Running Finance module (development)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/development.json

run-stage: ## Run in staging mode
	@echo "▶️  Running Finance module (staging)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/stage.json

run-prod: ## Run in production mode
	@echo "▶️  Running Finance module (production)..."
	@cd ../.. && ./$(BINARY) --config=$(CONFIG_DIR)/production.json

migrate: ## Run database migrations
	@echo "📦 Running Finance module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=finance --direction=up
	@echo "✅ Migrations complete"

migrate-down: ## Rollback database migrations
	@echo "⏪ Rolling back Finance module migrations..."
	@cd ../.. && go run ./infrastructure/database/migrate.go --module=finance --direction=down
	@echo "✅ Rollback complete"

test: ## Run tests
	@echo "🧪 Running Finance module tests..."
	@cd ../.. && go test -v ./app/finance/...
	@echo "✅ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "🧪 Running Finance module tests with coverage..."
	@cd ../.. && go test -v -coverprofile=coverage-finance.out ./app/finance/...
	@cd ../.. && go tool cover -html=coverage-finance.out
	@echo "✅ Coverage report generated"

docker-build: ## Build Docker image
	@echo "🐳 Building Finance module Docker image..."
	@docker build -f Dockerfile -t evero-finance:latest ../..
	@echo "✅ Docker image built"

docker-up: ## Start Docker containers
	@echo "🐳 Starting Finance module containers..."
	@docker-compose up -d
	@echo "✅ Containers started"

docker-down: ## Stop Docker containers
	@echo "🐳 Stopping Finance module containers..."
	@docker-compose down
	@echo "✅ Containers stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f finance

deploy: build migrate ## Deploy the finance module (build + migrate)
	@echo "🚀 Deploying Finance module..."
	@echo "✅ Deployment complete"

clean: ## Clean build artifacts
	@echo "🧹 Cleaning Finance module artifacts..."
	@cd ../.. && rm -rf $(BINARY)
	@cd ../.. && rm -f coverage-finance.out
	@echo "✅ Clean complete"

health: ## Check module health
	@echo "💰 Checking Finance module health..."
	@curl -f http://localhost:3003/health || echo "❌ Health check failed"

logs: ## View application logs
	@echo "📋 Viewing Finance module logs..."
	@tail -f ../../logs/finance.log

.DEFAULT_GOAL := help
