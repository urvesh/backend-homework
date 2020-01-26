GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=backend-homework

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

test:
	$(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)
