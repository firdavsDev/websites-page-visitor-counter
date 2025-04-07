# Visitor Counter API

A Go (Fiber) application for tracking website visitors with Swagger documentation.

## Prerequisites

- Go 1.23+
- Docker 20.10+
- Docker Compose 2.0+
- PostgreSQL (for development, included in Docker setup)

## Installation

### Manual Setup (Without Docker)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/visitor-counter.git
   cd visitor-counter
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Install Air for live reload**:
   ```bash
   go install github.com/air-verse/air@latest
   ```

4. **Set up environment variables**:
   Create `.env` file:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` with your configuration:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=youruser
   DB_PASSWORD=yourpassword
   DB_NAME=visitor_counter
   SERVER_PORT=8080
   ```

5. **Start PostgreSQL**:
   ```bash
   docker run --name visitor-counter-db -e POSTGRES_USER=youruser -e POSTGRES_PASSWORD=yourpassword -e POSTGRES_DB=visitor_counter -p 5432:5432 -d postgres:15-alpine
   ```

6. **Run the application**:
   ```bash
   air
   ```

### Docker Setup (Recommended)

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/visitor-counter.git
   cd visitor-counter
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   ```
   Edit `.env` as needed.

3. **Build and start containers**:
   ```bash
   docker-compose up --build
   ```

4. **Access the application**:
   - API: `http://localhost:8080/api/track`
   - Swagger UI: `http://localhost:8080/swagger/index.html`

## Development

### With Docker (Live Reload)

1. Start development environment:
   ```bash
   docker-compose up
   ```

2. The application will automatically reload when you make changes to:
   - Go files (`*.go`)
   - Templates (`*.html`, `*.tpl`)
   - Configuration files

### Without Docker

1. Start Air for live reload:
   ```bash
   air
   ```

2. Make your code changes - the server will restart automatically.

## Project Structure

```
visitor-counter/
├── cmd/
│   └── main.go          # Application entry point
├── config/              # Configuration files
├── docs/                # Swagger documentation
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   └── storage/         # Database operations
├── models/              # Data models
├── .env.example         # Environment variables template
├── air.toml             # Live reload configuration
├── docker-compose.yml   # Docker setup
├── Dockerfile           # Production Dockerfile
├── Dockerfile.dev       # Development Dockerfile
├── go.mod               # Go dependencies
└── README.md            # This file
```

## API Documentation

Interactive Swagger documentation is available at:
`http://localhost:8080/swagger/index.html`

## Troubleshooting

### Common Issues

1. **Permission errors**:
   ```bash
   docker-compose down -v
   docker system prune -f
   docker-compose build --no-cache
   ```

2. **Database connection issues**:
   - Verify your `.env` file matches Docker compose settings
   - Check if PostgreSQL is running:
     ```bash
     docker-compose logs db
     ```

3. **Air not detecting changes**:
   - Ensure your files are saved with LF line endings
   - Check `air.toml` includes all relevant file extensions

## Deployment

For production deployment, use the production Dockerfile:
```bash
docker build -t visitor-counter .
docker run -d -p 8080:8080 --env-file .env visitor-counter
```

## License

[MIT](LICENSE)