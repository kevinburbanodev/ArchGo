# Go Hexagonal Template

English | [Español](README.es.md)

A robust template for building Go applications using hexagonal architecture (also known as ports and adapters). This template uses Gin as the web framework for handling HTTP requests and provides a solid foundation for building scalable and maintainable APIs.

## Project Structure

```
.
├── cmd/
│   └── server/         # Application entry point
├── internal/
│   ├── handlers/       # HTTP handlers
│   ├── infrastructure/ # Concrete implementations
│   ├── middleware/     # Application middleware
│   └── modules/        # Application modules
│       └── user/       # User module
│           ├── application/    # Use cases
│           ├── domain/         # Models and ports
│           └── infrastructure/ # Implementations
└── tests/              # Application tests
```

## Requirements

- Go 1.21 or higher
- PostgreSQL
- Make (optional, for make commands)
- Gin Web Framework (automatically installed via go.mod)

## Tech Stack

- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Architecture**: Hexagonal (Ports and Adapters)
- **Authentication**: JWT
- **Documentation**: Swagger
- **Containerization**: Docker & Docker Compose

## Configuration

1. Clone the repository
2. Copy `.env.example` to `.env` and configure environment variables
3. Run `go mod download` to install dependencies
4. Run `go run cmd/server/main.go` to start the server

## Environment Variables

### Database Configuration
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_hexagonal
DB_SSL_MODE=disable
```

### GORM Configuration
```
GORM_LOG_LEVEL=debug    # Levels: debug, info, warn, error, silent
GORM_AUTO_MIGRATE=true  # true/false
```

### Server Configuration
```
PORT=3000
JWT_SECRET=your-secret-key
ENV=development
```

### Specific Variables Explanation

#### DB_SSL_MODE
Configures SSL mode for PostgreSQL connection:
- `disable`: No SSL (recommended for local development)
- `require`: Requires SSL connection
- `verify-ca`: Verifies server certificate is signed by a trusted CA
- `verify-full`: Verifies certificate and hostname (most secure)

#### GORM_LOG_LEVEL
Controls GORM logging level:
- `debug`: Shows all SQL queries and details
- `info`: Shows general information
- `warn`: Shows only warnings
- `error`: Shows only errors
- `silent`: No logs

#### GORM_AUTO_MIGRATE
Enables/disables automatic database migration:
- `true`: GORM will automatically create/update tables
- `false`: No automatic migrations will be performed

## API Documentation (Swagger)

This project uses Swagger for API documentation. To generate and view the documentation:

1. Install the Swagger CLI tool:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Generate the Swagger documentation:
   ```bash
   swag init -g cmd/server/main.go
   ```

3. The documentation will be available at:
   ```
   http://localhost:3000/swagger/index.html
   ```

4. To update the documentation after making changes to the API:
   ```bash
   swag init -g cmd/server/main.go
   ```

Note: Make sure to add Swagger annotations to your handlers to keep the documentation up to date.

## Endpoints

### Users

#### Create User
```bash
curl --location 'http://localhost:3000/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123"
}'
```

#### Login
```bash
curl --location 'http://localhost:3000/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "email": "test@example.com",
    "password": "password123"
}'
```

#### Get User (requires authentication)
```bash
curl --location 'http://localhost:3000/api/users/1' \
--header 'Authorization: Bearer <token>'
```

## Security

### Rate Limiting

The API implements rate limiting to prevent brute force and DoS attacks. Features include:

- **IP-based Limits**: 100 requests per minute
- **Response Headers**:
  - `X-RateLimit-Limit`: Total request limit (100)
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Time until counter reset

#### Example Response with Rate Limit
```http
HTTP/1.1 200 OK
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1710288000
```

#### Example Response When Limit Exceeded
```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{
    "error": "You have exceeded the request limit. Please wait a moment."
}
```

### Other Security Measures

- **JWT Authentication**: All protected routes require a valid JWT token
- **SSL/TLS**: Configurable through `DB_SSL_MODE` for secure database connections
- **Input Validation**: All input data is validated before processing
- **Security Headers**: The API includes standard security headers

## Docker

### Requirements
- Docker
- Docker Compose

### Docker Structure

#### Dockerfile
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
EXPOSE 3000
CMD ["./main"]
```

#### docker-compose.yml
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME=go_hexagonal
      - DB_SSL_MODE=disable
      - GORM_LOG_LEVEL=debug
      - GORM_AUTO_MIGRATE=true
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=go_hexagonal
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```

### Docker Features

1. **Multi-stage Build**
   - Reduces final image size
   - Separates build process from runtime environment
   - Includes only necessary files

2. **Services**
   - **app**: Main application service
     - Built from Dockerfile
     - Exposes port 3000
     - Configured with environment variables
   - **postgres**: PostgreSQL database
     - Uses official PostgreSQL image
     - Persists data through volume
     - Configured with basic credentials

3. **Networks**
   - Dedicated `app-network`
   - Service isolation
   - Secure container communication

4. **Volumes**
   - `postgres_data`: Persists PostgreSQL data
   - Prevents data loss on container restart

### Docker Commands

1. **Build and Run**
```bash
# Build and start all services
docker-compose up --build

# Run in background
docker-compose up -d
```

2. **Container Management**
```bash
# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Restart services
docker-compose restart
```

3. **Maintenance**
```bash
# Clean unused containers
docker system prune

# View resource usage
docker stats
```

### Environment Variables

Environment variables can be configured in two ways:

1. **In docker-compose.yml**
```yaml
environment:
  - DB_HOST=postgres
  - DB_PORT=5432
  # ... other variables
```

2. **In .env file**
```env
DB_HOST=postgres
DB_PORT=5432
# ... other variables
```

### Security Considerations

1. **Credentials**
   - Don't include sensitive credentials in code
   - Use environment variables or Docker secrets
   - Change default credentials in production

2. **Networks**
   - Use Docker networks to isolate services
   - Expose only necessary ports
   - Configure restrictive network policies

3. **Volumes**
   - Use named volumes for persistence
   - Configure appropriate permissions
   - Regular data backup

## Git Hooks

This project includes Git hooks to ensure code quality. To set up the hooks:

1. Make sure you have the required tools installed:
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. Configure Git to use the project's hooks:
   ```bash
   git config core.hooksPath .githooks
   ```

The pre-push hook will run:
- Tests
- Linter
- Code formatting checks

## Contact

### Author
**Kevin Fernando Burbano Aragón**  
Systems Engineer and Senior Software Developer with extensive experience in software development.

### Contact Information
- **Email**: [burbanokevin1997@gmail.com](mailto:burbanokevin1997@gmail.com)
- **GitHub**: [@kevinburbanodev](https://github.com/kevinburbanodev)
- **LinkedIn**: [Kevin Fernando Burbano Aragón](https://www.linkedin.com/in/kevin-fernando-burbano-arag%C3%B3n-78b3871a0/)

### Experience
- Systems Engineer with senior-level expertise in software development
- Specialized in clean architectures and design patterns
- Expert in API and microservices development
- Strong background in multiple technologies and frameworks

## License
MIT 