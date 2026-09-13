.PHONY: migrate-up migrate-down migrate-version migrate-create run build test docker-up docker-down docker-logs docker-ps

# Load .env file properly (handles quotes, comments)
ENV_FILE := .env
ifneq ("$(wildcard $(ENV_FILE))","")
  include $(ENV_FILE)
  export
endif

# Database migration commands
migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

migrate-down-all:
	migrate -path ./migrations -database "$(DATABASE_URL)" down

migrate-version:
	migrate -path ./migrations -database "$(DATABASE_URL)" version

migrate-force:
	@read -p "Force version number: " version; \
	migrate -path ./migrations -database "$(DATABASE_URL)" force $$version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# Development
run:
	air

build:
	go build -o ./bin/server ./cmd/server

# Run built binary
run-bin:
	./bin/server

# Test
test:
	go test ./...

# Lint (if you have golangci-lint)
lint:
	golangci-lint run

# Generate (for sqlc, etc.)
generate:
	go generate ./...

# Docker commands (run from project root)
docker-up:
	docker-compose -f deploy/docker-compose.yml up -d

docker-down:
	docker-compose -f deploy/docker-compose.yml down

docker-logs:
	docker-compose -f deploy/docker-compose.yml logs -f

docker-ps:
	docker-compose -f deploy/docker-compose.yml ps

# Clean
clean:
	rm -rf ./bin ./tmp
