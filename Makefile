include .env
export

.PHONY: run build migr-create db-up db-down sqlc

run:
	go run ./internal/cmd/main.go

build:
	go build ./internal/cmd/main.go

migr-create:
	dbmate --env-file ".env" -d "internal/db/migrations" new "$(name)"

db-up:
	dbmate -d ./internal/config/migrations -url ${DATABASE_URL} up

db-down:
	dbmate -d ./internal/config/migrations -url ${DATABASE_URL} down

sqlc:
	sqlc -f internal/config/sqlc.yaml generate
