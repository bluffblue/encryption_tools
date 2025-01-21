.PHONY: build run test clean

build:
	go build -o bin/encryption-tool src/main.go

run:
	go run src/main.go

test:
	go test ./...

clean:
	rm -rf bin/
