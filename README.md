# 🚀 Go API Starter Template

A production-ready Go API starter template with standardized response structure, authentication, role-based permissions, and external API integration capabilities.

## ✨ Features

- **🔐 JWT Authentication** - Secure token-based authentication
- **👥 Role-Based Access Control (RBAC)** - Flexible permission system
- **🌐 External API Integration** - Built-in HTTP client with retry logic
- **📊 Standardized Response Structure** - Consistent API responses
- **🔄 API Versioning Support** - Easy version management
- **🗄️ Database Migration** - Structured database schema management
- **📝 Comprehensive Logging** - Configurable logging levels
- **⚙️ Environment Configuration** - Flexible configuration management
- **🛡️ Middleware Support** - Authentication and API middleware
- **🛠️ CLI Tool** - Easy project generation with `go-starter` command

## 🛠️ Installation & Setup

### 1. Install Go Starter CLI

First, install the `go-starter` CLI tool:

```bash
go install github.com/Palguna1121/go-starter@latest
```

> **Note**: Make sure your `$GOPATH/bin` is in your system PATH to use the `go-starter` command globally.

### 2. Create New Project

Generate a new project using the CLI:

```bash
go-starter new my-awesome-api
cd my-awesome-api
```

This will create a new directory with all the template files and automatically replace template placeholders with your project name.

### 3. Initialize Dependencies

```bash
go mod tidy
```

### 4. Environment Setup

Copy the example environment file and configure it:

```bash
cp .env.example .env
```

Edit `.env` file with your configuration:

```env
# Application Configuration
APP_PORT=5220
APP_NAME=my-awesome-api

# Database Configuration
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database_name

# JWT Configuration
JWT_SECRET=your_super_secret_jwt_key_here
JWT_EXPIRES_IN=24h

# API Configuration
API_VERSION=v1
API_BASE_URL=http://localhost:5220/api/v1

# External API Configuration
EXTERNAL_API_BASE_URL=https://api.example.com
EXTERNAL_API_KEY=your_api_key_here

# Performance & Logging
REQUEST_TIMEOUT=30s
MAX_RETRIES=3
RETRY_DELAY=200ms
ENABLE_LOGGING=true
LOG_LEVEL=debug
ENVIRONMENT=development
```

## 📦 Prerequisites

### Install golang-migrate

Before running the application, install `golang-migrate` for database migrations:

**Option 1: Using Go**
```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Option 2: Download Binary (recommended on windows)**
1. Download from [golang-migrate releases](https://github.com/golang-migrate/migrate/releases)
2. Place `migrate` binary in your `$GOPATH/bin/`
3. Verify installation:
```bash
migrate -version
```

## 🏃‍♂️ Running the Application

### 1. Run Database Migrations

```bash
# Using Makefile (recommended)
make migrate-up

# Or using Go command directly
go run cmd/migrate/migrate.go up
```

### 2. Seed Initial Data (Optional)

```bash
go run cmd/seed/seed.go
```

### 3. Start the Development Server

```bash
# Using Makefile
make run

# Or using Go command
go run main.go
```

Your API will be available at `http://localhost:5220/api/v1`

## 🏗️ Project Structure

```
my-awesome-api/
├── cmd/                    # Command line utilities
│   ├── migrate/           # Database migration commands
│   └── seed/              # Database seeding commands
├── config/                 # Configuration management
│   ├── config.go          # Main configuration
│   ├── database.go        # Database configuration
│   └── external_api.go    # External API configuration
├── core/                   # Core application logic
│   ├── handlers/          # HTTP handlers
│   ├── helper/            # Utility functions
│   ├── middleware/        # Custom middleware
│   ├── models/            # Data models & structures
│   ├── response/          # Response utilities
│   ├── router/            # Route registry
│   └── services/          # Business logic services
├── v1/                     # API version 1
│   ├── controllers/       # HTTP controllers
│   ├── database/          # Migrations and seeds
│   │   ├── migrations/    # SQL migration files
│   │   └── seeds/         # Seed data files
│   ├── middleware/        # Version-specific middleware
│   ├── models/            # Version-specific models
│   └── routes/            # API route definitions
├── v2/                     # API version 2 (for future use)
├── logs/                   # Application logs
├── .env.example           # Environment template
├── .gitignore            # Git ignore rules
├── go.mod                # Go module definition
├── go.sum                # Go dependencies
├── main.go               # Application entry point
├── Makefile              # Build and development commands
└── README.md             # This file
```

## 📚 Usage Examples

### 🔧 Standardized Response Structure

All API responses follow a consistent format:

```json
{
  "status": "success|error",
  "code": 200,
  "message": "Operation completed successfully",
  "data": {
    // Response data here
  }
}
```

### 🎯 Available Response Methods

```go
// Success Responses
response.Success(c, "Operation successful", data)      // 200 OK
response.Created(c, "Resource created", data)          // 201 Created
response.Accepted(c, "Request accepted", data)         // 202 Accepted
response.NoContent(c)                                  // 204 No Content

// Error Responses
response.Error(c, 500, "Internal server error")       // Custom error code
response.BadRequest(c, "Invalid input provided")      // 400 Bad Request
response.Unauthorized(c, "Authentication required")   // 401 Unauthorized
response.Forbidden(c, "Access denied")                // 403 Forbidden
response.NotFound(c, "Resource not found")            // 404 Not Found
response.Conflict(c, "Resource already exists")       // 409 Conflict
response.ValidationError(c, "Validation failed")      // 422 Unprocessable Entity
```

### 🎯 Controller Implementation Example

```go
package controllers

import (
    "github.com/gin-gonic/gin"
    "your-project/core/response"
    "your-project/core/models"
)

type AuthController struct {
    // dependencies
}

func (a *AuthController) Me(c *gin.Context) {
    // Get authenticated user from middleware
    user, exists := c.Get("user")
    if !exists {
        response.Unauthorized(c, "User not authenticated")
        return
    }

    // Type assertion
    u, ok := user.(models.User)
    if !ok {
        response.Error(c, 500, "Invalid user data")
        return
    }

    // Prepare response data
    userData := gin.H{
        "id":         u.ID,
        "name":       u.Name,
        "email":      u.Email,
        "roles":      u.Roles,
        "created_at": u.CreatedAt,
        "updated_at": u.UpdatedAt,
    }

    response.Success(c, "User profile retrieved successfully", userData)
}
```

## 🌐 External API Integration

### 1. Define Route (v1/routes/api.go)
```go
func SetupRoutes(router *gin.Engine) {
    v1 := router.Group("/api/v1")
    {
        // External API routes
        external := v1.Group("/external")
        {
            external.GET("/users/:id", apiHandler.GetExternalUser)
            external.POST("/users", apiHandler.CreateExternalUser)
        }
    }
}
```

### 2. Handler Implementation (core/handlers/api_handler.go)
```go
func (h *APIHandler) GetExternalUser(c *gin.Context) {
    userID := c.Param("id")
    if userID == "" {
        response.BadRequest(c, "User ID is required")
        return
    }

    // Get authorization header
    authToken := c.GetHeader("Authorization")
    
    // Call external service
    user, err := h.userService.GetUserFromExternalAPI(userID, authToken)
    if err != nil {
        response.Error(c, 500, fmt.Sprintf("Failed to fetch user: %s", err.Error()))
        return
    }

    response.Success(c, "User retrieved successfully from external API", user)
}
```

### 3. Service Implementation (core/services/user_service.go)
```go
func (s *UserService) GetUserFromExternalAPI(userID string, token string) (*models.ExternalUser, error) {
    // Prepare API request
    url := fmt.Sprintf("%s/users/%s", s.config.ExternalAPIURL, userID)
    
    apiRequest := &models.APIRequest{
        Method: "GET",
        URL:    url,
        Headers: map[string]string{
            "Authorization": token,
            "Content-Type":  "application/json",
        },
        Timeout: time.Second * 30,
    }

    // Execute request using API client
    response := s.apiClient.ExecuteRequest(apiRequest)
    if !response.Success {
        return nil, fmt.Errorf("external API error: %s", response.Error)
    }

    // Parse response
    var user models.ExternalUser
    if err := mapstructure.Decode(response.Data, &user); err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }

    return &user, nil
}
```

## 🔄 API Versioning

The template supports easy API versioning:

### Environment Configuration
```env
API_VERSION=v1  # Routes to v1/ directory
# Change to v2, v3, etc. for different versions
```

### Version Structure
- `v1/` - Current stable version
- `v2/` - Future version (ready for implementation)
- Each version has its own controllers, routes, and models

### Accessing Different Versions
```bash
# Version 1
curl http://localhost:5220/api/v1/users

# Version 2 (when implemented)
curl http://localhost:5220/api/v2/users
```

## 🗄️ Database Management

### Migration Commands (using Makefile)

```bash
# Create new migration
make migrate-create name=create_products_table

# Run all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Rollback specific number of migrations
make migrate-down-to version=20250614000100

# Check current migration status
make migrate-status

# Force migration version (use with caution)
make migrate-force version=20250614000101
```

### Migration File Example
```sql
-- 20250614000100_create_products_table.up.sql
CREATE TABLE products (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
```

```sql
-- 20250614000100_create_products_table.down.sql
DROP INDEX idx_products_name ON products;
DROP TABLE products;
```

## 🔐 Authentication & Authorization

### JWT Authentication Flow

1. **Login**: User provides credentials, receives JWT token
2. **Authorization**: Include token in `Authorization: Bearer <token>` header
3. **Middleware**: Validates token and injects user data into context
4. **Access Control**: Check user roles/permissions for protected routes

### Example Authentication Usage

```go
// Protected route example
func (h *UserController) GetProfile(c *gin.Context) {
    // User data is automatically available from auth middleware
    user := c.MustGet("user").(models.User)
    
    response.Success(c, "Profile retrieved", gin.H{
        "user": user,
    })
}
```

### Role-Based Access Control (RBAC)

```go
// Middleware usage in routes
authorized := v1.Group("/")
authorized.Use(middleware.AuthMiddleware())
{
    // Requires authentication only
    authorized.GET("/profile", userController.GetProfile)
    
    // Requires specific permission
    authorized.Use(middleware.RequirePermission("users.manage"))
    authorized.GET("/admin/users", userController.ListAllUsers)
}
```

## 🚦 Available Middleware

### Built-in Middleware
- **AuthMiddleware**: JWT token validation
- **APIMiddleware**: External API request handling
- **CORS**: Cross-origin resource sharing
- **Logger**: Request/response logging
- **RateLimit**: API rate limiting (customizable)
- **Recovery**: Panic recovery with logging

### Custom Middleware Example
```go
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Pre-processing
        start := time.Now()
        
        // Process request
        c.Next()
        
        // Post-processing
        duration := time.Since(start)
        log.Printf("Request processed in %v", duration)
    }
}
```

## 📊 Logging Configuration

### Log Levels
- `debug`: Detailed debugging information
- `info`: General application information
- `warn`: Warning messages
- `error`: Error conditions
- `fatal`: Fatal errors causing program termination
- `panic`: Panic-level errors

### Environment Configuration
```env
ENABLE_LOGGING=true
LOG_LEVEL=debug
LOG_FORMAT=json  # or "text"
LOG_OUTPUT=file  # or "console" or "both"
```

### Usage in Code
```go
import "your-project/core/services"

logger := services.GetLogger()

logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error occurred", "error", err)
```

## 🧪 Development Commands

### Makefile Commands
```bash
# Development
make run          # Start development server
make build        # Build production binary
make test         # Run all tests
make test-cover   # Run tests with coverage
make lint         # Run linter

# Database
make migrate-up   # Run migrations
make migrate-down # Rollback migrations
make db-seed      # Seed database

# Docker (if implemented)
make docker-build # Build Docker image
make docker-run   # Run in Docker container
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./v1/controllers -v
```

## 🐳 Docker Support (Optional)

### Dockerfile Example
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 5220
CMD ["./main"]
```

### Docker Compose
```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "5220:5220"
    environment:
      - DB_HOST=mysql
    depends_on:
      - mysql
      
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: your_database
    ports:
      - "3306:3306"
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Changelog

### v0.1.2
- Initial release with CLI tool
- JWT authentication system
- RBAC implementation
- External API integration
- Database migration system
- Comprehensive logging
- Standardized response structure

## 🚀 Roadmap

- [ ] GraphQL support
- [ ] WebSocket integration
- [ ] Redis caching layer
- [ ] Swagger/OpenAPI documentation
- [ ] Health check endpoints
- [ ] Metrics and monitoring
- [ ] Docker containerization
- [ ] Kubernetes deployment configs

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- Database migrations by [golang-migrate](https://github.com/golang-migrate/migrate)
- JWT implementation using [golang-jwt](https://github.com/golang-jwt/jwt)

---

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/Palguna1121/go-starter/issues) page
2. Create a new issue with detailed information
3. Join our community discussions

**Happy coding! 🎉**

Made with ❤️ for the Go community.