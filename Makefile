.PHONY: help build up down logs clean dev dev-down test

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the production image
	docker compose build goexpert-lab-deploy-dev

up: ## Start the production service
	docker compose up goexpert-lab-deploy-dev -d

down: ## Stop and remove containers
	docker compose down

logs: ## View logs from the production service
	docker compose logs -f goexpert-lab-deploy-dev

test: ## Run tests
	go test ./...

run-local: ## Run the application locally
	go run cmd/api/main.go
