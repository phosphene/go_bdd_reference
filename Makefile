.PHONY: test lint build clean coverage

test:
	go test -v ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

build:
	go build -o bin/server ./cmd/server/main.go

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -testcache
