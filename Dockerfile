# Build stage
FROM golang:1.22-alpine3.21 AS builder

WORKDIR /usr/src/app

COPY . .

RUN go build -v -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.21.3

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/.env .env
COPY --from=builder /usr/src/app/migrate  ./migrate
COPY --from=builder /usr/src/app/db/migrations ./migrations
COPY --from=builder /usr/src/app/start.sh ./start.sh

RUN chmod +x /usr/src/app/start.sh

EXPOSE 8000

CMD [ "/usr/src/app/main" ]
ENTRYPOINT [ "/usr/src/app/start.sh" ]