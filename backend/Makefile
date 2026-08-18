include ../.env
export

.PHONY: run dev build migrate-diff migrate-apply migrate-status sqlc generate swagger

## Run the server locally
run:
	go run ./cmd/server

## Run the server with live reload
dev:
	POPISOMATOR_DEBUG=true go tool air -c .air.toml

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

## Regenerate Swagger/OpenAPI docs from handler annotations. Run this after adding/changing
## an endpoint's @... comment block (see internal/controller for examples) and commit the result.
swagger:
	go tool swag init -g cmd/server/main.go -o internal/docs --parseDependency --parseInternal
