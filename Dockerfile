FROM golang:1.19.1-alpine3.16 as builder

WORKDIR /app

COPY . .

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

RUN go build -o main cmd/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .
COPY templates ./templates
COPY --from=builder /app/migrate ./migrate
COPY migrations ./migrations

EXPOSE 8080

CMD [ "/app/main" ]