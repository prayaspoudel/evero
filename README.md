# Evero

Evero is a multi-domain, modular Go backend platform designed for building enterprise-grade applications across healthcare, insurance, and banking sectors. Built with clean architecture principles, it provides a robust foundation for domain-specific services with shared infrastructure components.

## 🏗️ Architecture

Evero follows a **modular monolith architecture** with clear separation of concerns:

- **Multi-Domain Support**: Healthcare, Insurance, and Banking modules
- **Clean Architecture**: Domain-driven design with layers (Entity, Use Case, Repository, Delivery)
- **Shared Infrastructure**: Reusable components across all modules
- **Event-Driven**: Kafka integration for asynchronous messaging
- **API-First**: RESTful APIs with comprehensive validation

## 🚀 Features

### Core Infrastructure
- **Configuration Management**: Environment-specific configurations (local, development, staging, production)
- **Database Support**: PostgreSQL with GORM ORM
- **Caching**: Redis integration for high-performance caching
- **Message Broker**: Kafka and RabbitMQ support
- **Logging**: Structured logging with Logrus and Zap
- **Validation**: Request validation with go-playground/validator
- **Routing**: Multiple router support (Fiber, Gorilla Mux, Gin)

### Healthcare Module

The healthcare module provides comprehensive patient and contact management capabilities:

#### Features
- **User Management**
  - User registration and authentication
  - JWT-based authentication
  - Session management
  - User profile updates

- **Contact Management**
  - Create, read, update, and delete contacts
  - Contact search and pagination
  - Email and phone validation
  - User-specific contact isolation

- **Address Management**
  - Multiple addresses per contact
  - Full CRUD operations
  - Address validation
  - Hierarchical data structure (User → Contact → Address)

### Access Module (SSO)

The access module provides comprehensive authentication and authorization services:

#### Features
- **Authentication**
  - User registration and login
  - JWT-based access and refresh tokens
  - Automatic token rotation
  - Session management

- **Authorization**
  - Multi-company/tenant support
  - OAuth2 authorization flows
  - Two-factor authentication (TOTP, SMS)
  - Role-based access control

- **Security**
  - Password reset and verification
  - Email verification
  - Audit logging
  - Account lockout protection

**📖 Documentation**: See [docs/ACCESS_README.md](docs/ACCESS_README.md)

**Quick Start**:
```bash
# Build and run the access module
go build -o bin/access app/access/main.go
./bin/access
```

**API**: http://localhost:8080 (configurable)

#### API Endpoints

**User Management**
```
POST   /api/register       - Register new user
POST   /api/login          - User login
GET    /api/users/current  - Get current user
PATCH  /api/users/current  - Update current user
POST   /api/logout         - User logout
```

**Contact Management**
```
POST   /api/contacts               - Create contact
GET    /api/contacts               - List contacts (with pagination)
GET    /api/contacts/:id           - Get contact by ID
PATCH  /api/contacts/:id           - Update contact
DELETE /api/contacts/:id           - Delete contact
```

**Address Management**
```
POST   /api/contacts/:contactId/addresses           - Create address
GET    /api/contacts/:contactId/addresses           - List addresses
GET    /api/contacts/:contactId/addresses/:id       - Get address
PATCH  /api/contacts/:contactId/addresses/:id       - Update address
DELETE /api/contacts/:contactId/addresses/:id       - Delete address
```

## 📁 Project Structure

```
evero/
├── app/                          # Application entry points
│   ├── healthcare/               # Healthcare application
│   ├── insurance/                # Insurance application
│   ├── banking/                  # Banking application
│   └── access/                   # Access/SSO application
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
├── modules/                      # Domain modules
│   ├── healthcare/
│   │   ├── app/                  # Application setup
│   │   ├── delivery/             # HTTP handlers/controllers
│   │   ├── entity/               # Domain entities
│   │   ├── features/             # Business logic (use cases)
│   │   ├── gateway/              # External integrations (Kafka, etc.)
│   │   ├── model/                # Request/Response models
│   │   ├── repository/           # Data access layer
│   │   └── test/                 # Unit and integration tests
│   ├── access/                   # Authentication & Authorization
│   │   ├── app/                  # Application setup
│   │   ├── delivery/             # HTTP controllers
│   │   ├── entity/               # User, Session, OAuth entities
│   │   ├── features/             # Auth use cases
│   │   ├── middleware/           # Auth middleware
│   │   ├── model/                # Request/Response DTOs
│   │   └── repository/           # Data access layer
│   └── user/                     # User management module
│
├── config/                       # Configuration files
│   ├── healthcare/
│   ├── insurance/
│   ├── banking/
│   └── access/                   # Access module configs
│
├── database/                     # Database migrations and seeds
│   ├── healthcare/
│   ├── insurance/
│   ├── banking/
│   └── access/                   # SSO database schemas
│
├── deployment/                   # Deployment configurations
│   └── access/                   # Docker files for access module
│
├── packages/                     # External service integrations
│   ├── lib/                      # Shared libraries
│   ├── sendgrid/                 # Email service
│   └── twilio/                   # SMS service
│
└── docs/                         # Documentation
    ├── ACCESS_README.md          # Access module documentation
    └── access/                   # Detailed access module docs
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
- **PostgreSQL**: 12 or higher
- **Redis**: 6 or higher (optional, for caching)
- **Kafka**: 2.8 or higher (optional, for event streaming)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/prayaspoudel/evero.git
   cd evero
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up the database**
   
   For Healthcare module:
   ```bash
   cd database/healthcare
   go run migrate.go
   ```

4. **Configure the application**
   
   Update the configuration file for your environment:
   ```
   config/healthcare/local.json
   ```

5. **Run the application**
   
   For Healthcare module:
   ```bash
   go run app/healthcare/main.go
   ```
   
   The server will start on `http://localhost:3000` (default)

## 📝 Configuration

Evero uses environment-specific JSON configuration files. Each module has its own configuration directory:

```
config/
├── healthcare/
│   ├── local.json          # Local development
│   ├── development.json    # Development environment
│   ├── stage.json          # Staging environment
│   └── production.json     # Production environment
```

### Sample Configuration Structure

```json
{
  "app": {
    "name": "Evero Healthcare API"
  },
  "web": {
    "port": 3000,
    "prefork": false
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "postgres",
    "name": "healthcare_db",
    "pool": {
      "idle": 10,
      "max": 100,
      "lifetime": 300
    }
  },
  "kafka": {
    "bootstrap.servers": "localhost:9092",
    "producer.enabled": false,
    "group.id": "healthcare-service"
  },
  "log": {
    "level": "info"
  }
}
```

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for a specific module
go test ./modules/healthcare/...

# Run tests with coverage
go test -cover ./...
```

### Postman Testing

The Healthcare module includes a comprehensive Postman collection with 50+ tests:

1. **Import the collection**
   ```
   modules/healthcare/Evero_Healthcare_API.postman_collection.json
   ```

2. **Start the server**
   ```bash
   go run app/healthcare/main.go
   ```

3. **Run the tests**
   - Start with "Seeded Data Tests" folder
   - Test user credentials are pre-seeded in the database
   - All environment variables are auto-populated

For detailed testing instructions, see [POSTMAN_TESTING_GUIDE.md](POSTMAN_TESTING_GUIDE.md)

## 🏛️ Infrastructure Components

### Configuration Manager
Centralized configuration loading with support for:
- Environment-specific files
- Module-specific overrides
- Type-safe access methods
- Hot-reload capability

See [infrastructure/config/README.md](infrastructure/config/README.md)

### Cache Manager
Multi-backend caching support:
- Redis
- In-memory cache
- Factory pattern for easy switching

See [infrastructure/cache/README.md](infrastructure/cache/README.md)

### Message Broker
Event-driven messaging with:
- Kafka producer/consumer
- RabbitMQ support
- Async event publishing

See [infrastructure/message-broker/README.md](infrastructure/message-broker/README.md)

### Setup Package
Reusable infrastructure bootstrapping:
- Database initialization
- Logger setup
- Validator configuration
- Web framework setup
- Message broker connections

See [infrastructure/setup/README.md](infrastructure/setup/README.md)

## 🔒 Security

- **Authentication**: JWT-based token authentication
- **Password Security**: bcrypt hashing
- **Input Validation**: Comprehensive request validation
- **SQL Injection**: Protected via ORM parameterized queries
- **Authorization**: User-specific data isolation

## 🚀 Deployment

### Docker Support

Each module can be deployed independently using Docker. Configuration files are environment-specific.

### Environment Variables

Set the following environment variable to specify the environment:
```bash
export EVERO_ENV=production  # Options: local, development, stage, production
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
GET /api/contacts?page=1&size=10
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
