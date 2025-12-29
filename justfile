test pkg="./...":
    - go test {{pkg}}

build:
    - go build -o erc .

run image:
    - go run . run --speed 2 {{image}}

quick image:
    - go run . run --speed 5 {{image}}

debug image:
    - go run . run --speed 2 --debug-image {{image}}

lint:
    - golangci-lint \
      --disable=unused \
      run

# Produce hexdumps that we can use to compute a diff of all the data that
# changed before and after some execution
hexa image:
    - xxd {{image}} {{image}}-a

hexb image:
    - xxd {{image}} {{image}}-b

hexdiff image:
    - diff -u {{image}}-a {{image}}-b > {{image}}.hex.diff

asma image:
    - cp {{image}}.asm {{image}}.asm-a

asmb image:
    - cp {{image}}.asm {{image}}.asm-b

asmdiff image:
    - diff -u {{image}}.asm-a {{image}}.asm-b > {{image}}.asm.diff
