.PHONY: docs

test:
	@go test -v ./...

docs:
	@swag init --dir app/api --generalInfo api.go
