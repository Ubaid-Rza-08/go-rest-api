# Go Microservices API — JWT Authentication + Post Service

A production-structured Golang microservices project using:

- Gin
- PostgreSQL
- JWT Authentication
- bcrypt Password Hashing
- Layered Architecture
- Repository Pattern
- Middleware
- REST APIs
- Microservice Communication using Shared JWT

---

# Architecture

```text
Client
   │
   ├── Auth Service (:8080)
   │      ├── Register
   │      ├── Login
   │      ├── JWT Generation
   │      └── User Management
   │
   └── Post Service (:8081)
          ├── JWT Validation
          ├── Create Post
          ├── Update Post
          ├── Delete Post
          └── Fetch Posts
```

---

# Services

| Service | Port | Responsibility |
|---|---|---|
| Auth Service | 8080 | Authentication & User Management |
| Post Service | 8081 | Post CRUD Operations |

---

# Tech Stack

| Technology | Purpose |
|---|---|
| Golang | Backend Language |
| Gin | HTTP Framework |
| PostgreSQL | Database |
| JWT | Authentication |
| bcrypt | Password Hashing |
| sqlx | Database Queries |
| UUID | Unique IDs |
| Middleware | Authorization & Logging |

---

# Features

## Auth Service

- User Registration
- User Login
- JWT Token Generation
- Password Hashing with bcrypt
- Role-Based Access
- Protected Routes
- User CRUD
- Profile Endpoint

---

## Post Service

- JWT Validation
- Create Posts
- Update Posts
- Delete Posts
- Get Single Post
- Get All Posts
- Ownership Authorization
- Soft Delete Support

---

# Project Structure

## Auth Service

```text
go-rest-api/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── utils/
│
├── migrations/
│   └── 001_init.sql
│
├── .env
├── go.mod
└── README.md
```

---

## Post Service

```text
post-service/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── service/
│   └── utils/
│
├── migrations/
│   └── 001_init.sql
│
├── .env
├── go.mod
└── README.md
```

---

# Database Setup

## Create Database

```sql
CREATE DATABASE go_rest_api;
```

---

# Run User Service Migration

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('admin', 'user');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name VARCHAR(100) NOT NULL,

    email VARCHAR(255) NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,

    role user_role NOT NULL DEFAULT 'user',

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ
);
```

---

# Run Post Service Migration

```sql
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL,

    title VARCHAR(255) NOT NULL,

    content TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ
);
```

---

# Environment Variables

## Auth Service `.env`

```env
SERVER_PORT=8080

GIN_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=go_rest_api
DB_SSLMODE=disable

JWT_SECRET=your-super-secret-key-change-in-production

JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

APP_ENV=development

BCRYPT_COST=12
```

---

## Post Service `.env`

```env
PORT=8081

DB_URL=postgres://postgres:yourpassword@localhost:5432/go_rest_api?sslmode=disable

JWT_SECRET=your-super-secret-key-change-in-production
```

---

# Install Dependencies

## Auth Service

```bash
go mod tidy
```

---

## Post Service

```bash
go mod tidy
```

---

# Run Services

## Start Auth Service

```bash
go run cmd/server/main.go
```

Runs on:

```text
http://localhost:8080
```

---

## Start Post Service

```bash
go run cmd/server/main.go
```

Runs on:

```text
http://localhost:8081
```

---

# API Endpoints

# Auth Service

## Register User

```http
POST /api/v1/auth/register
```

### Request

```json
{
  "name": "Ubaid",
  "email": "ubaid@example.com",
  "password": "Secret@123"
}
```

---

## Login User

```http
POST /api/v1/auth/login
```

### Request

```json
{
  "email": "ubaid@example.com",
  "password": "Secret@123"
}
```

### Response

```json
{
  "message": "login successful",
  "data": {
    "token": "JWT_TOKEN"
  }
}
```

---

## Get Profile

```http
GET /api/v1/profile
```

### Headers

```http
Authorization: Bearer JWT_TOKEN
```

---

# Post Service

## Create Post

```http
POST /api/v1/posts
```

### Headers

```http
Authorization: Bearer JWT_TOKEN
Content-Type: application/json
```

### Body

```json
{
  "title": "Learning Golang Microservices",
  "content": "JWT authentication between services is working."
}
```

---

## Get All Posts

```http
GET /api/v1/posts
```

---

## Get Post By ID

```http
GET /api/v1/posts/{id}
```

---

## Update Post

```http
PUT /api/v1/posts/{id}
```

### Body

```json
{
  "title": "Updated Title",
  "content": "Updated Content"
}
```

---

## Delete Post

```http
DELETE /api/v1/posts/{id}
```

---

# JWT Authentication Flow

1. User logs into Auth Service
2. Auth Service generates JWT token
3. Client sends JWT to Post Service
4. Post Service validates JWT
5. Protected operations are allowed

---

# Security Features

- JWT Authentication
- Password Hashing with bcrypt
- Protected Routes
- Ownership Authorization
- Soft Deletes
- Middleware-Based Security
- UUID Primary Keys

---

# Example Flow

## Step 1 — Register

```http
POST :8080/api/v1/auth/register
```

---

## Step 2 — Login

```http
POST :8080/api/v1/auth/login
```

Copy JWT token.

---

## Step 3 — Create Post

```http
POST :8081/api/v1/posts
Authorization: Bearer JWT_TOKEN
```

---

# Future Improvements

- Swagger/OpenAPI Documentation
- Docker Support
- Docker Compose
- Redis Caching
- API Gateway
- Rate Limiting
- Refresh Tokens
- CI/CD Pipeline
- Unit Testing
- Integration Testing
- Kafka/RabbitMQ
- Kubernetes Deployment
- Observability & Metrics

---

# Important Notes

This project is designed for:

- Learning Golang Backend Development
- Understanding JWT Authentication
- Understanding Microservice Architecture
- Practicing Repository Pattern
- Practicing Layered Architecture
- Learning Service-to-Service Authentication

This is not fully production-ready yet.

Missing:
- Distributed tracing
- Centralized logging
- Rate limiting
- Full test coverage
- API gateway
- Container orchestration
- Secret management
- Monitoring

---

# Author

Ubaid Rza

---

# License

MIT License
