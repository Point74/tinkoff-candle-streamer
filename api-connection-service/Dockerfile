FROM golang:1.24.3-alpine AS builder
WORKDIR /app

COPY api-connection-service/go.mod api-connection-service/go.sum ./api-connection-service/
COPY config/ ./config/
COPY contracts/ ./contracts/
COPY prometheus/metrics/ ./prometheus/metrics/

WORKDIR /app/api-connection-service
RUN go mod download
COPY api-connection-service/ ./
RUN go build -o api-connection-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-connection-service/api-connection-service .
RUN apk add --no-cache bash
ENTRYPOINT ["/app/api-connection-service"]