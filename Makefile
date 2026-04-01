.PHONY: test lint build clean

test:
	go test -v ./test/...

lint:
	golangci-lint run ./...

build:
	go build -o bin/server ./cmd/server/main.go

clean:
	rm -rf bin/
	go clean -testcache
