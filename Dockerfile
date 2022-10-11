FROM golang:1.19-alpine3.16 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/avgalaida/textBoard

COPY go.mod go.sum ./
COPY util util
COPY event event
COPY db db
COPY search search
COPY schema schema
COPY post-service post-service
COPY query-service query-service
COPY pusher-service pusher-service

RUN GO111MODULE=on go install ./...

FROM alpine:3.16
WORKDIR /usr/bin
COPY --from=build /go/bin .
