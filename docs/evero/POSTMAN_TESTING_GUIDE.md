# Postman Testing Guide

## 🚀 Quick Start (30 seconds to first API call!)

1. **Import Collection**: `modules/healthcare/Evero_Healthcare_API.postman_collection.json`

2. **Start App**:
   ```bash
   cd /Users/leapfrog/prayas_personal/union-products/evero
   go run app/healthcare/main.go
   ```
   Server starts on http://localhost:3000

3. **Run First Test**: 
   - Open "Seeded Data Tests" folder
   - Click "Login as john.doe"
   - Click Send
   - ✅ Done! Token saved, all other tests ready to run!

4. **Test Seeded Contacts**:
   - Run "Get Seeded Contact - Alice Johnson" → See real contact data
   - Run "List Addresses for Alice Johnson" → See 2 addresses (NY & Brooklyn)

## Database Setup ✅

The database has been migrated and seeded with test data. You're ready to test!

## 📊 Seeded Data Reference

### Users
| Username | Password | Has Contacts |
|----------|----------|--------------|
| john.doe | password123 | ✅ 3 contacts |
| jane.smith | SecurePass456! | ✅ 2 contacts |
| bob.wilson | password123 | ❌ 0 contacts |

### Contacts for john.doe
| Name | Email | Phone | ID | Addresses |
|------|-------|-------|----|-----------| 
| Alice Johnson | alice.johnson@example.com | +1-555-0101 | `550e8400-e29b-41d4-a716-446655440001` | 2 |
| Bob Williams | bob.williams@example.com | +1-555-0102 | `550e8400-e29b-41d4-a716-446655440002` | 1 |
| Charlie Brown | charlie.brown@example.com | +1-555-0103 | `550e8400-e29b-41d4-a716-446655440003` | 1 |

### Contacts for jane.smith
| Name | Email | Phone | ID | Addresses |
|------|-------|-------|----|-----------| 
| David Miller | david.miller@example.com | +1-555-0201 | `550e8400-e29b-41d4-a716-446655440004` | 1 |
| Emma Davis | emma.davis@example.com | +1-555-0202 | `550e8400-e29b-41d4-a716-446655440005` | 1 |

### Addresses
| Contact | Street | City | Province | ID |
|---------|--------|------|----------|----| 
| Alice Johnson | 123 Main Street | New York | NY | `660e8400-e29b-41d4-a716-446655440001` |
| Alice Johnson | 456 Oak Avenue | Brooklyn | NY | `660e8400-e29b-41d4-a716-446655440002` |
| Bob Williams | 789 Pine Road | Boston | MA | `660e8400-e29b-41d4-a716-446655440003` |
| Charlie Brown | 321 Elm Street | Chicago | IL | `660e8400-e29b-41d4-a716-446655440004` |
| David Miller | 555 Beach Boulevard | Miami | FL | `660e8400-e29b-41d4-a716-446655440005` |
| Emma Davis | 777 Valley Drive | Los Angeles | CA | `660e8400-e29b-41d4-a716-446655440006` |

## Postman Collection

Import: `modules/healthcare/Evero_Healthcare_API.postman_collection.json`

### 📁 Collection Structure

```
Evero Healthcare API Collection (50+ tests)
├── 🆕 Seeded Data Tests (8 tests) ⭐ START HERE
│   ├── Login as john.doe
│   ├── Login as jane.smith
│   ├── Get Seeded Contact - Alice Johnson
│   ├── Get Seeded Contact - Bob Williams
│   ├── Get Seeded Contact - Charlie Brown
│   ├── List Addresses for Alice Johnson
│   ├── Get Seeded Address - Alice's NY Address
│   └── Update Seeded Contact - Alice Johnson
│
├── User Management (13 tests)
│   ├── Register User - Success (creates test.user)
│   ├── Register User - Duplicate ID (tests against seeded john.doe)
│   ├── Login - Success (uses seeded credentials)
│   ├── Login - Invalid Credentials
│   ├── Get Current User - Success
│   ├── Get Current User - No Token
│   ├── Update Current User - Success
│   ├── Logout - Success
│   └── ... more edge cases
│
├── Contact Management (15 tests)
│   ├── Create Contact - Success
│   ├── Create Contact - Missing Fields
│   ├── Create Contact - Invalid Email
│   ├── List Contacts - Success
│   ├── List Contacts - With Pagination
│   ├── Get Contact - Success
│   ├── Update Contact - Success
│   ├── Delete Contact - Success
│   └── ... more edge cases
│
└── Address Management (14 tests)
    ├── Create Address - Success
    ├── Create Address - Max Length Exceeded
    ├── List Addresses - Success
    ├── Get Address - Success
    ├── Update Address - Success
    ├── Delete Address - Success
    └── ... more edge cases
```

### Collection Variables
Auto-configured variables:
- `base_url`: http://localhost:3000
- `auth_token`: Auto-saved after login
- `contact_id`: Auto-saved after creating a contact
- `address_id`: Auto-saved after creating an address

## Testing Workflows

## Testing Workflows

### 🎯 Recommended Testing Flows

#### Flow 1: Quick Test with Seeded Data (Fastest!)
1. **Seeded Data Tests** → Login as john.doe
2. **Seeded Data Tests** → Get Seeded Contact - Alice Johnson
3. **Seeded Data Tests** → List Addresses for Alice Johnson
4. **Contact Management** → Create Contact (test with new data)
5. **Contact Management** → List Contacts (see both seeded + new contacts)

#### Flow 2: New User Registration
1. **User Management** → Register User - Success (creates test.user)
2. **User Management** → Login - Success (with test.user credentials)
3. **User Management** → Get Current User - Success
4. **Contact Management** → Create Contact - Success
5. **Contact Management** → List Contacts - Success

#### Flow 3: Complete CRUD Operations
1. Login (seeded or new user)
2. Create Contact → Save ID automatically
3. Update Contact → Verify changes
4. Create Address for Contact → Save ID automatically
5. Update Address → Verify changes
6. Delete Address → Verify deletion
7. Delete Contact → Verify deletion

#### Flow 4: Test All Edge Cases
1. **User Management** → Register User - Duplicate ID (fails against seeded john.doe)
2. **User Management** → Login - Invalid Credentials
3. **User Management** → Get Current User - No Token (401 error)
4. **Contact Management** → Create Contact - Invalid Email
5. **Contact Management** → Get Contact - Invalid UUID
6. **Address Management** → Create Address - Max Length Exceeded

### Detailed Step-by-Step Testing

### 1. Authentication Testing

#### Option A: Use Seeded User (Recommended)
Run: **Seeded Data Tests > Login as john.doe**
- Uses seeded credentials: `john.doe` / `password123`
- ✅ Token auto-saved, ready to test!

#### Option B: Create New User
1. **Register**: **User Management > Register User - Success**
   ```json
   {
     "id": "test.user",
     "password": "TestPass123!",
     "name": "Test User"
   }
   ```

2. **Login**: **User Management > Login - Success**
   ```json
   {
     "id": "john.doe",
     "password": "password123"
   }
   ```
   ✅ The `auth_token` will be automatically saved

3. **Verify**: **User Management > Get Current User - Success**
   - Confirms your token is working

### 2. Contact Management Testing

#### Create Contact
Run: **Contact Management > Create Contact - Success**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@example.com",
  "phone": "+1-555-0123"
}
```
✅ The `contact_id` will be automatically saved

#### List Contacts
Run: **Contact Management > List Contacts - Success**
- Shows all your contacts (including seeded ones)

#### Get, Update, Delete Contact
- **Get**: **Contact Management > Get Contact - Success** (uses saved `contact_id`)
- **Update**: **Contact Management > Update Contact - Success**
- **Delete**: **Contact Management > Delete Contact - Success**

### 3. Address Management Testing

#### Create Address
Run: **Address Management > Create Address - Success**
```json
{
  "street": "123 Main Street",
  "city": "New York",
  "province": "NY",
  "postal_code": "10001",
  "country": "USA"
}
```
✅ The `address_id` will be automatically saved

#### List, Get, Update, Delete Address
- **List**: **Address Management > List Addresses - Success** (for current contact)
- **Get**: **Address Management > Get Address - Success** (uses saved `address_id`)
- **Update**: **Address Management > Update Address - Success**
- **Delete**: **Address Management > Delete Address - Success**

### 4. Edge Case Testing

The collection includes comprehensive edge case tests:

#### Validation Errors
- **Register User - Missing Fields**: Tests required field validation
- **Create Contact - Invalid Email**: Tests email format validation
- **Create Address - Max Length Exceeded**: Tests max length constraints

#### Authentication Errors
- **Get Current User - No Token**: Should return 401 Unauthorized
- **Get Current User - Invalid Token**: Should return authentication error
- **Create Contact - Unauthorized**: Tests protected routes without auth

#### Not Found Errors
- **Get Contact - Not Found**: Tests 404 response
- **Update Contact - Not Found**: Tests updating non-existent resource
- **Delete Contact - Not Found**: Tests deleting non-existent resource

#### Invalid Data
- **Get Contact - Invalid UUID**: Tests UUID format validation
- **Create Address - Invalid Contact ID**: Tests foreign key validation

## Expected Response Format

### Success Response
```json
{
  "data": {
    // Response data here
  }
}
```

### Error Response
```json
{
  "errors": "Error message here"
}
```

## Re-seeding Database

If you need to reset the database to initial state:

```bash
go run database/healthcare/migrate.go
```

This will recreate all test users, contacts, and addresses if they don't exist.

## Troubleshooting

### Issue: Connection Refused
**Solution**: Make sure the app is running on port 3000
```bash
go run app/healthcare/main.go
```

### Issue: Database Connection Failed
**Solution**: Verify PostgreSQL is running and database exists
```bash
# Check PostgreSQL status
brew services list | grep postgresql
# Or
pg_isready

# Create database if needed
createdb healthcare
```

### Issue: Authentication Failed
**Solution**: 
1. Make sure you've logged in successfully
2. Check that `auth_token` is saved in environment variables
3. Token might have expired - login again

### Issue: Contact/Address Not Found
**Solution**:
1. Verify you're using the correct IDs
2. Make sure you're logged in as the owner
3. Re-run the seeding script if data was deleted:
   ```bash
   go run database/healthcare/migrate.go
   ```

---

## Happy Testing! 🚀
