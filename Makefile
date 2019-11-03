test:
	./bin/test $(T)

lint:
	./bin/lint $(T)

coverage:
	COVERAGE=1 ./bin/test $(T)

build:
	./bin/build-font --font=assets/fonts/a2s.png --package=font > pkg/font/a2s.go
	./bin/build-font --font=assets/fonts/a2i.png --package=font > pkg/font/a2i.go
	go build ./cmd/erc

.PHONY: test lint coverage
