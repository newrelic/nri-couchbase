WORKDIR         := $(shell pwd)
INTEGRATION     := couchbase
BINARY_NAME      = nri-$(INTEGRATION)
GO_FILES        := ./src/
GOFLAGS          = -mod=readonly

all: build

build: clean test compile

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

test: 
	@echo "=== $(INTEGRATION) === [ test ]: Running unit tests..."
	@go test -race ./... -count=1

compile: 
	@echo "=== $(INTEGRATION) === [ compile ]: Building $(BINARY_NAME)..."
	@go build -o bin/$(BINARY_NAME) $(GO_FILES) 

# Include thematic Makefiles
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk

.PHONY: all build clean compile test 
