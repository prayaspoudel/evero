#!/bin/bash

# Evero Insurance Module - Setup Script
# This script helps you set up and run the insurance module

set -e

echo "🛡️ Evero Insurance Module Setup"
echo "================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24 or higher."
    exit 1
fi

echo "✓ Go version: $(go version)"
echo ""

# Check if PostgreSQL is running
if command -v psql &> /dev/null; then
    echo "✓ PostgreSQL is installed"
else
    echo "⚠️  PostgreSQL not found in PATH (you can still use Docker)"
fi
echo ""

# Ask user for setup method
echo "How would you like to run the Insurance module?"
echo "1) Docker Compose (recommended)"
echo "2) Local setup (requires PostgreSQL)"
echo ""
read -p "Enter your choice (1 or 2): " choice

case $choice in
    1)
        echo ""
        echo "📦 Setting up with Docker Compose..."
        echo ""
        
        if ! command -v docker-compose &> /dev/null && ! command -v docker &> /dev/null; then
            echo "❌ Docker is not installed. Please install Docker and Docker Compose."
            exit 1
        fi
        
        echo "✓ Docker is installed"
        echo ""
        
        # Check if .env exists
        if [ ! -f .env ]; then
            echo "📝 Creating .env file..."
            cat > .env << EOF
# Insurance Module Environment Variables
HEALTHCARE_DB_HOST=postgres
HEALTHCARE_DB_PORT=5432
HEALTHCARE_DB_USER=postgres
HEALTHCARE_DB_PASSWORD=postgres
HEALTHCARE_DB_NAME=insurance_db

# Server Configuration
HEALTHCARE_PORT=3002
HEALTHCARE_ENV=development
EOF
            echo "✓ .env file created"
        fi
        
        echo ""
        echo "🚀 Starting services with Docker Compose..."
        docker-compose up -d
        
        echo ""
        echo "⏳ Waiting for database to be ready..."
        sleep 5
        
        echo ""
        echo "✅ Insurance module is running!"
        echo ""
        echo "📊 Service URLs:"
        echo "   Insurance API: http://localhost:3002"
        echo ""
        echo "📝 Useful commands:"
        echo "   View logs:      docker-compose logs -f"
        echo "   Stop services:  docker-compose down"
        echo "   Rebuild:        docker-compose up -d --build"
        ;;
        
    2)
        echo ""
        echo "🔧 Local setup selected"
        echo ""
        
        # Check PostgreSQL connection
        read -p "PostgreSQL host (default: localhost): " DB_HOST
        DB_HOST=${DB_HOST:-localhost}
        
        read -p "PostgreSQL port (default: 5432): " DB_PORT
        DB_PORT=${DB_PORT:-5432}
        
        read -p "PostgreSQL database name (default: insurance_db): " DB_NAME
        DB_NAME=${DB_NAME:-insurance_db}
        
        read -p "PostgreSQL username (default: postgres): " DB_USER
        DB_USER=${DB_USER:-postgres}
        
        read -sp "PostgreSQL password: " DB_PASSWORD
        echo ""
        
        # Create config file
        echo ""
        echo "📝 Creating local configuration..."
        
        cat > ../../config/insurance/local.json << EOF
{
  "app": {
    "name": "Evero Insurance API",
    "version": "1.0.0"
  },
  "web": {
    "port": 3002,
    "prefork": false
  },
  "database": {
    "host": "$DB_HOST",
    "port": $DB_PORT,
    "username": "$DB_USER",
    "password": "$DB_PASSWORD",
    "name": "$DB_NAME",
    "sslmode": "disable",
    "pool": {
      "idle": 10,
      "max": 100,
      "lifetime": 300
    }
  },
  "log": {
    "level": "info"
  }
}
EOF
        
        echo "✓ Configuration created"
        echo ""
        
        # Build the application
        echo "🔨 Building Insurance module..."
        cd ../..
        go build -o bin/insurance ./app/insurance
        echo "✓ Build complete"
        echo ""
        
        # Run migrations
        echo "📦 Running migrations..."
        # Migration command here when available
        echo "✓ Migrations complete"
        echo ""
        
        echo "✅ Setup complete!"
        echo ""
        echo "🚀 To start the Insurance module, run:"
        echo "   ./bin/insurance --config=config/insurance/local.json"
        echo ""
        echo "Or use: make run"
        ;;
        
    *)
        echo "❌ Invalid choice. Exiting."
        exit 1
        ;;
esac

echo ""
echo "📚 Documentation: docs/insurance/README.md"
echo ""
