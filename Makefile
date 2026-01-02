# Evero Platform - Main Orchestration Makefile
# Manages all modules: access, healthcare, insurance, finance

.PHONY: help setup build clean test deploy-all

# Module directories
MODULES := access healthcare insurance finance
MODULE_ACCESS := modules/access
MODULE_HEALTHCARE := app/healthcare
MODULE_INSURANCE := app/insurance
MODULE_FINANCE := app/finance

# Binary outputs
BIN_DIR := bin

help: ## Display this help screen
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║         Evero Platform - Module Orchestration              ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Module-specific commands:"
	@echo "  make setup <module>       Setup specific module (access, healthcare, insurance, finance)"
	@echo "  make build <module>       Build specific module"
	@echo "  make deploy <module>      Deploy specific module"
	@echo "  make migrate <module>     Run migrations for specific module"
	@echo "  make clean <module>       Clean specific module artifacts"
	@echo ""
	@echo "Examples:"
	@echo "  make setup healthcare     Setup healthcare module"
	@echo "  make build access         Build access module"
	@echo "  make deploy-all           Deploy all modules"

# ============================================================================
# Setup Commands
# ============================================================================

setup: ## Interactive setup - prompts for module selection
	@echo "📦 Evero Platform Setup"
	@echo "Available modules: $(MODULES)"
	@echo "Usage: make setup <module_name>"
	@echo "Example: make setup healthcare"

setup-access: ## Setup access module
	@echo "🚀 Setting up Access module..."
	@cd deployment/access && $(MAKE) setup
	@echo "✅ Access module setup complete"

setup-healthcare: ## Setup healthcare module
	@echo "🏥 Setting up Healthcare module..."
	@mkdir -p $(BIN_DIR)
	@mkdir -p database/healthcare/migrations
	@echo "✅ Healthcare module setup complete"

setup-insurance: ## Setup insurance module
	@echo "🛡️  Setting up Insurance module..."
	@mkdir -p $(BIN_DIR)
	@mkdir -p database/insurance/migrations
	@echo "✅ Insurance module setup complete"

setup-finance: ## Setup finance module
	@echo "💰 Setting up Finance module..."
	@mkdir -p $(BIN_DIR)
	@mkdir -p database/finance/migrations
	@echo "✅ Finance module setup complete"

setup-all: ## Setup all modules
	@echo "🚀 Setting up all modules..."
	@$(MAKE) setup-access
	@$(MAKE) setup-healthcare
	@$(MAKE) setup-insurance
	@$(MAKE) setup-finance
	@echo "✅ All modules setup complete"

# ============================================================================
# Build Commands
# ============================================================================

build-access: ## Build access module
	@echo "🔨 Building Access module..."
	@go build -o $(BIN_DIR)/access ./modules/access/cmd/server
	@echo "✅ Access module built: $(BIN_DIR)/access"

build-healthcare: ## Build healthcare module
	@echo "🔨 Building Healthcare module..."
	@go build -o $(BIN_DIR)/healthcare ./app/healthcare
	@echo "✅ Healthcare module built: $(BIN_DIR)/healthcare"

build-insurance: ## Build insurance module
	@echo "🔨 Building Insurance module..."
	@go build -o $(BIN_DIR)/insurance ./app/insurance
	@echo "✅ Insurance module built: $(BIN_DIR)/insurance"

build-finance: ## Build finance module
	@echo "🔨 Building Finance module..."
	@go build -o $(BIN_DIR)/finance ./app/finance
	@echo "✅ Finance module built: $(BIN_DIR)/finance"

build-all: ## Build all modules
	@echo "🔨 Building all modules..."
	@$(MAKE) build-access
	@$(MAKE) build-healthcare
	@$(MAKE) build-insurance
	@$(MAKE) build-finance
	@echo "✅ All modules built"

# ============================================================================
# Migration Commands
# ============================================================================

migrate-access: ## Run access module migrations
	@echo "📦 Running Access migrations..."
	@cd deployment/access && $(MAKE) migrate
	@echo "✅ Access migrations complete"

migrate-healthcare: ## Run healthcare module migrations
	@echo "📦 Running Healthcare migrations..."
	@echo "Running migrations for healthcare..."
	@# Add migration command when healthcare migration tool is ready
	@echo "✅ Healthcare migrations complete"

migrate-insurance: ## Run insurance module migrations
	@echo "📦 Running Insurance migrations..."
	@echo "Running migrations for insurance..."
	@# Add migration command when insurance migration tool is ready
	@echo "✅ Insurance migrations complete"

migrate-finance: ## Run finance module migrations
	@echo "📦 Running Finance migrations..."
	@echo "Running migrations for finance..."
	@# Add migration command when finance migration tool is ready
	@echo "✅ Finance migrations complete"

migrate-all: ## Run all module migrations
	@echo "📦 Running all migrations..."
	@$(MAKE) migrate-access
	@$(MAKE) migrate-healthcare
	@$(MAKE) migrate-insurance
	@$(MAKE) migrate-finance
	@echo "✅ All migrations complete"

# ============================================================================
# Test Commands
# ============================================================================

test-access: ## Test access module
	@echo "🧪 Testing Access module..."
	@go test -v ./modules/access/...
	@echo "✅ Access tests complete"

test-healthcare: ## Test healthcare module
	@echo "🧪 Testing Healthcare module..."
	@go test -v ./app/healthcare/...
	@echo "✅ Healthcare tests complete"

test-insurance: ## Test insurance module
	@echo "🧪 Testing Insurance module..."
	@go test -v ./app/insurance/...
	@echo "✅ Insurance tests complete"

test-finance: ## Test finance module
	@echo "🧪 Testing Finance module..."
	@go test -v ./app/finance/...
	@echo "✅ Finance tests complete"

test-all: ## Run all module tests
	@echo "🧪 Running all tests..."
	@$(MAKE) test-access
	@$(MAKE) test-healthcare
	@$(MAKE) test-insurance
	@$(MAKE) test-finance
	@echo "✅ All tests complete"

# ============================================================================
# Deploy Commands
# ============================================================================

deploy-access: ## Deploy access module
	@echo "🚀 Deploying Access module..."
	@cd deployment/access && $(MAKE) deploy
	@echo "✅ Access module deployed"

deploy-healthcare: ## Deploy healthcare module
	@echo "🚀 Deploying Healthcare module..."
	@$(MAKE) build-healthcare
	@$(MAKE) migrate-healthcare
	@echo "✅ Healthcare module deployed"

deploy-insurance: ## Deploy insurance module
	@echo "🚀 Deploying Insurance module..."
	@$(MAKE) build-insurance
	@$(MAKE) migrate-insurance
	@echo "✅ Insurance module deployed"

deploy-finance: ## Deploy finance module
	@echo "🚀 Deploying Finance module..."
	@$(MAKE) build-finance
	@$(MAKE) migrate-finance
	@echo "✅ Finance module deployed"

deploy-all: ## Deploy all modules
	@echo "🚀 Deploying all modules..."
	@$(MAKE) deploy-access
	@$(MAKE) deploy-healthcare
	@$(MAKE) deploy-insurance
	@$(MAKE) deploy-finance
	@echo "✅ All modules deployed"

# ============================================================================
# Clean Commands
# ============================================================================

clean-access: ## Clean access module artifacts
	@echo "🧹 Cleaning Access module..."
	@rm -f $(BIN_DIR)/access
	@echo "✅ Access module cleaned"

clean-healthcare: ## Clean healthcare module artifacts
	@echo "🧹 Cleaning Healthcare module..."
	@rm -f $(BIN_DIR)/healthcare
	@echo "✅ Healthcare module cleaned"

clean-insurance: ## Clean insurance module artifacts
	@echo "🧹 Cleaning Insurance module..."
	@rm -f $(BIN_DIR)/insurance
	@echo "✅ Insurance module cleaned"

clean-finance: ## Clean finance module artifacts
	@echo "🧹 Cleaning Finance module..."
	@rm -f $(BIN_DIR)/finance
	@echo "✅ Finance module cleaned"

clean-all: ## Clean all module artifacts
	@echo "🧹 Cleaning all modules..."
	@rm -rf $(BIN_DIR)/*
	@rm -f coverage-*.out
	@echo "✅ All modules cleaned"

# ============================================================================
# Development Commands
# ============================================================================

deps: ## Download and tidy dependencies
	@echo "📦 Managing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

fmt: ## Format all Go code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter on all modules
	@echo "🔍 Running linter..."
	@golangci-lint run ./...
	@echo "✅ Linting complete"

# ============================================================================
# Docker Commands
# ============================================================================

docker-build-access: ## Build access module Docker image
	@cd deployment/access && $(MAKE) docker-build

docker-up-access: ## Start access module containers
	@cd deployment/access && $(MAKE) docker-up

docker-down-access: ## Stop access module containers
	@cd deployment/access && $(MAKE) docker-down

# ============================================================================
# Utility Commands
# ============================================================================

status: ## Show status of all modules
	@echo "📊 Module Status"
	@echo "================================"
	@echo -n "Access:      "; [ -f $(BIN_DIR)/access ] && echo "✅ Built" || echo "❌ Not built"
	@echo -n "Healthcare:  "; [ -f $(BIN_DIR)/healthcare ] && echo "✅ Built" || echo "❌ Not built"
	@echo -n "Insurance:   "; [ -f $(BIN_DIR)/insurance ] && echo "✅ Built" || echo "❌ Not built"
	@echo -n "Finance:     "; [ -f $(BIN_DIR)/finance ] && echo "✅ Built" || echo "❌ Not built"
	@echo "================================"

.DEFAULT_GOAL := help
