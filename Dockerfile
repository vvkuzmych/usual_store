# Base image for building the application
FROM golang:1.23 AS builder

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy Go modules manifests
COPY go.mod go.sum ./

COPY .env .env

# Download dependencies
RUN go mod download

# Copy the application source code
COPY . .

# Build the front-end binary
RUN make build_front

# Build the invoice binary
RUN make build_invoice

# Build the back-end binary
RUN make build_back

# Runtime stage
FROM debian:bullseye-slim

# Install required dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the built binaries from the builder stage
COPY --from=builder /app/dist/invoice ./invoice
COPY --from=builder /app/dist/usualstore ./usualstore
COPY --from=builder /app/dist/usualstore_api ./usualstore_api

# Expose ports for the front-end and back-end services
EXPOSE 8080 8081

# Start both front-end and back-end services
CMD ["sh", "-c", "./usualstore -port=${USUAL_STORE_PORT} & ./usualstore_api -port=${API_PORT} & ./invoice -port=${INVOICE_PORT}"]
