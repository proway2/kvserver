fmt:
	go fmt ./...

vet:
	go vet ./...
test:
	go test -v -cover ./...

build:
	go build -v .

client:
	@go build ./tools/client.go

.PHONY: client
