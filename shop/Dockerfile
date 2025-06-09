# Use an official Go runtime as a parent image
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app/shop-service

# Copy go mod and sum files
COPY shop/go.mod shop/go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY shop/ ./
COPY pkg/ ./pkg/

# Build the Go app
# Ensure the output binary name is distinct, e.g., shop-service
RUN go build -o /app/shop-service/shop-service ./cmd/main.go

# Start a new stage from scratch for a smaller image
FROM alpine:latest
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/shop-service/shop-service .
# Copy migrations (if they need to be run by the service or a separate tool)
COPY --from=builder /app/shop-service/migrations ./migrations

# Expose port 8081 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./shop-service"]
