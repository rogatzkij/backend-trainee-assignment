FROM golang:1.14.2-alpine3.11 AS builder
RUN apk add --update git
WORKDIR /go/src/service

COPY go.mod go.sum /go/src/service/
RUN GO111MODULE=on go mod download

COPY ./src .
RUN GO111MODULE=on go install all

FROM alpine:3.11
COPY --from=builder /go/bin/backend-trainee-assignment /usr/bin/service

RUN apk add --no-cache ca-certificates && \
  adduser -DH service

EXPOSE 9000

USER service
CMD [ "/usr/bin/service" ]

