.PHONY: docs vendor

vendor:
	@go mod tidy
	@go mod vendor

fmt:
	@go fmt ./internal/...

test:
	@go test ./internal/...

docs:
	@swag init --dir internal/api --generalInfo api.go

run: docs
	@go run cmd/moonbase/main.go

docker: docs
	@docker build -t moonbase .
