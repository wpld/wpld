GO           = @go

PROJECTROOT  = $(shell pwd)
D_BIN        = $(PROJECTROOT)/bin

MAKEFLAGS   += --silent

.PHONY: clean build install test

build: clean
	@echo ">  Building the binary..."
	$(GO) build -o $(D_BIN)/wpld $(PROJECTROOT)/main.go
install: clean
	@echo ">  Installing the binary..."
	$(GO) install
clean:
	@echo ">  Cleaning the build cache..."
	$(GO) clean .
test:
	@echo ">  Testing the project..."
	$(GO) test ./...
