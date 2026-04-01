VERSION ?= dev

## Build the binary
build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o lazymake .

## Run tests
test:
	go test ./...

## Install to GOPATH/bin
install:
	go install -ldflags "-s -w -X main.version=$(VERSION)" .

## Remove build artifacts
clean:
	rm -f lazymake
	rm -rf dist/

## Release with goreleaser (requires goreleaser)
release:
	goreleaser release --clean
