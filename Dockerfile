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
RUN make build -o main .
WORKDIR /out
RUN cp /app/main .

FROM alpine:latest
RUN apk add ca-certificates
COPY --from=builder /out/main /
COPY .env /
EXPOSE 8080
ENTRYPOINT ["/main"]