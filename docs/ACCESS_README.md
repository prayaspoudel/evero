# Access Module - Authentication & Authorization

The Access module is Evero's integrated SSO (Single Sign-On) authentication and authorization service, migrated from the standalone SSO project.

## 📍 Location in Evero

```
evero/
├── modules/access/          # Access module implementation
│   ├── entity/             # Database entities (User, Session, OAuth, etc.)
│   ├── model/              # Request/Response DTOs
│   ├── repository/         # Data access layer
│   ├── features/           # Business logic (auth use cases)
│   ├── delivery/           # HTTP controllers and routes
│   ├── middleware/         # Authentication middleware
│   └── app/                # Bootstrap and setup
├── config/access/          # Configuration files (local, dev, prod)
├── database/access/        # Migrations and schemas
├── deployment/access/      # Docker files
├── docs/access/            # Detailed documentation
└── app/access/             # Entry point (main.go)
```

## 🚀 Quick Start

### Running the Access Module

```bash
# Build the access module
go build -o bin/access app/access/main.go

# Run the access module
./bin/access
```

The service will start on port 8080 (configurable in `config/access/local.json`).

### Using Docker

```bash
cd deployment/access
docker-compose up -d
```

## 🔑 Key Features

- 🔐 **JWT Authentication**: Access and refresh token management
- 🔄 **Token Rotation**: Automatic refresh token rotation
- 👥 **Multi-Company Support**: Users can belong to multiple organizations
- 📊 **Session Management**: Track active sessions across devices
- 🔍 **Audit Logging**: Comprehensive security audit trail
- 🔑 **Password Management**: Reset, change, and verification
- ✉️ **Email Verification**: Built-in verification system
- 🛡️ **Two-Factor Auth**: TOTP and SMS 2FA support
- 🔌 **OAuth2 Support**: Standard OAuth2 authorization flows

## 📚 API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login user |
| POST | `/api/auth/logout` | Logout user (requires auth) |
| POST | `/api/auth/refresh` | Refresh access token |
| GET | `/health` | Health check |

### Example: Register User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure123",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Example: Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secure123"
  }'
```

Response:
```json
{
  "status": "success",
  "data": {
    "accessToken": "eyJhbGc...",
    "refreshToken": "def...",
    "expiresIn": 3600,
    "tokenType": "Bearer",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "firstName": "John",
      "lastName": "Doe"
    },
    "companies": []
  }
}
```

## ⚙️ Configuration

Configuration files are located in `config/access/`:

- `local.json` - Local development
- `development.json` - Development environment
- `production.json` - Production environment

### Example Configuration

```json
{
  "app": {
    "name": "evero-access",
    "environment": "local"
  },
  "web": {
    "host": "0.0.0.0",
    "port": 8080,
    "prefork": false
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "postgres",
    "password": "postgres",
    "database": "evero_access",
    "pool": {
      "min": 5,
      "max": 20,
      "idle": 10
    }
  },
  "jwt": {
    "secret": "your-secret-key-change-in-production",
    "expiration": 3600
  }
}
```

## 🗄️ Database

### Running Migrations

```bash
# Connect to database
psql -U postgres -d evero_access

# Run migrations
\i database/access/migrations/001_sso_schema.up.sql
```

### Tables

The access module uses 14 tables with `sso_` prefix:

- `sso_users` - User accounts
- `sso_companies` - Organizations/tenants
- `sso_user_companies` - User-company relationships
- `sso_sessions` - Active sessions
- `sso_refresh_tokens` - Refresh tokens
- `sso_user_two_factors` - 2FA settings
- `sso_backup_codes` - 2FA backup codes
- `sso_oauth_clients` - OAuth2 clients
- `sso_oauth_authorization_codes` - OAuth2 auth codes
- `sso_oauth_tokens` - OAuth2 tokens
- `sso_password_reset_tokens` - Password reset tokens
- `sso_email_verification_tokens` - Email verification tokens
- `sso_audit_logs` - Security audit trail
- `sso_account_lockouts` - Account lockout tracking

## 🔗 Integration with Other Modules

### Using the Auth Middleware

```go
import (
    "github.com/prayaspoudel/modules/access/middleware"
    "github.com/prayaspoudel/modules/access/features/auth"
)

// In your module's route setup
authMiddleware := middleware.NewAuthMiddleware(authUseCase)
protectedRoute.Use(authMiddleware.Authenticate)

// Get authenticated user
authCtx := middleware.GetAuth(ctx)
userID := authCtx.UserID
email := authCtx.Email
```

## 📖 Detailed Documentation

For comprehensive documentation, see:

- [docs/access/](../docs/access/) - Complete implementation guides
- [docs/access/QUICK_REFERENCE.md](../docs/access/QUICK_REFERENCE.md) - API quick reference
- [docs/access/MIGRATION_GUIDE.md](../docs/access/MIGRATION_GUIDE.md) - Migration from standalone SSO

## 🔒 Security

- **JWT Secret**: Always change the JWT secret in production
- **Database**: Use SSL/TLS for database connections in production
- **Passwords**: Bcrypt hashing with default cost
- **Rate Limiting**: Implement rate limiting for auth endpoints
- **CORS**: Configure CORS for your specific domains

## 🧪 Testing

```bash
# Run tests
cd modules/access
go test ./... -v

# Run with coverage
go test ./... -cover
```

## 🚢 Deployment

See [deployment/access/](../deployment/access/) for Docker and deployment configurations.

### Production Checklist

- [ ] Change JWT secret
- [ ] Enable database SSL
- [ ] Configure proper CORS origins
- [ ] Set up rate limiting
- [ ] Enable audit logging
- [ ] Configure email service (for verification/reset)
- [ ] Set up monitoring and alerting
- [ ] Configure backup strategy

## 📞 Support

For issues or questions about the Access module:

1. Check the [docs/access/](../docs/access/) folder
2. Review the migration guide
3. Check Evero main README

## 🔄 Migration from Standalone SSO

This module was migrated from the standalone SSO project. The backend code has been fully integrated into Evero's architecture. The original SSO repository now only contains:

- **admin-dashboard/** - Admin UI for managing users and settings
- **sdk/** - TypeScript SDK for frontend integration

Refer to those repositories for frontend-specific components.
