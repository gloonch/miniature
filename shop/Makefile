.PHONY: build run lint test docker clean

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOLINT=golangci-lint run
BINARY_NAME=shop-service
MAIN_PATH=./cmd/main.go

# Docker parameters
DOCKER_CMD=docker
IMAGE_NAME=shop-service
TAG=latest

# Default target
all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Ensure golangci-lint is installed: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
lint:
	@echo "Linting code..."
	$(GOLINT) ./...

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

docker-build:
	@echo "Building Docker image $(IMAGE_NAME):$(TAG)..."
	$(DOCKER_CMD) build -t $(IMAGE_NAME):$(TAG) -f Dockerfile .

docker-run:
	@echo "Running Docker container $(IMAGE_NAME):$(TAG)..."
	$(DOCKER_CMD) run -p 8081:8081 --env-file .env $(IMAGE_NAME):$(TAG)

clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

help:
	@echo "Available targets:"
	@echo "  build          Build the application"
	@echo "  run            Build and run the application"
	@echo "  lint           Lint the Go code"
	@echo "  test           Run unit tests"
	@echo "  docker-build   Build the Docker image"
	@echo "  docker-run     Run the Docker container (requires .env file)"
	@echo "  clean          Remove build artifacts"
	@echo "  help           Show this help message"
