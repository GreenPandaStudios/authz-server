# Use the official Golang image
FROM golang:1.21 AS builder

# Set the working directory
WORKDIR /app

# Copy application files
COPY . .

# Remove and recreate keys directory
RUN rm -rf ./.keys && mkdir -p ./.keys

# Disable CGO to prevent issues with cross-compilation
ENV CGO_ENABLED=0

# Install Go dependencies
RUN go mod tidy

# Build the Go application
RUN go build -o main .

# Use a minimal runtime image for final execution
FROM debian:bullseye-slim

# Create a non-root user
RUN useradd -m appuser
RUN chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Use the non-root user
USER appuser

# Set environment variables
ENV PORT=8080
ENV CLIENTS="client1|MyClientSecret1|http://localhost|client2|MyClientSecret2|http://localhost"

# Expose the application port
EXPOSE ${PORT}

# Command to run the application
CMD ["./main"]
