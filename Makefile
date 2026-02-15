BINARY=codemint
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-X github.com/codemint/codemint-cli/cmd.version=$(VERSION) -X github.com/codemint/codemint-cli/cmd.commit=$(COMMIT) -X github.com/codemint/codemint-cli/cmd.date=$(DATE) -X github.com/codemint/codemint-cli/cmd.builtBy=make

.PHONY: build test fmt release-dry-run

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) .

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal ./test

release-dry-run:
	mkdir -p dist
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)
	tar -czf dist/$(BINARY)_$(VERSION)_darwin_arm64.tar.gz -C dist $(BINARY)
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)
	tar -czf dist/$(BINARY)_$(VERSION)_darwin_amd64.tar.gz -C dist $(BINARY)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)
	tar -czf dist/$(BINARY)_$(VERSION)_linux_amd64.tar.gz -C dist $(BINARY)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY)
	tar -czf dist/$(BINARY)_$(VERSION)_linux_arm64.tar.gz -C dist $(BINARY)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/$(BINARY).exe
	(cd dist && zip $(BINARY)_$(VERSION)_windows_amd64.zip $(BINARY).exe)
	rm -f dist/$(BINARY) dist/$(BINARY).exe
	cp scripts/install.sh dist/install.sh
	chmod +x dist/install.sh
	./scripts/checksums.sh dist "$(BINARY)_$(VERSION)_*"
