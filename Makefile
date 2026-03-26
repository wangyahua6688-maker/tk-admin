APP_NAME=go-admin
BUILD_DIR=build

.PHONY: all build run docker-build docker-run

all: build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/godash

run:
	go run ./cmd/godash

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker-compose up --build
