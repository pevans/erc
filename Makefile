test:
	./bin/test $(T)

lint:
	./bin/lint $(T)

coverage:
	COVERAGE=1 ./bin/test $(T)

build:
	go build ./cmd/erc

.PHONY: test lint coverage
