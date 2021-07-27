# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR=bin/
BINARY_NAME=prcd
BINARY_MAC=prcd-mac

LOCAL_RUN_FLAG=

all: test build build-linux

build: build-prepare build-linux build-mac

run: build
	$(GOCMD) run . $(LOCAL_RUN_FLAG) &

test:
	$(GOTEST) -v 

clean:
	$(GOCLEAN)
	rm -f $(BINARY_DIR)$(BINARY_NAME)
	rm -f $(BINARY_DIR)$(BINARY_MAC)

build-prepare:
	mkdir -p $(BINARY_DIR)

# Cross compilation
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_NAME) -v

build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DIR)$(BINARY_MAC) -v
