# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /usr/local/bin ./...

FROM alpine
COPY --from=builder /usr/local/bin/* /usr/local/bin/
