.PHONY: docs vendor

vendor:
	@go mod tidy
	@go mod vendor

docs:
	@swag init --dir internal/api --generalInfo api.go

test:
	@go test -v ./...

run: docs
	@go run cmd/moonbase/main.go
