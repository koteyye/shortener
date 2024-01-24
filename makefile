SHORTENER_PORT := 8080
DATABASE_SHORTENER := postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable

.DEFAULT_GOAL := all

.PHONY: run
run:
	@go run cmd/shortener/main.go \
		-a=localhost:8081 \
		-b=http://localhost:8081 \
		-d=${DATABASE_SHORTENER} \
		-j=supersecret \
		-p=localhost:9090

.PHONY: all
all: test link autotest

.PHONY: generate
generate:
	@go generate ./...

.PHONY: up
up:
	@docker-compose -f ./scripts/docker-compose.yaml up -d

.PHONY: down
down:
	@docker-compose -f ./scripts/docker-compose.yaml down

.PHONY: lint
lint:
	@go vet -vettool=$(shell which statictest) ./...

.PHONY: test
test:
	@go test -short -race -timeout=30s -count=1 -cover ./...

.PHONY: build
build:
	@go build -buildvcs=false -o ./cmd/shortener/shortener ./cmd/shortener

.PHONY: autotest
autotest: build
	@./shortenertestbeta \
		-test.v -test.run=^TestIteration1 \
		-binary-path=cmd/shortener/shortener \
		-server-host=localhost \
		-server-port=8081 \
		-server-base-url=localhost:8081 \
		-source-path=cmd/shortener \
		-database-dsn=${DATABASE_SHORTENER} \
