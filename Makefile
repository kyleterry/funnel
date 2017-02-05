build:
	go generate
	go build

test:
	@./scripts/run-tests

.PHONY: build test
