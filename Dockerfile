# Builder
FROM golang:1.16-alpine as build

RUN go version

WORKDIR /cmd/api

ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o image-service ./cmd/api/main.go

# Distribution
FROM alpine:latest

RUN go version

WORKDIR /cmd/api

EXPOSE 8080

COPY --from-builder /cmd/api/main /cmd/api

CMD /cmd/api/main