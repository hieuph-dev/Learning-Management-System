.PHONY: help build up down restart logs clean ps

# Hiển thị help
help:
	@echo "LMS Docker Commands:"
	@echo "  make build    - Build Docker images"
	@echo "  make up       - Start all containers"
	@echo "  make down     - Stop all containers"
	@echo "  make restart  - Restart all containers"
	@echo "  make logs     - View logs"
	@echo "  make clean    - Remove all containers, volumes, and images"
	@echo "  make ps       - List running containers"

# Build Docker images
build:
	docker-compose build

# Start containers
up:
	docker-compose up -d

# Stop containers
down:
	docker-compose down

# Restart containers
restart:
	docker-compose restart

# View logs
logs:
	docker-compose logs -f

# View logs của app only
logs-app:
	docker-compose logs -f app

# View logs của postgres only
logs-db:
	docker-compose logs -f postgres

# Clean everything
clean:
	docker-compose down -v --rmi all

# List containers
ps:
	docker-compose ps

# Access app container shell
shell-app:
	docker-compose exec app sh

# Access postgres container shell
shell-db:
	docker-compose exec postgres psql -U postgres -d lms_db

# Rebuild and restart
rebuild:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d