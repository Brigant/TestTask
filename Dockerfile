########-- Build stage --########

FROM golang:1.21-alpine AS builder

WORKDIR /opt

COPY . /opt

RUN go build -o ./runner ./cmd

########-- Deploy stage --########

FROM alpine:3.18

WORKDIR /opt 

COPY --from=builder /opt/runner /opt/

ENTRYPOINT /opt/runner