.PHONY: make_doc
make_doc:
	@docker run --rm -v $(pwd):/code ghcr.io/swaggo/swag:latest

.PHONY: build
build:
	@CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/server ./cmd/main.go

.PHONY: run
run:
	@air

.PHONY: migration
migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration name=your_migration_name"; \
		exit 1; \
	fi
	@mkdir -p repository/migrations
	@timestamp=$$(date +%s); \
	touch "repository/migrations/$${timestamp}_$(name).up.sql"; \
	touch "repository/migrations/$${timestamp}_$(name).down.sql"; \
	echo "Created migration files:"; \
	echo "  repository/migrations/$${timestamp}_$(name).up.sql"; \
	echo "  repository/migrations/$${timestamp}_$(name).down.sql"

