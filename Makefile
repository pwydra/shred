# Makefile for the go applications 

APP_NAME := shred-service
DOCKER_IMAGE_NAME := shred-app
DOCKER_COMPOSE_FILE := docker-compose.yml
BIN_DIR := bin

.PHONY: all build docker-build test coverage run clean

all: build image

build:
	@echo "Building Go application..."
	go build -o $(BIN_DIR) ./cmd/shred-service
	@echo "Go application built."

image: build
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME) .
	@echo "Docker image built: $(DOCKER_IMAGE_NAME)"

test:
	@echo "Running unit tests..."
	go test -v ./...
	@echo "Unit tests finished."

coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	@echo "Coverage report generated: coverage.out.html"

up:
	@echo "Starting application with Docker Compose..."
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Application started (see logs with 'docker compose logs -f')"

clean:
	@echo "Stopping and removing Docker containers..."
	docker compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "Removing Go binary..."
	rm -f $(BIN_DIR)/$(APP_NAME)
	@echo "Clean complete."

# Helper target to view Docker Compose logs
logs:
	docker compose -f $(DOCKER_COMPOSE_FILE) logs -f

# Helper target to stop Docker Compose
down:
	docker compose -f $(DOCKER_COMPOSE_FILE) stop

# Helper target to restart Docker Compose
restart: down up
