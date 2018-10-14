test:
	go test ./...

lint:
	golint ./...

coverage:
	go test -cover ./...

.PHONY: test lint coverage
