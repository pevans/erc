test pkg="./...":
    - go test {{pkg}}

build:
    - go build -o erc .

run image:
    - go run . --speed 2 {{image}}

quick image:
    - go run . --speed 5 {{image}}

debug image:
    - go run . --speed 2 --debug-image {{image}}

lint:
    - golangci-lint \
      --disable=unused \
      run
