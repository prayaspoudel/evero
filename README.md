# Evero Platform

Enterprise-grade modular platform for building scalable applications across healthcare, insurance, finance, and banking sectors. Built with clean architecture principles and Go, Evero provides a robust foundation for domain-specific services with shared infrastructure components.

## 🏗️ Architecture

Evero follows a **modular monolith architecture** with clear separation of concerns:

- **Multi-Domain Support**: Access (SSO), Healthcare, Insurance, Finance, and Banking modules
- **Clean Architecture**: Domain-driven design with layers (Entity, Use Case, Repository, Controller)
- **Shared Infrastructure**: Reusable components across all modules
- **Event-Driven**: Kafka integration for asynchronous messaging
- **API-First**: RESTful APIs with comprehensive validation
- **Modular Deployment**: Each module can be deployed independently

## 🚀 Features

### Core Infrastructure
- **Configuration Management**: Environment-specific configurations (local, development, staging, production)
- **Database Support**: PostgreSQL with GORM ORM
- **Caching**: Redis integration for high-performance caching
- **Message Broker**: Kafka and RabbitMQ support
- **Logging**: Structured logging with Logrus and Zap
- **Validation**: Request validation with go-playground/validator
- **Routing**: Fiber v2 (primary framework)

### Access Module (Authentication & Authorization)

Enterprise-grade single sign-on and access control system.

#### Features
- **Authentication**
  - User registration and login
  - JWT-based access and refresh tokens
  - Token rotation and refresh
  - Session management with device tracking
  - Email verification workflow

- **Authorization**
  - Multi-company/tenant support
  - OAuth 2.0 authorization code flow
  - Two-factor authentication (TOTP)
  - Role-based access control (RBAC)
  - Fine-grained permissions

- **Security**
  - Password reset and verification
  - Account lockout protection
  - Audit logging
  - bcrypt password hashing
  - Secure token management

**📖 Documentation**: [docs/access/IMPLEMENTATION_SUMMARY.md](docs/access/IMPLEMENTATION_SUMMARY.md)  
**Migration Guide**: [docs/SSO_MIGRATION_SUMMARY.md](docs/SSO_MIGRATION_SUMMARY.md)  
**Quick Start**: [docs/access/QUICK_REFERENCE.md](docs/access/QUICK_REFERENCE.md)

**Build & Deploy**:
```bash
make build-access    # Build the module
make deploy-access   # Deploy the module
```

**API**: http://localhost:3000 (configurable)

### Healthcare Module

Comprehensive healthcare management system for patient care, appointments, and medical records.

#### Features
- **Patient Management**
  - Patient registration and demographics
  - Medical history tracking
  - Insurance information
  - Emergency contacts

- **Appointments & Scheduling**
  - Appointment booking
  - Provider availability
  - Appointment reminders
  - Waitlist management

- **Electronic Medical Records (EMR)**
  - Clinical notes
  - Diagnosis recording (ICD-10)
  - Prescription management
  - Lab results integration

- **Billing Integration**
  - Insurance claim generation
  - Payment processing
  - Finance module integration

**📖 Documentation**: [docs/healthcare/README.md](docs/healthcare/README.md)

**Build & Deploy**:
```bash
make build-healthcare    # Build the module
make deploy-healthcare   # Deploy the module
```

**API**: http://localhost:3001 (configurable)

### Insurance Module

Complete insurance management platform for policies, claims, and underwriting.

#### Features
- **Policy Management**
  - Policy creation and issuance
  - Premium calculation
  - Renewals and endorsements
  - Coverage management

- **Claims Processing**
  - Claim submission and intake
  - Assessment and adjudication
  - Payment processing
  - Fraud detection

- **Underwriting**
  - Risk assessment
  - Quote generation
  - Automated underwriting rules
  - Manual workflow support

- **Agent & Commission Management**
  - Agent registration
  - Commission calculation
  - Performance analytics

**📖 Documentation**: [docs/insurance/README.md](docs/insurance/README.md)

**Build & Deploy**:
```bash
make build-insurance    # Build the module
make deploy-insurance   # Deploy the module
```

**API**: http://localhost:3002 (configurable)

### Finance Module

Comprehensive financial management system with general ledger, AR/AP, and budgeting.

#### Features
- **General Ledger**
  - Chart of accounts
  - Double-entry bookkeeping
  - Journal entries
  - Period closing
  - Multi-currency support

- **Accounts Receivable/Payable**
  - Customer invoicing
  - Payment tracking
  - Aging reports
  - Vendor management

- **Budgeting & Forecasting**
  - Budget creation
  - Variance analysis
  - Cash flow forecasting

- **Financial Reporting**
  - Balance sheet
  - Income statement
  - Cash flow statement
  - Custom reports

**📖 Documentation**: [docs/finance/README.md](docs/finance/README.md)

**Build & Deploy**:
```bash
make build-finance    # Build the module
make deploy-finance   # Deploy the module
```

**API**: http://localhost:3003 (configurable)

### Banking Module

Core banking platform for account management, transactions, and lending (in planning).

**📖 Documentation**: [docs/banking/README.md](docs/banking/README.md)  
**Status**: 📋 Planning Phase

## 📁 Project Structure

```
evero/
├── Makefile                      # Root orchestration for all modules
├── bin/                          # Compiled binaries (gitignored)
│   ├── access
│   ├── healthcare
│   ├── insurance
│   └── finance
│
├── app/                          # Application entry points
│   ├── healthcare/               # Healthcare application
│   ├── insurance/                # Insurance application
│   └── finance/                  # Finance application
│
├── modules/                      # Domain modules
│   ├── access/                   # Authentication & Authorization module
│   │   ├── cmd/server/           # Entry point
│   │   ├── entities/             # Domain entities
│   │   ├── models/               # Request/Response models
│   │   ├── repositories/         # Data access layer
│   │   ├── usecases/             # Business logic
│   │   ├── controllers/          # HTTP handlers
│   │   └── route.go              # Route definitions
│   └── healthcare/               # Healthcare domain module
│       ├── delivery/             # HTTP handlers/controllers
│       ├── entity/               # Domain entities
│       ├── features/             # Business logic (use cases)
│       ├── gateway/              # External integrations
│       ├── model/                # Request/Response models
│       ├── repository/           # Data access layer
│       └── test/                 # Unit and integration tests
│
├── infrastructure/               # Shared infrastructure components
│   ├── cache/                    # Cache management (Redis, In-memory)
│   ├── config/                   # Configuration management
│   ├── database/                 # Database connections
│   ├── logger/                   # Logging utilities
│   ├── message-broker/           # Kafka/RabbitMQ integration
│   ├── router/                   # HTTP routers (Fiber, Gin, Mux)
│   ├── setup/                    # Infrastructure bootstrapping
│   └── validator/                # Request validation
│
├── config/                       # Configuration files
│   ├── access/                   # Access module configs
│   │   ├── local.json
│   │   ├── development.json
│   │   ├── stage.json
│   │   └── production.json
│   ├── healthcare/               # Healthcare configs
│   ├── insurance/                # Insurance configs
│   └── finance/                  # Finance configs
│
├── database/                     # Database migrations and seeds
│   ├── access/migrations/        # SSO database schemas
│   ├── healthcare/migrations/    # Healthcare schemas
│   ├── insurance/migrations/     # Insurance schemas
│   └── finance/migrations/       # Finance schemas (10 tables)
│
├── deployment/                   # Deployment configurations
│   ├── access/                   # Access module deployment
│   │   ├── Makefile              # Deployment tasks
│   │   ├── setup.sh              # Setup script
│   │   ├── Dockerfile
│   │   └── docker-compose.yml
│   ├── healthcare/               # Healthcare deployment
│   ├── insurance/                # Insurance deployment
│   └── finance/                  # Finance deployment
│
├── packages/                     # External service integrations
│   ├── lib/                      # Shared libraries
│   ├── sendgrid/                 # Email service
│   └── twilio/                   # SMS service
│
└── docs/                         # Documentation
    ├── evero/                    # Platform technical docs
    ├── access/                   # Access module documentation
    ├── healthcare/               # Healthcare module documentation
    ├── insurance/                # Insurance module documentation
    ├── finance/                  # Finance module documentation
    └── banking/                  # Banking module documentation
```

## 🛠️ Technology Stack

### Core
- **Go**: 1.24.4
- **Web Framework**: Fiber v2 (primary), Gin, Gorilla Mux
- **ORM**: GORM v1.30.3
- **Database**: PostgreSQL (primary), MySQL support via drivers
- **Validation**: go-playground/validator v10

### Infrastructure
- **Cache**: Redis (go-redis v9)
- **Message Broker**: Apache Kafka (Sarama v1.46), RabbitMQ
- **Logging**: Logrus v1.9.3, Zap v1.27.0
- **Configuration**: Viper v1.20.1
- **Database Drivers**: 
  - PostgreSQL: lib/pq v1.10.9
  - MySQL: gorm.io/driver/mysql v1.6.0
- **Security**: bcrypt (golang.org/x/crypto v0.41.0)
- **UUID**: google/uuid v1.6.0

## 🚦 Getting Started

### Prerequisites

- **Go**: 1.24.4 or higher
- **PostgreSQL**: 14 or higher
- **Redis**: 6 or higher (optional, for caching)
- **Kafka**: 2.8 or higher (optional, for event streaming)
- **Make**: For using the Makefile commands

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/prayaspoudel/evero.git
   cd evero
   ```

2. **View available commands**
   ```bash
   make help
   ```

3. **Setup a module** (e.g., healthcare)
   ```bash
   make setup-healthcare
   ```

4. **Build a module**
   ```bash
   make build-healthcare
   ```

5. **Run migrations**
   ```bash
   make migrate-healthcare
   ```

6. **Deploy a module**
   ```bash
   make deploy-healthcare
   ```

### Module-Specific Setup

Each module can be set up and deployed independently:

**Access Module (SSO)**:
```bash
make setup-access      # Setup access module
make build-access      # Build binary
make deploy-access     # Deploy with migrations
```

**Healthcare Module**:
```bash
make setup-healthcare
make build-healthcare
make deploy-healthcare
```

**Finance Module**:
```bash
make setup-finance
make build-finance
make deploy-finance
```

**All Modules**:
```bash
make setup-all         # Setup all modules
make build-all         # Build all modules
make deploy-all        # Deploy all modules
```

### Check Module Status

```bash
make status
```

Output:
```
📊 Module Status
================================
Access:      ✅ Built
Healthcare:  ✅ Built
Insurance:   ❌ Not built
Finance:     ✅ Built
================================
```

## 📝 Configuration

Evero uses environment-specific JSON configuration files. Each module has its own configuration directory:

```
config/
├── access/
│   ├── local.json          # Local development
│   ├── development.json    # Development environment
│   ├── stage.json          # Staging environment
│   └── production.json     # Production environment
├── healthcare/
├── insurance/
└── finance/
```

### Configuration Structure

Each module follows a consistent configuration structure:

```json
{
  "app": {
    "name": "Evero Healthcare API",
    "version": "1.0.0"
  },
  "web": {
    "port": 3000,
    "prefork": false,
    "cors_enabled": true
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "postgres",
    "name": "evero_db",
    "sslmode": "disable",
    "pool": {
      "idle": 10,
      "max": 100,
      "lifetime": 300
    }
  },
  "kafka": {
    "bootstrap.servers": "localhost:9092",
    "producer.enabled": false,
    "group.id": "evero-service"
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "log": {
    "level": "info",
    "format": "json"
  }
}
```

### Environment Selection

Set the environment using:
```bash
export EVERO_ENV=production  # Options: local, development, stage, production
```

Or specify when running:
```bash
./bin/access --config=config/access/production.json
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Test specific module
make test-access
make test-healthcare
make test-finance

# All module tests
make test-all

# Run with coverage
go test -cover ./...
```

### Module Status

Check which modules are built:
```bash
make status
```

## 🏛️ Infrastructure Components

### Configuration Manager
Centralized configuration loading with support for:
- Environment-specific files  
- Module-specific overrides
- Type-safe access methods
- Hot-reload capability

See [infrastructure/config/example_usage.md](infrastructure/config/example_usage.md)

### Cache Manager
Multi-backend caching support:
- Redis
- In-memory cache
- Factory pattern for easy switching

### Database Manager
Database connection management:
- PostgreSQL support with GORM
- Connection pooling
- Migration support
- Multi-database support

### Message Broker
Event-driven messaging:
- Kafka producer/consumer
- RabbitMQ support
- Async event publishing

### Router
HTTP routing with multiple framework support:
- Fiber v2 (primary)
- Gin
- Gorilla Mux

### Logger
Structured logging:
- Logrus
- Zap
- Configurable log levels
- JSON formatting

See [infrastructure/](infrastructure/) for detailed documentation.

## 🔒 Security

- **Authentication**: JWT-based token authentication
- **Password Security**: bcrypt hashing
- **Input Validation**: Comprehensive request validation
- **SQL Injection**: Protected via ORM parameterized queries
- **Authorization**: User-specific data isolation

## 🚀 Deployment

### Docker Support

Each module can be deployed independently using Docker:

```bash
# Build Docker image
make docker-build-access

# Start containers
make docker-up-access

# Stop containers
make docker-down-access
```

### Deployment Files

Each module has its own deployment configuration:

```
deployment/
├── access/
│   ├── Makefile              # Deployment commands
│   ├── setup.sh              # Setup script
│   ├── Dockerfile            # Docker image
│   └── docker-compose.yml    # Orchestration
├── healthcare/
├── insurance/
└── finance/
```

### Environment Configuration

Set the environment variable:
```bash
export EVERO_ENV=production
```

Configuration files are loaded based on this variable:
- `local` → config/[module]/local.json
- `development` → config/[module]/development.json
- `stage` → config/[module]/stage.json
- `production` → config/[module]/production.json

### Production Deployment

1. Build the module:
   ```bash
   make build-access
   ```

2. Run migrations:
   ```bash
   make migrate-access
   ```

3. Start the service:
   ```bash
   ./bin/access --config=config/access/production.json
   ```

Or use the combined deploy command:
```bash
make deploy-access  # Builds + migrates
```

## 📚 API Documentation

### Response Format

All API responses follow a consistent format:

**Success Response:**
```json
{
  "code": 200,
  "status": "success",
  "data": { ... }
}
```

**Error Response:**
```json
{
  "code": 400,
  "status": "error",
  "message": "Validation failed",
  "errors": ["field: first_name is required"]
}
```

### Pagination

List endpoints support pagination:

```
GET /api/v1/[resource]?page=1&size=10
```

Response includes metadata:
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "size": 10,
    "total_items": 50,
    "total_pages": 5
  }
}
```

### Module-Specific APIs

- **Access Module**: Authentication, authorization, OAuth 2.0, 2FA
  - See [docs/access/QUICK_REFERENCE.md](docs/access/QUICK_REFERENCE.md)
  
- **Healthcare Module**: Patient management, appointments, EMR
  - See [docs/healthcare/README.md](docs/healthcare/README.md)
  
- **Insurance Module**: Policies, claims, underwriting
  - See [docs/insurance/README.md](docs/insurance/README.md)
  
- **Finance Module**: General ledger, invoicing, budgeting
  - See [docs/finance/README.md](docs/finance/README.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style

- Follow Go standard formatting (`gofmt`)
- Write meaningful commit messages
- Add tests for new features
- Update documentation as needed
- Use the Makefile for common tasks

### Development Workflow

1. Create a new module or feature
2. Write tests first (TDD approach)
3. Implement the feature
4. Run tests: `make test-[module]`
5. Format code: `make fmt`
6. Run linter: `make lint`
7. Update documentation
8. Submit PR

## 📖 Documentation

### Module Documentation
- [Access Module](docs/access/IMPLEMENTATION_SUMMARY.md) - Authentication & authorization
- [Healthcare Module](docs/healthcare/README.md) - Healthcare management
- [Insurance Module](docs/insurance/README.md) - Insurance operations
- [Finance Module](docs/finance/README.md) - Financial management
- [Banking Module](docs/banking/README.md) - Banking services (planned)

### Platform Documentation
- [Platform Architecture](docs/evero/ARCHITECTURE.md) - Platform architecture and design principles
- [Deployment Guide](docs/evero/DEPLOYMENT.md) - Comprehensive deployment instructions
- [Infrastructure Guide](docs/evero/INFRASTRUCTURE.md) - Shared infrastructure components
- [Postman Testing Guide](docs/evero/POSTMAN_TESTING_GUIDE.md) - API testing with Postman
- [SSO Migration Summary](docs/SSO_MIGRATION_SUMMARY.md) - Access module migration details

### Quick References
- [Access Quick Reference](docs/access/QUICK_REFERENCE.md)
- [Makefile Commands](#-getting-started) - Use `make help`

## 📄 License

This project is private and proprietary.

## 👥 Authors

- **Prayas Poudel** - [@prayaspoudel](https://github.com/prayaspoudel)

## 🙏 Acknowledgments

- Built with Go and the amazing Go ecosystem
- Inspired by clean architecture principles
- Designed for scalability and maintainability

## 📞 Support

For questions or issues, please contact the development team or open an issue in the repository.

---

**Note**: This is a private repository. Please ensure you have proper authorization before accessing or using this codebase.
