# Use the official Golang image for building the app
FROM golang:1.22 AS builder

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

# Create a non-root user in the builder stage
RUN useradd -m appuser

# Use a minimal runtime image for final execution
FROM debian:bullseye-slim AS runtime
RUN apt-get update && apt-get install -y ca-certificates

# Create the same non-root user in the runtime image
RUN useradd -m appuser

# Ensure the /app directory exists and set correct ownership
RUN mkdir -p /app && chown -R appuser:appuser /app

# Set working directory
WORKDIR /app

# Copy the compiled binary and service account file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/service-account.json .

# Use the non-root user
USER appuser

# Set environment variables
ENV PORT=8080
ENV CLIENTS="client1|MyClientSecret1|http://localhost|client2|MyClientSecret2|http://localhost"
# Expose the application port
EXPOSE ${PORT}

# Command to run the application
CMD ["./main"]