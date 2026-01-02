# Evero Platform Architecture

## Overview

Evero is a modular monolith platform designed to support multiple business domains through a unified codebase while maintaining clear boundaries between modules. The platform follows clean architecture principles and domain-driven design.

## Architecture Principles

### 1. Modular Monolith
- **Single Codebase**: All modules share the same repository and Go module
- **Clear Boundaries**: Each module has its own handlers, services, repositories, and models
- **Shared Infrastructure**: Common infrastructure components (database, cache, logger, etc.) are shared
- **Independent Deployment**: Modules can be deployed independently using module-specific builds

### 2. Clean Architecture
```
┌─────────────────────────────────────────────────────────┐
│                    HTTP/REST API                         │
│                     (Handlers)                           │
├─────────────────────────────────────────────────────────┤
│                   Business Logic                         │
│                    (Services)                            │
├─────────────────────────────────────────────────────────┤
│                  Data Access Layer                       │
│                   (Repositories)                         │
├─────────────────────────────────────────────────────────┤
│                    Infrastructure                        │
│         (Database, Cache, Logger, Message Queue)         │
└─────────────────────────────────────────────────────────┘
```

### 3. Domain-Driven Design
Each module represents a bounded context:
- **Access**: Authentication, authorization, user management, multi-tenancy
- **Healthcare**: Patient management, appointments, medical records, billing
- **Insurance**: Policy management, claims processing, underwriting
- **Finance**: General ledger, accounts payable/receivable, budgeting
- **Banking**: Account management, transactions, deposits, loans (planned)

## Module Structure

```
evero/
├── app/                        # Application layer
│   ├── access/                 # Access module application
│   ├── healthcare/             # Healthcare module application
│   ├── insurance/              # Insurance module application
│   ├── finance/                # Finance module application
│   └── banking/                # Banking module application (planned)
│
├── config/                     # Configuration files
│   ├── access/                 # Access module config
│   ├── healthcare/             # Healthcare module config
│   └── ...
│
├── database/                   # Database layer
│   ├── access/                 # Access module migrations
│   ├── healthcare/             # Healthcare module migrations
│   └── ...
│
├── deployment/                 # Deployment configurations
│   ├── access/                 # Access deployment files
│   ├── healthcare/             # Healthcare deployment files
│   └── ...
│
├── docs/                       # Documentation
│   ├── evero/                  # Platform-level documentation
│   ├── access/                 # Module-specific documentation
│   └── ...
│
├── infrastructure/             # Shared infrastructure
│   ├── cache/                  # Redis cache implementation
│   ├── config/                 # Configuration management
│   ├── database/               # PostgreSQL database
│   ├── logger/                 # Structured logging
│   ├── message-broker/         # RabbitMQ/Kafka integration
│   ├── router/                 # HTTP router setup
│   ├── setup/                  # Infrastructure initialization
│   └── validator/              # Request validation
│
└── modules/                    # Module implementations
    ├── access/                 # Access module
    │   ├── handlers/           # HTTP handlers
    │   ├── services/           # Business logic
    │   ├── repositories/       # Data access
    │   └── models/             # Domain models
    └── ...
```

## Infrastructure Components

### Database
- **PostgreSQL**: Primary database for all modules
- **Per-Module Databases**: Each module has its own database for isolation
- **Migration Management**: SQL migrations in `database/[module]/migrations/`
- **Connection Pooling**: Configured per module

### Cache
- **Redis**: Distributed caching layer
- **Use Cases**: Session storage, rate limiting, temporary data
- **Configuration**: Per-module cache namespacing

### Logger
- **Structured Logging**: JSON-formatted logs
- **Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Context Propagation**: Request ID tracking across services

### Message Broker
- **RabbitMQ/Kafka**: Asynchronous communication
- **Event-Driven**: Inter-module communication via events
- **Queue-Based**: Background job processing

### Router
- **Gin Framework**: HTTP routing and middleware
- **Module Routing**: Each module registers its own routes
- **Middleware Stack**: Authentication, logging, CORS, rate limiting

### Validator
- **go-playground/validator**: Request validation
- **Custom Validators**: Domain-specific validation rules
- **Automatic Binding**: JSON/Form data binding

## Data Flow

### Request Flow
```
1. HTTP Request → Router
2. Router → Middleware (Auth, Logging, Validation)
3. Middleware → Handler (Module-specific)
4. Handler → Service (Business Logic)
5. Service → Repository (Data Access)
6. Repository → Database
7. Database → Repository (Response)
8. Repository → Service
9. Service → Handler
10. Handler → HTTP Response
```

### Event Flow
```
1. Service A → Event Publisher
2. Event Publisher → Message Broker
3. Message Broker → Event Consumer (Service B)
4. Event Consumer → Service B Logic
5. Service B → Repository (if needed)
```

## Module Communication

### Synchronous Communication
- **Direct Service Calls**: Allowed within the same process
- **HTTP APIs**: For external/future microservice communication
- **Shared Repositories**: Access to shared data models (with caution)

### Asynchronous Communication
- **Events**: Domain events published to message broker
- **Event Handlers**: Modules subscribe to relevant events
- **Decoupling**: Modules don't directly depend on each other

## Security Architecture

### Authentication
- **JWT Tokens**: Stateless authentication
- **Session Management**: Redis-backed sessions
- **OAuth Integration**: Support for external providers
- **2FA**: Two-factor authentication support

### Authorization
- **RBAC**: Role-based access control
- **Resource Permissions**: Fine-grained permissions
- **Company Isolation**: Multi-tenant data separation
- **Middleware**: Authorization checks at handler level

### Data Protection
- **Encryption at Rest**: Database encryption
- **Encryption in Transit**: TLS/HTTPS
- **Password Hashing**: bcrypt
- **Secrets Management**: Environment variables, vault integration

## Scalability

### Horizontal Scaling
- **Stateless Handlers**: Support multiple instances
- **Load Balancing**: Distribute requests across instances
- **Session Sharing**: Redis-backed sessions

### Database Scaling
- **Read Replicas**: For read-heavy modules
- **Connection Pooling**: Efficient connection management
- **Query Optimization**: Indexed queries, efficient joins

### Caching Strategy
- **L1 Cache**: In-memory (application level)
- **L2 Cache**: Redis (distributed)
- **Cache Invalidation**: Event-driven invalidation

## Deployment Architecture

### Module Deployment Options

#### Option 1: Single Binary (All Modules)
```bash
# Build all modules into one binary
make build

# Run all modules on different ports
make run-access    # Port 3000
make run-healthcare # Port 3001
make run-insurance  # Port 3002
make run-finance    # Port 3003
```

#### Option 2: Separate Binaries (Per Module)
```bash
# Build specific module
make build-healthcare

# Deploy independently
make deploy-healthcare
```

#### Option 3: Docker Containers
```bash
# Each module in its own container
docker-compose -f deployment/healthcare/docker-compose.yml up
```

### Infrastructure Deployment
- **Database**: PostgreSQL per module
- **Cache**: Shared Redis cluster
- **Message Broker**: Shared RabbitMQ/Kafka cluster
- **Reverse Proxy**: Nginx for routing
- **Load Balancer**: HAProxy/AWS ALB

## Development Workflow

### Local Development
1. Clone repository
2. Set up infrastructure (PostgreSQL, Redis)
3. Run migrations for required modules
4. Start module(s) using Makefile commands

### Testing Strategy
- **Unit Tests**: Service and repository layer tests
- **Integration Tests**: API endpoint tests
- **E2E Tests**: Full workflow tests
- **Load Tests**: Performance testing

### CI/CD Pipeline
1. **Build**: Compile and test all modules
2. **Test**: Run unit, integration, and E2E tests
3. **Security Scan**: Vulnerability scanning
4. **Deploy**: Deploy to staging/production
5. **Verify**: Health checks and smoke tests

## Monitoring and Observability

### Metrics
- **Application Metrics**: Request rate, latency, errors
- **Infrastructure Metrics**: CPU, memory, disk, network
- **Business Metrics**: Module-specific KPIs

### Logging
- **Centralized Logging**: ELK stack or similar
- **Structured Logs**: JSON format for parsing
- **Log Aggregation**: Correlation across modules

### Tracing
- **Distributed Tracing**: OpenTelemetry
- **Request Tracking**: Trace ID propagation
- **Performance Analysis**: Identify bottlenecks

### Health Checks
- **Liveness**: Is the service running?
- **Readiness**: Is the service ready to handle traffic?
- **Dependency Checks**: Database, cache, message broker

## Future Considerations

### Microservices Migration
- **Gradual Extraction**: Move modules to separate services
- **API Gateway**: Centralized routing
- **Service Mesh**: Istio/Linkerd for communication
- **Event Sourcing**: Event-driven architecture

### Multi-Region Support
- **Geographic Distribution**: Deploy in multiple regions
- **Data Replication**: Cross-region database replication
- **Latency Optimization**: Edge caching, CDN

### Advanced Features
- **GraphQL**: Alternative to REST APIs
- **WebSockets**: Real-time communication
- **gRPC**: High-performance inter-service communication
- **Service Discovery**: Consul/etcd for dynamic discovery

## Conclusion

Evero's architecture balances the simplicity of a monolith with the modularity of microservices. This approach provides:
- **Development Speed**: Shared infrastructure, simplified deployment
- **Clear Boundaries**: Module independence, domain isolation
- **Scalability**: Horizontal and vertical scaling options
- **Future-Proof**: Easy migration to microservices if needed

For module-specific architecture details, refer to the respective module documentation in `docs/[module]/`.
