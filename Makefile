.PHONY: run build test clean migrate-up migrate-down docker-up docker-down

# --- Load .env ---
ifneq (,$(wildcard ./.env))
	include .env
	export
endif

# --- Default DB variables (fallback nếu .env không có) ---
DB_HOST ?= localhost
DB_USER ?= postgres
DB_NAME ?= meobeo_talk
DB_PORT ?= 5432
DB_PASSWORD ?= 

# --- Default Goal ---
.DEFAULT_GOAL := run

# --- Run & Build ---
run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v -race ./...

clean:
	rm -rf bin/

# --- Migration helpers ---
MIGRATIONS_DIR := migrations

migrate-up:
	@echo "Running UP migrations..."
	@for file in $(MIGRATIONS_DIR)/*.up.sql; do \
		echo "Applying $$file"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -p $(DB_PORT) -f $$file; \
	done

migrate-down:
	@echo "Running DOWN migrations..."
	@for file in $(MIGRATIONS_DIR)/*.down.sql | sort -r; do \
		echo "Reverting $$file"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -U $(DB_USER) -d $(DB_NAME) -p $(DB_PORT) -f $$file; \
	done

# --- Docker ---
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down
