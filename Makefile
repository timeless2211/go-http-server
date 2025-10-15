.PHONY: run sqlc-generate goose-up goose-down

DB_URL ?= postgres://tuanbui.n9@localhost:5432/chirpy?sslmode=disable

run:
	go run .

sqlc-generate:
	sqlc generate

goose-up:
	goose -dir sql/schema postgres $(DB_URL) up

goose-down:
	goose -dir sql/schema postgres $(DB_URL) down

#goose -dir sql/schema create <migration_name> sql