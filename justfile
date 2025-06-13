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

# Run all checks (build, test, lint)
ci: build test

# Show available commands
help:
    @just --list 