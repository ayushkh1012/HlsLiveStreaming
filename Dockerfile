# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy necessary directories
COPY media/ ./media/
COPY ads/ ./ads/
COPY manifests/ ./manifests/
COPY config/ ./config/

# Create necessary directories with proper permissions
RUN mkdir -p /app/media/1080p \
    /app/media/720p \
    /app/media/480p \
    /app/media/360p \
    /app/ads/adv1/1080p \
    /app/ads/adv1/720p \
    /app/ads/adv1/480p \
    /app/ads/adv1/360p \
    /app/manifests && \
    chmod -R 755 /app

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 