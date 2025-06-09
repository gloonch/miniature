BINARY_NAME=product-service
PKG_PATH=./product/cmd

.PHONY: build run clean test docker-build docker-run

build:
	@echo "Building product service..."
	@go build -o $(BINARY_NAME) $(PKG_PATH)/main.go

run: build
	@echo "Running product service..."
	@./$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf ./coverage

test:
	@echo "Testing product service..."
	@go test -v ./product/... -coverprofile=./coverage/product.out
	@go tool cover -html=./coverage/product.out -o ./coverage/product.html

# Docker related targets
DOCKER_IMAGE_NAME=product-service
DOCKER_TAG=latest

docker-build:
	@echo "Building Docker image for product service..."
	@docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) -f ./product/Dockerfile .

# Example: docker-run might require network to connect to DB
# Ensure this is adjusted based on actual needs
docker-run:
	@echo "Running product service Docker container..."
	@docker run -p 8082:8082 --name $(DOCKER_IMAGE_NAME) --rm $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
