# Simple Makefile for a Go project

# Build the application
all: build run

build:
	@echo "Building..."
	@go build -o main main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@go test
