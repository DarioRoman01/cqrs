ARG GO_VERSION=1.18.4

FROM golang:${GO_VERSION}-alpine as builder

RUN go env -w GOPROXY=direct 
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src/

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY models models
COPY events events
COPY repository repository
COPY database database
COPY search search
COPY command-service command-service
COPY query-service query-service
COPY pusher-service pusher-service

RUN go install ./...

FROM alpine:3.11

WORKDIR /usr/bin

COPY --from=builder /go/bin .
