# Variables
APP_NAME=packulator
DOCKER_IMAGE=$(APP_NAME):latest

# Go commands
.PHONY: build run test clean

build:
	go build -o $(APP_NAME) ./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./...

test-coverage:
	go test -cover ./...

clean:
	rm -f $(APP_NAME)

# Docker commands
.PHONY: docker-build docker-run docker-stop docker-clean

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-clean:
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE) || true

# Development commands
.PHONY: dev-setup dev-up dev-down

dev-setup:
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

dev-up:
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	go run ./cmd/main.go

dev-down:
	docker-compose down

# Linting and formatting
.PHONY: fmt vet lint

fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet

# Database commands
.PHONY: db-up db-down db-reset

db-up:
	docker-compose up -d postgres

db-down:
	docker-compose stop postgres

db-reset:
	docker-compose down postgres
	docker volume rm packulator_postgres_data || true
	docker-compose up -d postgres

# Kubernetes commands
.PHONY: k8s-deploy k8s-delete k8s-status k8s-logs

k8s-deploy:
	kubectl apply -f k8s/

k8s-delete:
	kubectl delete -f k8s/

k8s-status:
	kubectl get all -n packulator

k8s-logs:
	kubectl logs -n packulator -l app=packulator --tail=50 -f

# Help
.PHONY: help

help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker containers"
	@echo "  dev-setup    - Setup development environment"
	@echo "  dev-up       - Start development mode"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  lint         - Run linting tools"
	@echo "  k8s-deploy   - Deploy to Kubernetes"
	@echo "  k8s-delete   - Delete from Kubernetes"
	@echo "  k8s-status   - Check Kubernetes status"
	@echo "  k8s-logs     - View application logs"