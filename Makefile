test:
	go test ./...

lint:
	golint ./...

coverage:
	go test -cover ./...

build:
	go build ./cmd/erc

.PHONY: test lint coverage
