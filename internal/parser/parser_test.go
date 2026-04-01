package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	content := `
VERSION := 1.0.0

## Build the binary
build: clean
	go build -o app .

test: ## Run all tests
	go test ./...

# Deploy to production
deploy:
	./deploy.sh

.PHONY: build test deploy

clean:
	rm -f app
`
	dir := t.TempDir()
	path := filepath.Join(dir, "Makefile")
	os.WriteFile(path, []byte(content), 0644)

	targets, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		name string
		desc string
	}{
		{"build", "Build the binary"},
		{"test", "Run all tests"},
		{"deploy", "Deploy to production"},
		{"clean", ""},
	}

	if len(targets) != len(expected) {
		t.Fatalf("got %d targets, want %d: %+v", len(targets), len(expected), targets)
	}

	for i, e := range expected {
		if targets[i].Name != e.name {
			t.Errorf("target %d: got name %q, want %q", i, targets[i].Name, e.name)
		}
		if targets[i].Description != e.desc {
			t.Errorf("target %d: got desc %q, want %q", i, targets[i].Description, e.desc)
		}
	}
}

func TestParseSkipsVariables(t *testing.T) {
	content := `
VERSION := 1.0.0
CC ?= gcc
LDFLAGS = -s

build:
	go build .
`
	dir := t.TempDir()
	path := filepath.Join(dir, "Makefile")
	os.WriteFile(path, []byte(content), 0644)

	targets, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(targets) != 1 || targets[0].Name != "build" {
		t.Fatalf("expected only 'build', got: %+v", targets)
	}
}
