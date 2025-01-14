.PHONY: build build-amd64 clean

# Default target
build:
	go build -o days-without-shit

# Build for amd64 Linux
build-amd64:
	GOOS=linux GOARCH=amd64 go build -o days-without-shit-amd64

# Clean builds
clean:
	rm -f days-without-shit days-without-shit-amd64
