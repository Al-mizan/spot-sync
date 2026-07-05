.PHONY: all build run test vet fmt

# Load environment variables
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

all: build

build:
	@echo "Building the application..."
	go build -o tmp/spotsync cmd/main.go

run:
	@echo "Running the application..."
	go run cmd/main.go

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...
