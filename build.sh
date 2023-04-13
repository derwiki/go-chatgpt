#!/bin/bash

# Build for Windows (64-bit)
echo "Building for Windows (64-bit)"
GOOS=windows GOARCH=amd64 go build -o release/chatgpt.exe chatgpt.go

# Build for Linux (64-bit)
echo "Building for Linux (64-bit)"
GOOS=linux GOARCH=amd64 go build -o release/chatgpt-linux chatgpt.go

# Build for macOS (64-bit)
echo "Building for macOS (64-bit)"
GOOS=darwin GOARCH=amd64 go build -o release/chatgpt-macos chatgpt.go

# Build for FreeBSD (64-bit)
echo "Building for FreeBSD (64-bit)"
GOOS=freebsd GOARCH=amd64 go build -o release/chatgpt-freebsd chatgpt.go

# Build for OpenBSD (64-bit)
echo "Building for OpenBSD (64-bit)"
GOOS=openbsd GOARCH=amd64 go build -o release/chatgpt-openbsd chatgpt.go

echo "Done!"

