# Build stage
FROM golang:1.22-alpine3.21 AS builder

WORKDIR /usr/src/app

COPY . .

RUN go build -v -o main main.go

# Run stage
FROM alpine:3.21.3

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/main .
COPY --from=builder /usr/src/app/.env .env

EXPOSE 8000

CMD [ "/usr/src/app/main" ]