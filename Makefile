.PHONY: docs vendor

vendor:
	@go mod tidy
	@go mod vendor

test:
	@go test -v ./...

docs:
	@swag init --dir app/api --generalInfo api.go
