all: build test

build:
	go build -v ./...

test:
	go test -v ./...

clean:
	go clean -v ./...