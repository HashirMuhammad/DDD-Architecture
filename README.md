# DDD User Service

A user management service built with Domain Driven Design (DDD) principles in Go.

## Features

- Create, Read, Update, Delete (CRUD) operations for users
- Domain-driven design architecture
- RESTful API with JSON responses
- In-memory storage (easily replaceable with database)
- Input validation and error handling

## Architecture

The project follows DDD principles with clean architecture:

```
cmd/
└── main.go                    # Application entry point

internal/
├── domain/                    # Domain layer (business logic)
│   ├── user.go               # User entity with business rules
│   └── repository.go         # Repository interface
├── application/              # Application layer (use cases)
│   ├── dto/                  # Data Transfer Objects
│   │   └── user_dto.go
│   └── service/              # Application services
│       └── user_service.go
├── infrastructure/           # Infrastructure layer (external concerns)
│   └── repository/
│       └── memory_user_repository.go
└── interfaces/              # Interface layer (controllers, HTTP)
    └── http/
        ├── handler/
        │   └── user_handler.go
        └── router/
            └── router.go
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Users
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

## User Model

```json
{
  "id": "string (UUID)",
  "name": "string",
  "email": "string (valid email format)",
  "username": "string (minimum 3 characters)"
}
```

## Database Setup

This application now uses MongoDB for data persistence.

### Prerequisites
- MongoDB installed and running on localhost:27017
- MongoDB can be downloaded from https://www.mongodb.com/try/download/community

### Database Setup Steps

1. **Start MongoDB**:
   ```bash
   # On Windows (if MongoDB is installed as a service)
   net start MongoDB
   
   # Or run MongoDB manually
   mongod
   ```

2. **Verify MongoDB Connection**:
   ```bash
   # Connect using MongoDB shell
   mongosh mongodb://localhost:27017/
   ```

The application will automatically:
- Connect to MongoDB at `mongodb://localhost:27017/`
- Create database `UserServiceDB`
- Create collection `users` with unique indexes for email and username

## Running the Application

1. **Install dependencies**:
```bash
go mod tidy
```

2. **Start MongoDB** (ensure it's running on localhost:27017)

3. **Run the application**:
```bash
go run cmd/main.go
```

The server will start on port 8080 (or the port specified in the PORT environment variable).

### Database Configuration

The application connects to:
- **MongoDB URI**: `mongodb://localhost:27017/`
- **Database**: `UserServiceDB`
- **Collection**: `users`

Connection details are configured in `internal/infrastructure/config/mongodb.go`.

## Example Usage

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "username": "johndoe"
  }'
```

### Get All Users
```bash
curl http://localhost:8080/api/v1/users
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id}
```

## Domain Rules

- Name cannot be empty
- Email must be in valid format and unique
- Username must be at least 3 characters and unique
- All fields are automatically trimmed and normalized (email/username to lowercase)

## MongoDB Features

- **Automatic Indexing**: Unique indexes on email and username
- **Document Storage**: Users stored as MongoDB documents
- **Persistent Storage**: Data survives application restarts
- **Concurrent Access**: MongoDB handles multiple connections

## Error Handling

The API returns appropriate HTTP status codes:
- 200: OK
- 201: Created
- 204: No Content
- 400: Bad Request (validation errors)
- 404: Not Found
- 409: Conflict (duplicate email/username)
- 500: Internal Server Error
