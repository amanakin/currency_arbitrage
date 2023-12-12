BIN=$(shell realpath ./bin)

.PHONY: build
build:
	GOBIN=$(BIN) go install -v -race ./cmd/history_simulation