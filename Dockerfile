# syntax=docker/dockerfile:1

FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build ./cmd/yl-msgdump/...
RUN go build ./cmd/yl-pgwriter/...
RUN go build ./cmd/yl-request/...

FROM alpine
COPY --from=builder /app/yl-* /app/
ENV PATH "$PATH:/app"
