all: build

build:
	go build -ldflags="-s -w" -o bin/pb cmd/pb/main.go