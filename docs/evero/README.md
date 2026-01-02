# Evero Platform Documentation

This directory contains platform-level technical documentation for the Evero modular monolith application.

## Documentation Index

### Core Documentation

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Platform architecture, design principles, module structure, data flow, and scalability strategies
- **[DEPLOYMENT.md](DEPLOYMENT.md)** - Comprehensive deployment guide covering local development, Docker, production deployments, and monitoring
- **[INFRASTRUCTURE.md](INFRASTRUCTURE.md)** - Shared infrastructure components (database, cache, logger, message broker, router, validator)
- **[POSTMAN_TESTING_GUIDE.md](POSTMAN_TESTING_GUIDE.md)** - API testing guide using Postman

## Quick Links

### Module Documentation
- [Access Module](../access/README.md) - Authentication, authorization, user management
- [Healthcare Module](../healthcare/README.md) - Patient management, appointments, medical records
- [Insurance Module](../insurance/README.md) - Policy management, claims processing, underwriting
- [Finance Module](../finance/README.md) - General ledger, accounts payable/receivable, budgeting
- [Banking Module](../banking/README.md) - Account management, transactions, loans (planned)

### Getting Started
1. Read [ARCHITECTURE.md](ARCHITECTURE.md) to understand the platform design
2. Follow [DEPLOYMENT.md](DEPLOYMENT.md) to set up your development environment
3. Review [INFRASTRUCTURE.md](INFRASTRUCTURE.md) to learn about shared components
4. Refer to module-specific docs for detailed implementation guides

## Platform Overview

Evero is a modular monolith platform that supports multiple business domains:
- **Modular Structure**: Clear boundaries between modules
- **Shared Infrastructure**: Common database, cache, logging, and routing components
- **Independent Deployment**: Modules can be deployed separately
- **Unified Codebase**: All modules in a single repository

## Architecture Highlights

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

## Key Technologies

- **Language**: Go 1.24+
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7.0+
- **HTTP Framework**: Gin
- **Message Broker**: RabbitMQ/Kafka
- **Container**: Docker

## Directory Structure

```
docs/evero/
├── README.md                    # This file
├── ARCHITECTURE.md              # Platform architecture guide
├── DEPLOYMENT.md                # Deployment and operations guide
├── INFRASTRUCTURE.md            # Infrastructure components guide
└── POSTMAN_TESTING_GUIDE.md     # API testing guide
```

## Contributing

When adding new platform-level documentation:
1. Place technical docs in this directory (`docs/evero/`)
2. Place module-specific docs in `docs/[module]/`
3. Update this README with links to new documentation
4. Follow the existing documentation structure and formatting

## Support

For questions or assistance:
- Review the relevant documentation first
- Check module-specific docs for feature details
- Contact the platform team for architecture questions
- Refer to deployment docs for operational issues

---

**Last Updated**: January 2024  
**Maintained By**: Evero Platform Team
