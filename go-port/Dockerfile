# Multi-stage Docker build for WebGL Water Tutorial Go Port

# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -s /bin/sh webgl

# Set working directory
WORKDIR /home/webgl/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy static assets and web files
COPY --from=builder /app/web ./web
COPY --from=builder /app/assets ./assets

# Copy original assets from parent directory
COPY ../dudvmap.png ./assets/
COPY ../normalmap.png ./assets/
COPY ../stone-texture.png ./assets/

# Create directories for runtime
RUN mkdir -p ./assets ./web/static ./web/templates ./web/shaders

# Change ownership to non-root user
RUN chown -R webgl:webgl /home/webgl

# Switch to non-root user
USER webgl

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Set environment variables
ENV PORT=8080
ENV ASSETS_PATH=./assets
ENV STATIC_PATH=./web/static

# Run the server
CMD ["./server"]
