start:
	docker compose up -d
	
migrate:
	@eval "$$(cat .env)" && \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
	fi && \
	migrate -path db/migrations -database "$$DB_URL" -verbose up

migrate-one:
	@eval "$$(cat .env)" && \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
	fi && \
	migrate -path db/migrations -database "$$DB_URL" -verbose up 1

rollback:
	@eval "$$(cat .env)" && \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
	fi && \
	migrate -path db/migrations -database "$$DB_URL" -verbose down

rollback-one:
	@eval "$$(cat .env)" && \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
	fi && \
	migrate -path db/migrations -database "$$DB_URL" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

env:
	cp .env.example .env

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/varsilias/simplebank/db/sqlc Store


.PHONY: start migrate rollback migrate-one rollback-one sqlc test env server mock
