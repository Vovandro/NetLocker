# syntax=docker/dockerfile:1

## Build
FROM golang:1.23.4-alpine AS builder

WORKDIR /go/src/app

RUN apk add --no-cache make build-base

ENV GOPATH /go

COPY go.mod .
COPY go.sum .
COPY .go/pkg/mod /go/pkg/mod
RUN go mod download && go mod verify

COPY . .

ENV GOPATH /go
ENV PATH $PATH:/go/bin:$GOPATH/bin

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags musl -o ./NetLocker ./cmd/app/main.go

## Deploy
FROM alpine:3.19 as runtime

WORKDIR /app

COPY --from=builder /go/src/app/NetLocker ./

WORKDIR /app

ENTRYPOINT [ "./NetLocker" ]
