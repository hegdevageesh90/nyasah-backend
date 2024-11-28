# Nyasah AI Backend - Multi-tenant Social Proofing Platform

This is the backend implementation for the Nyasah AI multi-tenant social proofing platform.

## Features

- Multi-tenant architecture
- Tenant management and configuration
- User authentication per tenant
- Generic entity management
- Review and social proof tracking
- Analytics per tenant

## Getting Started

1. Install Go (1.21 or later)
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

### Tenant Management (Admin)
- POST /api/admin/tenants - Create a new tenant
- GET /api/admin/tenants/:id - Get tenant details
- PUT /api/admin/tenants/:id - Update tenant settings

### Per-Tenant Authentication
- POST /api/auth/register - Register a new user
- POST /api/auth/login - Login and get JWT token

### Reviews
- POST /api/reviews - Create a new review
- GET /api/reviews - List all reviews
- GET /api/reviews/:id - Get a specific review

### Social Proof
- POST /api/social-proof - Create new social proof
- GET /api/social-proof - List all social proofs
- GET /api/social-proof/analytics - Get social proof analytics

## Testing

To test the API endpoints, you can use the following curl commands:

1. Create a new tenant (admin):
```bash
curl -X POST http://localhost:8080/api/admin/tenants \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Store",
    "domain": "example-store.com",
    "type": "ecommerce",
    "settings": {
      "allowedReviewTypes": ["product", "service"],
      "moderationEnabled": true
    }
  }'
```

2. Register a user (tenant-specific):
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "Test User"
  }'
```

3. Create a review (tenant-specific):
```bash
curl -X POST http://localhost:8080/api/reviews \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "entity_id": "ENTITY_UUID",
    "rating": 5,
    "content": "Great product!",
    "metadata": {
      "verified_purchase": true
    }
  }'
```

4. Get social proof analytics (tenant-specific):
```bash
curl -X GET http://localhost:8080/api/social-proof/analytics \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```