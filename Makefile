APP_NAME ?= go-code-review
GO_CMD ?= go
DOCKERCMD ?= docker
DOCKERCOMPOSECMD ?= docker compose
GOLANGCICMD ?= golangci-lint

.PHONY: *

# docker/up: starts the service with docker-compose
docker/up: docker/check-engine
	$(DOCKERCOMPOSECMD) up -d

# docker/down: stops the service
docker/down:
	$(DOCKERCOMPOSECMD) down

# generate: runs go mod tidy and generates swagger documentation
generate: tidy swagger/generate

# docker/check-engine: checks if docker engine is running
docker/check-engine:
	@if ! docker info >/dev/null 2>&1; then \
		echo "ERROR: Docker engine is not running. Please start Docker manually."; \
		exit 1; \
	fi

# test: runs tests locally
test:
	$(GO_CMD) test -v -cover -short ./...

# swagger/generate: generates swagger documentation
swagger/generate:
	swag init --parseDependency --parseInternal -g ./internal/api/api.go --output ./docs

# tidy: tidy dependencies
tidy:
	$(GO_CMD) mod tidy