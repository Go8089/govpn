.PHONY: build test fmt vet run-server run-client clean

build:
	go build ./...
test:
	go test ./...
fmt:
	go fmt ./...
vet:
	go vet ./...
run-server:
	go run ./cmd/server
run-client:
	go run ./cmd/client
clean:
	go clean
