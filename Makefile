.PHONY: build
build:
	go build -v ./...

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: mocks
mocks:
	mockery --case underscore --dir ./internal/service/ --output ./internal/service/mocks --all --disable-version-string

.PHONY: lint
lint:
	goimports -w cmd/api internal/apiserver internal/log internal/model internal/repository internal/service internal/utils
	gofmt -w cmd/api internal/apiserver internal/log internal/model internal/repository internal/service internal/utils
	golint ./...

.PHONY: test
test: lint
	go test -v ./internal/apiserver ./internal/repository



.DEFAULT_GOAL := build