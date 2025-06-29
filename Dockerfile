FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and pre-generated docs
COPY cmd/ cmd/
COPY internal/ internal/
COPY web/ web/
COPY docs/ docs/

# Build the application
RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary
COPY --from=build /app/main /app/main

# Copy web assets and docs
COPY --from=build /app/web/ /app/web/
COPY --from=build /app/docs/ /app/docs/

EXPOSE ${PORT}
CMD ["./main"]
