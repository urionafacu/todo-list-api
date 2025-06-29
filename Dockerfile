FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install swag for generating documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/
COPY web/ web/

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go

# Build the application
RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app

# Copy binary
COPY --from=build /app/main /app/main

# Copy web assets and generated docs
COPY --from=build /app/web/ /app/web/
COPY --from=build /app/docs/ /app/docs/

EXPOSE ${PORT}
CMD ["./main"]
