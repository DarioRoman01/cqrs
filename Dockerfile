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
COPY cache cache
COPY search search
COPY helpers helpers
COPY middlewares middlewares
COPY pusher-service pusher-service
COPY feeds-command feeds-command
COPY feeds-query feeds-query
COPY users-command users-command
COPY users-query users-query
COPY auth-service auth-service

RUN go install ./...

FROM alpine:3.11

WORKDIR /usr/bin

COPY --from=builder /go/bin .
