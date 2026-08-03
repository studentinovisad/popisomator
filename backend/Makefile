include ../.env
export

.PHONY: run build migrate-diff migrate-apply migrate-status sqlc generate

## Run the server locally
run:
	go run ./cmd/server

## Build a binary
build:
	go build -o bin/server ./cmd/server

## Generate a new migration from schema.sql, e.g. make migrate-diff name=create_users
migrate-diff:
	atlas migrate diff $(name) --env local

## Apply pending migrations to the DB
migrate-apply:
	atlas migrate apply --env local

## Show current migration status
migrate-status:
	atlas migrate status --env local

## Regenerate sqlc Go code from db/queries + db/schema.sql
sqlc:
	sqlc generate

## Regenerate everything code-gen related (go:generate directives, incl. sqlc)
generate:
	go generate ./...
