# Todo List API

A simple and secure REST API for managing TODO tasks, built in Go with JWT authentication and PostgreSQL database.

## 🚀 Features

- **TODO CRUD**: Create, read, update and delete tasks
- **User Management**: User registration and authentication
- **JWT Authentication**: Secure access and refresh tokens
- **Categories & Priorities**: Organize your tasks by category and priority levels
- **Due Dates**: Set deadlines for your tasks
- **REST API**: Well-structured endpoints following REST standards
- **Swagger Documentation**: Interactive API documentation
- **PostgreSQL Database**: Reliable data persistence
- **Soft Delete**: Logical deletion of records
- **Graceful Shutdown**: Elegant server shutdown handling

## 🛠 Tech Stack

- **Language**: Go 1.24.4
- **Web Framework**: Chi Router v5
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Documentation**: Swagger (swaggo/swag)
- **Containers**: Docker & Docker Compose
- **Validation**: go-playground/validator
- **Testing**: Testcontainers for integration tests
- **Hot Reload**: Air for development

## 📋 Prerequisites

- [Go](https://golang.org/dl/) 1.24 or higher
- [Docker](https://www.docker.com/) and Docker Compose
- [Just](https://github.com/casey/just) - Command runner

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone <your-repository>
cd todo-list-api
```

### 2. Set up environment variables

```bash
# Copy the example environment file
cp .env.example .env
```

Edit the `.env` file with your configuration. Required variables:

```bash
PORT=8080
APP_ENV=development
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=todo_db
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=password
BLUEPRINT_DB_SCHEMA=public
API_KEY=your-api-key-here
JWT_SECRET=your-super-secret-jwt-key
```

### 3. Install dependencies

```bash
go mod download
```

### 4. Start development environment (Recommended)

```bash
just dev
```

**This is the recommended way to start the server!** The `just dev` command will:

- Start PostgreSQL database in Docker
- Wait for the database to be ready
- Start the application with hot reload using Air
- Automatically restart the app when you make code changes

The API will be available at `http://localhost:8080`

## 🔧 Available Commands

### Development Commands

```bash
just                # Build and test (default command)
just build          # Build the application
just run            # Run the application (requires manual DB setup)
just dev            # 🌟 RECOMMENDED: Start development environment
just watch          # Hot reload with Air only
just test           # Run unit tests
just itest          # Run integration tests
just clean          # Clean compiled binaries
```

### Docker Commands

```bash
just docker-run     # Start database container only
just docker-down    # Stop database container
just dev-down       # Stop all development services
```

### CI/CD Commands

```bash
just ci             # Run CI checks (build + test)
```

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

### Interactive Swagger Documentation

Once the application is running, access the interactive documentation at:

```
http://localhost:8080/docs
```

### Main Endpoints

#### Authentication

- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/refresh` - Refresh JWT token

#### TODOs

- `GET /api/todos` - Get all TODOs
- `POST /api/todos` - Create new TODO
- `GET /api/todos/{id}` - Get specific TODO
- `PUT /api/todos/{id}` - Update TODO
- `DELETE /api/todos/{id}` - Delete TODO

### Usage Examples

#### Create a TODO

```bash
curl -X POST http://localhost:8080/api/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-jwt-token>" \
  -d '{
    "title": "Learn Go",
    "description": "Complete Go tutorial",
    "priority": "high",
    "category": "learning",
    "dueDate": "2024-12-31T23:59:59Z"
  }'
```

#### Get all TODOs

```bash
curl -H "Authorization: Bearer <your-jwt-token>" \
  http://localhost:8080/api/todos
```

## 🏗 Project Structure

```
todo-list-api/
├── cmd/api/               # Application entry point
├── internal/
│   ├── controller/        # HTTP controllers
│   ├── database/          # Database configuration
│   ├── middleware/        # HTTP middlewares
│   ├── models/           # Data models (GORM)
│   ├── repository/       # Data access layer
│   ├── server/           # Server configuration
│   ├── service/          # Business logic
│   └── utils/            # Utilities
├── docs/                 # Generated Swagger documentation
├── web/                  # Static files (if any)
├── docker-compose.yml    # Docker services
├── Dockerfile           # Application Docker image
├── justfile            # Development commands
├── .air.toml           # Hot reload configuration
├── .env.example        # Environment variables template
└── go.mod              # Go dependencies
```

## 🧪 Testing

### Unit Tests

```bash
just test
```

### Integration Tests

Integration tests use Testcontainers to create temporary PostgreSQL instances:

```bash
just itest
```

### Local CI Pipeline

```bash
just ci
```

## 🐳 Docker

### Development with Docker

```bash
# Recommended: Hybrid development (DB in Docker + local app with hot reload)
just dev

# Start only the database
just docker-run

# Stop development services
just dev-down
```

### Production

```bash
# Full build with Docker Compose
docker compose up --build
```

## 🔧 Development

### Hot Reload Development

For development with automatic reload (recommended):

```bash
just dev
```

This command will:

1. Start PostgreSQL in Docker
2. Wait for the database to be ready
3. Start the application with Air hot reload
4. Automatically restart when you make changes

### Environment Variables

The main environment variables include:

- `APP_ENV` - Application environment
- `PORT` - Server port
- `BLUEPRINT_DB_HOST` - Database host
- `BLUEPRINT_DB_PORT` - PostgreSQL port
- `BLUEPRINT_DB_DATABASE` - Database name
- `BLUEPRINT_DB_USERNAME` - PostgreSQL username
- `BLUEPRINT_DB_PASSWORD` - PostgreSQL password
- `BLUEPRINT_DB_SCHEMA` - Database schema
- `API_KEY` - API key for additional security
- `JWT_SECRET` - Secret key for JWT token signing

Make sure to copy `.env.example` to `.env` and fill in the appropriate values.

## 🤝 Contributing

1. Fork the project
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request

Make sure to run tests before submitting:

```bash
just ci
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
