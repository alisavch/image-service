GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=image_service
LINTER=golangci-lint

.PHONY: build
build:
	$(GOBUILD) -x ./cmd/api/main.go

.PHONY: run
run:
	$(GORUN) ./cmd/api/main.go

.PHONY: mocks
mocks:
	mockery --case underscore --dir ./internal/service/ --output ./internal/service/mocks --all --disable-version-string
	mockery --case underscore --dir ./internal/broker/ --output ./internal/broker/mocks --all --disable-version-string
	mockery --case underscore --dir ./internal/bucket/ --output ./internal/bucket/mocks --all --disable-version-string


.PHONY: lint
lint:
	$(LINTER) run --config .golangci.yaml

.PHONY: test
test:
	$(GOTEST) ./... -v

.PHONY: generate-spec
generate-spec:
	swagger generate spec -m -o swagger.yaml

.DEFAULT_GOAL := test