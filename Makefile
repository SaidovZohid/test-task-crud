-include .env
DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

run:
	@go run cmd/main.go

tidy:
	@go mod tidy
	@go mod vendor

cache:
	@go clean -testcache

test:
	# @go test -v -cover ./storage/...
	@go test -v -cover ./...

swag-init:
	@swag init -g api/api.go -o api/docs

# database migrations
#
# create migration name=users
create-migration:
	@migrate create -ext sql -dir migrations -seq $(name)
#
# up all migrations
migrateup:
	@migrate -path migrations -database "$(DB_URL)" -verbose up
#
# up migration last one
migrateup1:
	@migrate -path migrations -database "$(DB_URL)" -verbose up 1
#
# down migrations all
migratedown:
	@migrate -path migrations -database "$(DB_URL)" -verbose down
#
# down the migration last
migratedown1:
	@migrate -path migrations -database "$(DB_URL)" -verbose down 1

up:
	@docker compose --env-file ./.env.docker up -d

down:
	@docker compose down

build:
	@docker build --platform linux/amd64 --tag $(DOCKER_USERNAME)/test-task-crud:latest .

push: build
	@docker image push $(DOCKER_USERNAME)/test-task-crud:latest
	