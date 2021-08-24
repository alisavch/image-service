FROM golang:alpine as builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN apk --no-cache add build-base git gcc make

RUN go get github.com/gorilla/mux
RUN go get github.com/sirupsen/logrus
RUN go get github.com/streadway/amqp
RUN go get github.com/lib/pq

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api/main.go
WORKDIR /out
RUN cp /app/main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /out/main /
COPY .env /
EXPOSE 8080
ENTRYPOINT ["/main"]