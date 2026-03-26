FROM golang:1.24-alpine AS builder

ENV GOPROXY=https://proxy.golang.org,direct
ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# ✅ 强制指定 arm64
RUN GOARCH=arm64 GOOS=linux go build -o tk-admin-server ./cmd/godash

RUN chmod +x tk-admin-server

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/tk-admin-server .
COPY config.yaml .

RUN mkdir -p /app/uploads

EXPOSE 8080

CMD ["./tk-admin-server"]


