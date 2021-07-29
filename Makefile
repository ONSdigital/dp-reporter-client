.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: lint
lint:
	exit
