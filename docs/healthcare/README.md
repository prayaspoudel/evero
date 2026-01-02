# Healthcare Module

## Overview

The Healthcare module provides comprehensive healthcare management functionality for the Evero platform, including patient management, appointments, medical records, and billing integration.

## Features

### Patient Management
- Patient registration and demographics
- Medical history tracking
- Insurance information management
- Emergency contact management
- Patient portal access

### Appointments & Scheduling
- Appointment booking and management
- Provider availability scheduling
- Appointment reminders (email/SMS)
- Waitlist management
- Recurring appointments

### Electronic Medical Records (EMR)
- Clinical notes and documentation
- Diagnosis recording (ICD-10)
- Prescription management
- Lab results integration
- Medical imaging references

### Billing & Claims
- Insurance claim generation
- Payment processing
- Balance tracking
- Integration with finance module

## Architecture

```
app/healthcare/
├── entities/           # Domain models
├── models/            # Database models
├── repositories/      # Data access layer
├── usecases/         # Business logic
├── controllers/      # HTTP handlers
├── middleware/       # Healthcare-specific middleware
└── routes.go         # Route definitions
```

## Database Schema

### Core Tables

1. **healthcare_patients** - Patient demographics and information
2. **healthcare_appointments** - Appointment scheduling
3. **healthcare_providers** - Healthcare provider information
4. **healthcare_medical_records** - Medical records and history
5. **healthcare_prescriptions** - Medication prescriptions
6. **healthcare_diagnoses** - Patient diagnoses
7. **healthcare_vitals** - Patient vital signs
8. **healthcare_insurance** - Insurance information

## API Endpoints

### Patient Management
- `GET /api/v1/healthcare/patients` - List patients
- `POST /api/v1/healthcare/patients` - Create patient
- `GET /api/v1/healthcare/patients/:id` - Get patient details
- `PUT /api/v1/healthcare/patients/:id` - Update patient
- `DELETE /api/v1/healthcare/patients/:id` - Delete patient

### Appointments
- `GET /api/v1/healthcare/appointments` - List appointments
- `POST /api/v1/healthcare/appointments` - Schedule appointment
- `GET /api/v1/healthcare/appointments/:id` - Get appointment
- `PUT /api/v1/healthcare/appointments/:id` - Update appointment
- `DELETE /api/v1/healthcare/appointments/:id` - Cancel appointment

### Medical Records
- `GET /api/v1/healthcare/records/:patientId` - Get patient records
- `POST /api/v1/healthcare/records` - Create medical record
- `GET /api/v1/healthcare/records/:id` - Get record details
- `PUT /api/v1/healthcare/records/:id` - Update record

### Prescriptions
- `GET /api/v1/healthcare/prescriptions/:patientId` - Get prescriptions
- `POST /api/v1/healthcare/prescriptions` - Create prescription
- `GET /api/v1/healthcare/prescriptions/:id` - Get prescription
- `PUT /api/v1/healthcare/prescriptions/:id` - Update prescription

## Configuration

### Environment Variables
```env
HEALTHCARE_DB_HOST=localhost
HEALTHCARE_DB_PORT=5432
HEALTHCARE_DB_NAME=evero
HEALTHCARE_PORT=3001
HEALTHCARE_ENABLE_HIPAA_AUDIT=true
```

### Config Files
- `config/healthcare/local.json`
- `config/healthcare/development.json`
- `config/healthcare/production.json`

## Security & Compliance

### HIPAA Compliance
- All patient data encrypted at rest and in transit
- Comprehensive audit logging of PHI access
- Role-based access control (RBAC)
- Automatic session timeout
- Data retention policies

### Access Control
- Healthcare provider authentication
- Patient consent management
- Emergency access protocols
- Audit trail for all data access

## Integration

### Finance Module
- Automatic billing integration
- Insurance claim submission
- Payment tracking
- Revenue cycle management

### Access Module
- Single sign-on (SSO)
- Role-based permissions
- Multi-factor authentication
- Session management

## Build & Deploy

### Build
```bash
make build-healthcare
```

### Setup
```bash
make setup-healthcare
```

### Run Migrations
```bash
make migrate-healthcare
```

### Deploy
```bash
make deploy-healthcare
```

## Testing

```bash
make test-healthcare
```

## Monitoring

### Key Metrics
- Patient registrations
- Appointment utilization
- Provider efficiency
- Claim processing time
- System uptime

### Health Checks
- Database connectivity
- External API availability
- Storage capacity
- Response times

## Development

### Getting Started

1. Setup the module:
   ```bash
   make setup-healthcare
   ```

2. Run migrations:
   ```bash
   make migrate-healthcare
   ```

3. Build and run:
   ```bash
   make build-healthcare
   ./bin/healthcare
   ```

### Code Structure

- Follow clean architecture principles
- Use dependency injection
- Implement repository pattern
- Write unit and integration tests

## Documentation

- [API Documentation](./API.md)
- [Database Schema](./SCHEMA.md)
- [Integration Guide](./INTEGRATION.md)
- [Security & Compliance](./SECURITY.md)

## Support

For issues or questions:
1. Check module documentation
2. Review error logs
3. Contact healthcare module team

---

**Last Updated:** January 2, 2026  
**Version:** 1.0.0  
**Status:** In Development
