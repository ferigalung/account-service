FROM golang:1.24-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create a working directory
WORKDIR /app

# Install build dependencies
RUN apk --no-cache add git

# Copy go.mod and go.sum to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Go application
RUN go build -o main ./cmd/http

# Create a minimal image for running the application
FROM alpine:latest

# Set environment variables
ENV APP_ENV=production

# Create a non-root user and group for running the application
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create a directory for the application
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

# Set the user to appuser
USER appuser

EXPOSE 3000
CMD ["./main"]
