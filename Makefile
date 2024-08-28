all: build test

build:
	go build -v ./...

test:
	go test -v ./...

install:
	go install -v ./...

clean:
	go clean -v ./...
	rm -f .cover.html .cover.out

cover:
	go test -v -coverprofile .cover.out ./...
	go tool cover -html .cover.out -o .cover.html
	#open .cover.html