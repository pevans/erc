test:
	if [ "$(PKG)" ]; then go test ./$(PKG); else go test ./...; fi

lint:
	golint ./...

coverage:
	go test -cover ./...

build:
	go build ./cmd/erc

.PHONY: test lint coverage
