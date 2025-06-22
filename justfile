# Todo List API - Development Commands

# Default recipe - build and test
default: build test

# Build the application
build:
    @echo "Building..."
    @go build -o main cmd/api/main.go

# Run the application
run:
    @go run cmd/api/main.go

# Create and run DB container
docker-run:
    #!/usr/bin/env bash
    if docker compose up --build 2>/dev/null; then
        :
    else
        echo "Falling back to Docker Compose V1"
        docker-compose up --build
    fi

# Shutdown DB container
docker-down:
    #!/usr/bin/env bash
    if docker compose down 2>/dev/null; then
        :
    else
        echo "Falling back to Docker Compose V1"
        docker-compose down
    fi

# Test the application
test:
    @echo "Testing..."
    @go test ./... -v

# Integration tests for the application
itest:
    @echo "Running integration tests..."
    @go test ./internal/database -v

# Clean the binary
clean:
    @echo "Cleaning..."
    @rm -f main

# Live reload with Air
watch:
    #!/usr/bin/env bash
    if command -v air > /dev/null; then
        air
        echo "Watching..."
    else
        echo "Go's 'air' is not installed on your machine."
        read -p "Do you want to install it? [Y/n] " choice
        if [ "$choice" != "n" ] && [ "$choice" != "N" ]; then
            go install github.com/air-verse/air@latest
            air
            echo "Watching..."
        else
            echo "You chose not to install air. Exiting..."
            exit 1
        fi
    fi

# Development with Docker DB + Hot Reload (Hybrid approach)
dev:
    #!/usr/bin/env bash
    echo "🐘 Starting PostgreSQL in Docker..."
    if docker compose up psql_bp -d 2>/dev/null; then
        echo "✅ PostgreSQL started with Docker Compose v2"
    else
        echo "📦 Falling back to Docker Compose v1"
        docker-compose up psql_bp -d
    fi
    
    echo "⏳ Waiting for database to be ready..."
    sleep 5
    
    echo "🔥 Starting app with hot reload..."
    if command -v air > /dev/null; then
        air
    else
        echo "Go's 'air' is not installed on your machine."
        read -p "Do you want to install it? [Y/n] " choice
        if [ "$choice" != "n" ] && [ "$choice" != "N" ]; then
            go install github.com/air-verse/air@latest
            air
        else
            echo "You chose not to install air. Stopping database..."
            just dev-down
            exit 1
        fi
    fi

# Stop development services
dev-down:
    #!/usr/bin/env bash
    echo "🛑 Stopping development services..."
    if docker compose down 2>/dev/null; then
        echo "✅ Services stopped with Docker Compose v2"
    else
        echo "📦 Falling back to Docker Compose v1"
        docker-compose down
        echo "✅ Services stopped with Docker Compose v1"
    fi

# Run all checks (build, test, lint)
ci: build test

# Show available commands
help:
    @just --list 