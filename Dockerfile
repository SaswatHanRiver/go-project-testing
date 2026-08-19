# ── Stage 1: Build ───────────────────────────────────────────────────────────
# Like a Maven build container - compiles the Go binary
# We use the full Go image only for building, not for running
FROM golang:1.21-alpine AS builder

# Install git (needed by some Go modules)
RUN apk add --no-cache git

# Set working directory inside container
WORKDIR /app

# Copy dependency files first - Docker caches this layer
# Like copying pom.xml before source code so dependencies are cached
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the binary - CGO_ENABLED=0 makes a fully static binary (no C dependencies)
# -o app = output file named "app"
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# ── Stage 2: Run ─────────────────────────────────────────────────────────────
# Minimal Alpine image - only ~5MB, no Go toolchain needed
# Like running a Spring Boot JAR in a JRE-only image (not JDK)
FROM alpine:latest

# Install CA certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/app .

# Expose port 8080 - same as Spring Boot's server.port=8080
EXPOSE 8080

# Run the binary
# Note: .env is NOT copied - use Docker env vars or docker-compose instead
CMD ["./app"]
