# Use the official Python image from the Docker Hub
FROM python:3.10

# Set the working directory in the container
WORKDIR /app

# Copy the rest of the application code into the container
COPY . .

RUN rm -rf ./.keys
RUN mkdir -p ./.keys
# Create a non-root user and switch to it
RUN useradd -m appuser
RUN chown -R appuser:appuser ./.keys && chmod -R 700 ./.keys
USER appuser

# Install any dependencies specified in requirements.txt
RUN pip install --no-cache-dir -r requirements.txt



# Set environment variable for the port
ENV PORT 8080
ENV DOMAIN "localhost"
ENV PROTOCOL "http"
ENV CLIENTS "client1|MyClientSecret1|http://localhost|client2|MyClientSecret2|http://localhost"


# Expose the port the app runs on
EXPOSE 8080
# Command to run the application
CMD ["python3", "http_server.py"]