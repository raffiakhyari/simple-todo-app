# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Download dependencies first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /todo-api \
    ./cmd/api

# Runtime stage
FROM alpine:3.22

WORKDIR /app

# Install CA certificates for HTTPS/TLS connections
RUN apk --no-cache add ca-certificates

# Copy application binary
COPY --from=builder /todo-api /app/todo-api

# Application port
EXPOSE 8080

# Run application
ENTRYPOINT ["/app/todo-api"]