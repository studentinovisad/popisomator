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

## Apply migrations to the DB
migrate-apply:
	migrate -database $(POPISOMATOR_MIGRATION_DSN) -path ./db/migrations up

## Creates migration files (up/down)
migrate-create:
	migrate create -ext sql -dir ./db/migrations $(name)

## Show current migration status
migrate-status:
	migrate -database $(POPISOMATOR_MIGRATION_DSN) -path ./db/migrations version

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
