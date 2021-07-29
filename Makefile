.PHONY: build
build:
	go build -v ./...

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: mocks
mocks:
	mockery --case underscore --dir ./internal/service/ --output ./internal/service/mocks --all --disable-version-string

.PHONY : test
test:
	go test -v ./internal/apiserver ./internal/repository

.DEFAULT_GOAL := build