LINT = golangci-lint -E golint run

all: test

build:
	./bin/build-font --font=assets/fonts/a2s.png --package=font > pkg/font/a2s.go
	./bin/build-font --font=assets/fonts/a2i.png --package=font > pkg/font/a2i.go
	go build ./cmd/erc

coverage:
	COVERAGE=1 ./bin/test $(T)

lint:
	./bin/analyze "$(LINT)" $(T)

test:
	./bin/test $(T)

.PHONY: build coverage lint test
