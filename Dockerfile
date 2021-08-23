FROM golang:alpine as builder
RUN apk --no-cache add build-base git gcc make
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build -o /out/image-service .

FROM alpine:latest
RUN apk add ca-certificates
COPY --from=builder /out/image-service /app
#WORKDIR /out
#COPY --from=builder /app/main .
#COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["/app"]