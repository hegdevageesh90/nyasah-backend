# Nyasah AI Backend - Multi-tenant Social Proofing Platform

This is the backend implementation for the Nyasah AI multi-tenant social proofing platform.

## Features

- Multi-tenant architecture
- Tenant management and configuration
- User authentication per tenant
- Generic entity management
- Review and social proof tracking
- Analytics per tenant
- Multiple AI Provider Support (OpenAI, Claude, HuggingFace, Llama)

## Getting Started

1. Install Go (1.21 or later)
2. Clone the repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Set up environment variables in `.env`:
   ```
   PORT=8080
   JWT_SECRET=your-secure-secret-key
   DATABASE_URL=nyasah.db
   OPENAI_API_KEY=your-openai-api-key
   CLAUDE_API_KEY=your-claude-api-key
   HUGGINGFACE_API_KEY=your-huggingface-api-key
   LLAMA_SERVER_URL=http://localhost:8000
   ```
5. Run the server:
   ```bash
   go run main.go
   ```

## API Documentation

### Authentication

#### Register User
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

#### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Reviews

#### Create Review
```bash
curl -X POST http://localhost:8080/api/reviews \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "PRODUCT_UUID",
    "rating": 5,
    "content": "Great product!",
    "verified": true
  }'
```

#### List Reviews
```bash
curl -X GET http://localhost:8080/api/reviews \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```

#### Get Review
```bash
curl -X GET http://localhost:8080/api/reviews/REVIEW_UUID \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```

### Social Proof

#### Create Social Proof
```bash
curl -X POST http://localhost:8080/api/social-proof \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "purchase",
    "product_id": "PRODUCT_UUID",
    "content": "John D. just purchased this item!",
    "media_type": "text"
  }'
```

#### Get Analytics
```bash
curl -X GET http://localhost:8080/api/social-proof/analytics \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```

### AI Features

#### Query AI
```bash
curl -X POST http://localhost:8080/api/ai/query \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "What are the top trending products this week?"
  }'
```

#### Get Product Insights
```bash
curl -X GET http://localhost:8080/api/ai/insights/product/PRODUCT_UUID \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```

#### Get Recommendations
```bash
curl -X GET http://localhost:8080/api/ai/insights/recommendations \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```

#### Get Trend Analysis
```bash
curl -X GET http://localhost:8080/api/ai/insights/trends \
  -H "X-API-Key: TENANT_API_KEY" \
  -H "Authorization: Bearer USER_TOKEN"
```

## Postman Collection

[Download Postman Collection](./nyasah_api.json)

## AI Provider Configuration

The platform supports multiple AI providers:

1. OpenAI (GPT-3.5/4)
2. HuggingFace Models
3. Local Llama Deployment

Configure your preferred provider in the environment variables.