FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ 正确入口（根目录 main.go）
RUN go build -o tk-admin-server .

RUN chmod +x tk-admin-server

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/tk-admin-server .
COPY config.yaml ./config.yaml

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./tk-admin-server"]