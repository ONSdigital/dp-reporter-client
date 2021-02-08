MAIN=dp-reporter-client

BUILD=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
BIN_DIR?=.

PHONY: all
all: audit test build

PHONY: audit
audit:
	go list -m all | nancy sleuth

PHONY: build
build:
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build $(LDFLAGS) -o $(BUILD_ARCH)/$(BIN_DIR)/$(MAIN)

PHONY: test
test:
	go test -race -cover ./...


.PHONY: build debug test