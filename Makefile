APP_NAME=go-admin
BUILD_DIR=build

.PHONY: all build run docker-build docker-run backfill-draw-results

all: build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/godash

run:
	go run ./cmd/godash

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker-compose up --build

backfill-draw-results:
	go run ./tools/backfill_draw_results -config config.yaml -force
