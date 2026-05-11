<<<<<<< HEAD
# Go Microservices API — JWT Authentication + Post Service

A production-structured Golang microservices project using:
=======
# Golang Microservices Architecture — Auth Service + Post Service + API Gateway

A production-structured Golang microservices backend using:
>>>>>>> 0b72c8f (add api gateway)

- Gin
- PostgreSQL
- JWT Authentication
- bcrypt Password Hashing
<<<<<<< HEAD
- Layered Architecture
- Repository Pattern
- Middleware
- REST APIs
- Microservice Communication using Shared JWT
=======
- API Gateway
- Reverse Proxy
- Repository Pattern
- Middleware
- Layered Architecture
- Service Isolation
>>>>>>> 0b72c8f (add api gateway)

---

# Architecture

```text
<<<<<<< HEAD
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
=======
                        ┌─────────────────┐
                        │     Client      │
                        └────────┬────────┘
                                 │
                                 ▼
                    ┌────────────────────────┐
                    │      API Gateway       │
                    │      Port : 8000       │
                    └────────┬───────┬───────┘
                             │       │
                 ┌───────────┘       └───────────┐
                 ▼                               ▼

      ┌──────────────────┐         ┌──────────────────┐
      │   Auth Service   │         │   Post Service   │
      │    Port : 8080   │         │    Port : 8081   │
      └──────────────────┘         └──────────────────┘
                 │                               │
                 └──────────────┬────────────────┘
                                ▼
                      ┌─────────────────┐
                      │   PostgreSQL    │
                      └─────────────────┘
>>>>>>> 0b72c8f (add api gateway)
```

---

# Services

| Service | Port | Responsibility |
|---|---|---|
<<<<<<< HEAD
=======
| API Gateway | 8000 | Reverse Proxy & Routing |
>>>>>>> 0b72c8f (add api gateway)
| Auth Service | 8080 | Authentication & User Management |
| Post Service | 8081 | Post CRUD Operations |

---

<<<<<<< HEAD
=======
# Features

---

## API Gateway

- Reverse Proxy
- Centralized Entry Point
- Route Forwarding
- Header Forwarding
- JWT Forwarding
- CORS Middleware
- Request Logging
- Service Isolation

---

## Auth Service

- User Registration
- User Login
- JWT Generation
- Password Hashing
- Protected Routes
- User CRUD
- Role-Based Access
- Profile Endpoint

---

## Post Service

- JWT Validation
- Create Posts
- Get Posts
- Update Posts
- Delete Posts
- Ownership Authorization
- Soft Deletes
- Protected APIs

---

>>>>>>> 0b72c8f (add api gateway)
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
<<<<<<< HEAD
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
=======
| Reverse Proxy | API Gateway |
| Middleware | Logging & Authorization |
>>>>>>> 0b72c8f (add api gateway)

---

# Project Structure

<<<<<<< HEAD
## Auth Service
=======
# API Gateway

```text
api-gateway/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── gateway/
│   │   └── proxy.go
│   │
│   └── middleware/
│       ├── cors.go
│       └── logger.go
│
├── .env
├── go.mod
└── README.md
```

---

# Auth Service
>>>>>>> 0b72c8f (add api gateway)

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

<<<<<<< HEAD
## Post Service
=======
# Post Service
>>>>>>> 0b72c8f (add api gateway)

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

<<<<<<< HEAD
# Run User Service Migration
=======
# User Service Migration
>>>>>>> 0b72c8f (add api gateway)

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('admin', 'user');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS users (
<<<<<<< HEAD
=======

>>>>>>> 0b72c8f (add api gateway)
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

<<<<<<< HEAD
# Run Post Service Migration

```sql
CREATE TABLE IF NOT EXISTS posts (
=======
# Post Service Migration

```sql
CREATE TABLE IF NOT EXISTS posts (

>>>>>>> 0b72c8f (add api gateway)
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

<<<<<<< HEAD
## Auth Service `.env`
=======
# API Gateway `.env`

```env
PORT=8000

AUTH_SERVICE_URL=http://localhost:8080

POST_SERVICE_URL=http://localhost:8081
```

---

# Auth Service `.env`
>>>>>>> 0b72c8f (add api gateway)

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

<<<<<<< HEAD
## Post Service `.env`
=======
# Post Service `.env`
>>>>>>> 0b72c8f (add api gateway)

```env
PORT=8081

DB_URL=postgres://postgres:yourpassword@localhost:5432/go_rest_api?sslmode=disable

JWT_SECRET=your-super-secret-key-change-in-production
```

---

# Install Dependencies

<<<<<<< HEAD
=======
## API Gateway

```bash
go mod tidy
```

---

>>>>>>> 0b72c8f (add api gateway)
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

<<<<<<< HEAD
# Run Services

## Start Auth Service
=======
# Running Services

---

# Start Auth Service
>>>>>>> 0b72c8f (add api gateway)

```bash
go run cmd/server/main.go
```

Runs on:

```text
http://localhost:8080
```

---

<<<<<<< HEAD
## Start Post Service
=======
# Start Post Service
>>>>>>> 0b72c8f (add api gateway)

```bash
go run cmd/server/main.go
```

Runs on:

```text
http://localhost:8081
```

---

<<<<<<< HEAD
# API Endpoints

# Auth Service

## Register User

```http
POST /api/v1/auth/register
=======
# Start API Gateway

```bash
go run cmd/server/main.go
```

Runs on:

```text
http://localhost:8000
```

---

# Final API Endpoints

IMPORTANT:

Client should ONLY communicate with API Gateway.

Do NOT directly expose internal services to frontend.

---

# Auth APIs

## Register

```http
POST http://localhost:8000/api/v1/auth/register
>>>>>>> 0b72c8f (add api gateway)
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

<<<<<<< HEAD
## Login User

```http
POST /api/v1/auth/login
=======
## Login

```http
POST http://localhost:8000/api/v1/auth/login
>>>>>>> 0b72c8f (add api gateway)
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
<<<<<<< HEAD
GET /api/v1/profile
=======
GET http://localhost:8000/api/v1/profile
>>>>>>> 0b72c8f (add api gateway)
```

### Headers

```http
Authorization: Bearer JWT_TOKEN
```

---

<<<<<<< HEAD
# Post Service
=======
# Post APIs
>>>>>>> 0b72c8f (add api gateway)

## Create Post

```http
<<<<<<< HEAD
POST /api/v1/posts
=======
POST http://localhost:8000/api/v1/posts
>>>>>>> 0b72c8f (add api gateway)
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
<<<<<<< HEAD
GET /api/v1/posts
=======
GET http://localhost:8000/api/v1/posts
>>>>>>> 0b72c8f (add api gateway)
```

---

<<<<<<< HEAD
## Get Post By ID

```http
GET /api/v1/posts/{id}
=======
## Get Single Post

```http
GET http://localhost:8000/api/v1/posts/{id}
>>>>>>> 0b72c8f (add api gateway)
```

---

## Update Post

```http
<<<<<<< HEAD
PUT /api/v1/posts/{id}
=======
PUT http://localhost:8000/api/v1/posts/{id}
>>>>>>> 0b72c8f (add api gateway)
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
<<<<<<< HEAD
DELETE /api/v1/posts/{id}
=======
DELETE http://localhost:8000/api/v1/posts/{id}
>>>>>>> 0b72c8f (add api gateway)
```

---

# JWT Authentication Flow

<<<<<<< HEAD
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
=======
```text
1. User logs into Auth Service
2. Auth Service generates JWT token
3. Client stores JWT token
4. Client sends token to API Gateway
5. Gateway forwards request to Post Service
6. Post Service validates JWT
7. Protected operation succeeds
```
>>>>>>> 0b72c8f (add api gateway)

---

# Example Flow

<<<<<<< HEAD
## Step 1 — Register

```http
POST :8080/api/v1/auth/register
=======
# Step 1 — Register

```http
POST :8000/api/v1/auth/register
>>>>>>> 0b72c8f (add api gateway)
```

---

<<<<<<< HEAD
## Step 2 — Login

```http
POST :8080/api/v1/auth/login
=======
# Step 2 — Login

```http
POST :8000/api/v1/auth/login
>>>>>>> 0b72c8f (add api gateway)
```

Copy JWT token.

---

<<<<<<< HEAD
## Step 3 — Create Post

```http
POST :8081/api/v1/posts
=======
# Step 3 — Create Post

```http
POST :8000/api/v1/posts
```

Headers:

```http
>>>>>>> 0b72c8f (add api gateway)
Authorization: Bearer JWT_TOKEN
```

---

<<<<<<< HEAD
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
=======
# Step 4 — Get Posts

```http
GET :8000/api/v1/posts
```

---

# Security Features

- JWT Authentication
- bcrypt Password Hashing
- Middleware Authorization
- Ownership Validation
- Protected Routes
- Reverse Proxy Isolation
- Soft Deletes
- UUID Primary Keys

---

# Infrastructure Concepts Implemented

- Layered Architecture
- Repository Pattern
- API Gateway
- Reverse Proxy
- Microservices
- JWT Trust Between Services
- Service Isolation
- Middleware Chaining
- Structured Routing

---

# Future Improvements

- Docker
- Docker Compose
- Redis Caching
- Swagger/OpenAPI
- Refresh Tokens
- Rate Limiting
- API Versioning
- Service Discovery
- Circuit Breakers
- Kafka/RabbitMQ
- Kubernetes
- Centralized Logging
- Distributed Tracing
- CI/CD Pipeline
- Unit Testing
- Integration Testing
>>>>>>> 0b72c8f (add api gateway)

---

# Important Notes

<<<<<<< HEAD
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
=======
This project is for:

- Learning Golang Backend Development
- Understanding JWT Authentication
- Learning API Gateway Architecture
- Understanding Reverse Proxying
- Practicing Microservices
- Practicing Layered Architecture
- Understanding Middleware

This is NOT fully production-ready yet.

Missing:
- observability
- metrics
- tracing
- retries
- load balancing
- centralized authentication
- distributed transactions
- service discovery
- TLS termination
- container orchestration
>>>>>>> 0b72c8f (add api gateway)

---

# Author

<<<<<<< HEAD
Ubaid Rza
=======
Ubaid Raza
>>>>>>> 0b72c8f (add api gateway)

---

# License

<<<<<<< HEAD
MIT License
=======
MIT License
>>>>>>> 0b72c8f (add api gateway)
