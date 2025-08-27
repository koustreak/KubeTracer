# Multi-stage build for smaller final image
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kubetracer ./cmd/kubetracer

# Final stage - minimal image
FROM alpine:3.18

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -S kubetracer && adduser -S kubetracer -G kubetracer

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/kubetracer .

# Change ownership
RUN chown -R kubetracer:kubetracer /app

# Switch to non-root user
USER kubetracer

# Run the binary
ENTRYPOINT ["./kubetracer"]
