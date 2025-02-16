# Use the official Golang image from the Docker Hub
FROM golang:1.21

# Set the working directory in the container
WORKDIR /app

# Copy the rest of the application code into the container
COPY . .

# Create new directory for keys
RUN rm -rf ./.keys
RUN mkdir -p ./.keys

# Install Go dependencies before switching user
RUN go mod tidy


# Create a non-root user and switch to it
RUN useradd -m appuser
RUN chown -R appuser:appuser . && chmod -R 700 .
USER appuser

# Build the Go application
RUN go build -o main .

# Set environment variable for the port
ENV PORT 8080
ENV CLIENTS "client1|MyClientSecret1|http://localhost|client2|MyClientSecret2|http://localhost"

# Expose the port the app runs on
EXPOSE ${PORT}

# Command to run the application
CMD ["./main"]