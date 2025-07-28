# Task Management API Documentation

## Overview

This API is designed using Clean Architecture principles, ensuring separation of concerns, maintainability, and testability. The codebase is organized into the following layers:

- **Domain**: Core business entities and logic (User, Task).
- **Usecases**: Application-specific business rules (user and task usecases).
- **Repositories**: Data access abstraction (interfaces and MongoDB implementations).
- **Infrastructure**: External services (MongoDB connection, JWT, password hashing, auth middleware).
- **Delivery**: HTTP layer (controllers, routers, main entrypoint).

## MongoDB Usage

- The API uses MongoDB as its data store.
- Connection is established directly in `Delivery/main.go` using the official MongoDB Go driver.
- Collections used: `users`, `tasks` in the `task_manager` database.
- MongoDB URI is read from the `MONGODB_URI` environment variable (or from a `.env` file if present, defaults to `mongodb://localhost:27017`).

## Authorization

- JWT-based authentication is used.
- Register and login endpoints return a JWT token.
- **Protected endpoints require the `Authorization: Bearer <token>` header.**
- Middleware in `Infrastructure/auth_middleware.go` validates JWT and injects claims into the request context.
- Only users with the `admin` role can access certain endpoints (e.g., promote user).

## Endpoints

### Auth & User

- `POST /register` — Register a new user. Returns JWT and role. _(No auth required)_
- `POST /login` — Login with username/email and password. Returns JWT and role. _(No auth required)_
- `POST /promote` — Promote a user to admin (**Requires Authorization header, must be admin**)

### Tasks (all require authentication)

- `GET /tasks` — List all tasks. **Requires Authorization header**
- `GET /tasks/:id` — Get a task by ID. **Requires Authorization header**
- `POST /tasks` — Create a new task. **Requires Authorization header**
- `PUT /tasks/:id` — Update a task by ID. **Requires Authorization header**
- `DELETE /tasks/:id` — Delete a task by ID. **Requires Authorization header**

## Example Usage

### Register

`POST /register`

```json
{
  "username": "user",
  "email": "user@example.com",
  "password": "password123"
}
```

### Login

`POST /login`

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:

```json
{
  "message": "User logged in successfully",
  "token": "<JWT>",
  "role": "user"
}
```

### Authenticated Request Example

```
GET /tasks
Authorization: Bearer <JWT>
```

### Add a Task (Authenticated)

```POST /tasks
Authorization: Bearer <JWT>
Content-Type: application/json
```

```json
{
  "id": "1",
  "title": "My Task",
  "description": "task 1",
  "due_date": "2024-07-23T12:00:00Z",
  "status": "pending"
}
```

### Promote a User (Admin Only)

```POST /promote
Authorization: Bearer <JWT-of-admin>
Content-Type: application/json
```

```json
{
  "identifier": "user"
}
```

## Design Decisions

- **Clean Architecture**: Each layer is decoupled and only depends on abstractions.
- **MongoDB**: Used for persistence, with repository interfaces allowing for easy substitution.
- **JWT**: Used for stateless authentication and role-based access control.

## Architecture Patterns

### DAO vs DTO Pattern

- **DAO (Data Access Object)**: Used in the Repository layer for database serialization/deserialization with `bson` tags
  - `TaskDAO` and `UserDAO` in repositories for MongoDB operations
- **DTO (Data Transfer Object)**: Used in the Controller layer for HTTP request/response serialization with `json` tags
  - `TaskDTO` and `UserDTO` in controllers for API communication

### Clean Architecture Layers

- **Domain**: Pure business logic, no framework dependencies
- **Usecases**: Application business logic, orchestrates domain operations
- **Repositories**: Data access using DAOs, implements domain interfaces
- **Controllers**: HTTP handling using DTOs, calls usecases
- **Infrastructure**: External services (JWT, password hashing)

## Folder Structure

```
task-manager/
├── Delivery/                        # HTTP layer: request handling, controllers, routing
│   ├── main.go                      # App entrypoint: server setup, dependency wiring
│   ├── controllers/
│   │   └── controller.go            # HTTP controllers: handle API requests, call usecases
│   └── routers/
│       └── router.go                # Route definitions: Gin router setup
├── Domain/
│   └── domain.go                    # Core business entities (User, Task structs, interfaces)
├── Infrastructure/                  # External services: JWT, password, auth middleware
│   ├── auth_middleWare.go           # JWT authentication/authorization middleware
│   ├── jwt_service.go               # JWT token generation/validation
│   └── password_service.go          # Password hashing and verification
├── Repositories/                    # Data access abstraction: interfaces & MongoDB impls
│   ├── task_repository.go           # Task repository interface & MongoDB implementation
│   └── user_repository.go           # User repository interface & MongoDB implementation
└── Usecases/                        # Application business logic (orchestrates domain logic and rules)
    ├── task_usecases.go             # Task-related business logic
    └── user_usecases.go             # User-related business logic
```

**Descriptions:**

- **Delivery/**: Handles HTTP requests and responses (controllers, routing, main entrypoint).
- **Domain/**: Defines core business models and interfaces (e.g., User, Task, repository/service interfaces). Contains core domain logic.
- **Infrastructure/**: Implements external dependencies (JWT, password hashing, authentication).
- **Repositories/**: Implements data access, using domain interfaces. No business logic.
- **Usecases/**: Contains application business logic—coordinates domain logic, enforces application rules, and orchestrates interactions between domain and repositories.

This structure follows Clean Architecture, ensuring clear separation of concerns and maintainability.

## Running the API

1. Set up MongoDB and ensure it is running.
2. Set the `MONGODB_URI` environment variable if not using the default.
3. Run the API:
   ```
   go run Delivery/main.go
   ```
4. Use Postman or similar tools to interact with the endpoints.
