# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the messaging service
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/messaging-service ./cmd/messaging-service

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/messaging-service .

# Copy email templates
COPY --from=builder /app/cmd/messaging-service/templates ./templates

# Expose port (optional, messaging service doesn't need to expose ports)
EXPOSE 8081

# Run the messaging service
CMD ["./messaging-service"]

