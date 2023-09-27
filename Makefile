include .envrc

.PHONY: help
help:
	@echo "Usage"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n "Are you sure [y/N]" && read ans && [ $${ans:-N} = y ]

.PHONY: run/api
## run/api: run the cmd/api application
run/api:
	go run ./cmd/api

.PHONY: db/psql
## db/psql: connect to the database using psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

.PHONY: db/migarations/up
## db/migrations/up: apply all up database migrations
db/migrations/up: confirm
	@echo "Running up migrations"
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up

.PHONY: db/migarations/new
## db/migrations/new name=$1: create a new database migration
db/migrations/new: confirm
	@echo "Creating migarations files for ${name}..."
	migrate create -seq -ext-.sql -dir=./migrations ${name}

.PHONY: audit
## audit
audit: vendor
	@echo "Formatting code..."
	go fmt ./...
	@echo "Vetting code..."
	go vet ./...
	staticcheck ./...
	@echo "Running tests..."
	go test -race -vet=off ./...

.PHONY: vendor
## vendor
vendor:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Vendoring dependencies..."
	go mod vendor


