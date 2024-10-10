# Makefile for gotoolbox-service

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Main package path
MAIN_PACKAGE=.

# Binary name
BINARY_NAME=gotoolbox-service

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PACKAGE)

run:
	$(GORUN) $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

test:
	$(GOTEST) -v ./...

deps:
	$(GOGET) ./...

# Docker
docker-build:
	docker build -t gotoolbox-service .

docker-run:
	docker run -p 8080:8080 gotoolbox-service

docker-run-detached:
	docker run -d -p 8080:8080 gotoolbox-service

test-api:
	@echo "Testing /products endpoint:"
	@curl -s http://localhost:8080/products
	@echo "\n\nTesting /comments endpoint:"
	@curl -s http://localhost:8080/comments

.PHONY: all build run clean test deps docker-build docker-run docker-run-detached test-api
