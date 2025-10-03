test pkg="./...":
    - go test {{pkg}}

run image:
    - go run . --speed 2 {{image}}

quick image:
    - go run . --speed 5 {{image}}

lint:
    - golangci-lint \
      --disable=unused \
      run
