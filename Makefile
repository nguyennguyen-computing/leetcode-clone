# LeetCode Clone Makefile

.PHONY: help seed-db clean-db setup-dev

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

seed-db: ## Seed the database with sample data
	@echo "🌱 Seeding database with sample data..."
	@cd backend/scripts && ./seed.sh

clean-db: ## Clear all data from the database (WARNING: This will delete all data!)
	@echo "⚠️  WARNING: This will delete ALL data from the database!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "🧹 Clearing database..."
	@cd backend/scripts && DB_CLEAR_ONLY=true go run seed_data.go

setup-dev: ## Set up development environment
	@echo "🔧 Setting up development environment..."
	@echo "Starting Docker services..."
	@docker-compose up -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	@echo "Running database migrations..."
	@cd backend && go run scripts/validate_schema.go
	@echo "Seeding database with sample data..."
	@make seed-db
	@echo "✅ Development environment is ready!"

.DEFAULT_GOAL := help