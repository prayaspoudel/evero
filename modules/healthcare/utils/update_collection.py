import json

# Read the collection
with open('modules/healthcare/Evero_Healthcare_API.postman_collection.json', 'r') as f:
    collection = json.load(f)

# Helper function to update request bodies
def update_requests(items):
    for item in items:
        if 'item' in item:
            # It's a folder, recurse
            update_requests(item['item'])
        elif 'request' in item:
            # It's a request
            request = item['request']
            name = item.get('name', '')
            
            # Update Register User - Success to use test.user
            if name == 'Register User - Success':
                if 'body' in request and 'raw' in request['body']:
                    request['body']['raw'] = '{\n  "id": "test.user",\n  "password": "TestPass123!",\n  "name": "Test User"\n}'
            
            # Update Register User - Duplicate ID to use john.doe (seeded user)
            elif name == 'Register User - Duplicate ID':
                if 'body' in request and 'raw' in request['body']:
                    request['body']['raw'] = '{\n  "id": "john.doe",\n  "password": "password123",\n  "name": "John Doe Duplicate"\n}'
                    request['description'] = 'Try to register with duplicate user ID (john.doe already exists in seed data)'
            
            # Update Login - Success with correct seeded password
            elif name == 'Login - Success':
                if 'body' in request and 'raw' in request['body']:
                    request['body']['raw'] = '{\n  "id": "john.doe",\n  "password": "password123"\n}'
            
            # Update Login - Invalid Credentials
            elif name == 'Login - Invalid Credentials':
                if 'body' in request and 'raw' in request['body']:
                    request['body']['raw'] = '{\n  "id": "john.doe",\n  "password": "WrongPassword"\n}'

# Update all items
update_requests(collection['item'])

# Add new folder for testing seeded data
seeded_data_folder = {
    "name": "Seeded Data Tests",
    "item": [
        {
            "name": "Login as john.doe",
            "event": [
                {
                    "listen": "test",
                    "script": {
                        "exec": [
                            "var jsonData = pm.response.json();",
                            "if (jsonData.data && jsonData.data.token) {",
                            "    pm.environment.set(\"auth_token\", jsonData.data.token);",
                            "}"
                        ],
                        "type": "text/javascript"
                    }
                }
            ],
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": '{\n  "id": "john.doe",\n  "password": "password123"\n}'
                },
                "url": {
                    "raw": "{{base_url}}/api/users/_login",
                    "host": ["{{base_url}}"],
                    "path": ["api", "users", "_login"]
                },
                "description": "Login with seeded user john.doe"
            },
            "response": []
        },
        {
            "name": "Login as jane.smith",
            "event": [
                {
                    "listen": "test",
                    "script": {
                        "exec": [
                            "var jsonData = pm.response.json();",
                            "if (jsonData.data && jsonData.data.token) {",
                            "    pm.environment.set(\"auth_token\", jsonData.data.token);",
                            "}"
                        ],
                        "type": "text/javascript"
                    }
                }
            ],
            "request": {
                "method": "POST",
                "header": [
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": '{\n  "id": "jane.smith",\n  "password": "SecurePass456!"\n}'
                },
                "url": {
                    "raw": "{{base_url}}/api/users/_login",
                    "host": ["{{base_url}}"],
                    "path": ["api", "users", "_login"]
                },
                "description": "Login with seeded user jane.smith"
            },
            "response": []
        },
        {
            "name": "Get Seeded Contact - Alice Johnson",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "{{auth_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/contacts/550e8400-e29b-41d4-a716-446655440001",
                    "host": ["{{base_url}}"],
                    "path": ["api", "contacts", "550e8400-e29b-41d4-a716-446655440001"]
                },
                "description": "Get Alice Johnson (seeded contact for john.doe)"
            },
            "response": []
        },
        {
            "name": "Get Seeded Contact - Bob Williams",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "{{auth_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/contacts/550e8400-e29b-41d4-a716-446655440002",
                    "host": ["{{base_url}}"],
                    "path": ["api", "contacts", "550e8400-e29b-41d4-a716-446655440002"]
                },
                "description": "Get Bob Williams (seeded contact for john.doe)"
            },
            "response": []
        },
        {
            "name": "Get Seeded Contact - Charlie Brown",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "{{auth_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/contacts/550e8400-e29b-41d4-a716-446655440003",
                    "host": ["{{base_url}}"],
                    "path": ["api", "contacts", "550e8400-e29b-41d4-a716-446655440003"]
                },
                "description": "Get Charlie Brown (seeded contact for john.doe)"
            },
            "response": []
        },
        {
            "name": "List Addresses for Alice Johnson",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "{{auth_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/contacts/550e8400-e29b-41d4-a716-446655440001/addresses",
                    "host": ["{{base_url}}"],
                    "path": ["api", "contacts", "550e8400-e29b-41d4-a716-446655440001", "addresses"]
                },
                "description": "List all addresses for Alice Johnson (should return 2 addresses)"
            },
            "response": []
        },
        {
            "name": "Get Seeded Address - Alice's NY Address",
            "request": {
                "method": "GET",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "{{auth_token}}"
                    }
                ],
                "url": {
                    "raw": "{{base_url}}/api/contacts/550e8400-e29b-41d4-a716-446655440001/addresses/660e8400-e29b-41d4-a716-446655440001",
                    "host": ["{{base_url}}"],
                    "path": ["api", "contacts", "550e8400-e29b-41d4-a716-446655440001", "addresses", "660e8400-e29b-41d4-a716-446655440001"]
                },
                "description": "Get Alice's NY address (123 Main Street, New York)"
            },
            "response": []
        },
        {
            "name": "Update Seeded Contact - Alice Johnson",
            "request": {
                "method": "PUT",
                "header": [
                    {
                        "key": "Authorization",
                        "value": "{{auth_token}}"
                    },
                    {
                        "key": "Content-Type",
                        "value": "application/json"
                    }
                ],
                "body": {
                    "mode": "raw",
                    "raw": '{\n  "first_name": "Alice Updated",\n  "last_name": "Johnson Updated",\n  "email": "alice.updated@example.com",\n  "phone": "+1-555-9999"\n}'
                },
                "url": {
                    "raw": "{{base_url}}/api/contacts/550e8400-e29b-41d4-a716-446655440001",
                    "host": ["{{base_url}}"],
                    "path": ["api", "contacts", "550e8400-e29b-41d4-a716-446655440001"]
                },
                "description": "Update Alice Johnson's information"
            },
            "response": []
        }
    ]
}

# Insert the new folder at the beginning
collection['item'].insert(0, seeded_data_folder)

# Write the updated collection
with open('modules/healthcare/Evero_Healthcare_API.postman_collection.json', 'w') as f:
    json.dump(collection, f, indent='\t')

print("✅ Postman collection updated successfully!")
print("✅ Added 'Seeded Data Tests' folder with 8 pre-configured tests")
print("✅ Updated login credentials to match seeded passwords")
print("✅ Register User now uses test.user instead of john.doe")
