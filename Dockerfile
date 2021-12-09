FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk --no-cache add build-base git gcc make

WORKDIR /app

RUN go get github.com/go-delve/delve/cmd/dlv

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -gcflags="all=-N -l"  -o api ./cmd/api/main.go
FROM alpine:latest
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/api /api
COPY --from=builder /go/bin/dlv /dlv

ENTRYPOINT ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/api"]
