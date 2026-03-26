FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o tk-admin-server .

FROM debian:stable-slim

WORKDIR /app
COPY --from=builder /app/tk-admin-server .
COPY config.yaml .

RUN chmod +x tk-admin-server

EXPOSE 8080

CMD ["./tk-admin-server"]