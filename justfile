test pkg="./...":
    - go test {{pkg}}

run image:
    - go run . --speed 2 {{image}}

lint:
    - golangci-lint \
      --disable=unused \
      run
