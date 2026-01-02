# Evero Deployment Guide

## Overview

This guide covers deployment strategies for the Evero platform, including local development, Docker-based deployment, and production deployment options.

## Prerequisites

### Required Software
- **Go**: 1.24 or higher
- **PostgreSQL**: 15 or higher
- **Redis**: 7.0 or higher (optional, for caching)
- **Docker**: 24.0 or higher (for containerized deployment)
- **Make**: For build automation

### Optional Software
- **RabbitMQ/Kafka**: For message broker functionality
- **Nginx**: For reverse proxy
- **Postman**: For API testing

## Deployment Options

### 1. Local Development Setup

#### Quick Start
```bash
# Clone the repository
git clone <repository-url>
cd evero

# Set up a specific module (e.g., healthcare)
cd deployment/healthcare
./setup.sh

# Choose option 2 for local development
# Follow the interactive prompts
```

#### Manual Setup
```bash
# Install dependencies
go mod download

# Set up database for a module
createdb evero_healthcare_dev
export DATABASE_URL="postgresql://localhost:5432/evero_healthcare_dev?sslmode=disable"

# Run migrations
make migrate-healthcare

# Build the module
make build-healthcare

# Run the module
make run-healthcare
```

### 2. Docker Deployment

Each module has its own Docker setup in `deployment/[module]/`.

#### Single Module Deployment
```bash
# Navigate to module deployment directory
cd deployment/healthcare

# Start with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f healthcare-server

# Stop services
docker-compose down
```

#### Multi-Module Deployment
```bash
# Start all modules (from root)
make docker-up-healthcare
make docker-up-insurance
make docker-up-finance

# Or use a custom docker-compose file
docker-compose -f deployment/docker-compose.all.yml up -d
```

### 3. Production Deployment

#### Environment Configuration

Create environment-specific config files:

**config/healthcare/production.json**
```json
{
  "server": {
    "port": 3001,
    "environment": "production",
    "cors": {
      "allowed_origins": ["https://app.example.com"]
    }
  },
  "database": {
    "host": "postgres.example.com",
    "port": 5432,
    "database": "evero_healthcare",
    "sslmode": "require",
    "max_open_conns": 25,
    "max_idle_conns": 5
  },
  "redis": {
    "host": "redis.example.com",
    "port": 6379,
    "db": 0,
    "password": "${REDIS_PASSWORD}"
  },
  "jwt": {
    "access_token_expiry": "15m",
    "refresh_token_expiry": "7d"
  },
  "logging": {
    "level": "info",
    "format": "json"
  }
}
```

#### Build for Production
```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s" \
  -o bin/healthcare-server \
  ./app/healthcare

# Verify the build
./bin/healthcare-server --version
```

#### Deployment Methods

##### Method 1: Direct Binary Deployment
```bash
# On the server
scp bin/healthcare-server user@server:/opt/evero/
scp -r config/healthcare user@server:/opt/evero/config/

# SSH to server
ssh user@server

# Run with systemd
sudo cp evero-healthcare.service /etc/systemd/system/
sudo systemctl enable evero-healthcare
sudo systemctl start evero-healthcare
```

**evero-healthcare.service**
```ini
[Unit]
Description=Evero Healthcare Service
After=network.target postgresql.service

[Service]
Type=simple
User=evero
WorkingDirectory=/opt/evero
ExecStart=/opt/evero/bin/healthcare-server
Restart=always
RestartSec=10

Environment="CONFIG_PATH=/opt/evero/config/healthcare/production.json"
Environment="DATABASE_URL=postgresql://user:pass@localhost/evero_healthcare"

[Install]
WantedBy=multi-user.target
```

##### Method 2: Docker Production Deployment
```bash
# Build production image
docker build -f deployment/healthcare/Dockerfile \
  --target production \
  -t evero-healthcare:latest .

# Run container
docker run -d \
  --name evero-healthcare \
  --restart unless-stopped \
  -p 3001:3001 \
  -e CONFIG_PATH=/app/config/production.json \
  -e DATABASE_URL=postgresql://... \
  -v /opt/evero/config:/app/config \
  evero-healthcare:latest
```

##### Method 3: Kubernetes Deployment
```yaml
# deployment/healthcare/k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: evero-healthcare
spec:
  replicas: 3
  selector:
    matchLabels:
      app: evero-healthcare
  template:
    metadata:
      labels:
        app: evero-healthcare
    spec:
      containers:
      - name: healthcare
        image: evero-healthcare:latest
        ports:
        - containerPort: 3001
        env:
        - name: CONFIG_PATH
          value: /app/config/production.json
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: evero-secrets
              key: healthcare-db-url
        livenessProbe:
          httpGet:
            path: /health
            port: 3001
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 3001
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: evero-healthcare-service
spec:
  selector:
    app: evero-healthcare
  ports:
  - protocol: TCP
    port: 80
    targetPort: 3001
  type: LoadBalancer
```

## Infrastructure Setup

### Database Setup

#### PostgreSQL Installation
```bash
# Ubuntu/Debian
sudo apt-get install postgresql-15

# macOS
brew install postgresql@15

# Create databases for each module
createdb evero_access
createdb evero_healthcare
createdb evero_insurance
createdb evero_finance
```

#### Database Configuration
```bash
# Create user
createuser evero_user -P

# Grant privileges
psql -c "GRANT ALL PRIVILEGES ON DATABASE evero_healthcare TO evero_user;"
```

#### Run Migrations
```bash
# Using Makefile
make migrate-healthcare

# Or manually
psql -U evero_user -d evero_healthcare -f database/healthcare/migrations/001_initial_schema.sql
```

### Redis Setup

```bash
# Ubuntu/Debian
sudo apt-get install redis-server

# macOS
brew install redis

# Start Redis
redis-server

# Configure Redis (optional)
sudo vim /etc/redis/redis.conf
```

### Reverse Proxy (Nginx)

**nginx.conf**
```nginx
upstream evero_access {
    server localhost:3000;
}

upstream evero_healthcare {
    server localhost:3001;
}

upstream evero_insurance {
    server localhost:3002;
}

upstream evero_finance {
    server localhost:3003;
}

server {
    listen 80;
    server_name api.example.com;

    location /api/v1/access/ {
        proxy_pass http://evero_access/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/healthcare/ {
        proxy_pass http://evero_healthcare/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/v1/insurance/ {
        proxy_pass http://evero_insurance/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api/v1/finance/ {
        proxy_pass http://evero_finance/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Module-Specific Deployment

### Access Module
```bash
# Deploy Access module (SSO, Auth)
cd deployment/access
./setup.sh

# Or with Docker
docker-compose up -d

# Check health
curl http://localhost:3000/health
```

### Healthcare Module
```bash
cd deployment/healthcare
./setup.sh

# Or
make deploy-healthcare

# Verify
curl http://localhost:3001/health
```

### Insurance Module
```bash
cd deployment/insurance
docker-compose up -d

curl http://localhost:3002/health
```

### Finance Module
```bash
cd deployment/finance
make docker-up

curl http://localhost:3003/health
```

## Monitoring and Health Checks

### Health Endpoints

Each module exposes the following health endpoints:

- `GET /health` - Basic health check
- `GET /health/ready` - Readiness check (includes DB connection)
- `GET /health/live` - Liveness check

### Monitoring with Prometheus

**prometheus.yml**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'evero-healthcare'
    static_configs:
      - targets: ['localhost:3001']
    metrics_path: /metrics
```

### Logging

#### Application Logs
```bash
# View logs from systemd service
journalctl -u evero-healthcare -f

# View Docker logs
docker logs -f evero-healthcare

# View logs in Kubernetes
kubectl logs -f deployment/evero-healthcare
```

#### Log Aggregation (ELK Stack)
```bash
# Send logs to Logstash
# Configure in config/[module]/production.json
{
  "logging": {
    "level": "info",
    "format": "json",
    "output": "logstash",
    "logstash_host": "logstash.example.com:5000"
  }
}
```

## Security Best Practices

### SSL/TLS Configuration
```bash
# Generate SSL certificates (Let's Encrypt)
sudo certbot --nginx -d api.example.com

# Or use existing certificates
sudo cp /path/to/cert.pem /etc/ssl/certs/
sudo cp /path/to/key.pem /etc/ssl/private/
```

### Environment Variables
```bash
# Never commit sensitive data
# Use environment variables or secret management

# .env file (local development only)
DATABASE_URL=postgresql://user:pass@localhost/evero_healthcare
JWT_SECRET=your-secret-key
REDIS_PASSWORD=redis-password

# Production: Use secret management
# AWS Secrets Manager, HashiCorp Vault, etc.
```

### Firewall Configuration
```bash
# Allow only necessary ports
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
sudo ufw enable

# Restrict database access
sudo ufw allow from <app-server-ip> to any port 5432
```

## Scaling Strategies

### Horizontal Scaling

#### Load Balancer Configuration
```bash
# Using HAProxy
frontend evero_frontend
    bind *:80
    default_backend evero_healthcare_backend

backend evero_healthcare_backend
    balance roundrobin
    server healthcare1 10.0.1.1:3001 check
    server healthcare2 10.0.1.2:3001 check
    server healthcare3 10.0.1.3:3001 check
```

### Database Scaling

#### Read Replicas
```bash
# Configure PostgreSQL replication
# On primary server
sudo vim /etc/postgresql/15/main/postgresql.conf
# wal_level = replica
# max_wal_senders = 3

# On replica server
# Configure recovery settings
```

#### Connection Pooling (PgBouncer)
```ini
[databases]
evero_healthcare = host=localhost port=5432 dbname=evero_healthcare

[pgbouncer]
listen_addr = 127.0.0.1
listen_port = 6432
auth_type = md5
pool_mode = transaction
max_client_conn = 100
default_pool_size = 20
```

## Troubleshooting

### Common Issues

#### Database Connection Errors
```bash
# Check database is running
sudo systemctl status postgresql

# Test connection
psql -U evero_user -d evero_healthcare -c "SELECT 1;"

# Check connection string
echo $DATABASE_URL
```

#### Port Already in Use
```bash
# Find process using the port
lsof -i :3001

# Kill the process
kill -9 <PID>
```

#### Migration Failures
```bash
# Rollback migration
make rollback-healthcare

# Re-run migration
make migrate-healthcare

# Check migration status
psql -U evero_user -d evero_healthcare -c "SELECT * FROM schema_migrations;"
```

## Rollback Procedures

### Application Rollback
```bash
# Systemd
sudo systemctl stop evero-healthcare
sudo cp bin/healthcare-server.backup bin/healthcare-server
sudo systemctl start evero-healthcare

# Docker
docker-compose down
docker pull evero-healthcare:previous-version
docker-compose up -d

# Kubernetes
kubectl rollout undo deployment/evero-healthcare
```

### Database Rollback
```bash
# Run down migration
make rollback-healthcare

# Or manually
psql -U evero_user -d evero_healthcare -f database/healthcare/migrations/down/002_rollback.sql
```

## Backup and Recovery

### Database Backup
```bash
# Automated backup script
#!/bin/bash
BACKUP_DIR=/opt/backups
DATE=$(date +%Y%m%d_%H%M%S)

pg_dump -U evero_user evero_healthcare | gzip > $BACKUP_DIR/evero_healthcare_$DATE.sql.gz

# Keep only last 7 days
find $BACKUP_DIR -name "evero_healthcare_*.sql.gz" -mtime +7 -delete
```

### Restore from Backup
```bash
# Restore database
gunzip < evero_healthcare_20240115_120000.sql.gz | psql -U evero_user evero_healthcare
```

## Performance Optimization

### Database Optimization
```sql
-- Add indexes
CREATE INDEX idx_patients_created_at ON patients(created_at);
CREATE INDEX idx_appointments_patient_id ON appointments(patient_id);

-- Analyze tables
ANALYZE patients;
ANALYZE appointments;
```

### Application Optimization
```bash
# Enable Go profiling
import _ "net/http/pprof"

# CPU profiling
go tool pprof http://localhost:3001/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:3001/debug/pprof/heap
```

## Conclusion

This deployment guide provides comprehensive instructions for deploying Evero modules in various environments. For module-specific details, refer to the README files in `docs/[module]/` and deployment configurations in `deployment/[module]/`.

For assistance, contact the DevOps team or refer to the [Architecture documentation](ARCHITECTURE.md).
