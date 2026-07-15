include .env
export

GOOSE_IMAGE ?= ghcr.io/kukymbr/goose-docker:3.27.2
GOOSE_MIGRATION_DIR ?= ./migrations

# Prefer local goose binary if installed, otherwise use Docker
GOOSE ?= $(shell command -v goose 2>/dev/null)

.PHONY: db-up db-down migrate-create migrate-up migrate-down migrate-status migrate-reset migrate-version

## Start PostgreSQL
db-up:
	docker compose up -d postgres

## Stop PostgreSQL
db-down:
	docker compose down

## Create a new migration: make migrate-create name=add_users
migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=migration_name" && exit 1)
ifneq ($(GOOSE),)
	$(GOOSE) -dir $(GOOSE_MIGRATION_DIR) create $(name) sql
else
	docker run --rm -v "$(CURDIR)/migrations:/migrations" --entrypoint goose $(GOOSE_IMAGE) \
		-dir /migrations create $(name) sql
endif

## Apply all pending migrations
migrate-up: db-up
ifneq ($(GOOSE),)
	$(GOOSE) -dir $(GOOSE_MIGRATION_DIR) up
else
	GOOSE_COMMAND=up docker compose --profile tools run --rm goose
endif

## Rollback the last migration
migrate-down: db-up
ifneq ($(GOOSE),)
	$(GOOSE) -dir $(GOOSE_MIGRATION_DIR) down
else
	GOOSE_COMMAND=down docker compose --profile tools run --rm goose
endif

## Show migration status
migrate-status: db-up
ifneq ($(GOOSE),)
	$(GOOSE) -dir $(GOOSE_MIGRATION_DIR) status
else
	GOOSE_COMMAND=status docker compose --profile tools run --rm goose
endif

## Rollback all migrations
migrate-reset: db-up
ifneq ($(GOOSE),)
	$(GOOSE) -dir $(GOOSE_MIGRATION_DIR) reset
else
	GOOSE_COMMAND=reset docker compose --profile tools run --rm goose
endif

## Show current DB version
migrate-version: db-up
ifneq ($(GOOSE),)
	$(GOOSE) -dir $(GOOSE_MIGRATION_DIR) version
else
	GOOSE_COMMAND=version docker compose --profile tools run --rm goose
endif
