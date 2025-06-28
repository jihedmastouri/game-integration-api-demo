.PHONY: make_doc
make_doc:
	@docker run --rm -v $(pwd):/code ghcr.io/swaggo/swag:latest

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/server ./cmd/main.go

.PHONY: run
run:
	@air
