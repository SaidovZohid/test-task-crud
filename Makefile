run:
	@go run cmd/main.go

tidy: vendor
	@go mod tidy

vendor:
	@go mod vendor