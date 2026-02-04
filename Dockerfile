# Stage 1: Build the Go binary
FROM golang:alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# We name the output binary "main"
RUN go build -o main ./cmd/fitclassmaster/main.go

# Stage 2: Create a minimal image to run the app
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# CRITICAL: Copy the templates and static assets
# Your code expects these folders to be in the current directory
COPY --from=builder /app/internal/templates ./internal/templates
COPY --from=builder /app/internal/static ./internal/static

# Expose the port your app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]