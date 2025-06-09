# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Ensure the product service is built
RUN CGO_ENABLED=0 GOOS=linux go build -o /product-service ./product/cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /product-service .
# Expose port (if product service has an HTTP server, e.g., 8082)
# EXPOSE 8082
CMD ["./product-service"]
