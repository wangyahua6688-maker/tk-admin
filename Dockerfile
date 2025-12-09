FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o go-admin ./cmd/godash

FROM debian:stable-slim
WORKDIR /app
COPY --from=builder /app/go-admin .
COPY config.yaml .
EXPOSE 8080
CMD ["./go-admin", "--config", "config.yaml"]
