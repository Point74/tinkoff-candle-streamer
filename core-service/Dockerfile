FROM golang:1.24.3-alpine AS builder
WORKDIR /app

COPY core-service/go.mod core-service/go.sum ./core-service/
COPY config/ ./config/
COPY contracts/ ./contracts/
COPY prometheus/metrics/ ./prometheus/metrics/

WORKDIR /app/core-service
RUN go mod download
COPY core-service/ ./
RUN go build -o core-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/core-service/core-service .
RUN apk add --no-cache bash
ENTRYPOINT ["/app/core-service"]