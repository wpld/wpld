GO           = @go

PROJECTROOT  = $(shell pwd)
D_BIN        = $(PROJECTROOT)/bin

MAKEFLAGS   += --silent

.PHONY: clean build

build: clean
	@echo ">  Building binary..."
	$(GO) build -o $(D_BIN)/wpld $(PROJECTROOT)/main.go
clean:
	@echo ">  Cleaning build cache..."
	$(GO) clean .
