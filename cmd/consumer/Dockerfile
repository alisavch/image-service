FROM golang:1.14-alpine as build

RUN apk add --no-cache git

WORKDIR /src

RUN go get github.com/sirupsen/logrus
RUN go get github.com/streadway/amqp

COPY main.go /src

RUN go build main.go


FROM alpine as runtime

COPY --from=build /src/consumer /app/consumer

CMD [ "/app/consumer" ]