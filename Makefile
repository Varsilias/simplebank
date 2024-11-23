startdb:
	docker compose up -d
	
migrate:
	@eval "$$(cat .env)" && \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
	fi && \
	migrate -path db/migrations -database "$$DB_URL" -verbose up

rollback:
	@eval "$$(cat .env)" && \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable"; \
	fi && \
	migrate -path db/migrations -database "$$DB_URL" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

env:
	cp .env.example .env

.PHONY: startdb migrate rollback sqlc test env
