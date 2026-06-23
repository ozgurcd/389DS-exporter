BINARY := 389DS-exporter

.PHONY: all build clean distclean fmt vet lint gosec osv-scanner staticcheck govulncheck test run verify

all: build

build:
	go build -o $(BINARY) .

clean:
	rm -f $(BINARY)

distclean: clean
	rm -f vendor/

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test ./...

run:
	go run .

install:
	go install .

gosec:
	gosec ./...

osv-scanner:
	osv-scanner scan -r go.mod

staticcheck:
	staticcheck ./...

govulncheck:
	govulncheck ./...

verify: lint gosec osv-scanner staticcheck govulncheck build test
