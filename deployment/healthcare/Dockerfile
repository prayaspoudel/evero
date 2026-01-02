# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o healthcare-server ./app/healthcare

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/healthcare-server .

# Copy config files
COPY --from=builder /app/config/healthcare /root/config/healthcare

# Expose port
EXPOSE 3001

# Run the application
CMD ["./healthcare-server", "--config=config/healthcare/production.json"]
