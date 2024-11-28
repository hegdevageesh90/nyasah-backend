# Nyasah AI Backend

This is the backend implementation for the Nyasah AI social proofing platform.

## Features

- User authentication (register/login)
- Review management
- Social proof tracking
- Analytics

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

### Authentication
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

1. Register a new user:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

2. Login:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

3. Create a review (replace TOKEN with the JWT token from login):
```bash
curl -X POST http://localhost:8080/api/reviews \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id":"PRODUCT_UUID","rating":5,"content":"Great product!"}'
```

4. Get social proof analytics:
```bash
curl -X GET http://localhost:8080/api/social-proof/analytics \
  -H "Authorization: Bearer TOKEN"
```