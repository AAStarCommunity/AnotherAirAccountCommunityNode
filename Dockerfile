## build
FROM golang:1.23.0-alpine3.20 AS build-env

RUN apk add build-base

ADD . /go/src/app

WORKDIR /go/src/app

RUN go env -w GO111MODULE=on \
    && go mod tidy \
    && go build -o cnode

## run
FROM alpine:3.20.2

RUN mkdir -p /aa && mkdir -p /aa/log

WORKDIR /aa

COPY --from=build-env /go/src/app /aa/

ENV PATH $PATH:/aa

EXPOSE 80
EXPOSE 7946/udp

CMD ["/aa/cnode"]
