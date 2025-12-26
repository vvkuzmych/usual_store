# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the support service
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/support-service ./cmd/support-service

# Production stage
FROM debian:bullseye-slim

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /usr/local/bin/support-service ./support-service

# Expose port
EXPOSE 5000

# Run the service
CMD ["./support-service"]

