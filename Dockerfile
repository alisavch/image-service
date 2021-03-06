FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk --no-cache add build-base git gcc make

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder ["/app/api", "/"]

ENTRYPOINT ["/api"]
