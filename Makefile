.PHONY: build run test clean install deps

# Build the application
build:
	go build -o bin/eol ./cmd/eol

# Run the application
run:
	go run ./cmd/eol

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install the application
install: build
	sudo mv bin/eol /usr/local/bin/

# Development
dev: deps run