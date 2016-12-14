build:
	go generate
	go build

test:
	@./scripts/test.sh
